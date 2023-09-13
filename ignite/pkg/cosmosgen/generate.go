package cosmosgen

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

const (
	moduleCacheNamespace       = "generate.setup.module"
	includeProtoCacheNamespace = "generator.includes.proto"
)

var protocGlobalInclude = xfilepath.List(
	xfilepath.JoinFromHome(xfilepath.Path("local/include")),
	xfilepath.JoinFromHome(xfilepath.Path(".local/include")),
)

type ModulesInPath struct {
	Path           string
	ModuleIncludes ModuleIncludes
}

var (
	fileFound    = errors.New("Searched file found")     // Hacky way to terminate dir walk early
	fileNotFound = errors.New("Searched file not found") // Hacky way to terminate dir walk early
)

func (g *generator) setup() (err error) {
	// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
	// blockchain apps. Generate should be aware of these and make them available to the blockchain
	// app that wants to generate code for its own proto.
	//
	// blockchain apps may use different versions of the SDK. following code first makes sure that
	// app's dependencies are download by 'go mod' and cached under the local filesystem.
	// and then, it determines which version of the SDK is used by the app and what is the absolute path
	// of its source code.
	var errb bytes.Buffer
	if err := cmdrunner.
		New(
			cmdrunner.DefaultStderr(&errb),
			cmdrunner.DefaultWorkdir(g.appPath),
		).Run(g.ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return errors.Wrap(err, errb.String())
	}

	modFile, err := gomodule.ParseAt(g.appPath)
	if err != nil {
		return err
	}

	g.sdkImport = cosmosver.CosmosModulePath

	// Check if the Cosmos SDK import path points to a different path
	// and if so change the default one to the new location.
	for _, r := range modFile.Replace {
		if r.Old.Path == cosmosver.CosmosModulePath {
			g.sdkImport = r.New.Path
			break
		}
	}

	// Read the dependencies defined in the `go.mod` file
	g.deps, err = gomodule.ResolveDependencies(modFile, false)
	if err != nil {
		return err
	}

	// Discover any custom modules defined by the user's app
	g.appModules, err = g.discoverModules(g.appPath, g.protoDir)
	if err != nil {
		return err
	}
	g.appIncludes, err = g.resolveIncludes(g.appPath)
	// Go through the Go dependencies of the user's app within go.mod, some of them might be hosting Cosmos SDK modules
	// that could be in use by user's blockchain.
	//
	// Cosmos SDK is a dependency of all blockchains, so it's absolute that we'll be discovering all modules of the
	// SDK as well during this process.
	//
	// Even if a dependency contains some SDK modules, not all of these modules could be used by user's blockchain.
	// this is fine, we can still generate TS clients for those non modules, it is up to user to use (import in typescript)
	// not use generated modules.
	//
	// TODO: we can still implement some sort of smart filtering to detect non used modules by the user's blockchain
	// at some point, it is a nice to have.
	moduleCache := cache.New[ModulesInPath](g.cacheStorage, moduleCacheNamespace)
	for _, dep := range g.deps {
		// Try to get the cached list of modules for the current dependency package
		cacheKey := cache.Key(dep.Path, dep.Version)
		modulesInPath, err := moduleCache.Get(cacheKey)
		if err != nil && !errors.Is(err, cache.ErrorNotFound) {
			return err
		}

		// Discover the modules of the dependency package when they are not cached
		if errors.Is(err, cache.ErrorNotFound) {
			// Get the absolute path to the package's directory
			path, err := gomodule.LocatePath(g.ctx, g.cacheStorage, g.appPath, dep)
			if err != nil {
				return err
			}

			// Discover any modules defined by the package
			modules, err := g.discoverModules(path, "")
			if err != nil {
				return err
			}
			includes := []string{}
			if len(modules) > 0 {

				includes, err = g.resolveIncludes(path) // For versioning issues, we do dependency/includes resolution per module
				if err != nil {
					return err
				}
			}
			modulesInPath = ModulesInPath{
				Path: path, // Each go module
				ModuleIncludes: ModuleIncludes{ // Has a set of sdk modules and a set of includes to buidl them
					Includes: includes,
					Modules:  modules,
				},
			}

			if err := moduleCache.Put(cacheKey, modulesInPath); err != nil {
				return err
			}
		}

		g.thirdModules[modulesInPath.Path] = modulesInPath.ModuleIncludes
	}

	return nil
}

