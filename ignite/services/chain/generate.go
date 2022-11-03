package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/ignite/pkg/events"
)

type generateOptions struct {
	isGoEnabled          bool
	isTSClientEnabled    bool
	isComposablesEnabled bool
	isHooksEnabled       bool
	isVuexEnabled        bool
	isOpenAPIEnabled     bool
	tsClientPath         string
	vuexPath             string
	composablesPath      string
	hooksPath            string
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

// GenerateVuex enables generating proto based Typescript Client and Vuex Stores.
func GenerateVuex(path string) GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.isVuexEnabled = true
		o.vuexPath = path
	}
}

// GenerateComposables enables generating proto based Typescript Client and Vue 3 composables.
func GenerateComposables(path string) GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.isComposablesEnabled = true
		o.composablesPath = path
	}
}

// GenerateHooks enables generating proto based Typescript Client and React composables.
func GenerateHooks(path string) GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.isHooksEnabled = true
		o.hooksPath = path
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

	//nolint:staticcheck //ignore SA1019 until vuex config option is removed
	if p := conf.Client.Vuex.Path; p != "" {
		additionalTargets = append(additionalTargets, GenerateVuex(p))
	}

	if p := conf.Client.Composables.Path; p != "" {
		additionalTargets = append(additionalTargets, GenerateComposables(p))
	}

	if p := conf.Client.Hooks.Path; p != "" {
		additionalTargets = append(additionalTargets, GenerateHooks(p))
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

	c.ev.Send("Building proto...", events.ProgressUpdate())

	options := []cosmosgen.Option{
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	if targetOptions.isGoEnabled {
		options = append(options, cosmosgen.WithGoGeneration(c.app.ImportPath))
	}

	enableThirdPartyModuleCodegen := !c.protoBuiltAtLeastOnce && c.options.isThirdPartyModuleCodegenEnabled

	var openAPIPath, tsClientPath, vuexPath, composablesPath, hooksPath string

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
		//nolint:staticcheck //ignore SA1019 until vuex config option is removed
		vuexPath = targetOptions.vuexPath
		if vuexPath == "" {
			vuexPath = chainconfig.VuexPath(conf)
		}

		// Non absolute Vuex output paths must be treated as relative to the app directory
		if !filepath.IsAbs(vuexPath) {
			vuexPath = filepath.Join(c.app.Path, vuexPath)
		}

		vuexPath = c.joinGeneratedPath(vuexPath)
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

	if targetOptions.isComposablesEnabled {
		composablesPath = targetOptions.composablesPath

		if composablesPath == "" {
			composablesPath = chainconfig.ComposablesPath(conf)
		}

		// Non absolute Composables output paths must be treated as relative to the app directory
		if !filepath.IsAbs(composablesPath) {
			composablesPath = filepath.Join(c.app.Path, composablesPath)
		}

		if err := os.MkdirAll(composablesPath, 0o766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithComposablesGeneration(
				enableThirdPartyModuleCodegen,
				cosmosgen.ComposableModulePath(composablesPath),
				composablesPath,
			),
		)
	}

	if targetOptions.isHooksEnabled {
		hooksPath = targetOptions.hooksPath
		if hooksPath == "" {
			hooksPath = chainconfig.HooksPath(conf)
		}

		// Non absolute Hooks output paths must be treated as relative to the app directory
		if !filepath.IsAbs(hooksPath) {
			hooksPath = filepath.Join(c.app.Path, hooksPath)
		}

		if err := os.MkdirAll(hooksPath, 0o766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithHooksGeneration(
				enableThirdPartyModuleCodegen,
				cosmosgen.ComposableModulePath(hooksPath),
				hooksPath,
			),
		)
	}

	if targetOptions.isOpenAPIEnabled {
		openAPIPath = conf.Client.OpenAPI.Path
		if openAPIPath == "" {
			openAPIPath = chainconfig.DefaultOpenAPIPath
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

	if c.options.printGeneratedPaths {
		if targetOptions.isTSClientEnabled {
			c.ev.Send(
				fmt.Sprintf("Typescript client path: %s", tsClientPath),
				events.Icon(icons.Bullet),
				events.ProgressFinish(),
			)
		}

		if targetOptions.isComposablesEnabled {
			c.ev.Send(
				fmt.Sprintf("Vue composables path: %s", composablesPath),
				events.Icon(icons.Bullet),
				events.ProgressFinish(),
			)
		}

		if targetOptions.isHooksEnabled {
			c.ev.Send(
				fmt.Sprintf("React hooks path: %s", hooksPath),
				events.Icon(icons.Bullet),
				events.ProgressFinish(),
			)
		}

		if targetOptions.isVuexEnabled {
			c.ev.Send(
				fmt.Sprintf("Vuex stores path: %s", vuexPath),
				events.Icon(icons.Bullet),
				events.ProgressFinish(),
			)
		}

		if targetOptions.isOpenAPIEnabled {
			c.ev.Send(
				fmt.Sprintf("OpenAPI path: %s", openAPIPath),
				events.Icon(icons.Bullet),
				events.ProgressFinish(),
			)
		}
	}

	return nil
}

func (c Chain) joinGeneratedPath(rootPath string) string {
	if filepath.IsAbs(rootPath) {
		return filepath.Join(rootPath, "generated")
	}

	return filepath.Join(c.app.Path, rootPath, "generated")
}
