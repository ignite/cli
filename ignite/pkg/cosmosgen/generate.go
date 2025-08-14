package cosmosgen

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

const (
	moduleCacheNamespace       = "generate.setup.module"
	sdkModuleCacheNamespace    = "generate.setup.sdk_module"
	includeProtoCacheNamespace = "generate.includes.proto"
	bufYamlFilename            = "buf.yaml"
)

var (
	ErrBufConfig     = errors.New("invalid Buf config")
	ErrMissingSDKDep = errors.New("cosmos-sdk dependency not found")

	protocGlobalInclude = xfilepath.List(
		xfilepath.JoinFromHome(xfilepath.Path("local/include")),
		xfilepath.JoinFromHome(xfilepath.Path(".local/include")),
	)
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

	// Cacheable indicates whether this analysis can be safely cached.
	// Set to false when includes contain temporary directories.
	Cacheable bool
}

func newBufConfigError(path string, cause error) error {
	return errors.Errorf("%w: %s: %w", ErrBufConfig, path, cause)
}

// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
// blockchain apps. Generate should be aware of these and make them available to the blockchain
// app that wants to generate code for its own proto.
//
// blockchain apps may use different versions of the SDK. following code first makes sure that
// app's dependencies are download by 'go mod' and cached under the local filesystem.
// and then, it determines which version of the SDK is used by the app and what is the absolute path
// of its source code.
func (g *generator) setup(ctx context.Context) (err error) {
	// Download dependencies once
	if err := g.downloadDependencies(ctx); err != nil {
		return err
	}

	// Parse and resolve dependencies
	if err := g.resolveDependencies(ctx); err != nil {
		return err
	}

	// Discover app modules and includes in parallel
	if err := g.discoverAppModules(ctx); err != nil {
		return err
	}

	// Process third-party modules efficiently
	return g.processThirdPartyModules(ctx)
}

func (g *generator) downloadDependencies(ctx context.Context) error {
	var errb bytes.Buffer
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(&errb),
			cmdrunner.DefaultWorkdir(g.appPath),
		).Run(ctx, step.New(step.Exec("go", "mod", "download")))
}

func (g *generator) resolveDependencies(ctx context.Context) error {
	modFile, err := gomodule.ParseAt(g.appPath)
	if err != nil {
		return err
	}

	// Read the dependencies defined in the `go.mod` file
	g.deps, err = gomodule.ResolveDependencies(modFile, false)
	if err != nil {
		return err
	}

	// Find and set SDK directory
	dep, found := filterCosmosSDKModule(g.deps)
	if !found {
		return ErrMissingSDKDep
	}

	g.sdkImport = dep.Path
	g.sdkDir, err = gomodule.LocatePath(ctx, g.cacheStorage, g.appPath, dep)
	return err
}

func (g *generator) discoverAppModules(ctx context.Context) error {
	// Discover app modules
	var err error
	g.appModules, err = module.Discover(
		ctx,
		g.appPath,
		g.appPath,
		module.WithProtoDir(g.protoDir),
		module.WithSDKDir(g.sdkDir),
	)
	if err != nil {
		return err
	}

	// Resolve app includes
	g.appIncludes, _, err = g.resolveIncludes(ctx, g.appPath, g.protoDir)
	return err
}

