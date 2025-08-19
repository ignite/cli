package xast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type (
	// globalOpts represent the options for globals.
	globalOpts struct {
		globals []global
	}

	// GlobalOptions configures code generation.
	GlobalOptions func(*globalOpts)

	global struct {
		name, varType, value string
	}

	// GlobalType represents the global type.
	GlobalType string
)

const (
	GlobalTypeVar   GlobalType = "var"
	GlobalTypeConst GlobalType = "const"
)

// WithGlobal add a new global.
func WithGlobal(name, varType, value string) GlobalOptions {
	return func(c *globalOpts) {
		c.globals = append(c.globals, global{
			name:    name,
			varType: varType,
			value:   value,
		})
	}
}

func newGlobalOptions() globalOpts {
	return globalOpts{
		globals: make([]global, 0),
	}
}

// InsertGlobal inserts global variables or constants into the provided Go source code content after the import section.
// The function parses the provided content, locates the import section, and inserts the global declarations immediately after it.
// The type of globals (variables or constants) is specified by the globalType parameter.
// Each global declaration is defined by calling WithGlobal function with appropriate arguments.
// The function returns the modified content with the inserted global declarations.
func InsertGlobal(fileContent string, globalType GlobalType, globals ...GlobalOptions) (modifiedContent string, err error) {
	// apply global options.
	opts := newGlobalOptions()
	for _, o := range globals {
		o(&opts)
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	// Find the index of the import declaration or package declaration if no imports.
	var insertIndex int
	for i, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			insertIndex = i + 1
			break
		} else if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			insertIndex = i
			if funcDecl.Doc == nil {
				insertIndex++
			}
			break
		}
	}

	// Create global variable/constant declarations.
	for _, global := range opts.globals {
		// Create an identifier for the global.
		ident := ast.NewIdent(global.name)

		// Create a value expression if provided.
		var valueExpr ast.Expr
		if global.value != "" {
			valueExpr, err = parser.ParseExprFrom(fileSet, "", []byte(global.value), parser.ParseComments)
			if err != nil {
				return "", err
			}
		}

		// Create a declaration based on the global type.
		var spec ast.Spec
		switch globalType {
		case GlobalTypeVar:
			spec = &ast.ValueSpec{
				Names:  []*ast.Ident{ident},
				Type:   ast.NewIdent(global.varType),
				Values: []ast.Expr{valueExpr},
			}
		case GlobalTypeConst:
			spec = &ast.ValueSpec{
				Names:  []*ast.Ident{ident},
				Type:   ast.NewIdent(global.varType),
				Values: []ast.Expr{valueExpr},
			}
		default:
			return "", errors.Errorf("unsupported global type: %s", string(globalType))
		}

		// Insert the declaration after the import section or package declaration if no imports.
		f.Decls = append(
			f.Decls[:insertIndex],
			append([]ast.Decl{
				&ast.GenDecl{
					TokPos: 1,
					Tok:    token.Lookup(string(globalType)),
					Specs:  []ast.Spec{spec},
				},
			}, f.Decls[insertIndex:]...)...)
		insertIndex++
	}

	f.Comments = cmap.Filter(f).Comments()

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	// Return the modified content.
	return buf.String(), nil
}

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
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	// Append the function declaration to the file's declarations.
	f.Decls = append(f.Decls, funcDecl)

	f.Comments = cmap.Filter(f).Comments()

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type (
	// structOpts represent the options for structs.
	structOpts struct {
		values []structValue
	}

	// StructOpts configures struct changes.
	StructOpts func(*structOpts)

	structValue struct {
		value     string
		valueType string
	}
)

// AppendStructValue add a new value inside struct. For instances,
// the struct have only one field 'test struct{ test1 string }' and we want to add
// the `test2 int` the result will be 'test struct{ test1 string, test int }'.
func AppendStructValue(value, valueType string) StructOpts {
	return func(c *structOpts) {
		c.values = append(c.values, structValue{
			value:     value,
			valueType: valueType,
		})
	}
}

