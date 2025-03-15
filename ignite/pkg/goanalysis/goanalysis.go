// Package goanalysis provides a toolset for statically analysing Go applications
package goanalysis

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
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
		return nil, errors.Errorf("genDecl is not a var token: %v", genDecl.Tok)
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
		return nil, errors.Errorf("empty method declarations")
	}
	return decls, nil
}

// getCallExprName returns the method declaration inside the ast.Expr.
func getCallExprName(expr ast.Expr) (string, error) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", errors.Errorf("expression is not a *ast.CallExpr: %v", expr)
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", errors.Errorf("expression function is not a *ast.SelectorExpr: %v", call.Fun)
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", errors.Errorf("selector expression function is not a *ast.Ident: %v", sel.X)
	}

	return x.String() + "." + sel.Sel.String(), nil
}

// FindBlankImports find all blank imports ('_') into a file.
func FindBlankImports(node *ast.File) []string {
	// Iterate through the import declarations and find the blank imports.
	blankImports := make([]string, 0)
	for _, imp := range node.Imports {
		if imp.Name != nil && imp.Name.Name == "_" {
			blankImports = append(blankImports, strings.ReplaceAll(imp.Path.Value, `"`, ""))
		}
	}
	return blankImports
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

// AddOrRemoveTools helper function to add or remove tools from the go.mod file.
func AddOrRemoveTools(f *modfile.File, writer io.Writer, importsToAdd, importsToRemove []string) error {
	for _, imp := range importsToAdd {
		_ = f.AddTool(imp)
	}

	for _, imp := range importsToRemove {
		_ = f.DropTool(imp)
	}

	data, err := f.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod file: %w", err)
	}

	_, err = writer.Write(data)
	return err
}

// ReplaceCode replace a function implementation into a package path. The method will find
// the method signature and re-write the method implementation based in the new function.
func ReplaceCode(pkgPath, oldFunctionName, newFunction string) (err error) {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	all, err := parser.ParseDir(fileSet, absPath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range all {
		for _, f := range pkg.Files {
			found := false
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					// Check if the function has the name you want to replace.
					if funcDecl.Name.Name == oldFunctionName {
						// Replace the function body with the replacement code.
						replacementExpr, err := parser.ParseExpr(newFunction)
						if err != nil {
							return false
						}
						funcDecl.Body = &ast.BlockStmt{List: []ast.Stmt{
							&ast.ExprStmt{X: replacementExpr},
						}}
						found = true
						return false
					}
				}
				return true
			})
			if !found {
				continue
			}
			filePath := fileSet.Position(f.Package).Filename
			outFile, err := os.Create(filePath)
			if err != nil {
				return err
			}

			// Format and write the modified AST to the output file.
			if err := format.Node(outFile, fileSet, f); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// HasAnyStructFieldsInPkg finds the struct within a package folder and checks
// if any of the fields are defined in the struct.
func HasAnyStructFieldsInPkg(pkgPath, structName string, fields []string) (bool, error) {
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return false, err
	}
	fileSet := token.NewFileSet()
	all, err := parser.ParseDir(fileSet, absPath, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}

	fieldsNames := make(map[string]struct{})
	for _, field := range fields {
		fieldsNames[strings.ToLower(field)] = struct{}{}
	}

	exist := false
	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if _, ok := typeSpec.Type.(*ast.StructType); !ok ||
					typeSpec.Name.Name != structName ||
					typeSpec.Type == nil {
					return true
				}

				// Check if the struct has fields.
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// Iterate through the fields of the struct.
				for _, field := range structType.Fields.List {
					for _, fieldName := range field.Names {
						if _, ok := fieldsNames[strings.ToLower(fieldName.Name)]; !ok {
							continue
						}
						exist = true
						return false
					}
				}
				return true
			})
		}
	}
	return exist, nil
}
