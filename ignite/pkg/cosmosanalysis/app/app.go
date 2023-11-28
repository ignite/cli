package app

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/goenv"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/xast"
)

const registerRoutesMethod = "RegisterAPIRoutes"

// CheckKeeper checks for the existence of the keeper with the provided name in the app structure.
func CheckKeeper(path, keeperName string) error {
	// find app type
	appImpl, err := cosmosanalysis.FindImplementation(path, cosmosanalysis.AppImplementation)
	if err != nil {
		return err
	}
	if len(appImpl) != 1 {
		return fmt.Errorf("app.go should contain a single app (got %d)", len(appImpl))
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
		return fmt.Errorf("app doesn't contain %s", keeperName)
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

	// Loop on package's files
	var blankImports, discovered []string
	for _, f := range appPkg.Files {
		blankImports = append(blankImports, goanalysis.FindBlankImports(f)...)
		fileImports := goanalysis.FormatImports(f)
		d, err := DiscoverModules(f, chainRoot, fileImports)
		if err != nil {
			return nil, err
		}
		discovered = append(discovered, d...)
	}
	return mergeImports(blankImports, discovered), nil
}

// mergeImports merge all discovered imports into the blank imports found in the app files.
func mergeImports(blankImports, discovered []string) []string {
	imports := make([]string, len(blankImports))
	copy(imports, blankImports)
	for i, m := range discovered {
		split := strings.Split(m, "/")

		j := len(split)
		maxTrim := len(split) - 3
	LoopBack:
		for j > maxTrim {
			j--
			// x path means we are reaching the root of the module
			if split[j] == "x" {
				j = maxTrim
				goto LoopBack
			}
			for _, imp := range blankImports {
				// check if the import exist into the blank imports
				if strings.Contains(imp, m) {
					j = -1
					goto LoopBack
				}
			}
			m = strings.TrimSuffix(m, "/"+split[j])
		}
		if j == maxTrim {
			imports = append(imports, discovered[i])
		}
	}
	return imports
}

// DiscoverModules find a map of import modules based on the configured app.
func DiscoverModules(n ast.Node, chainRoot string, fileImports map[string]string) ([]string, error) {
	// find app type
	appImpl := cosmosanalysis.FindImplementationInFile(n, cosmosanalysis.AppImplementation)
	appTypeName := "App"
	switch {
	case len(appImpl) > 1:
		return nil, fmt.Errorf("app.go should contain only a single app (got %d)", len(appImpl))
	case len(appImpl) == 1:
		appTypeName = appImpl[0]
	}

	file, ok := n.(*ast.File)
	if !ok {
		return nil, nil
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

	// Add discovered modules to a list without duplicates
	var (
		modules []string
		skip    = make(map[string]struct{})
	)
	for _, name := range discovered {
		if _, ok := skip[name]; ok {
			continue
		}

		skip[name] = struct{}{}
		modules = append(modules, name)
	}
	return modules, nil
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
		if dep.Path == cosmosver.CosmosModulePath {
			pkg = dep.String()
			break
		}
	}

	if pkg == "" {
		return "", errors.New("Cosmos SDK package version not found")
	}

	// Check path of the package directory within Go's module cache
	path := filepath.Join(goenv.GoModCache(), pkg)
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return "", errors.New("local path to Cosmos SDK package not found")
	}
	return path, nil
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
