package modulecreate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
)

func addIBCModuleRoute(content, moduleName string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}
	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	registerIBCModules := findFunctionByName(file, "registerIBCModules")
	if registerIBCModules == nil {
		return "", errors.New(`function "registerIBCModules" not found`)
	}

	moduleNameExprText := fmt.Sprintf("%smoduletypes.ModuleName", moduleName)
	moduleConstructorExprText := fmt.Sprintf(
		"%smodule.NewIBCModule(app.appCodec, app.%sKeeper)",
		moduleName,
		xstrings.Title(moduleName),
	)

	hasRoute, err := hasAddRoute(registerIBCModules, moduleNameExprText, fileSet)
	if err != nil {
		return "", err
	}
	if !hasRoute {
		insertionIndex := -1
		for i, stmt := range registerIBCModules.Body.List {
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if !ok || len(assignStmt.Lhs) != 1 {
				continue
			}

			lhs, ok := assignStmt.Lhs[0].(*ast.Ident)
			if !ok || lhs.Name != "ibcv2Router" {
				continue
			}

			insertionIndex = i
			break
		}
		if insertionIndex == -1 {
			return "", errors.New(`assignment to "ibcv2Router" not found`)
		}

		stmtCode := fmt.Sprintf(
			"ibcRouter = ibcRouter.AddRoute(%s, %s)",
			moduleNameExprText,
			moduleConstructorExprText,
		)
		statements, err := parseStatements(stmtCode)
		if err != nil {
			return "", err
		}
		if len(statements) != 1 {
			return "", errors.New("unexpected number of statements while creating ibc route assignment")
		}

		registerIBCModules.Body.List = append(
			registerIBCModules.Body.List[:insertionIndex],
			append([]ast.Stmt{statements[0]}, registerIBCModules.Body.List[insertionIndex:]...)...,
		)
	}

	file.Comments = commentMap.Filter(file).Comments()

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, file); err != nil {
		return "", err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func hasAddRoute(funcDecl *ast.FuncDecl, moduleNameExpr string, fileSet *token.FileSet) (bool, error) {
	var (
		found bool
		err   error
	)

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		if found || err != nil {
			return false
		}

		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "AddRoute" || len(callExpr.Args) == 0 {
			return true
		}

		argText, argErr := exprString(fileSet, callExpr.Args[0])
		if argErr != nil {
			err = argErr
			return false
		}

		if normalizedExpr(argText) == normalizedExpr(moduleNameExpr) {
			found = true
			return false
		}

		return true
	})

	return found, err
}

func parseStatements(code string) ([]ast.Stmt, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", fmt.Sprintf("package p\nfunc _(){\n%s\n}", code), 0)
	if err != nil {
		return nil, err
	}

	funcDecl, ok := file.Decls[0].(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil, errors.New("failed to parse statements")
	}

	return funcDecl.Body.List, nil
}

func findFunctionByName(file *ast.File, funcName string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == funcName {
			return funcDecl
		}
	}
	return nil
}