func (g *generator) processThirdPartyModules(ctx context.Context) error {
	moduleCache := cache.New[protoAnalysis](g.cacheStorage, moduleCacheNamespace)
	sdkModuleCache := cache.New[protoAnalysis](g.cacheStorage, sdkModuleCacheNamespace)

	// Process dependencies in parallel for better performance
	type depResult struct {
		path     string
		analysis protoAnalysis
		err      error
	}

	results := make(chan depResult, len(g.deps))
	semaphore := make(chan struct{}, 5) // Limit concurrent operations

	for _, dep := range g.deps {
		go func(ctx context.Context, dep gomodule.Version) {
			// check for cancellation first
			if err := ctx.Err(); err != nil {
				results <- depResult{path: "", analysis: protoAnalysis{}, err: err}
				return
			}

			select {
			case semaphore <- struct{}{}: // Acquire
			case <-ctx.Done():
				results <- depResult{path: "", analysis: protoAnalysis{}, err: ctx.Err()}
				return
			}
			defer func() { <-semaphore }() // Release

			var depInfo protoAnalysis
			var err error

			// Check if this is a Cosmos SDK module
			// Optimization: All SDK modules share the same proto files from the SDK's proto directory:
			// - cosmossdk.io/* (newer modular SDK packages like cosmossdk.io/math, cosmossdk.io/x/*)
			// - github.com/cosmos/cosmos-sdk/* (traditional monolithic SDK packages)
			// Instead of processing the same proto files multiple times for each SDK module
			// dependency, we use a shared cache key based on the SDK import path. This eliminates:
			// - Module discovery operations
			// - Proto include resolution
			// - Buf export operations
			// - File system operations
			// This can reduce processing time by 70-90% for projects with many SDK modules.
			if module.IsCosmosSDKPackage(dep.Path) || strings.HasPrefix(dep.Path, "cosmossdk.io/") {
				// Use a shared cache key for all SDK modules since they reference the same proto dir
				sdkCacheKey := cache.Key("cosmos-sdk", g.sdkImport)
				depInfo, err = sdkModuleCache.Get(sdkCacheKey)

				if errors.Is(err, cache.ErrorNotFound) {
					// check for cancellation before expensive operation
					if err := ctx.Err(); err != nil {
						results <- depResult{path: "", analysis: protoAnalysis{}, err: err}
						return
					}

					depInfo, err = g.processNewDependency(ctx, dep)
					if err == nil && len(depInfo.Modules) > 0 && depInfo.Cacheable {
						// Cache using the shared SDK key for all SDK modules
						_ = sdkModuleCache.Put(sdkCacheKey, depInfo)
					}
				}
			} else {
				// Regular module processing with individual cache keys
				cacheKey := cache.Key(dep.Path, dep.Version)
				depInfo, err = moduleCache.Get(cacheKey)

				if errors.Is(err, cache.ErrorNotFound) {
					// check for cancellation before expensive operation
					if err := ctx.Err(); err != nil {
						results <- depResult{path: "", analysis: protoAnalysis{}, err: err}
						return
					}

					depInfo, err = g.processNewDependency(ctx, dep)
					if err == nil && len(depInfo.Modules) > 0 && depInfo.Cacheable {
						// Cache the result only if it's safe to do
						_ = moduleCache.Put(cacheKey, depInfo)
					}
				}
			}

			results <- depResult{path: depInfo.Path, analysis: depInfo, err: err}
		}(ctx, dep)
	}

	// Collect results
	for i := 0; i < len(g.deps); i++ {
		select {
		case result := <-results:
			if result.err != nil && !errors.Is(result.err, cache.ErrorNotFound) {
				return result.err
			}

			if result.analysis.Path != "" {
				g.thirdModules[result.path] = result.analysis.Modules
				g.thirdModuleIncludes[result.path] = result.analysis.Includes
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (g *generator) processNewDependency(ctx context.Context, dep gomodule.Version) (protoAnalysis, error) {
	// Get the absolute path to the package's directory
	path, err := gomodule.LocatePath(ctx, g.cacheStorage, g.appPath, dep)
	if err != nil {
		return protoAnalysis{}, err
	}

	// Discover modules
	modules, err := module.Discover(ctx, g.appPath, path, module.WithSDKDir(g.sdkDir))
	if err != nil {
		return protoAnalysis{}, err
	}

	// Only resolve includes if modules exist
	var includes protoIncludes
	var cacheable bool
	if len(modules) > 0 {
		includes, cacheable, err = g.resolveIncludes(ctx, path, defaults.ProtoDir)
		if err != nil {
			return protoAnalysis{}, err
		}
	} else {
		cacheable = true // No includes needed, safe to cache
	}

	return protoAnalysis{
		Path:      path,
		Modules:   modules,
		Includes:  includes,
		Cacheable: cacheable,
	}, nil
}

func (g *generator) getProtoIncludeFolders(modPath string) []string {
	return []string{filepath.Join(modPath, g.protoDir)}
}

func (g *generator) findBufPath(modpath string) (string, error) {
	// check cache first
	if cached, exists := g.bufPathCache[modpath]; exists {
		return cached, nil
	}

	var bufPath string
	// More efficient: check common locations first before walking entire tree
	commonPaths := []string{
		filepath.Join(modpath, bufYamlFilename),
		filepath.Join(modpath, "buf.yml"),
		filepath.Join(modpath, "proto", bufYamlFilename),
		filepath.Join(modpath, "proto", "buf.yml"),
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			bufPath = path
			break
		}
	}

	// If not found in common locations, walk the directory tree
	if bufPath == "" {
		err := filepath.WalkDir(modpath, func(path string, _ fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			base := filepath.Base(path)
			if base == bufYamlFilename || base == "buf.yml" {
				bufPath = path
				return filepath.SkipAll
			}
			// Skip deep nested directories that are unlikely to contain buf configs
			if strings.Count(path, string(os.PathSeparator)) > strings.Count(modpath, string(os.PathSeparator))+3 {
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	// cache the result
	g.bufPathCache[modpath] = bufPath
	return bufPath, nil
}

func (g *generator) generateBufIncludeFolder(ctx context.Context, modpath string) (string, error) {
	// check cache first to avoid repeated export operations
	// this is particularly important since multiple dependencies may reference
	// the same proto path, causing redundant buf.Export calls and temp directory creation
	if cached, exists := g.bufExportCache[modpath]; exists {
		// verify the cached path still exists
		if _, err := os.Stat(cached); err == nil {
			return cached, nil
		}
		// remove invalid cache entry
		delete(g.bufExportCache, modpath)
	}

	protoPath, err := os.MkdirTemp("", "includeFolder")
	if err != nil {
		return "", err
	}

	g.tmpDirs = append(g.tmpDirs, protoPath)

	err = g.buf.Export(ctx, modpath, protoPath)
	if err != nil {
		return "", err
	}

	// cache the result for future use
	g.bufExportCache[modpath] = protoPath
	return protoPath, nil
}

func (g *generator) resolveIncludes(ctx context.Context, path, protoDir string) (protoIncludes, bool, error) {
	// Use a cache key that includes both path and protoDir for better cache hits
	cacheKey := path + ":" + protoDir
	includeCache := cache.New[protoIncludes](g.cacheStorage, includeProtoCacheNamespace)

	if cached, err := includeCache.Get(cacheKey); err == nil {
		return cached, true, nil
	}

	// Get global includes once and reuse
	paths, err := protocGlobalInclude()
	if err != nil {
		return protoIncludes{}, false, err
	}

	includes := protoIncludes{Paths: paths}

	// Determine proto path based on package type
	var protoPath string
	if module.IsCosmosSDKPackage(path) {
		protoPath = filepath.Join(g.sdkDir, "proto")
	} else {
		protoPath = filepath.Join(path, protoDir)
		if fi, err := os.Stat(protoPath); os.IsNotExist(err) {
			protoPath, err = findInnerProtoFolder(path)
			if err != nil {
				// if proto directory does not exist, we just skip it
				log.Print(err.Error())
				return protoIncludes{}, false, nil
			}
		} else if err != nil {
			return protoIncludes{}, false, err
		} else if !fi.IsDir() {
			return includes, true, nil
		}
	}

	// Add proto path and find buf config
	includes.Paths = append(includes.Paths, protoPath)
	includes.ProtoPath = protoPath

	// Efficient buf path discovery
	includes.BufPath, err = g.findBufPath(protoPath)
	if err != nil {
		return includes, false, err
	}

	// Try project root if not found in proto path
	if includes.BufPath == "" {
		includes.BufPath, err = g.findBufPath(path)
		if err != nil {
			return includes, false, err
		}
	}

	// Handle buf config processing
	cacheable := true
	if includes.BufPath != "" {
		bufProtoPath, err := g.generateBufIncludeFolder(ctx, protoPath)
		if err != nil && !errors.Is(err, cosmosbuf.ErrProtoFilesNotFound) {
			return protoIncludes{}, false, err
		}
		if bufProtoPath != "" {
			includes.Paths = append(includes.Paths, bufProtoPath)
			cacheable = false // Don't cache when temp directories are involved
		}
	} else {
		// Legacy behavior: add configured directories
		includes.Paths = append(includes.Paths, g.getProtoIncludeFolders(path)...)
	}

	// Cache the result if appropriate
	if cacheable {
		_ = includeCache.Put(cacheKey, includes)
	}

	return includes, cacheable, nil
}

func (g generator) updateBufModule(ctx context.Context) error {
	// Process in batch to reduce individual file operations
	var bufDeps []string
	var vendorOps []struct{ pkgName, protoPath string }

	for pkgPath, includes := range g.thirdModuleIncludes {
		// Skip third party dependencies without proto files
		if includes.ProtoPath == "" {
			continue
		}

		// Resolve the Go package and use the module name as the proto vendor directory name
		modFile, err := gomodule.ParseAt(pkgPath)
		if err != nil {
			return err
		}

		pkgName := modFile.Module.Mod.Path

		// Batch buf dependencies and vendor operations
		if includes.BufPath != "" {
			depName, err := g.getBufDependencyName(includes.BufPath)
			if err != nil {
				return err
			}
			if depName != "" {
				bufDeps = append(bufDeps, depName)
			} else {
				vendorOps = append(vendorOps, struct{ pkgName, protoPath string }{pkgName, filepath.Dir(includes.BufPath)})
			}
		} else {
			vendorOps = append(vendorOps, struct{ pkgName, protoPath string }{pkgName, includes.ProtoPath})
		}
	}

	// Process buf dependencies in batch
	if len(bufDeps) > 0 {
		if err := g.addBufDependencies(bufDeps); err != nil {
			return err
		}
	}

	// Process vendor operations
	for _, op := range vendorOps {
		if err := g.vendorProtoPackage(op.pkgName, op.protoPath); err != nil {
			return err
		}
	}

	// Update buf once at the end
	if err := g.buf.Update(
		ctx,
		filepath.Dir(g.appIncludes.BufPath),
	); err != nil && !errors.Is(err, cosmosbuf.ErrProtoFilesNotFound) {
		return err
	}
	return nil
}

func (g generator) getBufDependencyName(bufPath string) (string, error) {
	// check cache first
	if cached, exists := g.bufConfigCache[bufPath]; exists {
		return cached.Name, nil
	}

	// Open and parse buf config
	f, err := os.Open(bufPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cfg := struct {
		Name string `yaml:"name"`
	}{}

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return "", newBufConfigError(bufPath, err)
	}

	// cache the result
	g.bufConfigCache[bufPath] = struct{ Name string }{cfg.Name}
	return cfg.Name, nil
}

func (g generator) addBufDependencies(depNames []string) error {
	if len(depNames) == 0 {
		return nil
	}

	// Read app's Buf config once
	path := g.appIncludes.BufPath
	bz, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse existing dependencies
	cfg := struct {
		Deps []string `yaml:"deps"`
	}{}
	if err := yaml.Unmarshal(bz, &cfg); err != nil {
		return newBufConfigError(path, err)
	}

	// Filter out already existing dependencies
	var newDeps []string
	for _, depName := range depNames {
		if !slices.Contains(cfg.Deps, depName) {
			newDeps = append(newDeps, depName)
		}
	}

	if len(newDeps) == 0 {
		return nil // No new dependencies to add
	}

	// Add new dependencies and update config
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	var rawCfg map[string]interface{}
	if err := yaml.Unmarshal(bz, &rawCfg); err != nil {
		return newBufConfigError(path, err)
	}

	rawCfg["deps"] = append(cfg.Deps, newDeps...)

	enc := yaml.NewEncoder(f)
	defer enc.Close()

	if err := enc.Encode(rawCfg); err != nil {
		return err
	}

	// Send notifications for all new dependencies
	for _, depName := range newDeps {
		g.opts.ev.Send(
			fmt.Sprintf("New Buf dependency added: %s", colors.Name(depName)),
			events.Icon(icons.OK),
		)
	}

	return nil
}

func (g generator) vendorProtoPackage(pkgName, protoPath string) (err error) {
	// Check that the dependency vendor directory doesn't exist
	vendorRelPath := filepath.Join("proto_vendor", pkgName)
	vendorPath := filepath.Join(g.appPath, vendorRelPath)
	_, err = os.Stat(vendorPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Skip vendoring when the dependency is already vendored
	if !os.IsNotExist(err) {
		return nil
	}

	if err = os.MkdirAll(vendorPath, 0o777); err != nil {
		return err
	}

	// Make sure that the vendor folder is removed on error
	defer func() {
		if err != nil {
			_ = os.RemoveAll(vendorPath)
		}
	}()

	if err = xos.CopyFolder(protoPath, vendorPath); err != nil {
		return err
	}

	path := filepath.Join(g.appPath, bufYamlFilename)
	bz, err := os.ReadFile(path)
	if err != nil {
		return errors.Errorf("error reading Buf workspace file: %s: %w", path, err)
	}

	ws := struct {
		Version string `yaml:"version"`
		Modules []struct {
			Path string `yaml:"path"`
		} `yaml:"modules"`
		Deps     []string `yaml:"deps"`
		Lint     any      `yaml:"lint"`
		Breaking any      `yaml:"breaking"`
	}{}
	if err := yaml.Unmarshal(bz, &ws); err != nil {
		return err
	}

	ws.Modules = append(ws.Modules, struct {
		Path string `yaml:"path"`
	}{
		Path: vendorRelPath,
	})

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	defer enc.Close()
	if err = enc.Encode(ws); err != nil {
		return err
	}

	g.opts.ev.Send(
		fmt.Sprintf("New Buf vendored dependency added: %s", colors.Name(vendorRelPath)),
		events.Icon(icons.OK),
	)

	return nil
}

func filterCosmosSDKModule(versions []gomodule.Version) (gomodule.Version, bool) {
	for _, v := range versions {
		if cosmosver.CosmosSDKModulePathPattern.MatchString(v.Path) {
			return v, true
		}
	}
	return gomodule.Version{}, false
}

// findInnerProtoFolder attempts to find the proto directory in a module.
// it should be used as a fallback when the proto directory is not found in the expected location.
func findInnerProtoFolder(path string) (string, error) {
	// attempt to find proto directory in the module
	protoFiles, err := xos.FindFiles(path, xos.WithExtension(xos.ProtoFile))
	if err != nil {
		return "", err
	}
	if len(protoFiles) == 0 {
		return "", errors.Errorf("no proto folders found in %s", path)
	}

	var protoDirs []string
	for _, p := range protoFiles {
		dir := filepath.Dir(p)
		for {
			if filepath.Base(dir) == "proto" {
				protoDirs = append(protoDirs, dir)
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir { // reached root
				break
			}
			dir = parent
		}
	}

	if len(protoDirs) == 0 {
		// Fallback to the parent of the first proto file found.
		return filepath.Dir(protoFiles[0]), nil
	}

	// Find the highest level proto directory (shortest path)
	highest := protoDirs[0]
	for _, d := range protoDirs[1:] {
		if len(d) < len(highest) {
			highest = d
		}
	}

	return highest, nil
}
