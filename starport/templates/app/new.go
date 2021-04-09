package app

import (
	"embed"
	"github.com/tendermint/starport/starport/pkg/templateutils"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	//go:embed launchpad/* launchpad/**/*
	fsLaunchpad embed.FS

	// these needs to be created in the compiler time, otherwise packr2 won't be
	// able to find boxes.
	templates = map[cosmosver.MajorVersion]packd.Walker{
		cosmosver.Stargate:  xgenny.NewEmbedWalker(fsStargate, "stargate/"),
		cosmosver.Launchpad: xgenny.NewEmbedWalker(fsLaunchpad, "launchpad/"),
	}
)

// New ...
func New(sdkVersion cosmosver.MajorVersion, opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(templates[sdkVersion]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("OwnerName", opts.OwnerName)
	ctx.Set("BinaryNamePrefix", opts.BinaryNamePrefix)
	ctx.Set("AddressPrefix", opts.AddressPrefix)
	ctx.Set("title", strings.Title)

	// Used for proto package name
	ctx.Set("formatOwnerName", templateutils.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))
	return g, nil
}
