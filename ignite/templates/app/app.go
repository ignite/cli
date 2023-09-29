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
	var (
		includePrefix = opts.IncludePrefixes
		includeSuffix []string
		excludePrefix []string
		excludeSuffix []string
	)
	if !opts.IsConsumerChain {
		// non consumer chain doesn't need "ccv_msg_filter_*" & "ante_handler.go" files
		excludePrefix = append(excludePrefix, "app/ccv_msg_filter_")
		excludePrefix = append(excludePrefix, "app/ante_handler.go")
	}
	g := genny.New()
	if err := g.SelectiveFS(subfs, includePrefix, includeSuffix, excludePrefix, excludeSuffix); err != nil {
		return g, fmt.Errorf("generator fs: %w", err)
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
