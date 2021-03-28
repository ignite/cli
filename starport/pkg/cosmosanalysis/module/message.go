package module

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// DiscoverMessages discovers sdk messages defined in a module that resides under modulePath.
func DiscoverMessages(modulePath string) (msgs []string, err error) {
	// parse go packages/files under modulePath.
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, modulePath, nil, 0)
	if err != nil {
		return nil, err
	}

	// collect all structs under modulePath to find out the ones that satisfy requirements.
	structs := make(map[string]requirements)

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				// look for struct methods.
				fdecl, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}

				// not a method.
				if fdecl.Recv == nil {
					return true
				}

				// fname is the name of method.
				fname := fdecl.Name.Name

				// find the struct name that method belongs to.
				t := fdecl.Recv.List[0].Type
				sident, ok := t.(*ast.Ident)
				if !ok {
					sexp, ok := t.(*ast.StarExpr)
					if !ok {
						return true
					}
					sident = sexp.X.(*ast.Ident)
				}
				sname := sident.Name

				// mark the requirement that this struct satisfies.
				if _, ok := structs[sname]; !ok {
					structs[sname] = newRequirements()
				}

				structs[sname][fname] = true

				return true
			})
		}
	}

	// checkRequirements checks if all requirements are satisfied.
	checkRequirements := func(r requirements) bool {
		for _, ok := range r {
			if !ok {
				return false
			}
		}
		return true
	}

	for name, reqs := range structs {
		if checkRequirements(reqs) {
			msgs = append(msgs, name)
		}
	}

	if len(msgs) == 0 {
		return nil, ErrModuleNotFound
	}

	return msgs, nil
}
