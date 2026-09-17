package gno

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ignite/cli/v30/ignite/pkg/errors"
)

// FuncSig describes one exported function of a gno package.
type FuncSig struct {
	// Name is the exported function name.
	Name string
	// Args are the function parameters.
	Args []FuncArg
	// Realm reports whether the function requires realm context (mutates
	// state): its first parameter has type realm.
	Realm bool
}

// FuncArg describes one function parameter.
type FuncArg struct {
	Name string
	Type string
}

// PackageSignatures holds the exported functions of a gno package.
type PackageSignatures struct {
	// Name is the package name.
	Name string
	// PkgPath is the gno.land module path.
	PkgPath string
	// Funcs are the exported functions, sorted by name.
	Funcs []FuncSig
}

// ExtractSignatures parses the gno package at dir and extracts its exported
// functions. The signatures are the input of client generators; the same
// data is available at runtime from a deployed chain via vm/qfuncs.
func ExtractSignatures(dir string) (*PackageSignatures, error) {
	pkgPath, err := parseGnoModDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading gnomod.toml in %s (scaffold one with `ignite scaffold realm`)", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	sig := &PackageSignatures{
		Name:    pathPkgName(pkgPath),
		PkgPath: pkgPath,
	}

	fset := token.NewFileSet()
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".gno") || strings.HasSuffix(name, "_test.gno") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %s", name)
		}
		sig.Name = f.Name.Name
		collectFuncSigs(sig, f)
	}

	sort.Slice(sig.Funcs, func(i, j int) bool {
		return sig.Funcs[i].Name < sig.Funcs[j].Name
	})
	if len(sig.Funcs) == 0 {
		return nil, errors.Errorf("no exported functions found in %s", dir)
	}
	return sig, nil
}

// collectFuncSigs appends the exported functions of a parsed file to sig.
func collectFuncSigs(sig *PackageSignatures, f *ast.File) {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !fn.Name.IsExported() {
			continue
		}
		funcSig := FuncSig{Name: fn.Name.Name}
		if fn.Type.Params != nil {
			for _, p := range fn.Type.Params.List {
				argType := exprString(p.Type)
				if argType == "realm" && len(funcSig.Args) == 0 {
					funcSig.Realm = true
					continue
				}
				for _, name := range p.Names {
					funcSig.Args = append(funcSig.Args, FuncArg{Name: name.Name, Type: argType})
				}
				if len(p.Names) == 0 {
					funcSig.Args = append(funcSig.Args, FuncArg{Name: fmt.Sprintf("arg%d", len(funcSig.Args)), Type: argType})
				}
			}
		}
		sig.Funcs = append(sig.Funcs, funcSig)
	}
}

// exprString renders a type expression back to source.
func exprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + exprString(t.Elt)
		}
		return "[" + exprString(t.Len) + "]" + exprString(t.Elt)
	case *ast.BasicLit:
		return t.Value
	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)
	case *ast.Ellipsis:
		return "..." + exprString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}
