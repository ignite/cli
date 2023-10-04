package app

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/xast"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
)

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

func FindRegisteredModules(chainRoot string) (modules []string, err error) {
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
	for _, f := range appPkg.Files {
		fileImports := goanalysis.FormatImports(f)
		discovered, err := FindKeepersModules(appPkg, fileImports)
		if err != nil {
			return nil, err
		}
		modules = append(modules, discovered...)
	}

	return modules, nil
}

func FindKeepersModules(n ast.Node, fileImports map[string]string) ([]string, error) {
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

	keeperParamsMap := make(map[string]struct{})
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
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
					if spec.Sel.Name != "Keeper" {
						continue
					}
					ident, ok := spec.X.(*ast.Ident)
					if !ok {
						continue
					}
					fileImport, ok := fileImports[ident.Name]
					if !ok {
						continue
					}
					keeperParamsMap[removeKeeperPkgPath(fileImport)] = struct{}{}
				}
			}
		}
	}

	keeperParams := make([]string, 0)
	for param := range keeperParamsMap {
		keeperParams = append(keeperParams, param)
	}

	return keeperParams, nil
}

func removeKeeperPkgPath(pkg string) string {
	path := strings.TrimSuffix(pkg, "/keeper")
	path = strings.TrimSuffix(path, "/controller")
	return strings.TrimSuffix(path, "/host")
}
