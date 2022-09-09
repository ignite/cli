package app

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/ignite/pkg/goanalysis"
)

var appImplementation = []string{
	"RegisterAPIRoutes",
	"RegisterTxService",
	"RegisterTendermintService",
}

// CheckKeeper checks for the existence of the keeper with the provided name in the app structure
func CheckKeeper(path, keeperName string) error {
	// find app type
	appImpl, err := cosmosanalysis.FindImplementation(path, appImplementation)
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

// FindRegisteredModules looks for all the registered modules in the App
// It finds activated modules by checking if imported modules are registered in the app and also checking if their query clients are registered
// It does so by:
// 1. Mapping out all the imports and named imports
// 2. Looking for the call to module.NewBasicManager and finds the modules registered there
// 3. Looking for the implementation of RegisterAPIRoutes and find the modules that call their RegisterGRPCGatewayRoutes
func FindRegisteredModules(chainRoot string) ([]string, error) {
	appFilePath, err := cosmosanalysis.FindAppFilePath(chainRoot)
	if err != nil {
		return nil, err
	}

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, appFilePath, nil, 0)
	if err != nil {
		return []string{}, err
	}

	packages, err := goanalysis.FindImportedPackages(appFilePath)
	if err != nil {
		return nil, err
	}

	basicManagerModule, err := findBasicManagerModule(packages)
	if err != nil {
		return nil, err
	}

	// The directory where the app file is located.
	// This is required to resolve references within the app package.
	appDir := filepath.Dir(appFilePath)

	var basicModules []string
	ast.Inspect(f, func(n ast.Node) bool {
		if pkgsReg := findBasicManagerRegistrations(n, basicManagerModule, appDir); pkgsReg != nil {
			for _, rp := range pkgsReg {
				importModule := packages[rp]
				if importModule == "" {
					// When the package is not defined in the same file use the package name as import
					importModule = rp
				}

				basicModules = append(basicModules, importModule)
			}

			return false
		}

		if pkgsReg := findRegisterAPIRoutersRegistrations(n); pkgsReg != nil {
			for _, rp := range pkgsReg {
				importModule := packages[rp]
				if importModule == "" {
					continue
				}
				basicModules = append(basicModules, importModule)
			}

			return false
		}

		return true
	})

	return basicModules, nil
}

func findBasicManagerRegistrations(n ast.Node, basicManagerModule, pkgDir string) []string {
	callExprType, ok := n.(*ast.CallExpr)
	if !ok {
		return nil
	}

	selectorExprType, ok := callExprType.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	identExprType, ok := selectorExprType.X.(*ast.Ident)
	if !ok || identExprType.Name != basicManagerModule || selectorExprType.Sel.Name != "NewBasicManager" {
		return nil
	}

	// TODO: Move parsing for each type of node to different functions
	var packagesRegistered []string
	for _, arg := range callExprType.Args {
		argAsCompositeLitType, ok := arg.(*ast.CompositeLit)
		if ok {
			compositeTypeSelectorExpr, ok := argAsCompositeLitType.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			compositeTypeX, ok := compositeTypeSelectorExpr.X.(*ast.Ident)
			if ok {
				packagesRegistered = append(packagesRegistered, compositeTypeX.Name)
				continue
			}
		}

		argAsCallType, ok := arg.(*ast.CallExpr)
		if ok {
			argAsFunctionType, ok := argAsCallType.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			argX, ok := argAsFunctionType.X.(*ast.Ident)
			if ok {
				packagesRegistered = append(packagesRegistered, argX.Name)
			}
		}

		// The list of modules are defined in a local variable
		if ident, ok := arg.(*ast.Ident); ok {
			packagesRegistered = append(packagesRegistered, parseAppModulesFromIdent(ident, pkgDir)...)
		}
	}

	return packagesRegistered
}

func findBasicManagerModule(pkgs map[string]string) (string, error) {
	for mod, pkg := range pkgs {
		if pkg == "github.com/cosmos/cosmos-sdk/types/module" {
			return mod, nil
		}
	}

	return "", errors.New("no module for BasicManager was found")
}

func findRegisterAPIRoutersRegistrations(n ast.Node) []string {
	funcLitType, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	if funcLitType.Name.Name != "RegisterAPIRoutes" {
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
