package xast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// AppendAutoCLIRPCCommand appends a new RPC command option to either Query or Tx
// in the AutoCLIOptions function.
func AppendAutoCLIRPCCommand(fileContent, serviceField, commandExpr string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, file, file.Comments)

	funcDecl := findFuncDecl(file, "AutoCLIOptions")
	if funcDecl == nil {
		return "", errors.Errorf(`function "AutoCLIOptions" not found`)
	}

	returnStmt, moduleOptions, err := findAutoCLIModuleOptions(funcDecl)
	if err != nil {
		return "", err
	}

	serviceDescriptorExpr, err := findCompositeLitFieldValue(moduleOptions, serviceField)
	if err != nil {
		return "", err
	}

	serviceDescriptor, ok := asCompositeLit(serviceDescriptorExpr)
	if !ok {
		return "", errors.Errorf("field %q is not a struct literal", serviceField)
	}

	rpcCommandOptionsExpr, err := findCompositeLitFieldValue(serviceDescriptor, "RpcCommandOptions")
	if err != nil {
		return "", err
	}

	newRPCCommand, err := parser.ParseExprFrom(fileSet, "", []byte(commandExpr), parser.AllErrors)
	if err != nil {
		return "", errors.Errorf("failed to parse RPC command expression: %w", err)
	}

	switch expr := rpcCommandOptionsExpr.(type) {
	case *ast.CompositeLit:
		expr.Elts = append(expr.Elts, newRPCCommand)
	case *ast.Ident:
		appendStmt, err := parseSingleStmt(
			fileSet,
			fmt.Sprintf("%[1]s = append(%[1]s, %s)", expr.Name, commandExpr),
		)
		if err != nil {
			return "", err
		}
		if err := insertStmtBeforeReturn(funcDecl, returnStmt, appendStmt); err != nil {
			return "", err
		}
	default:
		return "", errors.Errorf("field %q is not a slice literal or identifier", "RpcCommandOptions")
	}

	file.Comments = cmap.Filter(file).Comments()
	return formatNode(fileSet, file)
}

func findAutoCLIModuleOptions(funcDecl *ast.FuncDecl) (*ast.ReturnStmt, *ast.CompositeLit, error) {
	if funcDecl.Body == nil {
		return nil, nil, errors.New(`function "AutoCLIOptions" has no body`)
	}

	for _, stmt := range funcDecl.Body.List {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(returnStmt.Results) == 0 {
			continue
		}

		moduleOptions, ok := asCompositeLit(returnStmt.Results[0])
		if !ok {
			continue
		}

		return returnStmt, moduleOptions, nil
	}

	return nil, nil, errors.New(`return statement with module options not found in "AutoCLIOptions"`)
}

func findCompositeLitFieldValue(lit *ast.CompositeLit, fieldName string) (ast.Expr, error) {
	for _, elt := range lit.Elts {
		keyValueExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := keyValueExpr.Key.(*ast.Ident)
		if !ok || key.Name != fieldName {
			continue
		}

		return keyValueExpr.Value, nil
	}

	return nil, errors.Errorf("field %q not found", fieldName)
}

func asCompositeLit(expr ast.Expr) (*ast.CompositeLit, bool) {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		return expr, true
	case *ast.UnaryExpr:
		if expr.Op != token.AND {
			return nil, false
		}
		compositeLit, ok := expr.X.(*ast.CompositeLit)
		if !ok {
			return nil, false
		}
		return compositeLit, true
	default:
		return nil, false
	}
}

func parseSingleStmt(fileSet *token.FileSet, code string) (ast.Stmt, error) {
	body, err := codeToBlockStmt(fileSet, code)
	if err != nil {
		return nil, err
	}
	if len(body.List) != 1 {
		return nil, errors.New("expected a single statement")
	}

	return body.List[0], nil
}

func insertStmtBeforeReturn(funcDecl *ast.FuncDecl, returnStmt *ast.ReturnStmt, stmt ast.Stmt) error {
	for i := range funcDecl.Body.List {
		if funcDecl.Body.List[i] != returnStmt {
			continue
		}
		funcDecl.Body.List = append(funcDecl.Body.List[:i], append([]ast.Stmt{stmt}, funcDecl.Body.List[i:]...)...)
		return nil
	}

	return errors.New("return statement not found in function body")
}
