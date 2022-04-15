package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosgen"
	"github.com/ignite-hq/cli/ignite/pkg/giturl"
)

const (
	defaultTSClientPath = "vue/src/generated"
	defaultDartPath     = "flutter/lib"
	defaultOpenAPIPath  = "docs/static/openapi.yml"
)

type generateOptions struct {
	isGoEnabled       bool
	isTSClientEnabled bool
	isDartEnabled     bool
	isOpenAPIEnabled  bool
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
func GenerateTSClient() GenerateTarget {
	return func(o *generateOptions) {
		o.isTSClientEnabled = true
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

func (c *Chain) generateAll(ctx context.Context) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	var additionalTargets []GenerateTarget

	if conf.Client.Typescript.Path != "" {
		additionalTargets = append(additionalTargets, GenerateTSClient())
	}

	if conf.Client.Dart.Path != "" {
		additionalTargets = append(additionalTargets, GenerateDart())
	}

	if conf.Client.OpenAPI.Path != "" {
		additionalTargets = append(additionalTargets, GenerateOpenAPI())
	}

	return c.Generate(ctx, GenerateGo(), additionalTargets...)
}

// Generate makes code generation from proto files for given target and additionalTargets.
func (c *Chain) Generate(
	ctx context.Context,
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

	if err := cosmosgen.InstallDependencies(ctx, c.app.Path); err != nil {
		return err
	}

	fmt.Fprintln(c.stdLog().out, "🛠️  Building proto...")

	options := []cosmosgen.Option{
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	if targetOptions.isGoEnabled {
		options = append(options, cosmosgen.WithGoGeneration(c.app.ImportPath))
	}

	enableThirdPartyModuleCodegen := !c.protoBuiltAtLeastOnce && c.options.isThirdPartyModuleCodegenEnabled

	// generate Typescript Client code as well if it is enabled.
	if targetOptions.isTSClientEnabled {
		tsClientPath := conf.Client.Typescript.Path
		if tsClientPath == "" {
			tsClientPath = defaultTSClientPath
		}

		tsClientRootPath := filepath.Join(c.app.Path, tsClientPath)
		if err := os.MkdirAll(tsClientRootPath, 0766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithTSClientGeneration(
				enableThirdPartyModuleCodegen,
				func(m module.Module) string {
					parsedGitURL, _ := giturl.Parse(m.Pkg.GoImportName)
					return filepath.Join(tsClientRootPath, "client", parsedGitURL.UserAndRepo(), m.Pkg.Name)
				},
				tsClientRootPath,
			),
		)
	}

	if targetOptions.isDartEnabled {
		dartPath := conf.Client.Dart.Path

		if dartPath == "" {
			dartPath = defaultDartPath
		}

		rootPath := filepath.Join(c.app.Path, dartPath, "generated")
		if err := os.MkdirAll(rootPath, 0766); err != nil {
			return err
		}

		options = append(options,
			cosmosgen.WithDartGeneration(
				enableThirdPartyModuleCodegen,
				func(m module.Module) string {
					return filepath.Join(rootPath, m.Pkg.Name, "module")
				},
				rootPath,
			),
		)
	}

	if targetOptions.isOpenAPIEnabled {
		openAPIPath := conf.Client.OpenAPI.Path

		if openAPIPath == "" {
			openAPIPath = defaultOpenAPIPath
		}

		options = append(options, cosmosgen.WithOpenAPIGeneration(openAPIPath))
	}

	if err := cosmosgen.Generate(ctx, c.app.Path, conf.Build.Proto.Path, options...); err != nil {
		return &CannotBuildAppError{err}
	}

	c.protoBuiltAtLeastOnce = true

	return nil
}
