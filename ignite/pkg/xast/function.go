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
		newParams      []functionParam  // Parameters to add to the function.
		body           string           // New function body content.
		newLines       []functionLine   // Lines to insert at specific positions.
		insideCall     functionCalls    // Function calls to modify.
		insideStruct   functionStructs  // Struct literals to modify.
		appendTestCase []string         // Test cases to append.
		appendCode     []string         // Code to append at the end.
		returnVars     []string         // Return variables to modify.
		appendSwitch   functionSwitches // Switch cases to append.
		removeCalls    []string         // Function calls to remove.
	}

	// FunctionOptions configures code generation.
	FunctionOptions func(*functionOpts)

	// functionStruct represents a struct literal to modify.
	functionSwitch struct {
		condition  string // Condition to find.
		switchCase string // Switch case to insert.
		switchBody string // Code to insert.
	}

	functionSwitches    []functionSwitch
	functionSwitchesMap map[string]functionSwitches

	// functionStruct represents a struct literal to modify.
	functionStruct struct {
		name  string // Name of the struct type.
		param string // Name of the struct field.
		code  string // Code to insert.
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

// Map converts a slice of functionStructs to a map keyed by struct name.
func (s functionSwitches) Map() functionSwitchesMap {
	switchesMap := make(functionSwitchesMap)
	for _, c := range s {
		switches, ok := switchesMap[c.condition]
		if !ok {
			switches = make(functionSwitches, 0)
		}
		switchesMap[c.condition] = append(switches, c)
	}
	return switchesMap
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

func cloneFunctionCallsMap(src functionCallsMap) functionCallsMap {
	dst := make(functionCallsMap, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneFunctionStructsMap(src functionStructsMap) functionStructsMap {
	dst := make(functionStructsMap, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneFunctionSwitchesMap(src functionSwitchesMap) functionSwitchesMap {
	dst := make(functionSwitchesMap, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
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
// AppendInsideFuncCall appends code inside a function call.
// The callName parameter can be either:
//   - Simple name: "NewKeeper" matches any call to NewKeeper regardless of package/receiver
//   - Qualified name: "foo.NewKeeper" matches only calls to NewKeeper with foo as the package/receiver
//
// The code parameter is the argument to insert, and index specifies the position:
//   - index >= 0: insert at the specified position
//   - index == -1: append at the end
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
// the struct has only one parameter 'Params{Param1: param1}' and we want to add
// the param2 the result will be 'Params{Param1: param1, Param2: param2}'.
//
// The name parameter can be either:
//   - Simple name: "Keeper" matches any struct literal of type Keeper regardless of package
//   - Qualified name: "keeper.Keeper" matches only struct literals with keeper as the package
func AppendFuncStruct(name, param, code string) FunctionOptions {
	return func(c *functionOpts) {
		c.insideStruct = append(c.insideStruct, functionStruct{
			name:  name,
			param: param,
			code:  code,
		})
	}
}

// NewFuncReturn replaces return statements in a function.
func NewFuncReturn(returnVars ...string) FunctionOptions {
	return func(c *functionOpts) {
		c.returnVars = append(c.returnVars, returnVars...)
	}
}

// AppendSwitchCase inserts a new case with the code at a specific switch condition statement.
func AppendSwitchCase(condition, switchCase, switchBody string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendSwitch = append(c.appendSwitch, functionSwitch{
			condition:  condition,
			switchCase: switchCase,
			switchBody: switchBody,
		})
	}
}

// RemoveFuncCall removes function calls with the specified name from within a function.
// The callName can be either a simple function name like "doSomething" or a qualified
// name like "pkg.DoSomething".
func RemoveFuncCall(callName string) FunctionOptions {
	return func(c *functionOpts) {
		c.removeCalls = append(c.removeCalls, callName)
	}
}

// newFunctionOptions creates a new functionOpts with defaults.
func newFunctionOptions() functionOpts {
	return functionOpts{
		newParams:      make([]functionParam, 0),
		body:           "",
		newLines:       make([]functionLine, 0),
		insideCall:     make(functionCalls, 0),
		insideStruct:   make(functionStructs, 0),
		appendTestCase: make([]string, 0),
		appendCode:     make([]string, 0),
		returnVars:     make([]string, 0),
		removeCalls:    make([]string, 0),
	}
}

// ModifyFunction modifies a function in Go source code using functional options.
func ModifyFunction(content string, funcName string, functions ...FunctionOptions) (string, error) {
	// Collect all function options.
	opts := newFunctionOptions()
	for _, fn := range functions {
		fn(&opts)
	}

	// Parse source into AST
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return "", errors.Errorf("failed to parse file (%s): %w", funcName, err)
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
	code := make([]ast.Stmt, 0, len(appendCode))
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
	stmts := make([]ast.Expr, 0, len(returnVars))
	for _, returnVar := range returnVars {
		newRetExpr, err := parser.ParseExprFrom(fileSet, "", []byte(returnVar), parser.ParseComments)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, newRetExpr)
	}
	return stmts, nil
}

// appendSwitchCase appends a new case to a switch statement.
func appendSwitchCase(fileSet *token.FileSet, stmt ast.Node, fs functionSwitches) error {
	for _, f := range fs {
		// Parse the new case code
		newRetExpr, err := parser.ParseExprFrom(fileSet, "", []byte(f.switchCase), parser.ParseComments)
		if err != nil {
			return err
		}

		bodyStmt, err := codeToBlockStmt(fileSet, f.switchBody)
		if err != nil {
			return err
		}

		// Create a new case clause
		newCase := &ast.CaseClause{
			List:  []ast.Expr{newRetExpr},
			Body:  bodyStmt.List,
			Case:  token.NoPos, // Keep first item aligned with case keyword
			Colon: token.NoPos, // Keep colon aligned with case keyword
		}

		// Handle different types of switch statements
		switch statement := stmt.(type) {
		case *ast.TypeSwitchStmt:
			statement.Body.List = appendCaseToList(statement.Body.List, newCase)
		case *ast.SwitchStmt:
			statement.Body.List = appendCaseToList(statement.Body.List, newCase)
		default:
			return errors.Errorf("unsupported switch statement type: %T", stmt)
		}
	}
	return nil
}

// appendCaseToList handles inserting a case clause into a list of statements,
// placing it before any default case if one exists.
func appendCaseToList(list []ast.Stmt, newCase *ast.CaseClause) []ast.Stmt {
	if len(list) > 0 {
		lastCase, isDefault := list[len(list)-1].(*ast.CaseClause)
		if isDefault && len(lastCase.List) == 0 {
			// Insert before default.
			return append(list[:len(list)-1], newCase, list[len(list)-1])
		}
	}

	return append(list, newCase)
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
		maxLine := len(funcDecl.Body.List) - 1
		// Validate line number
		if int(newLine.number) > maxLine {
			return errors.Errorf("line number %d out of range (max %d)", newLine.number, maxLine)
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
	newContent, err := parser.ParseExprFrom(fSet, "", newFuncContent, parser.AllErrors)
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
		// Check if X is an identifier to get the package name
		if ident, ok := exp.X.(*ast.Ident); ok {
			// Return qualified name: package.Function
			return ident.Name + "." + exp.Sel.Name, true
		}
		// Fallback to just the selector name if X is not an identifier
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
func addStructs(fileSet *token.FileSet, expr *ast.CompositeLit, structs functionStructs) {
	// Find the current max offset to avoid reused positions
	file := fileSet.File(expr.Pos())
	maxOffset := file.Offset(expr.Rbrace)
	for _, elt := range expr.Elts {
		if pos := elt.End(); pos.IsValid() {
			offset := file.Offset(pos)
			if offset > maxOffset {
				maxOffset = offset
			}
		}
	}

	for i, s := range structs {
		// Advance position
		insertOffset := maxOffset + i
		insertPos := file.Pos(insertOffset)

		value := ast.NewIdent(s.code)
		value.NamePos = insertPos

		var newArg ast.Expr = value
		if s.param != "" {
			key := ast.NewIdent(s.param)
			key.NamePos = insertPos + token.Pos(i)

			newArg = &ast.KeyValueExpr{
				Key:   key,
				Value: value,
				Colon: insertPos,
			}
		}

		expr.Elts = append(expr.Elts, newArg)
		expr.Rbrace += token.Pos(i + 1)
	}

	// Ensure closing brace is on a new line
	if len(expr.Elts) > 0 {
		last := expr.Elts[len(expr.Elts)-1]
		if file.Line(expr.Rbrace) == file.Line(last.End()) {
			// Force a new line before Rbrace
			file.AddLine(file.Offset(expr.Rbrace))
		}
	}
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
func formatNode(fileSet *token.FileSet, n ast.Node) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, n); err != nil {
		return "", err
	}

	node := strings.TrimSpace(buf.String())
	return node, nil
}

// applyFunctionOptions applies all modifications to a function.
func applyFunctionOptions(fileSet *token.FileSet, f *ast.FuncDecl, opts *functionOpts) (err error) {
	var newFunctionBody *ast.BlockStmt
	if opts.body != "" {
		newFunctionBody, err = codeToBlockStmt(fileSet, opts.body)
		if err != nil {
			return err
		}
	}

	appendCode, err := addCode(fileSet, opts.appendCode)
	if err != nil {
		return err
	}

	returnStmts, err := modifyReturnVars(fileSet, opts.returnVars)
	if err != nil {
		return err
	}

	callMap := opts.insideCall.Map()
	callMapCheck := cloneFunctionCallsMap(callMap)
	structMap := opts.insideStruct.Map()
	structMapCheck := cloneFunctionStructsMap(structMap)
	switchesCasesMap := opts.appendSwitch.Map()
	switchesCasesMapCheck := cloneFunctionSwitchesMap(switchesCasesMap)

	if len(opts.removeCalls) > 0 {
		if err := removeFunctionCalls(f, opts.removeCalls); err != nil {
			return err
		}
	}

	if err := addParams(f, opts.newParams); err != nil {
		return err
	}

	if newFunctionBody != nil {
		f.Body = newFunctionBody
		f.Body.Rbrace = f.Body.Pos()
	}

	if err := addNewLine(fileSet, f, opts.newLines); err != nil {
		return err
	}

	if err := modifyReturn(f, returnStmts, appendCode); err != nil {
		return err
	}

	for _, bodyList := range f.Body.List {
		var (
			stmt ast.Stmt
			key  string
		)
		switch expr := bodyList.(type) {
		case *ast.TypeSwitchStmt:
			stmt = expr
			key, err = formatNode(fileSet, expr.Assign)
		case *ast.SwitchStmt:
			stmt = expr
			key, err = formatNode(fileSet, expr.Tag)
		default:
			continue
		}
		if err != nil {
			return err
		}

		switchCase, ok := switchesCasesMap[key]
		if !ok {
			continue
		}

		if err := appendSwitchCase(fileSet, stmt, switchCase); err != nil {
			return err
		}
		delete(switchesCasesMapCheck, key)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.CallExpr:
			name, exist := exprName(expr.Fun)
			if !exist {
				return true
			}

			var allCalls functionCalls
			if calls, ok := callMap[name]; ok {
				allCalls = append(allCalls, calls...)
				delete(callMapCheck, name)
			}

			if sel, isSel := expr.Fun.(*ast.SelectorExpr); isSel {
				simpleName := sel.Sel.Name
				if calls, ok := callMap[simpleName]; ok {
					allCalls = append(allCalls, calls...)
					delete(callMapCheck, simpleName)
				}
			}

			if len(allCalls) == 0 {
				return true
			}

			if err = addFunctionCall(expr, allCalls); err != nil {
				return false
			}

		case *ast.CompositeLit:
			name, exist := exprName(expr.Type)
			if !exist {
				return true
			}

			var allStructs functionStructs
			if structs, ok := structMap[name]; ok {
				allStructs = append(allStructs, structs...)
				delete(structMapCheck, name)
			}

			if sel, isSel := expr.Type.(*ast.SelectorExpr); isSel {
				simpleName := sel.Sel.Name
				if structs, ok := structMap[simpleName]; ok {
					allStructs = append(allStructs, structs...)
					delete(structMapCheck, simpleName)
				}
			}

			if len(allStructs) == 0 {
				return true
			}

			addStructs(fileSet, expr, allStructs)
		}
		return true
	})
	if err != nil {
		return err
	}

	if err := addTestCase(fileSet, f, opts.appendTestCase); err != nil {
		return err
	}

	if len(callMapCheck) > 0 {
		return errors.Errorf("function calls not found: %v", callMapCheck)
	}
	if len(structMapCheck) > 0 {
		return errors.Errorf("function structs not found: %v", structMapCheck)
	}
	if len(switchesCasesMapCheck) > 0 {
		return errors.Errorf("function switch not found: %v", switchesCasesMapCheck)
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

// RemoveFunction removes a function declaration from the file content.
func RemoveFunction(content, funcName string) (string, error) {
	// Parse source into AST.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return "", errors.Errorf("failed to parse file: %w", err)
	}

	cmap := ast.NewCommentMap(fset, file, file.Comments)

	// Find the function to remove.
	var found bool
	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == funcName {
			found = true
			// Remove comments associated with this function.
			delete(cmap, decl)
			continue // Skip this declaration to remove it.
		}
		newDecls = append(newDecls, decl)
	}

	if !found {
		return "", errors.Errorf("function %q not found", funcName)
	}

	// Update file declarations and comments.
	file.Decls = newDecls
	file.Comments = cmap.Filter(file).Comments()

	return formatNode(fset, file)
}

// removeFunctionCalls removes all function calls matching the specified names from a function.
func removeFunctionCalls(f *ast.FuncDecl, callNames []string) error {
	if f.Body == nil {
		return nil
	}

	callMap := make(map[string]bool, len(callNames))
	for _, name := range callNames {
		callMap[name] = true
	}

	matchesCall := func(callExpr *ast.CallExpr) bool {
		switch fun := callExpr.Fun.(type) {
		case *ast.Ident:
			return callMap[fun.Name]
		case *ast.SelectorExpr:
			if ident, ok := fun.X.(*ast.Ident); ok {
				return callMap[ident.Name+"."+fun.Sel.Name]
			}
		}
		return false
	}

	var filterStmts func([]ast.Stmt) []ast.Stmt
	filterStmts = func(stmts []ast.Stmt) []ast.Stmt {
		filtered := make([]ast.Stmt, 0, len(stmts))
		for _, stmt := range stmts {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok && matchesCall(callExpr) {
					continue
				}
			}

			switch typedStmt := stmt.(type) {
			case *ast.BlockStmt:
				typedStmt.List = filterStmts(typedStmt.List)
			case *ast.IfStmt:
				if typedStmt.Body != nil {
					typedStmt.Body.List = filterStmts(typedStmt.Body.List)
				}
				switch elseNode := typedStmt.Else.(type) {
				case *ast.BlockStmt:
					elseNode.List = filterStmts(elseNode.List)
				case *ast.IfStmt:
					elseNode.Body.List = filterStmts(elseNode.Body.List)
					if elseBlock, ok := elseNode.Else.(*ast.BlockStmt); ok {
						elseBlock.List = filterStmts(elseBlock.List)
					}
				}
			case *ast.ForStmt:
				if typedStmt.Body != nil {
					typedStmt.Body.List = filterStmts(typedStmt.Body.List)
				}
			case *ast.RangeStmt:
				if typedStmt.Body != nil {
					typedStmt.Body.List = filterStmts(typedStmt.Body.List)
				}
			case *ast.SwitchStmt:
				if typedStmt.Body != nil {
					for _, caseClause := range typedStmt.Body.List {
						if cc, ok := caseClause.(*ast.CaseClause); ok {
							cc.Body = filterStmts(cc.Body)
						}
					}
				}
			case *ast.TypeSwitchStmt:
				if typedStmt.Body != nil {
					for _, caseClause := range typedStmt.Body.List {
						if cc, ok := caseClause.(*ast.CaseClause); ok {
							cc.Body = filterStmts(cc.Body)
						}
					}
				}
			}
			filtered = append(filtered, stmt)
		}
		return filtered
	}

	f.Body.List = filterStmts(f.Body.List)
	return nil
}
