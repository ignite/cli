package module

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

// ErrModuleNotFound error returned when an sdk module cannot be found.
var ErrModuleNotFound = errors.New("sdk module not found")

// requirements holds a list of sdk.Msg's method names.
type requirements map[string]bool

// newRequirements creates a new list of requirements(method names) that needed by a sdk.Msg impl.
// TODO(low priority): dynamically get these from the source code of underlying version of the sdk.
func newRequirements() requirements {
	return requirements{
		"Route":         false,
		"Type":          false,
		"GetSigners":    false,
		"GetSignBytes":  false,
		"ValidateBasic": false,
	}
}

// Msgs is a module import path-sdk msgs pair.
type Msgs map[string][]string

// Module keeps metadata about a Cosmos SDK module.
type Module struct {
	// Name of the module.
	Name string

	// Pkg holds the proto package info.
	Pkg protoanalysis.Package

	// Msg is a list of sdk.Msg implementation of the module.
	Msgs []Msg
}

// Msg keeps metadata about an sdk.Msg implementation.
type Msg struct {
	// Name of the type.
	Name string

	// URI of the type.
	URI string

	// FilePath is the path of the .proto file where message is defined at.
	FilePath string
}

// Discover discovers and returns modules and their types that implements sdk.Msg.
// sourcePath is the root path of an sdk blockchain.
//
// discovery algorithm make use of proto definitions to discover modules inside the blockchain.
//
// checking whether a type implements sdk.Msg is done by running a simple algorithm of comparing method names
// of each type in a package with sdk.Msg's, which satisfies our needs for the time being.
// for a more opinionated check:
//   - go/types.Implements() might be utilized and as needed.
//   - instead of just comparing method names, their full signatures can be compared.
func Discover(sourcePath string) ([]Module, error) {
	// find out base Go import path of the blockchain.
	gm, err := gomodule.ParseAt(sourcePath)
	if err != nil {
		return nil, err
	}
	basegopath := gm.Module.Mod.Path

	// find proto packages that belong to modules under x/.
	pkgs, err := findModuleProtoPkgs(sourcePath, basegopath)
	if err != nil {
		return nil, err
	}

	var modules []Module

	// discover discovers and sdk module by a proto pkg.
	discover := func(pkg protoanalysis.Package) error {
		pkgrelpath := strings.TrimPrefix(pkg.GoImportPath(), basegopath)
		pkgpath := filepath.Join(sourcePath, pkgrelpath)

		msgs, err := DiscoverModule(pkgpath)
		if err == ErrModuleNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		var (
			spname = strings.Split(pkg.Name, ".")
			m      = Module{
				Name: spname[len(spname)-1],
				Pkg:  pkg,
			}
		)

		for _, msg := range msgs {
			pkgmsg, err := pkg.MessageByName(msg)
			if err != nil { // no msg found in the proto defs corresponds to discovered sdk message.
				return nil
			}

			m.Msgs = append(m.Msgs, Msg{
				Name:     msg,
				URI:      fmt.Sprintf("%s.%s", pkg.Name, msg),
				FilePath: pkgmsg.Path,
			})
		}

		modules = append(modules, m)

		return nil
	}

	for _, pkg := range pkgs {
		if err := discover(pkg); err != nil {
			return nil, err
		}
	}

	return modules, nil
}

// DiscoverModule discovers sdk messages defined in a module that resides under modulePath.
func DiscoverModule(modulePath string) (msgs []string, err error) {
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

func findModuleProtoPkgs(sourcePath, bpath string) ([]protoanalysis.Package, error) {
	// find out all proto packages inside blockchain.
	allprotopkgs, err := protoanalysis.DiscoverPackages(sourcePath)
	if err != nil {
		return nil, err
	}

	// filter out proto packages that do not represent x/ modules of blockchain.
	var xprotopkgs []protoanalysis.Package
	for _, pkg := range allprotopkgs {
		if !strings.HasPrefix(pkg.GoImportName, bpath) {
			continue
		}

		xprotopkgs = append(xprotopkgs, pkg)
	}

	return xprotopkgs, nil
}
