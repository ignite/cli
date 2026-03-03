package typed

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

const (
	autoCLIServiceQuery = "Query"
	autoCLIServiceTx    = "Tx"
)

// AppendAutoCLIQueryOptions appends options to the Query RpcCommandOptions in AutoCLIOptions.
func AppendAutoCLIQueryOptions(content string, options ...string) (string, error) {
	return appendAutoCLIOptions(content, autoCLIServiceQuery, options...)
}

// AppendAutoCLITxOptions appends options to the Tx RpcCommandOptions in AutoCLIOptions.
func AppendAutoCLITxOptions(content string, options ...string) (string, error) {
	return appendAutoCLIOptions(content, autoCLIServiceTx, options...)
}

func appendAutoCLIOptions(content, service string, options ...string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}
	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	autoCLIOptionsFunc := findFunctionByName(file, "AutoCLIOptions")
	if autoCLIOptionsFunc == nil {
		return "", errors.New(`function "AutoCLIOptions" not found`)
	}

	moduleOptionsLit, err := findModuleOptionsLiteral(autoCLIOptionsFunc)
	if err != nil {
		return "", err
	}

	serviceDescriptorField, found := findCompositeField(moduleOptionsLit, service)
	if !found {
		return "", errors.Errorf("field %q not found in ModuleOptions", service)
	}
	serviceDescriptorLit, found := resolveCompositeLiteral(serviceDescriptorField.Value)
	if !found {
		return "", errors.Errorf("field %q is not a composite literal in ModuleOptions", service)
	}

	rpcCommandOptionsField, found := findCompositeField(serviceDescriptorLit, "RpcCommandOptions")
	if !found {
		return "", errors.Errorf(`field "RpcCommandOptions" not found in %q service descriptor`, service)
	}
	rpcCommandOptionsLit, found := resolveCompositeLiteral(rpcCommandOptionsField.Value)
	if !found {
		return "", errors.Errorf(`field "RpcCommandOptions" in %q service descriptor is not a composite literal`, service)
	}

	existingRPCMethods := map[string]struct{}{}
	for _, elt := range rpcCommandOptionsLit.Elts {
		method, ok := rpcMethod(elt)
		if !ok {
			continue
		}
		existingRPCMethods[method] = struct{}{}
	}

	for _, option := range options {
		if strings.TrimSpace(option) == "" {
			continue
		}

		optionExpr, parseErr := parser.ParseExpr(option)
		if parseErr != nil {
			return "", errors.Errorf("failed to parse autocli option expression: %w", parseErr)
		}

		method, ok := rpcMethod(optionExpr)
		if ok {
			if _, exists := existingRPCMethods[method]; exists {
				continue
			}
			existingRPCMethods[method] = struct{}{}
		}

		rpcCommandOptionsLit.Elts = append(rpcCommandOptionsLit.Elts, optionExpr)
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

func findModuleOptionsLiteral(autoCLIOptionsFunc *ast.FuncDecl) (*ast.CompositeLit, error) {
	if autoCLIOptionsFunc.Body == nil {
		return nil, errors.New(`function "AutoCLIOptions" has no body`)
	}

	for _, stmt := range autoCLIOptionsFunc.Body.List {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(returnStmt.Results) != 1 {
			continue
		}

		moduleOptionsLit, found := resolveCompositeLiteral(returnStmt.Results[0])
		if !found {
			continue
		}

		if !isModuleOptionsLiteral(moduleOptionsLit) {
			continue
		}

		return moduleOptionsLit, nil
	}

	return nil, errors.New(`return statement with "autocliv1.ModuleOptions" literal not found in "AutoCLIOptions"`)
}

func isModuleOptionsLiteral(moduleOptionsLit *ast.CompositeLit) bool {
	selector, ok := moduleOptionsLit.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return pkgIdent.Name == "autocliv1" && selector.Sel.Name == "ModuleOptions"
}

func resolveCompositeLiteral(expr ast.Expr) (*ast.CompositeLit, bool) {
	switch typedExpr := expr.(type) {
	case *ast.CompositeLit:
		return typedExpr, true
	case *ast.UnaryExpr:
		if typedExpr.Op == token.AND {
			return resolveCompositeLiteral(typedExpr.X)
		}
	case *ast.ParenExpr:
		return resolveCompositeLiteral(typedExpr.X)
	}

	return nil, false
}

func findCompositeField(compLit *ast.CompositeLit, fieldName string) (*ast.KeyValueExpr, bool) {
	for _, elt := range compLit.Elts {
		keyValue, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyIdent, ok := keyValue.Key.(*ast.Ident)
		if !ok || keyIdent.Name != fieldName {
			continue
		}

		return keyValue, true
	}

	return nil, false
}

func rpcMethod(expr ast.Expr) (string, bool) {
	commandLit, found := resolveCompositeLiteral(expr)
	if !found {
		return "", false
	}

	rpcMethodField, found := findCompositeField(commandLit, "RpcMethod")
	if !found {
		return "", false
	}

	rpcMethodValue, ok := rpcMethodField.Value.(*ast.BasicLit)
	if !ok || rpcMethodValue.Kind != token.STRING {
		return "", false
	}

	method, err := strconv.Unquote(rpcMethodValue.Value)
	if err != nil {
		return "", false
	}

	return method, true
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
