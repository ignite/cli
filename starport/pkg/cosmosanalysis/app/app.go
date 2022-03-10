package app

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
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
		return errors.New("app.go should contain a single app")
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

// GetRegisteredModules looks for all the registered modules in the App
// It connects the imported go module to the registered module and returns a list of the go modules that are registered
// It does so by:
// 1. Mapping out all the imports and named imports
// 2. Looking for the call to module.NewBasicManager and finds the modules registered there
// 3. Looking for the implementation of RegisterAPIRoutes and find the modules that call their RegisterGRPCGatewayRoutes
func GetRegisteredModules(pathToApp string) ([]string, error) {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, pathToApp, nil, 0)
	if err != nil {
		return []string{}, err
	}

	basicModules := make([]string, 0)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			importMap := extractImports(f)

			ast.Inspect(f, func(n ast.Node) bool {
				if pkgsReg := lookForBasicManagerRegistrations(n); pkgsReg != nil {
					for _, rp := range pkgsReg {
						importModule := importMap[rp]
						basicModules = append(basicModules, importModule)
					}

					return false
				}

				if pkgsReg := lookForRegisterApiRouters(n); pkgsReg != nil {
					for _, rp := range pkgsReg {
						importModule := importMap[rp]
						if importModule == "" {
							continue
						}
						basicModules = append(basicModules, importModule)
					}

					return false
				}

				return true
			})
		}
	}

	return basicModules, nil
}

func extractImports(f *ast.File) map[string]string {
	importMap := make(map[string]string) // name -> import
	ast.Inspect(f, func(n ast.Node) bool {
		// look for struct methods.
		importType, ok := n.(*ast.ImportSpec)
		if ok {
			var importName string
			if importType.Name != nil {
				importName = importType.Name.Name
			} else {
				importParts := strings.Split(importType.Path.Value, "/")
				importName = importParts[len(importParts)-1]
			}

			importMap[strings.Trim(importName, "\"")] = strings.Trim(importType.Path.Value, "\"")
		}

		return true
	})

	return importMap
}

func lookForBasicManagerRegistrations(n ast.Node) []string {
	callExprType, ok := n.(*ast.CallExpr)
	if ok {
		selectorExprType, ok := callExprType.Fun.(*ast.SelectorExpr)
		if ok {
			identExprType, ok := selectorExprType.X.(*ast.Ident)
			if ok && identExprType.Name == "module" && selectorExprType.Sel.Name == "NewBasicManager" {
				packagesRegistered := make([]string, len(callExprType.Args))

				for i, arg := range callExprType.Args {
					argAsCompositeLitType, ok := arg.(*ast.CompositeLit)
					if ok {
						compositeTypeSelectorExpr, ok := argAsCompositeLitType.Type.(*ast.SelectorExpr)
						if ok {
							compositeTypeX, ok := compositeTypeSelectorExpr.X.(*ast.Ident)
							if ok {
								packagesRegistered[i] = compositeTypeX.Name
								continue
							}
						}
					}

					argAsCallType, ok := arg.(*ast.CallExpr)
					if ok {
						argAsFunctionType, ok := argAsCallType.Fun.(*ast.SelectorExpr)
						if ok {
							argX, ok := argAsFunctionType.X.(*ast.Ident)
							if ok {
								packagesRegistered[i] = argX.Name
							}
						}
					}
				}

				return packagesRegistered
			}
		}
	}

	return nil
}

func lookForRegisterApiRouters(n ast.Node) []string {
	funcLitType, ok := n.(*ast.FuncDecl)
	if ok && funcLitType.Name.Name == "RegisterAPIRoutes" {
		packagesRegistered := make([]string, 0)

		for _, stmt := range funcLitType.Body.List {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if ok {
				exprCall, ok := exprStmt.X.(*ast.CallExpr)
				if ok {
					exprFun, ok := exprCall.Fun.(*ast.SelectorExpr)
					if ok && exprFun.Sel.Name == "RegisterGRPCGatewayRoutes" {
						identType, ok := exprFun.X.(*ast.Ident)
						if ok {
							pkgName := identType.Name
							if pkgName == "" {
								continue
							}
							packagesRegistered = append(packagesRegistered, identType.Name)
						}
					}
				}
			}
		}

		return packagesRegistered
	}

	return nil
}
