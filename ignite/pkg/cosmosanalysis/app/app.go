package app

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

const registerRoutesMethod = "RegisterAPIRoutes"

// CheckKeeper checks for the existence of the keeper with the provided name in the app structure.
func CheckKeeper(path, keeperName string) error {
	// find app type
	appImpl, err := cosmosanalysis.FindEmbed(path, cosmosanalysis.AppEmbeddedTypes)
	if err != nil {
		return err
	}
	if len(appImpl) != 1 {
		return errors.Errorf("app.go should contain a single app (got %d)", len(appImpl))
	}
	appTypeName := appImpl[0]

	// Inspect the module for app struct
	var found bool
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, path, nil, 0)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				// look for struct methods.
				appType, ok := n.(*ast.TypeSpec)
				if !ok || appType.Name.Name != appTypeName {
					return true
				}

				appStruct, ok := appType.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// Search for the keeper specific field
				for _, field := range appStruct.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == keeperName {
							found = true
							return false
						}
					}
				}

				return false
			})
		}
	}

	if !found {
		return errors.Errorf("app doesn't contain %s", keeperName)
	}
	return nil
}

// FindRegisteredModules returns all registered modules into the chain root.
func FindRegisteredModules(chainRoot string) ([]string, error) {
	// Assumption: modules are registered in the app package
	appFilePath, err := cosmosanalysis.FindAppFilePath(chainRoot)
	if err != nil {
		return nil, err
	}
	// The directory where the app file is located.
	// This is required to resolve references within the app package.
	appDir := filepath.Dir(appFilePath)

	appPkg, _, err := xast.ParseDir(appDir)
	if err != nil {
		return nil, err
	}

	// Search the app for the imported SDK modules
	var discovered []string
	for _, f := range appPkg.Files {
		discovered = append(discovered, goanalysis.FindBlankImports(f)...)
		fileImports := goanalysis.FormatImports(f)
		d, err := DiscoverModules(f, chainRoot, fileImports)
		if err != nil {
			return nil, err
		}
		discovered = append(discovered, d...)
	}

	// Discover IBC wired modules
	// TODO: This can be removed once IBC modules use dependency injection
	ibcPath := filepath.Join(chainRoot, "app", "ibc.go")
	if _, err := os.Stat(ibcPath); err == nil {
		m, err := discoverIBCModules(ibcPath)
		if err != nil {
			return nil, err
		}

		discovered = append(discovered, m...)
	}

	return removeDuplicateEntries(discovered), nil
}

// DiscoverModules find a map of import modules based on the configured app.
func DiscoverModules(file *ast.File, chainRoot string, fileImports map[string]string) ([]string, error) {
	// find app type
	appImpl := cosmosanalysis.FindEmbedInFile(file, cosmosanalysis.AppEmbeddedTypes)
	appTypeName := "App"
	switch {
	case len(appImpl) > 1:
		return nil, errors.Errorf("app.go should contain only a single app (got %d)", len(appImpl))
	case len(appImpl) == 1:
		appTypeName = appImpl[0]
	}

	var discovered []string
	for _, decl := range file.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			discovered = append(discovered, discoverKeeperModules(x, appTypeName, fileImports)...)
		case *ast.FuncDecl:
			// The modules registered by Cosmos SDK `rumtime.App` are included
			// when the app registers API modules though the `App` instance.
			if isRuntimeAppCalled(x) {
				m, err := discoverRuntimeAppModules(chainRoot)
				if err != nil {
					return nil, err
				}
				discovered = append(discovered, m...)
			}
		}
	}
	return removeDuplicateEntries(discovered), nil
}

func removeDuplicateEntries(entries []string) (res []string) {
	seen := make(map[string]struct{})
	for _, e := range entries {
		if _, ok := seen[e]; ok {
			continue
		}

		seen[e] = struct{}{}
		res = append(res, e)
	}
	return
}

func discoverKeeperModules(d *ast.GenDecl, appTypeName string, imports map[string]string) []string {
	var modules []string
	for _, spec := range d.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		if typeSpec.Name.Name != appTypeName {
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		for _, field := range structType.Fields.List {
			f := field.Type
		CheckSpec:
			switch spec := f.(type) {
			case *ast.StarExpr:
				f, ok = spec.X.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				goto CheckSpec
			case *ast.SelectorExpr:
				if !strings.HasSuffix(spec.Sel.Name, "Keeper") {
					continue
				}
				ident, ok := spec.X.(*ast.Ident)
				if !ok {
					continue
				}
				fileImport, ok := imports[ident.Name]
				if !ok {
					continue
				}
				modules = append(modules, removeKeeperPkgPath(fileImport))
			}
		}
	}
	return modules
}

