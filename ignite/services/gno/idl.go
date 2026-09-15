package gno

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// IDL describes a gno package interface in an Anchor-like shape, so client
// generators and external tools can consume it.
type IDL struct {
	// Version is the IDL schema version.
	Version string `json:"version"`
	// Name is the package name.
	Name string `json:"name"`
	// PkgPath is the gno.land module path.
	PkgPath string `json:"pkgPath"`
	// Kind is "realm" or "package".
	Kind string `json:"kind"`
	// Instructions are the callable functions of the package.
	Instructions []IDLInstruction `json:"instructions"`
}

// IDLInstruction describes one exported function.
type IDLInstruction struct {
	// Name is the exported function name.
	Name string `json:"name"`
	// Args are the function parameters.
	Args []IDLArg `json:"args"`
	// Realm reports whether the function requires realm context (mutates
	// state): its first parameter has type realm.
	Realm bool `json:"realm"`
}

// IDLArg describes one function parameter.
type IDLArg struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// idlVersion is the IDL schema version.
const idlVersion = "1.0"

// GenerateIDL builds the IDL of the gno package at dir. The module path is
// read from gnomod.toml.
func GenerateIDL(dir string) (*IDL, error) {
	pkgPath, err := parseGnoModDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading gnomod.toml in %s (scaffold one with `ignite scaffold realm`)", dir)
	}
	kind := "package"
	if strings.Contains(pkgPath, "/r/") || strings.HasSuffix(filepath.Dir(pkgPath), "/r") {
		kind = "realm"
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	idl := &IDL{
		Version: idlVersion,
		Name:    pathPkgName(pkgPath),
		PkgPath: pkgPath,
		Kind:    kind,
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
		idl.Name = f.Name.Name
		collectIDLFuncs(idl, f)
	}

	sort.Slice(idl.Instructions, func(i, j int) bool {
		return idl.Instructions[i].Name < idl.Instructions[j].Name
	})
	if len(idl.Instructions) == 0 {
		return nil, errors.Errorf("no exported functions found in %s", dir)
	}
	return idl, nil
}

// collectIDLFuncs appends the exported functions of a parsed file to idl.
func collectIDLFuncs(idl *IDL, f *ast.File) {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !fn.Name.IsExported() {
			continue
		}
		instr := IDLInstruction{Name: fn.Name.Name}
		if fn.Type.Params != nil {
			for _, p := range fn.Type.Params.List {
				argType := exprString(p.Type)
				if argType == "realm" && len(instr.Args) == 0 {
					instr.Realm = true
					continue
				}
				for _, name := range p.Names {
					instr.Args = append(instr.Args, IDLArg{Name: name.Name, Type: argType})
				}
				if len(p.Names) == 0 {
					instr.Args = append(instr.Args, IDLArg{Name: fmt.Sprintf("arg%d", len(instr.Args)), Type: argType})
				}
			}
		}
		idl.Instructions = append(idl.Instructions, instr)
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

// WriteIDLJSON marshals the IDL to indented JSON.
func (i *IDL) WriteIDLJSON() ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}
