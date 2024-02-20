package xast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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
		newFuncContent := fmt.Sprintf("package p; func _() { %s }", strings.TrimSpace(opts.body))
		newContent, err := parser.ParseFile(fileSet, "", newFuncContent, parser.ParseComments)
		if err != nil {
			return "", err
		}
		newFunctionBody = newContent.Decls[0].(*ast.FuncDecl).Body
	}

	// Parse the content of the append code an ast.
	appendCode := make([]ast.Stmt, 0)
	for _, codeToInsert := range opts.appendCode {
		insertionExpr, err := parser.ParseExprFrom(fileSet, "", []byte(codeToInsert), parser.ParseComments)
		if err != nil {
			return "", err
		}
		appendCode = append(appendCode, &ast.ExprStmt{X: insertionExpr})
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
				Type:  &ast.Ident{Name: p.varType},
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
