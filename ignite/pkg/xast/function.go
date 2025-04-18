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
		newParams      []functionParam
		body           string
		newLines       []functionLine
		insideCall     []functionCall
		insideStruct   []functionStruct
		appendTestCase []string
		appendCode     []string
		returnVars     []string
	}

	// FunctionOptions configures code generation.
	FunctionOptions func(*functionOpts)

	functionStruct struct {
		name  string
		param string
		code  string
		index int
	}
	functionCall struct {
		name  string
		code  string
		index int
	}
	functionParam struct {
		name    string
		varType string
		index   int
	}
	functionLine struct {
		code   string
		number uint64
	}
)

// AppendFuncParams add a new param value.
func AppendFuncParams(name, varType string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.newParams = append(c.newParams, functionParam{
			name:    name,
			varType: varType,
			index:   index,
		})
	}
}

// ReplaceFuncBody replace all body of the function, the method will replace first and apply the other options after.
func ReplaceFuncBody(body string) FunctionOptions {
	return func(c *functionOpts) {
		c.body = body
	}
}

// AppendFuncTestCase append test a new test case, if exists, of a function in Go source code content.
func AppendFuncTestCase(testCase string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendTestCase = append(c.appendTestCase, testCase)
	}
}

// AppendFuncCode append code before the end or the return, if exists, of a function in Go source code content.
func AppendFuncCode(code string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendCode = append(c.appendCode, code)
	}
}

// AppendFuncAtLine append a new code at line.
func AppendFuncAtLine(code string, lineNumber uint64) FunctionOptions {
	return func(c *functionOpts) {
		c.newLines = append(c.newLines, functionLine{
			code:   code,
			number: lineNumber,
		})
	}
}

// AppendInsideFuncCall add code inside another function call. For instances, the method have a parameter a
// call 'New(param1, param2)' and we want to add the param3 the result will be 'New(param1, param2, param3)'.
func AppendInsideFuncCall(callName, code string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.insideCall = append(c.insideCall, functionCall{
			name:  callName,
			code:  code,
			index: index,
		})
	}
}

// AppendFuncStruct add code inside another function call. For instances,
// the struct have only one parameter 'Params{Param1: param1}' and we want to add
// the param2 the result will be 'Params{Param1: param1, Param2: param2}'.
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

// NewFuncReturn replaces return statements in a Go function with a new return statement.
func NewFuncReturn(returnVars ...string) FunctionOptions {
	return func(c *functionOpts) {
		c.returnVars = append(c.returnVars, returnVars...)
	}
}

func newFunctionOptions() functionOpts {
	return functionOpts{
		newParams:      make([]functionParam, 0),
		body:           "",
		newLines:       make([]functionLine, 0),
		insideCall:     make([]functionCall, 0),
		insideStruct:   make([]functionStruct, 0),
		appendTestCase: make([]string, 0),
		appendCode:     make([]string, 0),
		returnVars:     make([]string, 0),
	}
}

