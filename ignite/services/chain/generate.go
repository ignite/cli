package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
)

const (
	defaultVuexPath    = "vue/src/store"
	defaultDartPath    = "flutter/lib"
	defaultOpenAPIPath = "docs/static/openapi.yml"
)

type generateOptions struct {
	isGoEnabled       bool
	isTSClientEnabled bool
	isVuexEnabled     bool
	isDartEnabled     bool
	isOpenAPIEnabled  bool
	tsClientPath      string
}

// GenerateTarget is a target to generate code for from proto files.
type GenerateTarget func(*generateOptions)

// GenerateGo enables generating proto based Go code needed for the chain's source code.
func GenerateGo() GenerateTarget {
	return func(o *generateOptions) {
		o.isGoEnabled = true
	}
}

// GenerateTSClient enables generating proto based Typescript Client.
// The path assigns the output path to use for the generated Typescript client
// overriding the configured or default path. Path can be an empty string.
func GenerateTSClient(path string) GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.tsClientPath = path
	}
}

// GenerateTSClient enables generating proto based Typescript Client.
func GenerateVuex() GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.isVuexEnabled = true
	}
}

// GenerateDart enables generating Dart client.
func GenerateDart() GenerateTarget {
	return func(o *generateOptions) {
		o.isDartEnabled = true
	}
}

// GenerateOpenAPI enables generating OpenAPI spec for your chain.
func GenerateOpenAPI() GenerateTarget {
	return func(o *generateOptions) {
		o.isOpenAPIEnabled = true
	}
}

// generateFromConfig makes code generation from proto files from the given config
func (c *Chain) generateFromConfig(ctx context.Context, cacheStorage cache.Storage) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	var additionalTargets []GenerateTarget

	// parse config for additional target
	if p := conf.Client.Typescript.Path; p != "" {
		additionalTargets = append(additionalTargets, GenerateTSClient(p))
	}

	if conf.Client.Vuex.Path != "" {
		additionalTargets = append(additionalTargets, GenerateVuex())
	}

	if conf.Client.Dart.Path != "" {
		additionalTargets = append(additionalTargets, GenerateDart())
	}

	if conf.Client.OpenAPI.Path != "" {
		additionalTargets = append(additionalTargets, GenerateOpenAPI())
	}

	return c.Generate(ctx, cacheStorage, GenerateGo(), additionalTargets...)
}

// Generate makes code generation from proto files for given target and additionalTargets.
func (c *Chain) Generate(
	ctx context.Context,
	cacheStorage cache.Storage,
	target GenerateTarget,
	additionalTargets ...GenerateTarget,
) error {
	var targetOptions generateOptions

	for _, apply := range append(additionalTargets, target) {
		apply(&targetOptions)
	}

	conf, err := c.Config()
	if err != nil {
		return err
	}

	if err := cosmosgen.InstallDepTools(ctx, c.app.Path); err != nil {
		return err
	}

	c.ev.Send("🛠  Building proto...")

	options := []cosmosgen.Option{
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	if targetOptions.isGoEnabled {
		options = append(options, cosmosgen.WithGoGeneration(c.app.ImportPath))
	}

	enableThirdPartyModuleCodegen := !c.protoBuiltAtLeastOnce && c.options.isThirdPartyModuleCodegenEnabled

	var dartPath, openAPIPath, tsClientPath, vuexPath string

	if targetOptions.isTSClientEnabled {
		tsClientPath = targetOptions.tsClientPath
		if tsClientPath == "" {
			tsClientPath = chainconfig.TSClientPath(conf)
		}

		// Non absolute TS client output paths must be treated as relative to the app directory
		if !filepath.IsAbs(tsClientPath) {
			tsClientPath = filepath.Join(c.app.Path, tsClientPath)
		}

		if err := os.MkdirAll(tsClientPath, 0o766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithTSClientGeneration(
				cosmosgen.TypescriptModulePath(tsClientPath),
				tsClientPath,
			),
		)
	}

	if targetOptions.isVuexEnabled {
		vuexPath = conf.Client.Vuex.Path
		if vuexPath == "" {
			vuexPath = defaultVuexPath
		}

		if filepath.IsAbs(vuexPath) {
			// TODO: Should we always generate Vuex code inside a "generated" directory?
			vuexPath = filepath.Join(vuexPath, "generated")
		} else {
			vuexPath = filepath.Join(c.app.Path, vuexPath, "generated")
		}

		if err := os.MkdirAll(vuexPath, 0o766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithVuexGeneration(
				enableThirdPartyModuleCodegen,
				cosmosgen.TypescriptModulePath(vuexPath),
				vuexPath,
			),
		)
	}

	if targetOptions.isDartEnabled {
		dartPath = conf.Client.Dart.Path
		if dartPath == "" {
			dartPath = defaultDartPath
		}

		if filepath.IsAbs(dartPath) {
			// TODO: Should we always generate Dart code inside a "generated" directory?
			dartPath = filepath.Join(dartPath, "generated")
		} else {
			dartPath = filepath.Join(c.app.Path, dartPath, "generated")
		}

		if err := os.MkdirAll(dartPath, 0o766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithDartGeneration(
				enableThirdPartyModuleCodegen,
				func(m module.Module) string {
					return filepath.Join(dartPath, m.Pkg.Name, "module")
				},
				dartPath,
			),
		)
	}

	if targetOptions.isOpenAPIEnabled {
		openAPIPath = conf.Client.OpenAPI.Path
		if openAPIPath == "" {
			openAPIPath = defaultOpenAPIPath
		}

		// Non absolute OpenAPI paths must be treated as relative to the app directory
		if !filepath.IsAbs(openAPIPath) {
			openAPIPath = filepath.Join(c.app.Path, openAPIPath)
		}

		options = append(options, cosmosgen.WithOpenAPIGeneration(openAPIPath))
	}

	if err := cosmosgen.Generate(ctx, cacheStorage, c.app.Path, conf.Build.Proto.Path, options...); err != nil {
		return &CannotBuildAppError{err}
	}

	c.protoBuiltAtLeastOnce = true

	return nil
}
