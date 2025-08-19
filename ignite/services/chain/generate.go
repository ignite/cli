package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/events"
)

type generateOptions struct {
	useCache             bool
	isProtoVendorEnabled bool
	isGoEnabled          bool
	isTSClientEnabled    bool
	isComposablesEnabled bool
	isOpenAPIEnabled     bool
	tsClientPath         string
	composablesPath      string
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
func GenerateTSClient(path string, useCache bool) GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
		o.tsClientPath = path
		o.useCache = useCache
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

// GenerateOpenAPI enables generating OpenAPI spec for your chain.
func GenerateOpenAPI() GenerateTarget {
	return func(o *generateOptions) {
		o.isOpenAPIEnabled = true
	}
}

// GenerateProtoVendor enables `proto_vendor` folder generation.
// Proto vendor is generated from Go dependencies that contain proto files that
// are not included in the app's Buf config.
// Enabling proto vendoring might update Buf config with missing dependencies
// if a Go dependency contains proto files and a Buf config with a name that is
// not listed in the Buf dependencies.
func GenerateProtoVendor() GenerateTarget {
	return func(o *generateOptions) {
		o.isProtoVendorEnabled = true
	}
}

// generateFromConfig makes code generation from proto files from the given config.
func (c *Chain) generateFromConfig(ctx context.Context, cacheStorage cache.Storage, generateClients bool) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	// Additional code generation targets
	var targets []GenerateTarget

	if conf.Client.OpenAPI.Path != "" {
		targets = append(targets, GenerateOpenAPI())
	}

	if generateClients {
		if p := conf.Client.Typescript.Path; p != "" {
			targets = append(targets, GenerateTSClient(p, true))
		}

		if p := conf.Client.Composables.Path; p != "" {
			targets = append(targets, GenerateComposables(p))
		}
	}

	// Generate proto based code for Go and optionally for any optional targets
	return c.Generate(ctx, cacheStorage, GenerateGo(), targets...)
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

	c.ev.Send("Building proto...", events.ProgressUpdate())

	options := []cosmosgen.Option{cosmosgen.CollectEvents(c.ev)}

	if targetOptions.isGoEnabled {
		options = append(options, cosmosgen.WithGoGeneration())
	}

	if targetOptions.isProtoVendorEnabled {
		options = append(options, cosmosgen.UpdateBufModule())
	}

	var (
		openAPIPath, tsClientPath, composablesPath string
		updateConfig                               bool
	)

	if targetOptions.isOpenAPIEnabled {
		openAPIPath = conf.Client.OpenAPI.Path
		if openAPIPath == "" {
			openAPIPath = chainconfig.DefaultOpenAPIPath
		}

		// Non-absolute OpenAPI paths must be treated as relative to the app directory
		if !filepath.IsAbs(openAPIPath) {
			openAPIPath = filepath.Join(c.app.Path, openAPIPath)
		}

		options = append(options, cosmosgen.WithOpenAPIGeneration(openAPIPath))
	}

	if targetOptions.isTSClientEnabled {
		tsClientPath = targetOptions.tsClientPath
		if tsClientPath == "" {
			tsClientPath = chainconfig.TSClientPath(*conf)

			// When TS client is generated make sure the config is updated
			// with the output path when the client path option is empty.
			if conf.Client.Typescript.Path == "" {
				conf.Client.Typescript.Path = tsClientPath
				updateConfig = true
			}
		}

		// Non-absolute TS client output paths must be treated as relative to the app directory
		if !filepath.IsAbs(tsClientPath) {
			tsClientPath = filepath.Join(c.app.Path, tsClientPath)
		}

		options = append(options,
			cosmosgen.WithTSClientGeneration(
				cosmosgen.TypescriptModulePath(tsClientPath),
				tsClientPath,
				targetOptions.useCache,
			),
		)
	}

	if targetOptions.isComposablesEnabled {
		composablesPath = targetOptions.composablesPath

		if composablesPath == "" {
			composablesPath = chainconfig.ComposablesPath(conf)

			if conf.Client.Composables.Path == "" {
				conf.Client.Composables.Path = composablesPath
				updateConfig = true
			}
		}

		// Non-absolute Composables output paths must be treated as relative to the app directory
		if !filepath.IsAbs(composablesPath) {
			composablesPath = filepath.Join(c.app.Path, composablesPath)
		}

		options = append(options,
			cosmosgen.WithComposablesGeneration(
				cosmosgen.ComposableModulePath(composablesPath),
				composablesPath,
			),
		)
	}

	if err := cosmosgen.Generate(
		ctx,
		cacheStorage,
		c.app.Path,
		conf.Build.Proto.Path,
		c.app.ImportPath,
		chainconfig.DefaultVuePath,
		options...,
	); err != nil {
		return &CannotBuildAppError{err}
	}

	// Check if the client config options have to be updated with the paths of the generated code
	if updateConfig {
		if err := c.saveClientConfig(conf.Client); err != nil {
			return errors.Errorf("error adding generated paths to config file: %w", err)
		}
	}

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

func (c Chain) saveClientConfig(client base.Client) error {
	path := c.ConfigPath()
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	// Initialize the config to the file values ignoring empty
	// values that otherwise would be initialized to defaults.
	// Defaults must be ignored to avoid writing them to the
	// YAML config file when they are not present.
	var cfg chainconfig.Config
	if err := cfg.Decode(file); err != nil {
		return err
	}

	cfg.Client = client

	return chainconfig.Save(cfg, path)
}
