package dry

import (
	"embed"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/typed"
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
		)
	)
	return g, typed.Box(template, opts, g)
}
