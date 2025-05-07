package xast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type (
	// functionOpts represent the options for functions.
	functionOpts struct {
		funcName       string          // Name of the function to modify.
		newParams      []functionParam // Parameters to add to the function.
		body           string          // New function body content.
		newLines       []functionLine  // Lines to insert at specific positions.
		insideCall     functionCalls   // Function calls to modify.
		insideStruct   functionStructs // Struct literals to modify.
		appendTestCase []string        // Test cases to append.
		appendCode     []string        // Code to append at the end.
		returnVars     []string        // Return variables to modify.
	}

	// FunctionOptions configures code generation.
	FunctionOptions func(*functionOpts)

	// functionStruct represents a struct literal to modify.
	functionStruct struct {
		name  string // Name of the struct type.
		param string // Name of the struct field.
		code  string // Code to insert.
		index int    // Position to insert at.
	}
	functionStructs    []functionStruct
	functionStructsMap map[string]functionStructs

	// functionCall represents a function call to modify.
	functionCall struct {
		name  string // Name of the function.
		code  string // Code to insert.
		index int    // Position to insert at.
	}
	functionCalls    []functionCall
	functionCallsMap map[string]functionCalls

	// functionParam represents a parameter to add to a function.
	functionParam struct {
		name    string // Parameter name.
		varType string // Parameter type.
		index   int    // Position to insert at.
	}

	// functionLine represents a line of code to insert.
	functionLine struct {
		code   string // Code to insert.
		number uint64 // Line number to insert at.
	}
)

// Field creates an AST field node from the function parameter.
func (p functionParam) Field() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(p.name)},
		Type:  ast.NewIdent(p.varType),
	}
}

// Map converts a slice of functionStructs to a map keyed by struct name.
func (s functionStructs) Map() functionStructsMap {
	structMap := make(functionStructsMap)
	for _, c := range s {
		structs, ok := structMap[c.name]
		if !ok {
			structs = make(functionStructs, 0)
		}
		structMap[c.name] = append(structs, c)
	}
	return structMap
}

// Map converts a slice of functionCalls to a map keyed by function name.
func (c functionCalls) Map() functionCallsMap {
	callMap := make(functionCallsMap)
	for _, c := range c {
		calls, ok := callMap[c.name]
		if !ok {
			calls = make(functionCalls, 0)
		}
		callMap[c.name] = append(calls, c)
	}
	return callMap
}

// AppendFuncParams adds a new parameter to a function.
func AppendFuncParams(name, varType string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.newParams = append(c.newParams, functionParam{
			name:    name,
			varType: varType,
			index:   index,
		})
	}
}

// ReplaceFuncBody replaces the entire body of a function, the method will replace first and apply the other options after.
func ReplaceFuncBody(body string) FunctionOptions {
	return func(c *functionOpts) {
		c.body = body
	}
}

// AppendFuncTestCase adds a test case to a test function, if exists, of a function in Go source code content.
func AppendFuncTestCase(testCase string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendTestCase = append(c.appendTestCase, testCase)
	}
}

// AppendFuncCode adds code to the end of a function, if exists, of a function in Go source code content.
func AppendFuncCode(code string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendCode = append(c.appendCode, code)
	}
}

// AppendFuncCodeAtLine inserts code at a specific line number.
var AppendFuncCodeAtLine = AppendFuncAtLine

// AppendFuncAtLine inserts code at a specific line number.
func AppendFuncAtLine(code string, lineNumber uint64) FunctionOptions {
	return func(c *functionOpts) {
		c.newLines = append(c.newLines, functionLine{
			code:   code,
			number: lineNumber,
		})
	}
}

// AppendInsideFuncCall adds an argument to a function call. For instances, the method have a parameter a
// // call 'New(param1, param2)' and we want to add the param3 the result will be 'New(param1, param2, param3)'.
func AppendInsideFuncCall(callName, code string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.insideCall = append(c.insideCall, functionCall{
			name:  callName,
			code:  code,
			index: index,
		})
	}
}

