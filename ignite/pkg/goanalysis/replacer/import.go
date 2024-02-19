package replacer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
)

type (
	// importOpts represent the options for imp.
	importOpts struct {
		imports []imp
	}

	// ImportOptions configures code generation.
	ImportOptions func(*importOpts)

	imp struct {
		repo  string
		name  string
		index int
	}
)

// WithImport add a new import. If the index is -1 will append in the end of the imports.
func WithImport(repo, name string, index int) ImportOptions {
	return func(c *importOpts) {
		c.imports = append(c.imports, imp{
			repo:  repo,
			name:  name,
			index: index,
		})
	}
}

func newImportOptions() importOpts {
	return importOpts{
		imports: make([]imp, 0),
	}
}

// AppendImports appends import statements to the existing import block in Go source code content.
func AppendImports(fileContent string, imports ...ImportOptions) (string, error) {
	// apply global options.
	opts := newImportOptions()
	for _, o := range imports {
		o(&opts)
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Find the existing import declaration.
	var importDecl *ast.GenDecl
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT || len(genDecl.Specs) == 0 {
			continue
		}
		importDecl = genDecl
		break
	}

	if importDecl == nil {
		// If no existing import declaration found, create a new one.
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: make([]ast.Spec, 0),
		}
		f.Decls = append([]ast.Decl{importDecl}, f.Decls...)
	}

	// Check existing imports to avoid duplicates.
	existImports := make(map[string]struct{})
	for _, spec := range importDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		existImports[importSpec.Path.Value] = struct{}{}
	}

	// Add new import statements.
	for _, importStmt := range opts.imports {
		// Check if the import already exists.
		path := strconv.Quote(importStmt.repo)
		if _, ok := existImports[path]; ok {
			continue
		}
		// Create a new import spec.
		spec := &ast.ImportSpec{
			Name: &ast.Ident{
				Name: importStmt.name,
			},
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: path,
			},
		}

		switch {
		case importStmt.index == -1:
			// Append the new argument to the end
			importDecl.Specs = append(importDecl.Specs, spec)
		case importStmt.index >= 0 && importStmt.index <= len(importDecl.Specs):
			// Insert the new argument at the specified index
			importDecl.Specs = append(importDecl.Specs[:importStmt.index], append([]ast.Spec{spec}, importDecl.Specs[importStmt.index:]...)...)
		default:
			return "", fmt.Errorf("index out of range") // Stop the inspection, an error occurred
		}
	}

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}
