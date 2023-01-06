package cosmosgen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

const (
	defaultSDKImport     = "github.com/cosmos/cosmos-sdk"
	moduleCacheNamespace = "generate.setup.module"
)

var protocGlobalInclude = xfilepath.List(
	xfilepath.JoinFromHome(xfilepath.Path("local/include")),
	xfilepath.JoinFromHome(xfilepath.Path(".local/include")),
)

type ModulesInPath struct {
	Path    string
	Modules []module.Module
}

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

	g.sdkImport = defaultSDKImport

	// Check if the Cosmos SDK import path points to a different path
	// and if so change the default one to the new location.
	for _, r := range modFile.Replace {
		if r.Old.Path == defaultSDKImport {
			g.sdkImport = r.New.Path
			break
		}
	}

	// Read the dependencies defined in the `go.mod` file
	g.deps, err = gomodule.ResolveDependencies(modFile)
	if err != nil {
		return err
	}

	// Discover any custom modules defined by the user's app
	g.appModules, err = g.discoverModules(g.appPath, g.protoDir)
	if err != nil {
		return err
	}

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

			modulesInPath = ModulesInPath{
				Path:    path,
				Modules: modules,
			}

			if err := moduleCache.Put(cacheKey, modulesInPath); err != nil {
				return err
			}
		}

		g.thirdModules[modulesInPath.Path] = append(g.thirdModules[modulesInPath.Path], modulesInPath.Modules...)
	}

	return nil
}

func (g *generator) resolveDependencyInclude() ([]string, error) {
	// Init paths with the global include paths for protoc
	paths, err := protocGlobalInclude()
	if err != nil {
		return nil, err
	}

	// Relative paths to proto directories
	protoDirs := append([]string{g.protoDir}, g.o.includeDirs...)

	// Create a list of proto import paths for the dependencies.
	// These paths will be available to be imported from the chain app's proto files.
	for rootPath, m := range g.thirdModules {
		// Skip modules without proto files
		if m == nil {
			continue
		}

		// Check each one of the possible proto directory names for the
		// current module and append them only when the directory exists.
		for _, d := range protoDirs {
			p := filepath.Join(rootPath, d)
			f, err := os.Stat(p)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}

				return nil, err
			}

			if f.IsDir() {
				paths = append(paths, p)
			}
		}
	}

	return paths, nil
}

func (g *generator) resolveInclude(path string) (paths []string, err error) {
	// Append chain app's proto paths
	paths = append(paths, filepath.Join(path, g.protoDir))
	for _, p := range g.o.includeDirs {
		paths = append(paths, filepath.Join(path, p))
	}

	// Append paths for dependencies that have protocol buffer files
	includePaths, err := g.resolveDependencyInclude()
	if err != nil {
		return nil, err
	}

	paths = append(paths, includePaths...)

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