// AppendFuncStruct adds a field to a struct literal. For instance,
// // the struct has only one parameter 'Params{Param1: param1}' and we want to add
// // the param2 the result will be 'Params{Param1: param1, Param2: param2}'.
func AppendFuncStruct(name, param, code string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.insideStruct = append(c.insideStruct, functionStruct{
			name:  name,
			param: param,
			code:  code,
			index: index,
		})
	}
}

// NewFuncReturn replaces return statements in a function.
func NewFuncReturn(returnVars ...string) FunctionOptions {
	return func(c *functionOpts) {
		c.returnVars = append(c.returnVars, returnVars...)
	}
}

// newFunctionOptions creates a new functionOpts with defaults.
func newFunctionOptions(funcName string) functionOpts {
	return functionOpts{
		funcName:       funcName,
		newParams:      make([]functionParam, 0),
		body:           "",
		newLines:       make([]functionLine, 0),
		insideCall:     make(functionCalls, 0),
		insideStruct:   make(functionStructs, 0),
		appendTestCase: make([]string, 0),
		appendCode:     make([]string, 0),
		returnVars:     make([]string, 0),
	}
}

// ModifyFunction modifies a function in Go source code using functional options.
func ModifyFunction(content string, funcName string, functions ...FunctionOptions) (string, error) {
	// Collect all function options.
	opts := newFunctionOptions(funcName)
	for _, fn := range functions {
		fn(&opts)
	}

	// Parse source into AST
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return "", errors.Errorf("failed to parse file: %w", err)
	}

	cmap := ast.NewCommentMap(fset, file, file.Comments)

	// Find the target function.
	funcDecl := findFuncDecl(file, funcName)
	if funcDecl == nil {
		return "", errors.Errorf("function %q not found", funcName)
	}

	// Apply modifications.
	if err := applyFunctionOptions(fset, funcDecl, &opts); err != nil {
		return "", err
	}

	file.Comments = cmap.Filter(file).Comments()
	return formatNode(fset, file)
}

// findFuncDecl finds a function declaration in an AST by name.
func findFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return fd
		}
	}
	return nil
}

// addCode converts string code snippets into AST statements.
func addCode(fileSet *token.FileSet, appendCode []string) ([]ast.Stmt, error) {
	code := make([]ast.Stmt, 0)
	for _, codeToInsert := range appendCode {
		body, err := codeToBlockStmt(fileSet, codeToInsert)
		if err != nil {
			return nil, err
		}
		code = append(code, body.List...)
	}
	return code, nil
}

// modifyReturnVars converts return variable strings into AST expressions.
func modifyReturnVars(fileSet *token.FileSet, returnVars []string) ([]ast.Expr, error) {
	stmts := make([]ast.Expr, 0)
	for _, returnVar := range returnVars {
		newRetExpr, err := parser.ParseExprFrom(fileSet, "", []byte(returnVar), parser.ParseComments)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, newRetExpr)
	}
	return stmts, nil
}

// addParams adds new parameters to a function declaration.
func addParams(funcDecl *ast.FuncDecl, newParams []functionParam) error {
	for _, p := range newParams {
		switch {
		case p.index == -1:
			// Append at end
			funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, p.Field())
		case p.index >= 0 && p.index <= len(funcDecl.Type.Params.List):
			// Insert at index
			funcDecl.Type.Params.List = append(
				funcDecl.Type.Params.List[:p.index],
				append([]*ast.Field{p.Field()}, funcDecl.Type.Params.List[p.index:]...)...,
			)
		default:
			return errors.Errorf("params index %d out of range", p.index)
		}
	}
	return nil
}

// addNewLine inserts code at specific line numbers in a function body.
func addNewLine(fileSet *token.FileSet, funcDecl *ast.FuncDecl, newLines []functionLine) error {
	for _, newLine := range newLines {
		// Validate line number
		if newLine.number > uint64(len(funcDecl.Body.List))-1 {
			return errors.Errorf("line number %d out of range", newLine.number)
		}

		// Parse insertion code
		insertionExpr, err := codeToBlockStmt(fileSet, newLine.code)
		if err != nil {
			return err
		}

		// Insert code at the specified line number.
		funcDecl.Body.List = append(
			funcDecl.Body.List[:newLine.number],
			append(insertionExpr.List, funcDecl.Body.List[newLine.number:]...)...,
		)
	}
	return nil
}