func newStructOptions() structOpts {
	return structOpts{
		values: make([]structValue, 0),
	}
}

// ModifyStruct modifies a struct in the provided Go source code.
func ModifyStruct(fileContent, structName string, options ...StructOpts) (string, error) {
	// Apply struct options.
	opts := newStructOptions()
	for _, o := range options {
		o(&opts)
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	// Locate and modify the struct declaration.
	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true // Not a general declaration, continue searching.
		}

		for _, spec := range genDecl.Specs {
			// Check for type specification.
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			// Check if the type is a struct.
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, v := range opts.values {
				// Add the new field to the struct.
				newField := &ast.Field{
					Names: []*ast.Ident{ast.NewIdent(v.value)},
					Type:  ast.NewIdent(v.valueType),
				}
				structType.Fields.List = append(structType.Fields.List, newField)
			}

			found = true
			return false // Stop searching once we modify the struct.
		}
		return true
	})
	if !found {
		return "", errors.Errorf("struct %q not found in file content", structName)
	}

	f.Comments = cmap.Filter(f).Comments()

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	// Return the modified content.
	return buf.String(), nil
}

type (
	// globalArrayOpts represent the options for globar array variables.
	globalArrayOpts struct {
		values []string
	}

	// GlobalArrayOpts configures global array variable changes.
	GlobalArrayOpts func(*globalArrayOpts)
)

// AppendGlobalArrayValue add a new value inside a global array variable. For instances,
// the variable have only one field '[]]' and we want to add
// the `test2 int` the result will be 'test struct{ test1 string, test int }'.
func AppendGlobalArrayValue(value string) GlobalArrayOpts {
	return func(c *globalArrayOpts) {
		c.values = append(c.values, value)
	}
}

func newGlobalArrayOptions() globalArrayOpts {
	return globalArrayOpts{
		values: make([]string, 0),
	}
}

// ModifyGlobalArrayVar modifies an array global array variable in the provided Go source code by appending new values.
func ModifyGlobalArrayVar(fileContent, globalName string, options ...GlobalArrayOpts) (string, error) {
	opts := newGlobalArrayOptions()
	for _, o := range options {
		o(&opts)
	}

	fileSet := token.NewFileSet()

	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) == 0 || valueSpec.Names[0].Name != globalName {
				continue
			}

			if len(valueSpec.Values) == 0 {
				continue
			}

			compLit, ok := valueSpec.Values[0].(*ast.CompositeLit)
			if !ok {
				continue
			}

			file := fileSet.File(compLit.Pos())
			maxOffset := file.Offset(compLit.Rbrace)
			for _, elt := range compLit.Elts {
				if pos := elt.End(); pos.IsValid() {
					offset := file.Offset(pos)
					if offset > maxOffset {
						maxOffset = offset
					}
				}
			}

			for i, v := range opts.values {
				// Advance position
				insertOffset := maxOffset + i
				insertPos := file.Pos(insertOffset)

				value := ast.NewIdent(v)
				value.NamePos = insertPos

				compLit.Elts = append(compLit.Elts, value)
				compLit.Rbrace += token.Pos(i + 1)
			}

			// Ensure closing brace is on a new line and add comma after last element
			if len(compLit.Elts) > 0 {
				last := compLit.Elts[len(compLit.Elts)-1]
				if file.Line(compLit.Rbrace) == file.Line(last.End())-1 {
					// Add comma after last element
					file.AddLine(file.Offset(compLit.Rbrace))
					compLit.Rbrace += token.Pos(1)
				}
			}

			found = true
			return false // Stop searching once we modify the struct.
		}
		return true
	})

	if !found {
		return "", errors.Errorf("global array %q not found in file content", globalName)
	}

	f.Comments = cmap.Filter(f).Comments()

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}
