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

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

type (
	// functionOpts represent the options for functions.
	functionOpts struct {
		newParams    []param
		body         string
		newLines     []line
		insideCall   []call
		insideStruct []str
		appendCode   []string
		returnVars   []string
	}

	// FunctionOptions configures code generation.
	FunctionOptions func(*functionOpts)

	str struct {
		structName string
		paramName  string
		code       string
		index      int
	}
	call struct {
		name  string
		code  string
		index int
	}
	param struct {
		name    string
		varType string
		index   int
	}
	line struct {
		code   string
		number uint64
	}
)

// AppendFuncParams add a new param value.
func AppendFuncParams(name, varType string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.newParams = append(c.newParams, param{
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

// AppendFuncCode append code before the end or the return, if exists, of a function in Go source code content.
func AppendFuncCode(code string) FunctionOptions {
	return func(c *functionOpts) {
		c.appendCode = append(c.appendCode, code)
	}
}

// AppendFuncCodeAtLine append a new code at line.
var AppendFuncCodeAtLine = AppendFuncAtLine

// AppendFuncAtLine append a new code at line.
func AppendFuncAtLine(code string, lineNumber uint64) FunctionOptions {
	return func(c *functionOpts) {
		c.newLines = append(c.newLines, line{
			code:   code,
			number: lineNumber,
		})
	}
}

// AppendInsideFuncCall add code inside another function call. For instances, the method have a parameter a
// call 'New(param1, param2)' and we want to add the param3 the result will be 'New(param1, param2, param3)'.
func AppendInsideFuncCall(callName, code string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.insideCall = append(c.insideCall, call{
			name:  callName,
			code:  code,
			index: index,
		})
	}
}

// AppendInsideFuncStruct add code inside another function call. For instances,
// the struct have only one parameter 'Params{Param1: param1}' and we want to add
// the param2 the result will be 'Params{Param1: param1, Param2: param2}'.
func AppendInsideFuncStruct(structName, paramName, code string, index int) FunctionOptions {
	return func(c *functionOpts) {
		c.insideStruct = append(c.insideStruct, str{
			structName: structName,
			paramName:  paramName,
			code:       code,
			index:      index,
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
		newParams:    make([]param, 0),
		body:         "",
		newLines:     make([]line, 0),
		insideCall:   make([]call, 0),
		insideStruct: make([]str, 0),
		appendCode:   make([]string, 0),
		returnVars:   make([]string, 0),
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

	callMap := make(map[string][]call)
	callMapCheck := make(map[string][]call)
	for _, c := range opts.insideCall {
		calls, ok := callMap[c.name]
		if !ok {
			calls = []call{}
		}
		callMap[c.name] = append(calls, c)
		callMapCheck[c.name] = append(calls, c)
	}

	structMap := make(map[string][]str)
	structMapCheck := make(map[string][]str)
	for _, s := range opts.insideStruct {
		structs, ok := structMap[s.structName]
		if !ok {
			structs = []str{}
		}
		structMap[s.structName] = append(structs, s)
		structMapCheck[s.structName] = append(structs, s)
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
		}

		// Add the new code at line.
		for _, newLine := range opts.newLines {
			// Check if the function body has enough lines.
			if newLine.number > uint64(len(funcDecl.Body.List))-1 {
				errInspect = errors.Errorf("line number %d out of range", newLine.number)
				return false
			}

			// Parse the Go code to insert.
			body, err := codeToBlockStmt(fileSet, newLine.code)
			if err != nil {
				errInspect = err
				return false
			}

			// Insert code at the specified line number.
			funcDecl.Body.List = append(
				funcDecl.Body.List[:newLine.number],
				append(body.List, funcDecl.Body.List[newLine.number:]...)...,
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
					if s.paramName != "" {
						newArg = &ast.KeyValueExpr{
							Key:   ast.NewIdent(s.paramName),
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