// modifyReturn handles return statement modifications and code appending.
func modifyReturn(funcDecl *ast.FuncDecl, returnStmts []ast.Expr, appendCode []ast.Stmt) error {
	if len(funcDecl.Body.List) > 0 {
		lastStmt := funcDecl.Body.List[len(funcDecl.Body.List)-1]
		switch stmt := lastStmt.(type) {
		case *ast.ReturnStmt:
			// Modify the return statement
			if len(returnStmts) > 0 {
				stmt.Results = returnStmts
			}
			if len(appendCode) > 0 {
				// Insert before return
				appendCode = append(appendCode, stmt)
				funcDecl.Body.List = append(funcDecl.Body.List[:len(funcDecl.Body.List)-1], appendCode...)
			}
		default:
			if len(returnStmts) > 0 {
				return errors.New("return statement not found")
			}
			// Append at end
			if len(appendCode) > 0 {
				funcDecl.Body.List = append(funcDecl.Body.List, appendCode...)
			}
		}
	} else {
		if len(returnStmts) > 0 {
			return errors.New("return statement not found")
		}
		// Append to empty body
		if len(appendCode) > 0 {
			funcDecl.Body.List = append(funcDecl.Body.List, appendCode...)
		}
	}
	return nil
}

// addTestCase adds test cases to a test function.
func addTestCase(fSet *token.FileSet, funcDecl *ast.FuncDecl, testCase []string) error {
	// Find tests variable
	for _, stmt := range funcDecl.Body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok || len(assignStmt.Lhs) == 0 {
			continue
		}

		// Check for "tests" variable
		ident, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok || ident.Name != "tests" {
			continue
		}

		// Get composite literal
		compositeLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		// Add test cases
		for _, tc := range testCase {
			testCaseStmt, err := structToBlockStmt(fSet, tc)
			if err != nil {
				return err
			}
			compositeLit.Elts = append(compositeLit.Elts, testCaseStmt)
		}
	}
	return nil
}

// structToBlockStmt parses struct literal code into AST expression.
func structToBlockStmt(fSet *token.FileSet, code string) (ast.Expr, error) {
	newFuncContent := toStruct(code)
	newContent, err := parser.ParseExprFrom(fSet, "temp.go", newFuncContent, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	newCompositeList, ok := newContent.(*ast.CompositeLit)
	if !ok {
		return nil, errors.New("not a composite literal")
	}

	if len(newCompositeList.Elts) != 1 {
		return nil, errors.New("composite literal has more than one element or zero")
	}

	return newCompositeList.Elts[0], nil
}

// toStruct wraps code in an anonymous struct literal for parsing.
func toStruct(code string) string {
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "\n\t", "\n")
	code = strings.ReplaceAll(code, "\n	", "\n")
	return fmt.Sprintf("struct {}{ %s  }", code)
}

// exprName extracts the name from an AST expression.
func exprName(expr ast.Expr) (string, bool) {
	switch exp := expr.(type) {
	case *ast.Ident:
		return exp.Name, true
	case *ast.SelectorExpr:
		return exp.Sel.Name, true
	default:
		return "", false
	}
}

// addFunctionCall modifies function call arguments.
func addFunctionCall(expr *ast.CallExpr, calls functionCalls) error {
	for _, c := range calls {
		newArg := ast.NewIdent(c.code)
		newArg.NamePos = token.Pos(c.index)

		switch {
		case c.index == -1:
			// Append at end
			expr.Args = append(expr.Args, newArg)
		case c.index >= 0 && c.index <= len(expr.Args):
			// Insert at index
			expr.Args = append(expr.Args[:c.index], append([]ast.Expr{newArg}, expr.Args[c.index:]...)...)
		default:
			return errors.Errorf("function call index %d out of range", c.index)
		}
	}
	return nil
}

