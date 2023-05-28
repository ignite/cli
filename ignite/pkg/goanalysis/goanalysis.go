// Package goanalysis provides a toolset for statically analysing Go applications
package goanalysis

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	mainPackage     = "main"
	goFileExtension = ".go"
	toolsBuildTag   = "//go:build tools\n\n"
)

// ErrMultipleMainPackagesFound is returned when multiple main packages found while expecting only one.
var ErrMultipleMainPackagesFound = errors.New("multiple main packages found")

// DiscoverMain finds main Go packages under path.
func DiscoverMain(path string) (pkgPaths []string, err error) {
	uniquePaths := make(map[string]struct{})

	err = filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(filePath, goFileExtension) {
			return err
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}

		if mainPackage == parsed.Name.Name {
			dir := filepath.Dir(filePath)
			uniquePaths[dir] = struct{}{}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for path := range uniquePaths {
		pkgPaths = append(pkgPaths, path)
	}

	return pkgPaths, nil
}

// DiscoverOneMain tries to find only one main Go package under path.
func DiscoverOneMain(path string) (pkgPath string, err error) {
	pkgPaths, err := DiscoverMain(path)
	if err != nil {
		return "", err
	}

	count := len(pkgPaths)
	if count == 0 {
		return "", errors.New("main package cannot be found")
	}
	if count > 1 {
		return "", ErrMultipleMainPackagesFound
	}

	return pkgPaths[0], nil
}

// FuncVarExists finds a genesis variable goImport into the go file.
func FuncVarExists(f *ast.File, goImport, methodSignature string) bool {
	var (
		importAlias = ""
		imports     = FormatImports(f)
	)
	for alias, imp := range imports {
		if imp == goImport {
			importAlias = alias
		}
	}
	if importAlias == "" {
		return false
	}
	methodDecl := importAlias + "." + methodSignature

	for _, d := range f.Decls {
		if declVarExists(d, methodDecl) {
			return true
		}
	}
	return false
}

// declVarExists find a variable declaration into a ast.Decl.
func declVarExists(decl ast.Decl, methodDecl string) bool {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		for _, stmt := range d.Body.List {
			switch v := stmt.(type) {
			case *ast.DeclStmt:
				if declVarExists(v.Decl, methodDecl) {
					return true
				}
			case *ast.AssignStmt:
				if len(v.Rhs) == 0 {
					continue
				}
				decl, err := getCallExprName(v.Rhs[0])
				if err != nil {
					continue
				}
				if decl == methodDecl {
					return true
				}
			case *ast.IfStmt:
				stmt, ok := v.Init.(*ast.AssignStmt)
				if !ok || len(stmt.Rhs) == 0 {
					continue
				}
				decl, err := getCallExprName(stmt.Rhs[0])
				if err != nil {
					continue
				}
				if decl == methodDecl {
					return true
				}
			}
		}
	case *ast.GenDecl:
		decls, err := getGenDeclNames(d)
		if err != nil {
			return false
		}
		for _, decl := range decls {
			if decl == methodDecl {
				return true
			}
		}
	}
	return false
}

// getGenDeclNames returns a list of the method declaration inside the ast.GenDecl.
func getGenDeclNames(genDecl *ast.GenDecl) ([]string, error) {
	if genDecl.Tok != token.VAR {
		return nil, fmt.Errorf("genDecl is not a var token: %v", genDecl.Tok)
	}
	var decls []string
	for _, spec := range genDecl.Specs {
		valueDecl, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, id := range valueDecl.Names {
			vSpec, ok := id.Obj.Decl.(*ast.ValueSpec)
			if !ok || len(vSpec.Values) == 0 {
				continue
			}

			cursorDecl, err := getCallExprName(vSpec.Values[0])
			if err != nil {
				continue
			}
			decls = append(decls, cursorDecl)
		}
	}
	if len(decls) == 0 {
		return nil, fmt.Errorf("empty method declarations")
	}
	return decls, nil
}

// getGenDeclNames returns the method declaration inside the ast.Expr.
func getCallExprName(expr ast.Expr) (string, error) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", fmt.Errorf("expression is not a *ast.CallExpr: %v", expr)
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", fmt.Errorf("expression function is not a *ast.SelectorExpr: %v", call.Fun)
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("selector expression function is not a *ast.Ident: %v", sel.X)
	}

	return x.String() + "." + sel.Sel.String(), nil
}

// FormatImports translate f.Imports into a map where name -> package.
// Name is the alias if declared, or the last element of the package path.
func FormatImports(f *ast.File) map[string]string {
	m := make(map[string]string) // name -> import
	for _, imp := range f.Imports {
		var importName string
		if imp.Name != nil && imp.Name.Name != "_" && imp.Name.Name != "." {
			importName = imp.Name.Name
		} else {
			importParts := strings.Split(imp.Path.Value, "/")
			importName = importParts[len(importParts)-1]
		}

		name := strings.Trim(importName, "\"")
		m[name] = strings.Trim(imp.Path.Value, "\"")
	}
	return m
}

// UpdateInitImports helper function to remove and add underscore (init) imports to an *ast.File.
func UpdateInitImports(file *ast.File, writer io.Writer, importsToAdd, importsToRemove []string) error {
	// Create a map for faster lookup of items to remove
	importMap := make(map[string]bool)
	for _, astImport := range file.Imports {
		value, err := strconv.Unquote(astImport.Path.Value)
		if err != nil {
			return err
		}
		importMap[value] = true
	}
	for _, removeImport := range importsToRemove {
		importMap[removeImport] = false
	}
	for _, addImport := range importsToAdd {
		importMap[addImport] = true
	}

	// Add the imports
	for _, d := range file.Decls {
		if dd, ok := d.(*ast.GenDecl); ok {
			if dd.Tok == token.IMPORT {
				file.Imports = make([]*ast.ImportSpec, 0)
				dd.Specs = make([]ast.Spec, 0)
				for imp, exist := range importMap {
					if exist {
						spec := createUnderscoreImport(imp)
						file.Imports = append(file.Imports, spec)
						dd.Specs = append(dd.Specs, spec)
					}
				}
			}
		}
	}

	if _, err := writer.Write([]byte(toolsBuildTag)); err != nil {
		return fmt.Errorf("failed to write the build tag: %w", err)
	}
	return format.Node(writer, token.NewFileSet(), file)
}

// createUnderscoreImports helper function to create an AST underscore import with given path.
func createUnderscoreImport(imp string) *ast.ImportSpec {
	return &ast.ImportSpec{
		Name: ast.NewIdent("_"),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(imp),
		},
	}
}
