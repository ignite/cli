package replacer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// AppendFunction appends a new function to the end of the Go source code content.
func AppendFunction(fileContent string, function string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the function body as a separate file.
	funcFile, err := parser.ParseFile(fileSet, "", "package main\n"+function, parser.AllErrors)
	if err != nil {
		return "", err
	}

	// Extract the first declaration, assuming it's a function declaration.
	var funcDecl *ast.FuncDecl
	for _, decl := range funcFile.Decls {
		if fDecl, ok := decl.(*ast.FuncDecl); ok {
			funcDecl = fDecl
			break
		}
	}
	if funcDecl == nil {
		return "", errors.Errorf("no function declaration found in the provided function body")
	}

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Append the function declaration to the file's declarations.
	f.Decls = append(f.Decls, funcDecl)

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// AppendCodeToFunction inserts code before the end or the return, if exists, of a function in Go source code content.
func AppendCodeToFunction(fileContent, functionName, codeToInsert string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Parse the Go code to insert.
	insertionExpr, err := parser.ParseExpr(codeToInsert)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the code you want to replace.
			if funcDecl.Name.Name == functionName {
				// Check if there is a return statement in the function.
				if len(funcDecl.Body.List) > 0 {
					lastStmt := funcDecl.Body.List[len(funcDecl.Body.List)-1]
					switch lastStmt.(type) {
					case *ast.ReturnStmt:
						// If there is a return, insert before it.
						funcDecl.Body.List = append(funcDecl.Body.List[:len(funcDecl.Body.List)-1], &ast.ExprStmt{X: insertionExpr}, lastStmt)
					default:
						// If there is no return, insert at the end of the function body.
						funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{X: insertionExpr})
					}
				} else {
					// If there are no statements in the function body, insert at the end of the function body.
					funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{X: insertionExpr})
				}
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", errors.Errorf("function %s not found", functionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReplaceReturnVars replaces return statements in a Go function with a new return statement.
func ReplaceReturnVars(fileContent, functionName string, returnVars ...string) (string, error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	returnStmts := make([]ast.Expr, 0)
	for _, returnVar := range returnVars {
		// Parse the new return var to expression.
		newRetExpr, err := parser.ParseExpr(returnVar)
		if err != nil {
			return "", err
		}
		returnStmts = append(returnStmts, newRetExpr)
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the code you want to replace.
			if funcDecl.Name.Name == functionName {
				// Replace the return statements.
				for _, stmt := range funcDecl.Body.List {
					if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
						// Remove existing return statements.
						retStmt.Results = nil
						// Add the new return statement.
						retStmt.Results = append(retStmt.Results, returnStmts...)
					}
				}
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", errors.Errorf("function %s not found", functionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReplaceFunctionContent replaces a function implementation content in Go source code.
func ReplaceFunctionContent(fileContent, oldFunctionName, newFunction string) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Parse the content of the new function into an ast.File.
	newFuncContent := fmt.Sprintf("package p; func _() { %s }", strings.TrimSpace(newFunction))
	newFile, err := parser.ParseFile(fileSet, "", newFuncContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the code you want to replace.
			if funcDecl.Name.Name == oldFunctionName {
				// Take the body of the new function from the parsed file.
				newFunctionBody := newFile.Decls[0].(*ast.FuncDecl).Body
				// Replace the function body with the body of the new function.
				funcDecl.Body = newFunctionBody
				found = true
				return false
			}
		}
		return true
	})

	if !found {
		return "", errors.Errorf("function %s not found in file content", oldFunctionName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// AppendParamToFunctionCall inserts a parameter to a function call inside a function in Go source code content.
func AppendParamToFunctionCall(fileContent, functionName, functionCallName, paramToAdd string, index int) (modifiedContent string, err error) {
	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var (
		found      bool
		errInspect error
	)
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Check if the function has the name you want to replace.
			if funcDecl.Name.Name == functionName {
				ast.Inspect(funcDecl, func(n ast.Node) bool {
					callExpr, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					// Check if the call expression matches the function call name
					ident, ok := callExpr.Fun.(*ast.Ident)
					if !ok || ident.Name != functionCallName {
						selector, ok := callExpr.Fun.(*ast.SelectorExpr)
						if !ok || selector.Sel.Name != functionCallName {
							return true
						}
					}
					// Construct the new argument to be added
					newArg := &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(paramToAdd),
					}
					switch {
					case index == -1:
						// Append the new argument to the end
						callExpr.Args = append(callExpr.Args, newArg)
						found = true
					case index >= 0 && index <= len(callExpr.Args):
						// Insert the new argument at the specified index
						callExpr.Args = append(callExpr.Args[:index], append([]ast.Expr{newArg}, callExpr.Args[index:]...)...)
						found = true
					default:
						errInspect = fmt.Errorf("index out of range")
						return false // Stop the inspection, an error occurred
					}
					return true // Continue the inspection for duplicated calls
				})
				return false // Stop the inspection, we found what we needed
			}
		}
		return true // Continue inspecting
	})
	if errInspect != nil {
		return "", errInspect
	}
	if !found {
		return "", fmt.Errorf("function %s not found or no calls to %s inside the function", functionName, functionCallName)
	}

	// Write the modified AST to a buffer.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}
