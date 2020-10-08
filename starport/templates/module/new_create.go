package module

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
	cosmosver.Launchpad: packr.New("module/templates/launchpad", "./launchpad"),
	cosmosver.Stargate:  packr.New("module/templates/stargate", "./stargate"),
}

// New ...
func NewCreateLaunchpad(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(templates[cosmosver.Launchpad]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}

// New ...
func NewCreateStargate(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(templates[cosmosver.Stargate]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}