// addStructs modifies struct literal fields.
func addStructs(expr *ast.CompositeLit, structs functionStructs) error {
	for _, s := range structs {
		var newArg ast.Expr = ast.NewIdent(s.code)
		if s.param != "" {
			newArg = &ast.KeyValueExpr{
				Key:   ast.NewIdent(s.param),
				Colon: token.Pos(s.index),
				Value: ast.NewIdent(s.code),
			}
		}

		switch {
		case s.index == -1:
			// Append at end
			expr.Elts = append(expr.Elts, newArg)
		case s.index >= 0 && s.index <= len(expr.Elts):
			// Insert at index
			expr.Elts = append(expr.Elts[:s.index], append([]ast.Expr{newArg}, expr.Elts[s.index:]...)...)
		default:
			return errors.Errorf("function call index %d out of range", s.index)
		}
	}
	return nil
}

// codeToBlockStmt parses code string into AST block statement.
func codeToBlockStmt(fileSet *token.FileSet, code string) (*ast.BlockStmt, error) {
	newFuncContent := toCode(code)
	newContent, err := parser.ParseFile(fileSet, "", newFuncContent, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return newContent.Decls[0].(*ast.FuncDecl).Body, nil
}

// toCode wraps code in a function for parsing.
func toCode(code string) string {
	return fmt.Sprintf("package p; func _() { %s }", strings.TrimSpace(code))
}

// formatNode formats an AST node into Go source code.
func formatNode(fileSet *token.FileSet, f *ast.File) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// applyFunctionOptions applies all modifications to a function.
func applyFunctionOptions(fileSet *token.FileSet, f *ast.FuncDecl, opts *functionOpts) (err error) {
	// Parse new function body if provided
	var newFunctionBody *ast.BlockStmt
	if opts.body != "" {
		newFunctionBody, err = codeToBlockStmt(fileSet, opts.body)
		if err != nil {
			return err
		}
	}

	// Parse append code.
	appendCode, err := addCode(fileSet, opts.appendCode)
	if err != nil {
		return err
	}

	// Parse return variables.
	returnStmts, err := modifyReturnVars(fileSet, opts.returnVars)
	if err != nil {
		return err
	}

	// Create maps for tracking modifications.
	var (
		callMap        = opts.insideCall.Map()
		callMapCheck   = opts.insideCall.Map()
		structMap      = opts.insideStruct.Map()
		structMapCheck = opts.insideStruct.Map()
	)

	// Apply all modifications.
	var (
		found      bool
		errInspect error
	)
	ast.Inspect(f, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != opts.funcName {
			return true
		}

		// Add parameters.
		if err := addParams(funcDecl, opts.newParams); err != nil {
			errInspect = err
			return false
		}

		// Replace body if needed.
		if newFunctionBody != nil {
			funcDecl.Body = newFunctionBody
			funcDecl.Body.Rbrace = funcDecl.Body.Pos()
		}

		// Add new lines.
		if err := addNewLine(fileSet, funcDecl, opts.newLines); err != nil {
			errInspect = err
			return false
		}

		// Modify returns and append code.
		if err := modifyReturn(funcDecl, returnStmts, appendCode); err != nil {
			errInspect = err
			return false
		}

		// Modify function calls and struct literals.
		ast.Inspect(funcDecl, func(n ast.Node) bool {
			switch expr := n.(type) {
			case *ast.CallExpr:
				name, exist := exprName(expr.Fun)
				if !exist {
					return true
				}

				calls, ok := callMap[name]
				if !ok {
					return true
				}

				if err := addFunctionCall(expr, calls); err != nil {
					errInspect = err
					return false
				}
				delete(callMapCheck, name)

			case *ast.CompositeLit:
				name, exist := exprName(expr.Type)
				if !exist {
					return true
				}

				structs, ok := structMap[name]
				if !ok {
					return true
				}

				if err := addStructs(expr, structs); err != nil {
					errInspect = err
					return false
				}
				delete(structMapCheck, name)

			default:
				return true
			}
			return true
		})
		if errInspect != nil {
			return false
		}

		// Verify all modifications were applied.
		if len(callMapCheck) > 0 {
			errInspect = errors.Errorf("function calls not found: %v", callMapCheck)
			return false
		}
		if len(structMapCheck) > 0 {
			errInspect = errors.Errorf("function structs not found: %v", structMapCheck)
			return false
		}

		// Add test cases.
		if err := addTestCase(fileSet, funcDecl, opts.appendTestCase); err != nil {
			errInspect = err
			return false
		}

		found = true
		return false
	})

	if errInspect != nil {
		return errInspect
	}
	if !found {
		return errors.Errorf("function %s not found in file content", opts.funcName)
	}

	return nil
}

