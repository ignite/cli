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

//go:embed files/* files/**/*
var files embed.FS

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
	)
	if !opts.IsConsumerChain {
		// non-consumer chain doesn't need "consumer_*" & "ante_handler.go" files
		excludePrefix = append(excludePrefix, "app/consumer_")
		excludePrefix = append(excludePrefix, "app/ante_handler.go")
	}
	g := genny.New()
	if err := g.SelectiveFS(subfs, includePrefix, nil, excludePrefix, nil); err != nil {
		return g, errors.Errorf("generator fs: %w", err)
	}
	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("GitHubPath", opts.GitHubPath)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
	ctx.Set("IsConsumerChain", opts.IsConsumerChain)
	ctx.Set("DepTools", cosmosgen.DepTools())

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))

	return g, nil
}
