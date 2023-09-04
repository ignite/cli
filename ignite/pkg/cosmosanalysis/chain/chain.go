package chain

import (
	"bytes"
	"context"
	"fmt"
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
	gomod "golang.org/x/mod/module"
)

const (
	moduleCacheNamespace       = "analyze.setup.module"
	includeProtoCacheNamespace = "analyze.includes.proto"
	includeCacheNamespace      = "analyze.includes.module"
)

var protocGlobalInclude = xfilepath.List(
	xfilepath.JoinFromHome(xfilepath.Path("local/include")),
	xfilepath.JoinFromHome(xfilepath.Path(".local/include")),
)
var protoFound = errors.New("Proto file found") // Hacky way to terminate dir walk early

type ModulesInPath struct {
	Path     string          `json:"path,omitempty"`
	Modules  []module.Module `json:"modules,omitempty"`
	HasProto bool            `json:"has_proto,omitempty"`
	Includes []string        `json:"includes,omitempty"`
}
type AllModules struct {
	ModulePaths []ModulesInPath `json:"modules_in_path,omitempty"`
}

type Analyzer struct {
	ctx          context.Context
	appPath      string
	protoDir     string
	sdkImport    string
	appModules   []module.Module
	cacheStorage cache.Storage
	deps         []gomod.Version
	includeDirs  []string
	thirdModules map[string][]ModulesInPath // app dependency-modules pair.
}

func GetModuleList(ctx context.Context, appPath, protoPath string, thirdPartyPaths []string) (*AllModules, error) {
	cacheStorage, err := cache.NewStorage(filepath.Join(appPath, "analyzer_cache.db"))
	if err != nil {
		return nil, err
	}
	var errb bytes.Buffer
	if err := cmdrunner.
		New(
			cmdrunner.DefaultStderr(&errb),
			cmdrunner.DefaultWorkdir(appPath),
		).Run(ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return nil, errors.Wrap(err, errb.String())
	}

	modFile, err := gomodule.ParseAt(appPath)
	a := &Analyzer{
		ctx:          ctx,
		appPath:      appPath,
		protoDir:     protoPath,
		includeDirs:  thirdPartyPaths,
		thirdModules: make(map[string][]ModulesInPath),
		cacheStorage: cacheStorage,
	}
	if err != nil {
		return nil, err
	}

	a.sdkImport = cosmosver.CosmosModulePath

	// Check if the Cosmos SDK import path points to a different path
	// and if so change the default one to the new location.
	for _, r := range modFile.Replace {
		if r.Old.Path == cosmosver.CosmosModulePath {
			a.sdkImport = r.New.Path
			break
		}
	}

	// Read the dependencies defined in the `go.mod` file
	a.deps, err = gomodule.ResolveDependencies(modFile, true)
	if err != nil {
		return nil, err
	}

	// Discover any custom modules defined by the user's app
	a.appModules, err = a.discoverModules(a.appPath, a.protoDir)
	if err != nil {
		return nil, err
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
	moduleCache := cache.New[ModulesInPath](a.cacheStorage, moduleCacheNamespace)
	for _, dep := range a.deps {
		// Try to get the cached list of modules for the current dependency package
		cacheKey := cache.Key(dep.Path, dep.Version)
		modulesInPath, err := moduleCache.Get(cacheKey)
		if err != nil && !errors.Is(err, cache.ErrorNotFound) {
			return nil, err
		}

		// Discover the modules of the dependency package when they are not cached
		if errors.Is(err, cache.ErrorNotFound) {
			// Get the absolute path to the package's directory
			path, err := gomodule.LocatePath(a.ctx, a.cacheStorage, a.appPath, dep)
			if err != nil {
				return nil, err
			}
			hasProto, err := a.checkForProto(path)
			if err != nil {
				return nil, err
			}
			if hasProto {
				// Discover any modules defined by the package
				modules, err := a.discoverModules(path, "")
				if err != nil {
					return nil, err
				}
				includePaths, err := a.getProtoIncludeFolders(path)
				if err != nil {
					return nil, err
				}
				modulesInPath = ModulesInPath{
					Path:     path,
					Modules:  modules,
					HasProto: true,
					Includes: includePaths,
				}
			} else {
				modulesInPath = ModulesInPath{
					Path:     path,
					Modules:  []module.Module{},
					HasProto: false,
				}
			}
			if err := moduleCache.Put(cacheKey, modulesInPath); err != nil {
				return nil, err
			}
		}
		if modulesInPath.HasProto {
			a.thirdModules[modulesInPath.Path] = append(a.thirdModules[modulesInPath.Path], modulesInPath)
		}
	}

	if err != nil {
		return nil, err
	}
	var modulelist []ModulesInPath
	modulelist = append(modulelist, ModulesInPath{Path: appPath, Modules: a.appModules})
	for _, modules := range a.thirdModules {
		{
			modulelist = append(modulelist, modules...)
		}
	}
	allModules := &AllModules{
		ModulePaths: modulelist,
	}
	fmt.Println(allModules)
	return allModules, nil
}

func (a *Analyzer) getProtoIncludeFolders(modPath string) ([]string, error) {
	// Read the mod file for this module
	modFile, err := gomodule.ParseAt(modPath)
	if err != nil {
		return nil, err
	}
	includePaths := []string{}
	// Get the imports/deps from the mod file (include indirect)
	deps, err := gomodule.ResolveDependencies(modFile, true)
	if err != nil {
		return nil, err
	}
	// Initialize include cache for proto checking (lots of common includes across modules. No need to traverse repeatedly)
	includeProtoCache := cache.New[bool](a.cacheStorage, includeProtoCacheNamespace)

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
			path, err := gomodule.LocatePath(a.ctx, a.cacheStorage, a.appPath, dep)
			if err != nil {
				return nil, err
			}
			hasProto, err = a.checkForProto(path)
			if err != nil {
				return nil, err
			}
		}

		if hasProto {
			path, err := gomodule.LocatePath(a.ctx, a.cacheStorage, a.appPath, dep)
			if err != nil {
				return nil, err
			}
			includePaths = append(includePaths, path)
		}
	}
	return includePaths, nil
}

func (a *Analyzer) checkForProto(modpath string) (bool, error) {
	err := filepath.Walk(modpath,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".proto" {
				return protoFound
			}
			return nil
		})
	if err == protoFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func (a *Analyzer) discoverModules(path, protoDir string) ([]module.Module, error) {
	var filteredModules []module.Module

	modules, err := module.Discover(a.ctx, a.appPath, path, protoDir)
	if err != nil {
		return nil, err
	}

	protoPath := filepath.Join(path, a.protoDir)

	for _, m := range modules {
		if !strings.HasPrefix(m.Pkg.Path, protoPath) {
			continue
		}

		filteredModules = append(filteredModules, m)
	}
	return filteredModules, nil
}
