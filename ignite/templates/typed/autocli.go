package typed

import (
	"fmt"
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

	optionsToInsert := make([]string, 0, len(options))
	for _, option := range options {
		optionExpr, optionText, parseErr := parseRPCOption(option)
		if parseErr != nil {
			return "", parseErr
		}
		if optionExpr == nil {
			continue
		}

		method, ok := rpcMethod(optionExpr)
		if ok {
			if _, exists := existingRPCMethods[method]; exists {
				continue
			}
			existingRPCMethods[method] = struct{}{}
		}

		optionsToInsert = append(optionsToInsert, optionText)
	}

	if len(optionsToInsert) == 0 {
		return content, nil
	}

	content, err = insertAutoCLIOptions(content, fileSet, rpcCommandOptionsLit, optionsToInsert)
	if err != nil {
		return "", err
	}

	formatted, err := format.Source([]byte(content))
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func parseRPCOption(option string) (ast.Expr, string, error) {
	option = normalizeOption(option)
	if option == "" {
		return nil, "", nil
	}

	code := fmt.Sprintf("package p\nvar _ = []*autocliv1.RpcCommandOptions{\n%s,\n}\n", option)
	file, err := parser.ParseFile(token.NewFileSet(), "", code, 0)
	if err != nil {
		return nil, "", errors.Errorf("failed to parse autocli option expression: %w", err)
	}

	genDecl, ok := file.Decls[0].(*ast.GenDecl)
	if !ok || len(genDecl.Specs) == 0 {
		return nil, "", errors.New("failed to parse autocli option expression: generated declaration is invalid")
	}
	valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok || len(valueSpec.Values) == 0 {
		return nil, "", errors.New("failed to parse autocli option expression: generated value spec is invalid")
	}
	optionsLit, ok := valueSpec.Values[0].(*ast.CompositeLit)
	if !ok || len(optionsLit.Elts) == 0 {
		return nil, "", errors.New("failed to parse autocli option expression: generated options literal is invalid")
	}

	return optionsLit.Elts[0], option, nil
}

func normalizeOption(option string) string {
	option = strings.TrimSpace(option)
	option = strings.TrimSuffix(option, ",")
	return strings.TrimSpace(option)
}

func insertAutoCLIOptions(
	content string,
	fileSet *token.FileSet,
	optionsLiteral *ast.CompositeLit,
	optionsToInsert []string,
) (string, error) {
	file := fileSet.File(optionsLiteral.Rbrace)
	if file == nil {
		return "", errors.New(`failed to find token file for "RpcCommandOptions"`)
	}

	insertOffset := file.Offset(optionsLiteral.Rbrace)
	if insertOffset < 0 || insertOffset > len(content) {
		return "", errors.New(`invalid insertion offset for "RpcCommandOptions"`)
	}

	closingIndentOffset := insertOffset
	for closingIndentOffset > 0 {
		char := content[closingIndentOffset-1]
		if char != '\t' && char != ' ' {
			break
		}
		closingIndentOffset--
	}

	closingIndent := content[closingIndentOffset:insertOffset]
	optionIndent := closingIndent + "\t"

	var insertion strings.Builder
	for _, option := range optionsToInsert {
		insertion.WriteString(indentOption(option, optionIndent))
		insertion.WriteString(",\n")
	}

	return content[:closingIndentOffset] + insertion.String() + content[closingIndentOffset:], nil
}

func indentOption(option, baseIndent string) string {
	lines := strings.Split(option, "\n")
	lines = trimEmptyLines(lines)
	minIndent := minIndentation(lines)

	indented := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmedLine := strings.TrimRight(line, " \t")
		indented = append(indented, baseIndent+removeIndent(trimmedLine, minIndent))
	}

	return strings.Join(indented, "\n")
}

func trimEmptyLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return lines[start:end]
}

func minIndentation(lines []string) int {
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := 0
		for indent < len(line) {
			if line[indent] != ' ' && line[indent] != '\t' {
				break
			}
			indent++
		}

		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent < 0 {
		return 0
	}

	return minIndent
}

func removeIndent(line string, indent int) string {
	if indent <= 0 {
		return line
	}

	i := 0
	for i < len(line) && i < indent {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
		i++
	}

	return line[i:]
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
