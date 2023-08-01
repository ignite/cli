package module

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/app"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
)

// Msgs is a module import path-sdk msgs pair.
type Msgs map[string][]string

// Module keeps metadata about a Cosmos SDK module.
type Module struct {
	// Name of the module.
	Name string

	// GoModulePath of the app where the module is defined.
	GoModulePath string

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

	// Paginated indicates that the query is using pagination.
	Paginated bool
}

// Type is a proto type that might be used by module.
type Type struct {
	Name string

	// FilePath is the path of the .proto file where message is defined at.
	FilePath string
}

type moduleDiscoverer struct {
	sourcePath        string
	protoPath         string
	basegopath        string
	registeredModules []string
}

// Discover discovers and returns modules and their types that are registered in the app
// chainRoot is the root path of the chain
// sourcePath is the root path of the go module which the proto dir is from
//
// Discovery algorithm make use of registered modules and proto definitions to find relevant
// registered modules. It does so by:
// 1. Getting all the registered Go modules from the app.
// 2. Parsing the proto files to find services and messages.
// 3. Check if the proto services are implemented in any of the registered modules.
func Discover(ctx context.Context, chainRoot, sourcePath, protoDir string) ([]Module, error) {
	// find out base Go import path of the blockchain.
	gm, err := gomodule.ParseAt(sourcePath)
	if err != nil {
		if errors.Is(err, gomodule.ErrGoModNotFound) {
			return []Module{}, nil
		}
		return nil, err
	}

	// Find all the modules registered by the app
	registeredModules, err := app.FindRegisteredModules(chainRoot)
	if err != nil {
		return nil, err
	}

	// Go import path of the app module
	basegopath := gm.Module.Mod.Path
	rootgopath := RootGoImportPath(basegopath)

	// Keep the custom app's modules and filter out the third
	// party ones that are not defined within the app.
	appModules := make([]string, 0)
	for _, m := range registeredModules {
		if strings.HasPrefix(m, rootgopath) {
			appModules = append(appModules, m)
		}
	}

	if len(appModules) == 0 {
		return []Module{}, nil
	}

	md := &moduleDiscoverer{
		protoPath:         filepath.Join(sourcePath, protoDir),
		sourcePath:        sourcePath,
		basegopath:        basegopath,
		registeredModules: appModules,
	}

	// Find proto packages that belong to modules under x/.
	pkgs, err := md.findModuleProtoPkgs(ctx)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return []Module{}, nil
	}

	var modules []Module

	for _, pkg := range pkgs {
		m, err := md.discover(pkg)
		if err != nil {
			return nil, err
		}

		if m.Name == "" {
			continue
		}

		modules = append(modules, m)
	}

	return modules, nil
}

// IsRootPath checks if a Go import path is a custom app module.
// Custom app modules are defined inside the "x" directory.
func IsRootPath(path string) bool {
	return filepath.Base(filepath.Dir(path)) == "x"
}

// RootPath returns the Go import path of a custom app module.
// An empty string is returned when the path doesn't belong to a custom module.
func RootPath(path string) string {
	for !IsRootPath(path) {
		if path = filepath.Dir(path); path == "." {
			return ""
		}
	}

	return path
}

// RootGoImportPath returns a Go import path with the version suffix removed.
func RootGoImportPath(importPath string) string {
	if p, v := path.Split(importPath); semver.IsValid(v) {
		return strings.TrimRight(p, "/")
	}

	return importPath
}

func extractRelPath(pkgGoImportPath, baseGoPath string) (string, error) {
	// Remove the import prefix to get the relative path
	if strings.HasPrefix(pkgGoImportPath, baseGoPath) {
		return strings.TrimPrefix(pkgGoImportPath, baseGoPath), nil
	}

	// When the import path prefix defined by the proto package
	// doesn't match the base Go import path it means that the
	// latter might have a version suffix and the former doesn't.
	if p, v := path.Split(baseGoPath); semver.IsValid(v) {
		// Use the import path without the version as prefix
		p = strings.TrimRight(p, "/")

		return strings.TrimPrefix(pkgGoImportPath, p), nil
	}

	return "", fmt.Errorf("proto go import %s is not relative to %s", pkgGoImportPath, baseGoPath)
}

