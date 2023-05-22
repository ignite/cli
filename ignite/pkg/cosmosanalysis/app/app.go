package app

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/xast"
)

const (
	appWiringImport     = "cosmossdk.io/depinject"
	appWiringCallMethod = "Inject"
)

var appImplementation = []string{
	"RegisterAPIRoutes",
	"RegisterTxService",
	"RegisterTendermintService",
}

// CheckKeeper checks for the existence of the keeper with the provided name in the app structure.
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
// 3. Looking for the implementation of RegisterAPIRoutes and find the modules that call their RegisterGRPCGatewayRoutes.
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
		err := xast.Inspect(f, func(n ast.Node) error {
			// Find modules in module.NewBasicManager call
			pkgs, err := findBasicManagerRegistrations(n, appDir, fileImports)
			if err != nil {
				return err
			}

			if pkgs != nil {
				for _, p := range pkgs {
					importModule := fileImports[p]
					if importModule == "" {
						// When the package is not defined in the same file use the package name as import
						importModule = p
					}
					modules = append(modules, importModule)
				}
				return xast.ErrStop
			}

			// Find modules in RegisterAPIRoutes declaration
			if pkgs := findRegisterAPIRoutesRegistrations(n); pkgs != nil {
				for _, p := range pkgs {
					importModule := fileImports[p]
					if importModule == "" {
						continue
					}
					modules = append(modules, importModule)
				}
				return xast.ErrStop
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return modules, nil
}

// CheckAppWiring check if the app wiring exists finding the `appconfig.Compose` method call.
func CheckAppWiring(chainRoot string) (bool, error) {
	// Assumption: modules are registered in the app package
	appFilePath, err := cosmosanalysis.FindAppFilePath(chainRoot)
	if err != nil {
		return false, err
	}
	// The directory where the app file is located.
	// This is required to resolve references within the app package.
	appDir := filepath.Dir(appFilePath)

	appPkg, _, err := xast.ParseDir(appDir)
	if err != nil {
		return false, err
	}

	// Loop on package's files
	for _, f := range appPkg.Files {
		exists := goanalysis.FuncVarExists(f, appWiringImport, appWiringCallMethod)
		if exists {
			return true, nil
		}
	}
	return false, nil
}

func exprToString(n ast.Expr) (string, error) {
	buf := bytes.Buffer{}
	fset := token.NewFileSet()

	// Convert the expression node to Go
	if err := format.Node(&buf, fset, n); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newExprError(msg string, n ast.Expr) error {
	s, err := exprToString(n)
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return fmt.Errorf("%s: %s", msg, s)
}

func newUnexpectedTypeErr(n any) error {
	return errors.Errorf("unexpected type %T", n)
}

func findBasicManagerRegistrations(n ast.Node, pkgDir string, fileImports map[string]string) (packages []string, err error) {
	callExprType, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	selectorExprType, ok := callExprType.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	identExprType, ok := selectorExprType.X.(*ast.Ident)
	if !ok {
		return
	}
	basicModulePkgName := findBasicManagerPkgName(fileImports)
	if basicModulePkgName == "" {
		// cosmos-sdk/types/module is not imported in this file, skip
		return
	}
	if identExprType.Name != basicModulePkgName || selectorExprType.Sel.Name != "NewBasicManager" {
		return
	}

	// Node "n" defines the call to NewBasicManager, let's loop on its args to discover modules
	for _, arg := range callExprType.Args {
		switch v := arg.(type) {

		case *ast.CompositeLit:
			// The arg is an app module
			ps, err := parsePkgNameFromCompositeLit(v, pkgDir)
			if err != nil {
				return nil, err
			}
			packages = append(packages, ps...)

		case *ast.CallExpr:
			// The arg is a function call that returns the app module
			ps, err := parsePkgNameFromCall(v, pkgDir)
			if err != nil {
				return nil, err
			}
			packages = append(packages, ps...)

		case *ast.Ident:
			// The list of modules are defined in a local variable
			ps, err := parseAppModulesFromIdent(v, pkgDir)
			if err != nil {
				return nil, err
			}

			packages = append(packages, ps...)
		case *ast.SelectorExpr:
			// The list of modules is defined in a variable of a different package
			ps, err := parseAppModulesFromSelectorExpr(v, pkgDir, fileImports)
			if err != nil {
				return nil, err
			}
			packages = append(packages, ps...)
		default:
			return nil, newExprError("unsupported NewBasicManager() argument format", arg)
		}
	}
	return packages, nil
}

func findBasicManagerPkgName(pkgs map[string]string) string {
	for mod, pkg := range pkgs {
		if pkg == "github.com/cosmos/cosmos-sdk/types/module" {
			return mod
		}
	}
	return ""
}

func findRegisterAPIRoutesRegistrations(n ast.Node) []string {
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
