package xast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
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
	GlobalType uint8
)

const (
	GlobalTypeVar GlobalType = iota
	GlobalTypeConst
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

	// Determine the declaration type based on GlobalType.
	var declType string
	if globalType == GlobalTypeVar {
		declType = "var"
	} else if globalType == GlobalTypeConst {
		declType = "const"
	} else {
		return "", errors.Errorf("unsupported global type: %d", globalType)
	}

	// Create global variable/constant declarations.
	for _, global := range opts.globals {
		// Create an identifier for the global.
		ident := &ast.Ident{Name: global.name}

		// Create a value expression if provided.
		var valueExpr ast.Expr
		if global.value != "" {
			valueExpr, err = parser.ParseExpr(global.value)
			if err != nil {
				return "", err
			}
		}

		// Create a declaration based on the global type.
		var spec ast.Spec
		if declType == "var" {
			spec = &ast.ValueSpec{
				Names:  []*ast.Ident{ident},
				Type:   ast.NewIdent(global.varType),
				Values: []ast.Expr{valueExpr},
			}
		} else if declType == "const" {
			spec = &ast.ValueSpec{
				Names:  []*ast.Ident{ident},
				Type:   ast.NewIdent(global.varType),
				Values: []ast.Expr{valueExpr},
			}
		}

		// Insert the declaration after the import section or package declaration if no imports.
		f.Decls = append(f.Decls[:insertIndex], append([]ast.Decl{&ast.GenDecl{Tok: token.Lookup(declType), Specs: []ast.Spec{spec}}}, f.Decls[insertIndex:]...)...)
		insertIndex++
	}

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	// Return the modified content.
	return buf.String(), nil
}
