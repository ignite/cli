package module

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

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

	// HTTPQueries is a list of module queries.
	HTTPQueries []HTTPQuery

	// Types is a list of proto types that might be used by module.
	Types []Type
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

// HTTPQuery is an sdk Query.
type HTTPQuery struct {
	// Name of the RPC func.
	Name string

	// FullName of the query with service name and rpc func name.
	FullName string

	// HTTPAnnotations keeps info about http annotations of query.
	Rules []protoanalysis.HTTPRule
}

// Type is a proto type that might be used by module.
type Type struct {
	Name string

	// FilePath is the path of the .proto file where message is defined at.
	FilePath string
}

type moduleDiscoverer struct {
	sourcePath string
	protoPath  string
	basegopath string
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
func Discover(ctx context.Context, sourcePath, protoDir string) ([]Module, error) {
	// find out base Go import path of the blockchain.
	gm, err := gomodule.ParseAt(sourcePath)
	if err != nil {
		if err == gomodule.ErrGoModNotFound {
			return []Module{}, nil
		}
		return nil, err
	}

	md := &moduleDiscoverer{
		protoPath:  filepath.Join(sourcePath, protoDir),
		sourcePath: sourcePath,
		basegopath: gm.Module.Mod.Path,
	}

	// find proto packages that belong to modules under x/.
	pkgs, err := md.findModuleProtoPkgs(ctx)
	if err != nil {
		return nil, err
	}

	var modules []Module

	for _, pkg := range pkgs {
		m, err := md.discover(pkg)
		if err != nil {
			return nil, err
		}

		modules = append(modules, m)
	}

	return modules, nil
}

// discover discovers and sdk module by a proto pkg.
func (d *moduleDiscoverer) discover(pkg protoanalysis.Package) (Module, error) {
	pkgrelpath := strings.TrimPrefix(pkg.GoImportPath(), d.basegopath)
	pkgpath := filepath.Join(d.sourcePath, pkgrelpath)

	msgs, err := cosmosanalysis.FindImplementation(pkgpath, messageImplementation)
	if err != nil {
		return Module{}, err
	}
	if len(msgs) == 0 {
		// No message means the module has not been found
		return Module{}, nil
	}

	namesplit := strings.Split(pkg.Name, ".")
	m := Module{
		Name: namesplit[len(namesplit)-1],
		Pkg:  pkg,
	}

	// fill sdk Msgs.
	for _, msg := range msgs {
		pkgmsg, err := pkg.MessageByName(msg)
		if err != nil {
			// no msg found in the proto defs corresponds to discovered sdk message.
			// if it cannot be found, nothing to worry about, this means that it is used
			// only internally and not open for actual use.
			continue
		}

		m.Msgs = append(m.Msgs, Msg{
			Name:     msg,
			URI:      fmt.Sprintf("%s.%s", pkg.Name, msg),
			FilePath: pkgmsg.Path,
		})
	}

	// isType whether if protomsg can be added as an any Type to Module.
	isType := func(protomsg protoanalysis.Message) bool {
		// do not use GenesisState type.
		if protomsg.Name == "GenesisState" {
			return false
		}

		// do not use if an SDK message.
		for _, msg := range msgs {
			if msg == protomsg.Name {
				return false
			}
		}

		// do not use if used as a request/return type type of an RPC.
		for _, s := range pkg.Services {
			for _, q := range s.RPCFuncs {
				if q.RequestType == protomsg.Name || q.ReturnsType == protomsg.Name {
					return false
				}
			}
		}

		return true
	}

	// fill types.
	for _, protomsg := range pkg.Messages {
		if !isType(protomsg) {
			continue
		}

		m.Types = append(m.Types, Type{
			Name:     protomsg.Name,
			FilePath: protomsg.Path,
		})
	}

	// fill queries.
	for _, s := range pkg.Services {
		for _, q := range s.RPCFuncs {
			if len(q.HTTPRules) == 0 {
				continue
			}
			m.HTTPQueries = append(m.HTTPQueries, HTTPQuery{
				Name:     q.Name,
				FullName: s.Name + q.Name,
				Rules:    q.HTTPRules,
			})
		}
	}

	return m, nil
}

func (d *moduleDiscoverer) findModuleProtoPkgs(ctx context.Context) ([]protoanalysis.Package, error) {
	// find out all proto packages inside blockchain.
	allprotopkgs, err := protoanalysis.Parse(ctx, d.protoPath)
	if err != nil {
		return nil, err
	}

	// filter out proto packages that do not represent x/ modules of blockchain.
	var xprotopkgs []protoanalysis.Package
	for _, pkg := range allprotopkgs {
		if !strings.HasPrefix(pkg.GoImportName, d.basegopath) {
			continue
		}

		xprotopkgs = append(xprotopkgs, pkg)
	}

	return xprotopkgs, nil
}
