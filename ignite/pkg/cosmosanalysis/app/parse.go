package app

import (
	"go/ast"
	"go/build"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/xast"
)

func parseAppModulesFromIdent(n *ast.Ident, pkgDir string) ([]string, error) {
	if n == nil {
		return nil, errors.Errorf("nil node")
	}
	if n.Obj == nil {
		// The variable is defined in another file within the same package
		return parseAppModulesFromPkgIdent(n.Name, pkgDir)
	}

	// Variable declaration within the same file
	decl, ok := n.Obj.Decl.(*ast.ValueSpec)
	if !ok {
		return nil, newUnexpectedTypeErr(n.Obj.Decl)
	}

	values, ok := decl.Values[0].(*ast.CompositeLit)
	if !ok {
		return nil, newUnexpectedTypeErr(decl.Values[0])
	}

	var pkgNames []string
	for _, e := range values.Elts {
		switch v := e.(type) {
		case *ast.CompositeLit:
			vt, ok := v.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			if pkg, ok := vt.X.(*ast.Ident); ok {
				pkgNames = append(pkgNames, pkg.Name)
			}
		case *ast.CallExpr:
			ps, err := parsePkgNameFromCall(v, pkgDir)
			if err != nil {
				return nil, err
			}
			pkgNames = append(pkgNames, ps...)
		default:
			return nil, newUnexpectedTypeErr(e)
		}
	}

	return pkgNames, nil
}

func parseAppModulesFromSelectorExpr(n *ast.SelectorExpr, pkgDir string, fileImports map[string]string) ([]string, error) {
	// Get the name of the package where the app modules variable is defined
	ident, ok := n.X.(*ast.Ident)
	if !ok {
		return nil, newUnexpectedTypeErr(n.X)
	}

	// Get the import path of that package and resolve the full path to it
	ctx := build.Default
	ctx.Dir = pkgDir

	pkg, err := ctx.Import(fileImports[ident.Name], "", build.FindOnly)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return parseAppModulesFromPkgIdent(n.Sel.Name, pkg.Dir)
}

func parseAppModulesFromPkgIdent(identName, pkgDir string) ([]string, error) {
	pkg, _, err := xast.ParseDir(pkgDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, f := range pkg.Files {
		ident := f.Scope.Objects[identName]
		if ident == nil {
			continue
		}

		var pkgNames []string
		switch decl := ident.Decl.(type) {
		case *ast.FuncDecl:
			pkgNames, err = parsePkgNameFromFuncDecl(decl, pkgDir)
			if err != nil {
				return nil, err
			}

		case *ast.ValueSpec:
			values, ok := decl.Values[0].(*ast.CompositeLit)
			if !ok {
				continue
			}

			for _, e := range values.Elts {
				switch v := e.(type) {
				case *ast.CompositeLit:
					// The app module is defined using a struct
					ps, err := parsePkgNameFromCompositeLit(v, pkgDir)
					if err != nil {
						return nil, err
					}
					pkgNames = append(pkgNames, ps...)

				case *ast.CallExpr:
					// The app module is defined using a function call that returns it
					ps, err := parsePkgNameFromCall(v, pkgDir)
					if err != nil {
						return nil, err
					}
					pkgNames = append(pkgNames, ps...)
				}
			}
		}
		imports := goanalysis.FormatImports(f)
		var packages []string
		for _, pkgName := range pkgNames {
			if p := imports[pkgName]; p != "" {
				packages = append(packages, p)
			}
		}
		return packages, nil
	}
	return nil, errors.Errorf("unable to find identifier %s in package %s", identName, pkgDir)
}

func parsePkgNameFromCompositeLit(n *ast.CompositeLit, pkgDir string) ([]string, error) {
	switch v := n.Type.(type) {
	case *ast.SelectorExpr:
		if pkg, ok := v.X.(*ast.Ident); ok {
			return []string{pkg.Name}, nil
		}
		return nil, newUnexpectedTypeErr(v.X)
	case *ast.ArrayType:
		var pkgs []string
		for _, elt := range n.Elts {
			switch elt := elt.(type) {
			case *ast.CompositeLit:
				ps, err := parsePkgNameFromCompositeLit(elt, pkgDir)
				if err != nil {
					return nil, err
				}
				pkgs = append(pkgs, ps...)
			case *ast.CallExpr:
				ps, err := parsePkgNameFromCall(elt, pkgDir)
				if err != nil {
					return nil, err
				}
				pkgs = append(pkgs, ps...)
			default:
				return nil, newUnexpectedTypeErr(elt)
			}
		}
		return pkgs, nil
	}
	return nil, newUnexpectedTypeErr(n.Type)
}

func parsePkgNameFromCall(n *ast.CallExpr, pkgDir string) ([]string, error) {
	switch f := n.Fun.(type) {

	case *ast.SelectorExpr:
		if pkg, ok := f.X.(*ast.Ident); ok {
			return []string{pkg.Name}, nil
		}
		return nil, newUnexpectedTypeErr(f.X)

	case *ast.Ident:
		if f.Name == "append" {
			var pkgs []string
			for _, arg := range n.Args {
				switch arg := arg.(type) {
				case *ast.CompositeLit:
					ps, err := parsePkgNameFromCompositeLit(arg, pkgDir)
					if err != nil {
						return nil, err
					}
					pkgs = append(pkgs, ps...)

				case *ast.CallExpr:
					ps, err := parsePkgNameFromCall(arg, pkgDir)
					if err != nil {
						return nil, err
					}
					pkgs = append(pkgs, ps...)

				default:
					return nil, newUnexpectedTypeErr(arg)
				}
			}
			return pkgs, nil
		}
		// read func return statement
		if f.Obj == nil {
			return parseAppModulesFromPkgIdent(f.Name, pkgDir)
		}
		fd, ok := f.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			return nil, newUnexpectedTypeErr(f.Obj.Decl)
		}
		return parsePkgNameFromFuncDecl(fd, pkgDir)
	}
	return nil, newUnexpectedTypeErr(n.Fun)
}

func parsePkgNameFromFuncDecl(f *ast.FuncDecl, pkgDir string) ([]string, error) {
	lastBodyStmt := f.Body.List[len(f.Body.List)-1]
	ret, ok := lastBodyStmt.(*ast.ReturnStmt)
	if !ok {
		return nil, newUnexpectedTypeErr(lastBodyStmt)
	}
	c, ok := ret.Results[0].(*ast.CompositeLit)
	if !ok {
		return nil, newUnexpectedTypeErr(ret.Results[0])
	}
	return parsePkgNameFromCompositeLit(c, pkgDir)
}
