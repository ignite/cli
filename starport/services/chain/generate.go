package chain

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
	"github.com/tendermint/starport/starport/pkg/giturl"
)

type generateOptions struct {
	isProtoBuildGoEnabled      bool
	isProtoBuildVuexEnabled    bool
	isProtoBuildOpenAPIEnabled bool
}

// GenerateTarget is a target to generate code for from proto files.
type GenerateTarget func(*generateOptions)

// GenerateGo enables generating proto based Go code needed for the chain's source code.
func GenerateGo() GenerateTarget {
	return func(o *generateOptions) {
		o.isProtoBuildGoEnabled = true
	}
}

// GenerateVuex enables generating proto based Vuex store.
func GenerateVuex() GenerateTarget {
	return func(o *generateOptions) {
		o.isProtoBuildVuexEnabled = true
	}
}

// GenerateOpenAPI enables generating OpenAPI spec for your chain.
func GenerateOpenAPI() GenerateTarget {
	return func(o *generateOptions) {
		o.isProtoBuildOpenAPIEnabled = true
	}
}

func (c *Chain) generateAll(ctx context.Context) error {
	return c.Generate(ctx, GenerateGo(), GenerateVuex(), GenerateOpenAPI())
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

	if targetOptions.isProtoBuildGoEnabled {
		options = append(options, cosmosgen.WithGoGeneration(c.app.ImportPath))
	}

	enableThirdPartyModuleCodegen := !c.protoBuiltAtLeastOnce && c.options.isThirdPartyModuleCodegenEnabled

	// generate Vuex code as well if it is enabled.
	if targetOptions.isProtoBuildVuexEnabled && conf.Client.Vuex.Path != "" {
		storeRootPath := filepath.Join(c.app.Path, conf.Client.Vuex.Path, "generated")
		options = append(options,
			cosmosgen.WithVuexGeneration(
				enableThirdPartyModuleCodegen,
				func(m module.Module) string {
					parsedGitURL, _ := giturl.Parse(m.Pkg.GoImportName)
					return filepath.Join(storeRootPath, parsedGitURL.UserAndRepo(), m.Pkg.Name, "module")
				},
				storeRootPath,
			),
		)
	}
	if targetOptions.isProtoBuildOpenAPIEnabled && conf.Client.OpenAPI.Path != "" {
		options = append(options, cosmosgen.WithOpenAPIGeneration(conf.Client.OpenAPI.Path))
	}

	if err := cosmosgen.Generate(ctx, c.app.Path, conf.Build.Proto.Path, options...); err != nil {
		return &CannotBuildAppError{err}
	}

	c.protoBuiltAtLeastOnce = true

	return nil
}
