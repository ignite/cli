package app

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var files embed.FS

var (
	ibcConfig        = "app/ibc.go"
	minimalAppConfig = "app/minimal_app_config.go"
	appConfig        = "app/app_config.go"
)

// NewGenerator returns the generator to scaffold a new Cosmos SDK app.
func NewGenerator(opts *Options) (*genny.Generator, error) {
	// Remove "files/" prefix
	subfs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, fmt.Errorf("generator sub: %w", err)
	}
	g := genny.New()

	// always exclude minimal app config it will be created later
	// minimal_app_config is only used for the minimal app template
	excludePrefix := []string{minimalAppConfig}
	if opts.IsChainMinimal {
		// minimal chain does not have ibc or classic app config
		excludePrefix = append(excludePrefix, ibcConfig, appConfig)
	}

	if err := g.SelectiveFS(subfs, opts.IncludePrefixes, nil, excludePrefix, nil); err != nil {
		return g, fmt.Errorf("generator fs: %w", err)
	}

	if opts.IsChainMinimal {
		file, err := subfs.Open(fmt.Sprintf("%s.plush", minimalAppConfig))
		if err != nil {
			return g, fmt.Errorf("open minimal app config: %w", err)
		}

		g.File(genny.NewFile(appConfig, file))
	}

	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("GitHubPath", opts.GitHubPath)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
	ctx.Set("DepTools", cosmosgen.DepTools())
	ctx.Set("IsChainMinimal", opts.IsChainMinimal)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))

	return g, nil
}
