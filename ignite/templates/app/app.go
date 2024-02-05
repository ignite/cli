package app

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

var (
	//go:embed files/* files/**/*
	files embed.FS

	//go:embed files-minimal/* files-minimal/**/*
	filesMinimal embed.FS
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
	g := genny.New()

	var excludePrefix []string
	if opts.IsChainMinimal {
		// minimal chain does not have ibc
		excludePrefix = append(excludePrefix, ibcConfig)
	}

	if err := g.SelectiveFS(subfs, opts.IncludePrefixes, nil, excludePrefix, nil); err != nil {
		return g, errors.Errorf("generator fs: %w", err)
	}

	if opts.IsChainMinimal {
		// Remove "files-minimal/" prefix
		subfs, err := fs.Sub(filesMinimal, "files-minimal")
		if err != nil {
			return nil, errors.Errorf("generator sub minimal: %w", err)
		}
		// Override files from "files" with the ones from "files-minimal"
		if err := g.FS(subfs); err != nil {
			return g, errors.Errorf("generator fs minimal: %w", err)
		}
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
