package xast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type (
	// importOpts represent the options for imp.
	importOpts struct {
		imports []imp
	}

	// ImportOptions configures code generation.
	ImportOptions func(*importOpts)

	imp struct {
		path string
		name string
	}
)

// WithImport add a new import at the end of the imports.
func WithImport(repo string) ImportOptions {
	return func(c *importOpts) {
		c.imports = append(c.imports, imp{
			path: repo,
			name: "",
		})
	}
}

// WithNamedImport add a new import with name at the end of the imports.
func WithNamedImport(name, repo string) ImportOptions {
	return func(c *importOpts) {
		c.imports = append(c.imports, imp{
			name: name,
			path: repo,
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
	if len(opts.imports) == 0 {
		return fileContent, nil
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	// Add new import statements.
	for _, importPath := range opts.imports {
		deleteImportsByPath(fileSet, f, importPath.path)

		if !astutil.AddNamedImport(fileSet, f, importPath.name, importPath.path) {
			if hasImport(f, importPath.name, importPath.path) {
				continue
			}
			return "", errors.Errorf("failed to add import %s - %s", importPath.name, importPath.path)
		}
	}
	ast.SortImports(fileSet, f)

	f.Comments = cmap.Filter(f).Comments()

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RemoveImports removes import statements from the existing import block in Go source code content.
func RemoveImports(fileContent string, imports ...ImportOptions) (string, error) {
	// apply global options.
	opts := newImportOptions()
	for _, o := range imports {
		o(&opts)
	}
	if len(opts.imports) == 0 {
		return fileContent, nil
	}

	fileSet := token.NewFileSet()

	// Parse the Go source code content.
	f, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", err
	}
	cmap := ast.NewCommentMap(fileSet, f, f.Comments)

	// Remove import statements.
	for _, importPath := range opts.imports {
		astutil.DeleteNamedImport(fileSet, f, importPath.name, importPath.path)
	}

	f.Comments = cmap.Filter(f).Comments()

	// Format the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, f); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func deleteImportsByPath(fileSet *token.FileSet, file *ast.File, path string) {
	names := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		if importPath(spec) == path {
			names = append(names, importName(spec))
		}
	}

	for _, name := range names {
		astutil.DeleteNamedImport(fileSet, file, name, path)
	}
}

func hasImport(file *ast.File, name, path string) bool {
	for _, spec := range file.Imports {
		if importName(spec) == name && importPath(spec) == path {
			return true
		}
	}
	return false
}

func importName(spec *ast.ImportSpec) string {
	if spec.Name == nil {
		return ""
	}
	return spec.Name.Name
}

func importPath(spec *ast.ImportSpec) string {
	value, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return ""
	}
	return value
}
