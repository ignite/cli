package app

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func parseImports(n []*ast.ImportSpec) map[string]string {
	if n == nil {
		return nil
	}

	imports := make(map[string]string)
	for _, imp := range n {
		var name string

		if imp.Name != nil {
			// Use the import alias as name
			name = imp.Name.Name
		} else {
			// Split the import path and get the last element to use as name
			parts := strings.Split(imp.Path.Value, "/")
			name = parts[len(parts)-1]
		}

		name = strings.Trim(name, `"`)
		imports[name] = strings.Trim(imp.Path.Value, `"`)
	}

	return imports
}

func parseAppModulesFromIdent(n *ast.Ident, pkgDir string) (pkgNames []string) {
	if n == nil {
		return
	} else if n.Obj == nil {
		// The variable is defined in another file within the same package
		return parseAppModulesFromPkgIdent(n.Name, pkgDir)
	}

	// Variable declaration within the same file
	decl, ok := n.Obj.Decl.(*ast.ValueSpec)
	if !ok {
		return
	}

	values, ok := decl.Values[0].(*ast.CompositeLit)
	if !ok {
		return
	}

	for _, e := range values.Elts {
		v, ok := e.(*ast.CompositeLit)
		if !ok {
			continue
		}

		vt, ok := v.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		if pkg, ok := vt.X.(*ast.Ident); ok {
			pkgNames = append(pkgNames, pkg.Name)
		}
	}

	return pkgNames
}

func parseAppModulesFromSelectorExpr(n *ast.SelectorExpr, pkgDir string, pkgs map[string]string) []string {
	if n == nil {
		return nil
	}

	// Get the name of the package where the app modules variable is defined
	ident, ok := n.X.(*ast.Ident)
	if !ok {
		return nil
	}

	// Get the import path of that package and resolve the full path to it
	ctx := build.Default
	ctx.Dir = pkgDir

	pkg, err := ctx.Import(pkgs[ident.Name], "", build.FindOnly)
	if err != nil {
		return nil
	}

	return parseAppModulesFromPkgIdent(n.Sel.Name, pkg.Dir)
}

func parseAppModulesFromPkgIdent(identName, pkgDir string) (packages []string) {
	files, err := os.ReadDir(pkgDir)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() || !strings.HasSuffix(fileName, ".go") {
			continue
		}

		fileName = filepath.Join(pkgDir, fileName)
		f, err := parser.ParseFile(fset, fileName, nil, 0)
		if err != nil {
			return
		}

		ident := f.Scope.Objects[identName]
		if ident == nil {
			continue
		}

		decl, ok := ident.Decl.(*ast.ValueSpec)
		if !ok {
			continue
		}

		values, ok := decl.Values[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		imports := parseImports(f.Imports)

		for _, e := range values.Elts {
			var pkgName string

			switch v := e.(type) {
			case *ast.CompositeLit:
				// The app module is defined using a struct
				pkgName = parsePkgNameFromCompositeLit(v)
			case *ast.CallExpr:
				// The app module is defined using a function call that returns it
				pkgName = parsePkgNameFromCall(v)
			}

			if p := imports[pkgName]; p != "" {
				packages = append(packages, p)
			}
		}
	}

	return packages
}

func parsePkgNameFromCompositeLit(n *ast.CompositeLit) string {
	s, ok := n.Type.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	if pkg, ok := s.X.(*ast.Ident); ok {
		return pkg.Name
	}

	return ""
}

func parsePkgNameFromCall(n *ast.CallExpr) string {
	s, ok := n.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	if pkg, ok := s.X.(*ast.Ident); ok {
		return pkg.Name
	}

	return ""
}