// discover discovers and sdk module by a proto pkg.
func (d *moduleDiscoverer) discover(pkg protoanalysis.Package) (Module, error) {
	// Check if the proto package services are implemented
	// by any of the modules registered by the app.
	if ok, err := d.isPkgFromRegisteredModule(pkg); err != nil || !ok {
		return Module{}, err
	}

	pkgRelPath, err := extractRelPath(pkg.GoImportPath(), d.basegopath)
	if err != nil {
		return Module{}, err
	}

	// Find the `sdk.Msg` interface implementation
	pkgPath := filepath.Join(d.sourcePath, pkgRelPath)
	msgs, err := cosmosanalysis.FindImplementation(pkgPath, messageImplementation)
	if err != nil {
		return Module{}, err
	}

	if len(pkg.Services)+len(msgs) == 0 {
		return Module{}, nil
	}

	namesplit := strings.Split(pkg.Name, ".")
	m := Module{
		Name:         namesplit[len(namesplit)-1],
		GoModulePath: d.basegopath,
		Pkg:          pkg,
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

		// do not use if used as a request/return type of RPC.
		for _, s := range pkg.Services {
			for i, q := range s.RPCFuncs {
				if q.RequestType == protomsg.Name || q.ReturnsType == protomsg.Name {
					// Check if the service response message is using pagination and
					// update the RPC function. This is done here to avoid extra loops
					// just to update the pagination property.
					s.RPCFuncs[i].Paginated = hasPagination(protomsg)

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
				Name:      q.Name,
				FullName:  s.Name + q.Name,
				Rules:     q.HTTPRules,
				Paginated: q.Paginated,
			})
		}
	}

	return m, nil
}

func (d *moduleDiscoverer) findModuleProtoPkgs(ctx context.Context) ([]protoanalysis.Package, error) {
	// find out all proto packages inside blockchain.
	allprotopkgs, err := protoanalysis.Parse(ctx, nil, d.protoPath)
	if err != nil {
		return nil, err
	}

	// Remove version suffix from the Go import path if it exists.
	// Proto files might omit the version in the Go import path even
	// when the app module is using versioning.
	basegopath := RootGoImportPath(d.basegopath)

	// filter out proto packages that do not represent x/ modules of blockchain.
	var xprotopkgs []protoanalysis.Package
	for _, pkg := range allprotopkgs {
		if !strings.HasPrefix(pkg.GoImportPath(), basegopath) {
			continue
		}

		xprotopkgs = append(xprotopkgs, pkg)
	}

	return xprotopkgs, nil
}

// Checks if the proto package is implemented by any of the modules registered by the app.
func (d moduleDiscoverer) isPkgFromRegisteredModule(pkg protoanalysis.Package) (bool, error) {
	// Get the Go module import defined by the proto package
	goModuleImport := pkg.GoImportPath()

	// Try to get the Go import path of the custom app module that should implement
	// the package RPC services. When the import path doesn't import a package
	// from the standard "x" folder use the path defined by the proto package.
	// Using the custom app module root path guarantees that if the RPC services
	// implementation exists in the module it will always be found.
	if p := RootPath(goModuleImport); p != "" {
		goModuleImport = p
	}

	// Get a Go import path with the version suffix removed
	rootGoPath := RootGoImportPath(d.basegopath)

	for _, m := range d.registeredModules {
		// Extract the relative module path from the Go import path
		implRelPath := strings.TrimPrefix(m, d.basegopath)

		// Handle the case where the Go module has a version
		// suffix and the registered module doesn't.
		if implRelPath == m {
			implRelPath = strings.TrimPrefix(m, rootGoPath)
		}

		// Absolute path to the app module
		implPath := filepath.Join(d.sourcePath, implRelPath)

		for _, s := range pkg.Services {
			// List of the RPC service method names defined by the current proto service
			methods := make([]string, len(s.RPCFuncs))
			for i, rpcFunc := range s.RPCFuncs {
				methods[i] = rpcFunc.Name
			}

			// Find the Go implementation of the service defined in the proto package
			found, err := cosmosanalysis.DeepFindImplementation(implPath, methods)
			if err != nil {
				return false, err
			}

			// Sometimes the registered module definition is located in a different
			// directory branch from where the RPC implementation is defined. In this
			// case search the RPC implementation in all custom app module files.
			if len(found) == 0 && strings.HasPrefix(m, goModuleImport) {
				altImplRelPath := strings.TrimPrefix(goModuleImport, d.basegopath)
				if altImplRelPath == goModuleImport {
					altImplRelPath = strings.TrimPrefix(goModuleImport, rootGoPath)
				}

				altImplPath := filepath.Join(d.sourcePath, altImplRelPath)

				found, err = cosmosanalysis.DeepFindImplementation(altImplPath, methods)
				if err != nil {
					return false, err
				}
			}

			if len(found) > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

func hasPagination(msg protoanalysis.Message) bool {
	for _, fieldType := range msg.Fields {
		// Message field type suffix check to match common pagination types:
		//    cosmos.base.query.v1beta1.PageRequest
		//    cosmos.base.query.v1beta1.PageResponse
		if strings.HasSuffix(fieldType, "PageRequest") || strings.HasSuffix(fieldType, "PageResponse") {
			return true
		}
	}

	return false
}
