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
	"github.com/ignite/cli/ignite/pkg/cosmosbuf"
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

// protoIncludes contains proto include paths for a package.
type protoIncludes struct {
	// Paths is a list of proto include paths.
	Paths []string

	// BufPath is the path to the Buf config file when it exists.
	BufPath string

	// ProtoPath contains the path to the package's proto directory.
	ProtoPath string
}

// protoAnalysis contains proto module analysis data for a Go package dependency.
type protoAnalysis struct {
	// Path is the full path to the Go dependency
	Path string

	// Modules contains the proto modules analysis data.
	// The list is empty when the Go package has no proto files.
	Modules []module.Module

	// Includes contain proto include paths.
	// These paths should be used when generating code.
	Includes protoIncludes
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
	if err != nil {
		return err
	}

	// Go through the Go dependencies of the user's app within go.mod, some of them
	// might be hosting Cosmos SDK modules that could be in use by user's blockchain.
	//
	// Cosmos SDK is a dependency of all blockchains, so it's absolute that we'll be
	// discovering all modules of the SDK as well during this process.
	//
	// Even if a dependency contains some SDK modules, not all of these modules could
	// be used by user's blockchain. This is fine, we can still generate TS clients
	// for those non modules, it is up to user to use (import in typescript) not use
	// generated modules.
	//
	// TODO: we can still implement some sort of smart filtering to detect non used
	// modules by the user's blockchain at some point, it is a nice to have.
	moduleCache := cache.New[protoAnalysis](g.cacheStorage, moduleCacheNamespace)
	for _, dep := range g.deps {
		// Try to get the cached list of modules for the current dependency package
		cacheKey := cache.Key(dep.Path, dep.Version)
		depInfo, err := moduleCache.Get(cacheKey)
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

			// Dependency/includes resolution per module is done to solve versioning issues
			var includes protoIncludes
			if len(modules) > 0 {
				includes, err = g.resolveIncludes(path)
				if err != nil {
					return err
				}
			}

			depInfo = protoAnalysis{
				Path:     path,
				Modules:  modules,
				Includes: includes,
			}

			if err := moduleCache.Put(cacheKey, depInfo); err != nil {
				return err
			}
		}

		g.thirdModules[depInfo.Path] = depInfo.Modules
		g.thirdModuleIncludes[depInfo.Path] = depInfo.Includes
	}

	return nil
}

func (g *generator) getProtoIncludeFolders(modPath string) []string {
	// Add default protoDir and default includeDirs
	includePaths := []string{filepath.Join(modPath, g.protoDir)}
	for _, dir := range g.opts.includeDirs {
		includePaths = append(includePaths, filepath.Join(modPath, dir))
	}
	return includePaths
}

func (g *generator) findBufPath(modpath string) (string, error) {
	var bufPath string
	err := filepath.WalkDir(modpath, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		if base == "buf.yaml" || base == "buf.yml" {
			bufPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return bufPath, nil
}

func (g *generator) generateBufIncludeFolder(modpath string) (string, error) {
	protoPath, err := os.MkdirTemp("", "includeFolder")
	if err != nil {
		return "", err
	}

	g.tmpDirs = append(g.tmpDirs, protoPath)

	err = g.buf.Export(g.ctx, modpath, protoPath)
	if err != nil {
		return "", err
	}
	return protoPath, nil
}

func (g *generator) resolveIncludes(path string) (protoIncludes, error) {
	// Init paths with the global include paths for protoc
	paths, err := protocGlobalInclude()
	if err != nil {
		return protoIncludes{}, err
	}

	includes := protoIncludes{Paths: paths}

	// Check that the app/package proto directory exists
	protoPath := filepath.Join(path, g.protoDir)
	fi, err := os.Stat(protoPath)
	if err != nil && !os.IsNotExist(err) {
		return protoIncludes{}, err
	} else if !fi.IsDir() {
		// Just return the global includes when a proto directory doesn't exist
		return includes, nil
	}

	// Add app's proto path to the list of proto paths
	includes.Paths = append(includes.Paths, protoPath)
	includes.ProtoPath = protoPath

	// Check if a Buf config file is present
	bufPath, err := g.findBufPath(protoPath)
	if err != nil {
		return includes, err
	}

	if bufPath != "" {
		includes.BufPath = bufPath

		// When a Buf config exists export all protos needed
		// to build the modules to a temporary include folder.
		// TODO: Should this be optional and not done for the app includes? Duplicates proto folder.
		bufProtoPath, err := g.generateBufIncludeFolder(protoPath)
		if err != nil && !errors.Is(err, cosmosbuf.ErrProtoFilesNotFound) {
			return protoIncludes{}, err
		}

		// Use exported files only when the path contains ".proto" files
		if bufProtoPath != "" {
			includes.Paths = append(includes.Paths, bufProtoPath)
			return includes, nil
		}
	}

	// When there is no Buf config add the configured directories
	// instead to keep the legacy (non Buf) behavior.
	includes.Paths = append(includes.Paths, g.getProtoIncludeFolders(path)...)

	return includes, nil
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
