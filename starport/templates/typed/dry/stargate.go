package dry

import (
	"embed"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/typed"
)

var (
	//go:embed stargate/component/* stargate/component/**/*
	fsStargateComponent embed.FS
)

// NewStargate returns the generator to scaffold a basic type in a Stargate module.
func NewStargate(opts *typed.Options) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsStargateComponent,
			"stargate/component/",
			opts.AppPath,
			false,
		)
	)
	return g, typed.Box(template, opts, g)
}