// ModifyFunction modify a function based in the options.
func ModifyFunction(fileContent, functionName string, functions ...FunctionOptions) (modifiedContent string, err error) {
	// Apply function options.
	opts := newFunctionOptions()
	for _, o := range functions {
		o(&opts)
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Parse the content of the new function into an ast.
	var newFunctionBody *ast.BlockStmt
	if opts.body != "" {
		newFunctionBody, err = codeToBlockStmt(fileSet, opts.body)
		if err != nil {
			return "", err
		}
	}

	// Parse the content of the append code an ast.
	appendCode := make([]ast.Stmt, 0)
	for _, codeToInsert := range opts.appendCode {
		body, err := codeToBlockStmt(fileSet, codeToInsert)
		if err != nil {
			return "", err
		}
		appendCode = append(appendCode, body.List...)
	}

	// Parse the content of the return vars into an ast.
	returnStmts := make([]ast.Expr, 0)
	for _, returnVar := range opts.returnVars {
		// Parse the new return var to expression.
		newRetExpr, err := parser.ParseExprFrom(fileSet, "", []byte(returnVar), parser.ParseComments)
		if err != nil {
			return "", err
		}
		returnStmts = append(returnStmts, newRetExpr)
	}

	callMap := make(map[string][]functionCall)
	callMapCheck := make(map[string][]functionCall)
	for _, c := range opts.insideCall {
		calls, ok := callMap[c.name]
		if !ok {
			calls = []functionCall{}
		}
		callMap[c.name] = append(calls, c)
		callMapCheck[c.name] = append(calls, c)
	}

	structMap := make(map[string][]functionStruct)
	structMapCheck := make(map[string][]functionStruct)
	for _, s := range opts.insideStruct {
		structs, ok := structMap[s.name]
		if !ok {
			structs = []functionStruct{}
		}
		structMap[s.name] = append(structs, s)
		structMapCheck[s.name] = append(structs, s)
	}

	// Parse the Go code to insert.
	var (
		found      bool
		errInspect error
	)
	ast.Inspect(f, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != functionName {
			return true
		}

		for _, p := range opts.newParams {
			fieldParam := &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(p.name)},
				Type:  ast.NewIdent(p.varType),
			}
			switch {
			case p.index == -1:
				// Append the new argument to the end
				funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, fieldParam)
			case p.index >= 0 && p.index <= len(funcDecl.Type.Params.List):
				// Insert the new argument at the specified index
				funcDecl.Type.Params.List = append(
					funcDecl.Type.Params.List[:p.index],
					append([]*ast.Field{fieldParam}, funcDecl.Type.Params.List[p.index:]...)...,
				)
			default:
				errInspect = errors.Errorf("params index %d out of range", p.index)
				return false
			}
		}

		// Check if the function has the code you want to replace.
		if newFunctionBody != nil {
			funcDecl.Body = newFunctionBody
			funcDecl.Body.Rbrace = funcDecl.Body.Pos() // Re-adjust positions if necessary.
		}

		// Add the new code at line.
		for _, newLine := range opts.newLines {
			// Check if the function body has enough lines.
			if newLine.number > uint64(len(funcDecl.Body.List))-1 {
				errInspect = errors.Errorf("line number %d out of range", newLine.number)
				return false
			}
			// Parse the Go code to insert.
			insertionExpr, err := parser.ParseExprFrom(fileSet, "", []byte(newLine.code), parser.ParseComments)
			if err != nil {
				errInspect = err
				return false
			}
			// Insert code at the specified line number.
			funcDecl.Body.List = append(
				funcDecl.Body.List[:newLine.number],
				append([]ast.Stmt{&ast.ExprStmt{X: insertionExpr}}, funcDecl.Body.List[newLine.number:]...)...,
			)
		}

		// Check if there is a return statement in the function.
		if len(funcDecl.Body.List) > 0 {
			lastStmt := funcDecl.Body.List[len(funcDecl.Body.List)-1]
			switch stmt := lastStmt.(type) {
			case *ast.ReturnStmt:
				// Replace the return statements.
				if len(returnStmts) > 0 {
					// Remove existing return statements.
					stmt.Results = nil
					// Add the new return statement.
					stmt.Results = append(stmt.Results, returnStmts...)
				}
				if len(appendCode) > 0 {
					// If there is a return, insert before it.
					appendCode = append(appendCode, stmt)
					funcDecl.Body.List = append(funcDecl.Body.List[:len(funcDecl.Body.List)-1], appendCode...)
				}
			default:
				if len(returnStmts) > 0 {
					errInspect = errors.New("return statement not found")
					return false
				}
				// If there is no return, insert at the end of the function body.
				if len(appendCode) > 0 {
					funcDecl.Body.List = append(funcDecl.Body.List, appendCode...)
				}
			}
		} else {
			if len(returnStmts) > 0 {
				errInspect = errors.New("return statement not found")
				return false
			}
			// If there are no statements in the function body, insert at the end of the function body.
			if len(appendCode) > 0 {
				funcDecl.Body.List = append(funcDecl.Body.List, appendCode...)
			}
		}

		// Add new code to the function callers.
		ast.Inspect(funcDecl, func(n ast.Node) bool {
			switch expr := n.(type) {
			case *ast.CallExpr: // Add a new parameter to a function call.
				// Check if the call expression matches the function call name.
				name := ""
				switch exp := expr.Fun.(type) {
				case *ast.Ident:
					name = exp.Name
				case *ast.SelectorExpr:
					name = exp.Sel.Name
				default:
					return true
				}

				calls, ok := callMap[name]
				if !ok {
					return true
				}

				// Construct the new argument to be added
				for _, c := range calls {
					newArg := ast.NewIdent(c.code)
					newArg.NamePos = token.Pos(c.index)
					switch {
					case c.index == -1:
						// Append the new argument to the end
						expr.Args = append(expr.Args, newArg)
					case c.index >= 0 && c.index <= len(expr.Args):
						// Insert the new argument at the specified index
						expr.Args = append(expr.Args[:c.index], append([]ast.Expr{newArg}, expr.Args[c.index:]...)...)
					default:
						errInspect = errors.Errorf("function call index %d out of range", c.index)
						return false // Stop the inspection, an error occurred
					}
				}
				delete(callMapCheck, name)
			case *ast.CompositeLit: // Add a new parameter to a literal struct.
				// Check if the call expression matches the function call name.
				name := ""
				switch exp := expr.Type.(type) {
				case *ast.Ident:
					name = exp.Name
				case *ast.SelectorExpr:
					name = exp.Sel.Name
				default:
					return true
				}

				structs, ok := structMap[name]
				if !ok {
					return true
				}

				// Construct the new argument to be added
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
						// Append the new argument to the end
						expr.Elts = append(expr.Elts, newArg)
					case s.index >= 0 && s.index <= len(expr.Elts):
						// Insert the new argument at the specified index
						expr.Elts = append(expr.Elts[:s.index], append([]ast.Expr{newArg}, expr.Elts[s.index:]...)...)
					default:
						errInspect = errors.Errorf("function call index %d out of range", s.index)
						return false // Stop the inspection, an error occurred
					}
				}
				delete(structMapCheck, name)
			default:
				return true
			}
			return true // Continue the inspection for duplicated calls
		})
		if errInspect != nil {
			return false
		}
		if len(callMapCheck) > 0 {
			errInspect = errors.Errorf("function calls not found: %v", callMapCheck)
			return false
		}
		if len(structMapCheck) > 0 {
			errInspect = errors.Errorf("function structs not found: %v", structMapCheck)
			return false
		}

		// Locate the `tests` variable inside the function
		for _, stmt := range funcDecl.Body.List {
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if !ok || len(assignStmt.Lhs) == 0 {
				continue
			}

			// Check if the `tests` variable is being declared
			ident, ok := assignStmt.Lhs[0].(*ast.Ident)
			if !ok || ident.Name != "tests" {
				continue
			}

			// Find the composite literal (slice) for the `tests` variable
			compositeLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit)
			if !ok {
				continue
			}

			for _, testCase := range opts.appendTestCase {
				// Parse the new test case into an AST expression
				testCaseStmt, err := structToBlockStmt(testCase)
				if err != nil {
					errInspect = err
					return false
				}
				// Append the new test case to the list
				compositeLit.Elts = append(compositeLit.Elts, testCaseStmt)
			}
		}

		// everything is ok, mark as found and stop the inspect
		found = true
		return false
	})
	if errInspect != nil {
		return "", errInspect
	}
	if !found {
		return "", errors.Errorf("function %s not found in file content", functionName)
	}

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	// Return the modified content.
	return buf.String(), nil
}

