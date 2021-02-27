package module

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

// requirements holds a list of sdk.Msg's method names.
type requirements map[string]bool

// newRequirements creates a new list of requirements(method names) that needed by a sdk.Msg impl.
// TODO(low priority): dynamically get these from the source code of underlying version of the sdk.
func newRequirements() requirements {
	return requirements{
		"Reset":         false,
		"String":        false,
		"ProtoMessage":  false,
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
	bpath := gm.Module.Mod.Path

	// find proto packages that belong to modules under x/.
	xprotopkgs, err := findModuleProtoPkgs(sourcePath, bpath)
	if err != nil {
		return nil, err
	}

	var modules []Module

	for _, xproto := range xprotopkgs {
		var (
			rxpath = strings.TrimPrefix(xproto.GoImportName, bpath)
			xpath  = filepath.Join(sourcePath, rxpath)
		)

		xmsgs, err := DiscoverModule(xpath)
		if err != nil {
			return nil, err
		}

		var msgs []Msg

		for _, xmsg := range xmsgs {
			msgs = append(msgs, Msg{
				Name:     xmsg,
				URI:      fmt.Sprintf("%s.%s", xproto.Name, xmsg),
				FilePath: xproto.MessageByName(xmsg).Path,
			})
		}

		var (
			spname = strings.Split(xproto.Name, ".")
			m      = Module{
				Name: spname[len(spname)-1],
				Pkg:  xproto,
				Msgs: msgs,
			}
		)

		modules = append(modules, m)
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
				sexp, ok := fdecl.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					return true
				}

				sname := sexp.X.(*ast.Ident).Name

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
