package app

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field/plushhelpers"
)

var (
	//go:embed files/* files/**/*
	files embed.FS

	//go:embed files-minimal/* files-minimal/**/*
	filesMinimal embed.FS

	//go:embed files-consumer/* files-consumer/**/*
	filesConsumer embed.FS
)

const (
	ibcConfig = "app/ibc.go"
)

// NewGenerator returns the generator to scaffold a new Cosmos SDK app.
func NewGenerator(opts *Options) (*genny.Generator, error) {
	// Remove "files/" prefix
	subfs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, errors.Errorf("generator sub: %w", err)
	}

	var (
		includePrefix = opts.IncludePrefixes
		excludePrefix []string
		overridesFS   = make(map[string]embed.FS)
	)

	if opts.IsChainMinimal {
		// minimal chain does not have ibc
		excludePrefix = append(excludePrefix, ibcConfig)
		overridesFS["files-minimal"] = filesMinimal
	}
	if opts.IsConsumerChain {
		overridesFS["files-consumer"] = filesConsumer
	}

	g := genny.New()
	if err := g.SelectiveFS(subfs, includePrefix, nil, excludePrefix, nil); err != nil {
		return g, errors.Errorf("generator fs: %w", err)
	}

	for prefix, embed := range overridesFS {
		// Remove  prefix
		subfs, err := fs.Sub(embed, prefix)
		if err != nil {
			return g, errors.Errorf("generator sub %s: %w", prefix, err)
		}
		// Override files from "files" with the ones from embed
		if err := g.FS(subfs); err != nil {
			return g, errors.Errorf("generator fs %s: %w", prefix, err)
		}
	}

	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("GitHubPath", opts.GitHubPath)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
<<<<<<< HEAD
	ctx.Set("IsConsumerChain", opts.IsConsumerChain)
=======
	ctx.Set("CoinType", opts.CoinType)
>>>>>>> 2b45eaa2 (feat: wire custom coin type and get bech32 prefix (#4569))
	ctx.Set("DepTools", cosmosgen.DepTools())
	ctx.Set("IsChainMinimal", opts.IsChainMinimal)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))

	return g, nil
}
