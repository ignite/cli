package cosmosgen

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	gomodule "golang.org/x/mod/module"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/cosmosbuf"
)

// generateOptions used to configure code generation.
type generateOptions struct {
	includeDirs []string
	useCache    bool

	isGoEnabled     bool
	isPulsarEnabled bool

	jsOut            func(module.Module) string
	tsClientRootPath string

	vuexOut      func(module.Module) string
	vuexRootPath string

	composablesOut      func(module.Module) string
	composablesRootPath string

	hooksOut      func(module.Module) string
	hooksRootPath string

	specOut string
}

// TODO add WithInstall.

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

func WithVuexGeneration(out ModulePathFunc, vuexRootPath string) Option {
	return func(o *generateOptions) {
		o.vuexOut = out
		o.vuexRootPath = vuexRootPath
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

// WithGoGeneration adds Go code generation.
func WithGoGeneration() Option {
	return func(o *generateOptions) {
		o.isGoEnabled = true
	}
}

// WithPulsarGeneration adds Go pulsar code generation.
func WithPulsarGeneration() Option {
	return func(o *generateOptions) {
		o.isPulsarEnabled = true
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

// generator generates code for sdk and sdk apps.
type generator struct {
	ctx          context.Context
	buf          cosmosbuf.Buf
	cacheStorage cache.Storage
	appPath      string
	protoDir     string
	gomodPath    string
	o            *generateOptions
	sdkImport    string
	deps         []gomodule.Version
	appModules   []module.Module
	thirdModules map[string][]module.Module // app dependency-modules pair.
}

// Generate generates code from protoDir of an SDK app residing at appPath with given options.
// protoDir must be relative to the projectPath.
func Generate(ctx context.Context, cacheStorage cache.Storage, appPath, protoDir, gomodPath string, options ...Option) error {
	b, err := cosmosbuf.New()
	if err != nil {
		return err
	}

	g := &generator{
		ctx:          ctx,
		buf:          b,
		appPath:      appPath,
		protoDir:     protoDir,
		gomodPath:    gomodPath,
		o:            &generateOptions{},
		thirdModules: make(map[string][]module.Module),
		cacheStorage: cacheStorage,
	}

	for _, apply := range options {
		apply(g.o)
	}

	if err := g.setup(); err != nil {
		return err
	}

	// Go generation must run first so the types are created before other
	// generated code that requires sdk.Msg implementations to be defined
	if g.o.isGoEnabled {
		if err := g.generateGo(); err != nil {
			return err
		}
	}
	if g.o.isPulsarEnabled {
		if err := g.generatePulsar(); err != nil {
			return err
		}
	}

	if g.o.jsOut != nil {
		if err := g.generateTS(); err != nil {
			return err
		}
	}

	if g.o.vuexOut != nil {
		if err := g.generateVuex(); err != nil {
			return err
		}

		// Update Vuex store dependencies when Vuex stores are generated.
		// This update is required to link the "ts-client" folder so the
		// package is available during development before publishing it.
		if err := g.updateVuexDependencies(); err != nil {
			return err
		}

		// Update Vue app dependencies when Vuex stores are generated.
		// This update is required to link the "ts-client" folder so the
		// package is available during development before publishing it.
		if err := g.updateVueDependencies(); err != nil {
			return err
		}

	}

	if g.o.composablesRootPath != "" {
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
	if g.o.hooksRootPath != "" {
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

	if g.o.specOut != "" {
		if err := g.generateOpenAPISpec(); err != nil {
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
