package app

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var files embed.FS

// NewGenerator returns the generator to scaffold a new Cosmos SDK app.
func NewGenerator(opts *Options) (*genny.Generator, error) {
	// Remove "files/" prefix
	subfs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, fmt.Errorf("generator sub: %w", err)
	}
	g := genny.New()
	if err := g.OnlyFS(subfs, opts.IncludePrefixes, nil); err != nil {
		return g, fmt.Errorf("generator fs: %w", err)
	}
	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("GitHubPath", opts.GitHubPath)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
	ctx.Set("DepTools", cosmosgen.DepTools())

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))

	return g, nil
}