func (g *generator) getProtoIncludeFolders(modPath string) ([]string, error) {
	/*
		// Read the mod file for this module
		modFile, err := gomodule.ParseAt(modPath)
		var errb bytes.Buffer
		if err != nil {
			return nil, nil
		}

		// Cache this go_module's dependencies locally
		if err := cmdrunner.
			New(
				cmdrunner.DefaultStderr(&errb),
				cmdrunner.DefaultWorkdir(modPath),
			).Run(g.ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
			return nil, nil
		}
	*/
	includePaths := []string{filepath.Join(modPath, g.protoDir)} // Add default protoDir and default includeDirs
	for _, dir := range g.o.includeDirs {
		includePaths = append(includePaths, filepath.Join(modPath, dir))
	}
	/*
		// Get the imports/deps from the mod file (exclude indirect)
		deps, err := gomodule.ResolveDependencies(modFile, false)
		if err != nil {
			return nil, err
		}
		// Initialize include cache for proto checking (lots of common includes across modules. No need to traverse repeatedly)
		includeProtoCache := cache.New[bool](g.cacheStorage, includeProtoCacheNamespace)

		for _, dep := range deps {

			// Check for proto file in this dependency
			cacheKey := cache.Key(dep.Path, dep.Version)
			hasProto, err := includeProtoCache.Get(cacheKey)

			// Return unexpected error
			if err != nil && !errors.Is(err, cache.ErrorNotFound) {
				return nil, err
			}

			// If result not already cached
			if errors.Is(err, cache.ErrorNotFound) {
				path, err := gomodule.LocatePath(g.ctx, g.cacheStorage, modPath, dep)
				if err != nil {
					return nil, err
				}
				hasProto, err = g.checkForProto(path)
				if err != nil {
					return nil, err
				}
			}
			if hasProto {
				path, err := gomodule.LocatePath(g.ctx, g.cacheStorage, modPath, dep)
				if err != nil {
					return nil, err
				}
				if !strings.Contains(path, "gogo/googleapis") {
					//	includePaths = append(includePaths, path)
				}
			}
		}
	*/
	return includePaths, nil
}

func (g *generator) checkForProto(modpath string) (bool, error) {
	err := filepath.WalkDir(modpath,
		func(path string, _ fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".proto" {
				return fileFound
			}
			return nil
		})
	if err == fileFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, fileNotFound
}

func (g *generator) checkForBuf(modpath string) (string, error) {
	var bufPath string
	err := filepath.WalkDir(modpath,
		func(path string, _ fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == "buf.yaml" {
				bufPath = path
				return fileFound
			}
			return nil
		})
	if err == fileFound {
		return bufPath, nil
	}
	if err != nil {
		return "", err
	}
	return "", fileNotFound
}

func (g *generator) generateBufIncludeFolder(modpath string) (string, error) {
	protoPath, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	err = g.buf.Export(g.ctx, modpath, protoPath)
	if err != nil {
		return "", err
	}
	return protoPath, nil
}

func (g *generator) resolveIncludes(path string) (paths []string, err error) {
	// Init paths with the global include paths for protoc
	paths, err = protocGlobalInclude()
	if err != nil {
		return nil, err
	}
	p := filepath.Join(path, g.protoDir) // Look for the default protoDir ("/proto")
	paths = append(paths, p)
	f, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if f.IsDir() {
		if !strings.Contains(p, "cosmos/cosmos-sdk") { // ignore cosmos-sdk
			bufPath, err := g.checkForBuf(p) // Look for a buf file to make our lives easier
			if err != nil && err != fileNotFound {
				return nil, err
			}
			if err == nil { // If it exists, use buf export to package all protos needed to build the modules in a temp folder
				protoPath, err := g.generateBufIncludeFolder(p)
				if err != nil {
					return nil, err
				}
				paths = append(paths, protoPath)
			}
			if bufPath == "" { // if it doesn't exist, get list of incldue folders the old-fashioned way
				protoPaths, err := g.getProtoIncludeFolders(path)
				if err != nil {
					return nil, err
				}
				paths = append(paths, protoPaths...)
			}
		}
	}

	return paths, nil
}

func (g *generator) discoverModules(path, protoDir string) ([]module.Module, error) {
	var filteredModules []module.Module

	modules, err := module.Discover(g.ctx, g.appPath, path, protoDir)
	if err != nil {
		return nil, err
	}

	protoPath := filepath.Join(path, g.protoDir)

	for _, m := range modules {
		if !strings.HasPrefix(m.Pkg.Path, protoPath) {
			continue
		}

		filteredModules = append(filteredModules, m)
	}

	return filteredModules, nil
}