func codeToBlockStmt(fileSet *token.FileSet, code string) (*ast.BlockStmt, error) {
	newFuncContent := toCode(code)
	newContent, err := parser.ParseFile(fileSet, "", newFuncContent, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return newContent.Decls[0].(*ast.FuncDecl).Body, nil
}

func toCode(code string) string {
	return fmt.Sprintf("package p; func _() { %s }", strings.TrimSpace(code))
}

func structToBlockStmt(code string) (ast.Expr, error) {
	newFuncContent := toStruct(code)
	newContent, err := parser.ParseExpr(newFuncContent)
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

func toStruct(code string) string {
	return fmt.Sprintf(`struct {}{ %s }`, strings.TrimSpace(code))
}

// ModifyCaller replaces all arguments of a specific function call in the given content.
// The callerExpr should be in the format "pkgname.FuncName" or just "FuncName".
// The modifiers function is called with the existing arguments and should return the new arguments.
func ModifyCaller(content, callerExpr string, modifiers func([]string) ([]string, error)) (string, error) {
	// parse the caller expression to extract package name and function name
	var pkgName, funcName string
	parts := strings.Split(callerExpr, ".")
	if len(parts) == 1 {
		funcName = parts[0]
	} else if len(parts) == 2 {
		pkgName = parts[0]
		funcName = parts[1]
	} else {
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
