// Package goanalysis provides a toolset for statically analysing Go applications
package goanalysis

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	mainPackage     = "main"
	goFileExtension = ".go"
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

// GenVarExists finds a genesis variable goImport into the go file.
func GenVarExists(f *ast.File, goImport, methodSignature string) bool {
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
		genDecl, ok := d.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
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

				call, ok := vSpec.Values[0].(*ast.CallExpr)
				if !ok {
					continue
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				x, ok := sel.X.(*ast.Ident)
				if !ok {
					continue
				}

				cursorDeclaration := x.String() + "." + sel.Sel.String()
				if cursorDeclaration == methodDecl {
					return true
				}
			}
		}

	}
	return false
}

// FormatImports translate f.Imports into a map where name -> package.
// Name is the alias if declared, or the last element of the package path.
func FormatImports(f *ast.File) map[string]string {
	m := make(map[string]string) // name -> import
	for _, imp := range f.Imports {
		var importName string
		if imp.Name != nil {
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