// ModifyCaller replaces all arguments of a specific function call in the given content.
// The callerExpr should be in the format "pkgname.FuncName" or just "FuncName".
// The modifiers function is called with the existing arguments and should return the new arguments.
func ModifyCaller(content, callerExpr string, modifiers func([]string) ([]string, error)) (string, error) {
	// parse the caller expression to extract package name and function name
	var pkgName, funcName string
	parts := strings.Split(callerExpr, ".")
	switch len(parts) {
	case 1:
		funcName = parts[0]
	case 2:
		pkgName = parts[0]
		funcName = parts[1]
	default:
		return "", errors.New("invalid caller expression format, use 'pkgname.FuncName' or 'FuncName'")
	}

	fileSet := token.NewFileSet()
	// preserve original source positions for maintaining whitespace
	fileSet.AddFile("", fileSet.Base(), len(content))

	f, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// track positions of all call expressions that need modification
	type callModification struct {
		node     *ast.CallExpr
		newArgs  []string
		startPos token.Pos
		endPos   token.Pos
	}

	var modifications []callModification

	errInspect := Inspect(f, func(n ast.Node) error {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return nil
		}

		// check if this call matches our target function
		match := false
		switch fun := callExpr.Fun.(type) {
		case *ast.Ident:
			// handle case of FuncName()
			if pkgName == "" && fun.Name == funcName {
				match = true
			}
		case *ast.SelectorExpr:
			// handle case of pkg.FuncName()
			if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == pkgName && fun.Sel.Name == funcName {
				match = true
			}
		}

		if !match {
			return nil
		}

		// extract current arguments as strings
		currentArgs := make([]string, len(callExpr.Args))
		for i, arg := range callExpr.Args {
			var buf bytes.Buffer
			if err := format.Node(&buf, fileSet, arg); err != nil {
				return err
			}
			currentArgs[i] = buf.String()
		}

		// apply the modifier function
		newArgs, err := modifiers(currentArgs)
		if err != nil {
			return err
		}

		// record this modification for later application
		modifications = append(modifications, callModification{
			node:     callExpr,
			newArgs:  newArgs,
			startPos: callExpr.Lparen + 1, // position right after the left parenthesis
			endPos:   callExpr.Rparen,     // position of the right parenthesis
		})

		return nil
	})

	if errInspect != nil {
		return "", errInspect
	}

	if len(modifications) == 0 {
		return "", errors.Errorf("function call %s not found in file content", callerExpr)
	}

	// apply modifications in reverse order to avoid position shifts
	sort.Slice(modifications, func(i, j int) bool {
		return modifications[i].startPos > modifications[j].startPos
	})

	// make modifications directly to the content string
	result := []byte(content)
	for _, mod := range modifications {
		// build the new arguments string
		newArgsStr := strings.Join(mod.newArgs, ", ")

		// replace the arguments in the original content
		startOffset := fileSet.Position(mod.startPos).Offset
		endOffset := fileSet.Position(mod.endPos).Offset

		result = append(
			result[:startOffset],
			append(
				[]byte(newArgsStr),
				result[endOffset:]...,
			)...,
		)
	}

	return string(result), nil
}
