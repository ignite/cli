package dry

import (
	"embed"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/typed"
)

//go:embed files/component/* files/component/**/*
var fsComponent embed.FS

// NewStargate returns the generator to scaffold a basic type in a Stargate module.
func NewStargate(opts *typed.Options) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsComponent,
			"files/component/",
			opts.AppPath,
		)
	)
	return g, typed.Box(template, opts, g)
}