func discoverRuntimeAppModules(chainRoot string) ([]string, error) {
	// Resolve the absolute path to the Cosmos SDK module
	cosmosPath, err := resolveCosmosPackagePath(chainRoot)
	if err != nil {
		return nil, err
	}

	var modules []string

	// When runtime package doesn't exists it means is an older Cosmos SDK version,
	// so all the module API registrations are defined within user's app.
	path := filepath.Join(cosmosPath, "runtime", "app.go")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return modules, nil
	}

	f, _, err := xast.ParseFile(path)
	if err != nil {
		return nil, err
	}

	imports := goanalysis.FormatImports(f)
	err = xast.Inspect(f, func(n ast.Node) error {
		if pkgs := findRegisterAPIRoutesRegistrations(n); pkgs != nil {
			for _, p := range pkgs {
				if m := imports[p]; m != "" {
					modules = append(modules, m)
				}
			}
			return xast.ErrStop
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func discoverIBCModules(ibcPath string) ([]string, error) {
	f, _, err := xast.ParseFile(ibcPath)
	if err != nil {
		return nil, err
	}

	var (
		names   []string
		imports = goanalysis.FormatImports(f)
	)
	err = xast.Inspect(f, func(n ast.Node) error {
		fn, _ := n.(*ast.FuncDecl)
		if fn == nil {
			return nil
		}

		if fn.Name.Name != "RegisterIBC" && fn.Name.Name != "AddIBCModuleManager" {
			return nil
		}

		for _, stmt := range fn.Body.List {
			x, _ := stmt.(*ast.AssignStmt)
			if x == nil {
				continue
			}

			if len(x.Rhs) == 0 {
				continue
			}

			c, _ := x.Rhs[0].(*ast.CompositeLit)
			if c == nil {
				continue
			}

			s, _ := c.Type.(*ast.SelectorExpr)
			if s == nil || s.Sel.Name != "AppModule" {
				continue
			}

			if m, _ := s.X.(*ast.Ident); m != nil {
				names = append(names, m.Name)
			}
		}

		return xast.ErrStop
	})
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, n := range names {
		modules = append(modules, imports[n])
	}
	return modules, nil
}

func resolveCosmosPackagePath(chainRoot string) (string, error) {
	modFile, err := gomodule.ParseAt(chainRoot)
	if err != nil {
		return "", err
	}

	deps, err := gomodule.ResolveDependencies(modFile, false)
	if err != nil {
		return "", err
	}

	var pkg string
	for _, dep := range deps {
		// dependencies are resolved, so we need to check for possible SDK forks
		if cosmosver.CosmosSDKModulePathPattern.MatchString(dep.Path) {
			pkg = dep.String()
			break
		}
	}

	if pkg == "" {
		return "", errors.New("cosmos-sdk package version not found")
	}

	m, err := gomodule.FindModule(context.Background(), chainRoot, pkg)
	if err != nil {
		return "", err
	}
	return m.Dir, nil
}

func findRegisterAPIRoutesRegistrations(n ast.Node) []string {
	funcLitType, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	if funcLitType.Name.Name != registerRoutesMethod {
		return nil
	}

	var packagesRegistered []string
	for _, stmt := range funcLitType.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		exprCall, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		exprFun, ok := exprCall.Fun.(*ast.SelectorExpr)
		if !ok || exprFun.Sel.Name != "RegisterGRPCGatewayRoutes" {
			continue
		}
		identType, ok := exprFun.X.(*ast.Ident)
		if !ok {
			continue
		}
		pkgName := identType.Name
		if pkgName == "" {
			continue
		}
		packagesRegistered = append(packagesRegistered, identType.Name)
	}
	return packagesRegistered
}

func removeKeeperPkgPath(pkg string) string {
	path := strings.TrimSuffix(pkg, "/keeper")
	path = strings.TrimSuffix(path, "/controller")
	return strings.TrimSuffix(path, "/host")
}

func isRuntimeAppCalled(fn *ast.FuncDecl) bool {
	if fn.Name.Name != registerRoutesMethod {
		return false
	}

	for _, stmt := range fn.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		exprCall, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		exprFun, ok := exprCall.Fun.(*ast.SelectorExpr)
		if !ok || exprFun.Sel.Name != registerRoutesMethod {
			continue
		}

		exprSel, ok := exprFun.X.(*ast.SelectorExpr)
		if !ok || exprSel.Sel.Name != "App" {
			continue
		}

		identType, ok := exprSel.X.(*ast.Ident)
		if !ok || identType.Name != "app" {
			continue
		}

		return true
	}

	return false
}
