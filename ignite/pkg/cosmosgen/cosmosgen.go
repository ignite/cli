package cosmosgen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	gomodule "golang.org/x/mod/module"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/events"
)

// generateOptions used to configure code generation.
type generateOptions struct {
	includeDirs     []string
	useCache        bool
	updateBufModule bool
	ev              events.Bus

	generateProtobuf bool

	jsOut            func(module.Module) string
	tsClientRootPath string

	composablesOut      func(module.Module) string
	composablesRootPath string

	hooksOut      func(module.Module) string
	hooksRootPath string

	specOut string
}

// ModulePathFunc defines a function type that returns a path based on a Cosmos SDK module.
type ModulePathFunc func(module.Module) string

// Option configures code generation.
type Option func(*generateOptions)

// WithTSClientGeneration adds Typescript Client code generation.
// The tsClientRootPath is used to determine the root path of generated Typescript classes.
func WithTSClientGeneration(out ModulePathFunc, tsClientRootPath string, useCache bool) Option {
	return func(o *generateOptions) {
		o.jsOut = out
		o.tsClientRootPath = tsClientRootPath
		o.useCache = useCache
	}
}

func WithComposablesGeneration(out ModulePathFunc, composablesRootPath string) Option {
	return func(o *generateOptions) {
		o.composablesOut = out
		o.composablesRootPath = composablesRootPath
	}
}

func WithHooksGeneration(out ModulePathFunc, hooksRootPath string) Option {
	return func(o *generateOptions) {
		o.hooksOut = out
		o.hooksRootPath = hooksRootPath
	}
}

// WithGoGeneration adds protobuf (gogoproto and pulsar) code generation.
func WithGoGeneration() Option {
	return func(o *generateOptions) {
		o.generateProtobuf = true
	}
}

// WithOpenAPIGeneration adds OpenAPI spec generation.
func WithOpenAPIGeneration(out string) Option {
	return func(o *generateOptions) {
		o.specOut = out
	}
}

// IncludeDirs configures the third party proto dirs that used by app's proto.
// relative to the projectPath.
func IncludeDirs(dirs []string) Option {
	return func(o *generateOptions) {
		o.includeDirs = dirs
	}
}

// UpdateBufModule enables Buf config proto dependencies update.
// This option updates app's Buf config when proto packages or
// Buf modules are found within the Go dependencies.
func UpdateBufModule() Option {
	return func(o *generateOptions) {
		o.updateBufModule = true
	}
}

// CollectEvents sets an event bus for sending generation feedback events.
func CollectEvents(ev events.Bus) Option {
	return func(c *generateOptions) {
		c.ev = ev
	}
}

// generator generates code for sdk and sdk apps.
type generator struct {
	buf                 cosmosbuf.Buf
	cacheStorage        cache.Storage
	appPath             string
	protoDir            string
	gomodPath           string
	opts                *generateOptions
	sdkImport           string
	sdkDir              string
	deps                []gomodule.Version
	appModules          []module.Module
	appIncludes         protoIncludes
	thirdModules        map[string][]module.Module
	thirdModuleIncludes map[string]protoIncludes
	tmpDirs             []string
}

func (g *generator) cleanup() {
	// Remove temporary directories created during generation
	for _, path := range g.tmpDirs {
		_ = os.RemoveAll(path)
	}
}

// Generate generates code from protoDir of an SDK app residing at appPath with given options.
// protoDir must be relative to the projectPath.
func Generate(ctx context.Context, cacheStorage cache.Storage, appPath, protoDir, gomodPath string, options ...Option) error {
	b, err := cosmosbuf.New()
	if err != nil {
		return err
	}

	defer func() {
		if err := b.Cleanup(); err != nil {
			fmt.Println("Cleanup error:", err)
		}
	}()

	g := &generator{
		buf:                 b,
		appPath:             appPath,
		protoDir:            protoDir,
		gomodPath:           gomodPath,
		opts:                &generateOptions{},
		thirdModules:        make(map[string][]module.Module),
		thirdModuleIncludes: make(map[string]protoIncludes),
		cacheStorage:        cacheStorage,
	}

	defer g.cleanup()

	for _, apply := range options {
		apply(g.opts)
	}

	if err := g.setup(ctx); err != nil {
		return err
	}

	// Update app's Buf config for third party discovered proto modules.
	// Go dependency packages might contain proto files which could also
	// optionally be using Buf, so for those cases the discovered proto
	// files should be available before code generation.
	if g.opts.updateBufModule {
		if err := g.updateBufModule(ctx); err != nil {
			return err
		}
	}

	// Go generation must run first so the types are created before other
	// generated code that requires sdk.Msg implementations to be defined
	if g.opts.generateProtobuf {
		if err := g.generateGoGo(ctx); err != nil {
			return err
		}

		if err := g.generatePulsar(ctx); err != nil {
			return err
		}
	}

	if g.opts.specOut != "" {
		if err := g.generateOpenAPISpec(ctx); err != nil {
			return err
		}
	}

	if g.opts.jsOut != nil {
		if err := g.generateTS(ctx); err != nil {
			return err
		}
	}

	if g.opts.composablesRootPath != "" {
		if err := g.generateComposables("vue"); err != nil {
			return err
		}

		// Update Vue app dependencies when Vue composables are generated.
		// This update is required to link the "ts-client" folder so the
		// package is available during development before publishing it.
		if err := g.updateComposableDependencies("vue"); err != nil {
			return err
		}
	}
	if g.opts.hooksRootPath != "" {
		if err := g.generateComposables("react"); err != nil {
			return err
		}

		// Update React app dependencies when React hooks are generated.
		// This update is required to link the "ts-client" folder so the
		// package is available during development before publishing it.
		if err := g.updateComposableDependencies("react"); err != nil {
			return err
		}
	}

	return nil
}

// TypescriptModulePath generates TS module paths for Cosmos SDK modules.
// The root path is used as prefix for the generated paths.
func TypescriptModulePath(rootPath string) ModulePathFunc {
	return func(m module.Module) string {
		return filepath.Join(rootPath, m.Pkg.Name)
	}
}

// ComposableModulePath generates useQuery hook/composable module paths for Cosmos SDK modules.
// The root path is used as prefix for the generated paths.
func ComposableModulePath(rootPath string) ModulePathFunc {
	return func(m module.Module) string {
		replacer := strings.NewReplacer("-", "_", ".", "_")
		modPath := strcase.ToCamel(replacer.Replace(m.Pkg.Name))
		return filepath.Join(rootPath, "use"+modPath)
	}
}
