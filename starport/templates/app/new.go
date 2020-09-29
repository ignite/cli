package app

import (
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var templates = map[cosmosver.MajorVersion]*packr.Box{
	cosmosver.Launchpad: packr.New("app/templates/launchpad", "./launchpad"),
	cosmosver.Stargate:  packr.New("app/templates/stargate", "./stargate"),
}

// New ...
func New(sdkVersion cosmosver.MajorVersion, opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(templates[sdkVersion]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
	ctx.Set("title", strings.Title)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))
	return g, nil
}
