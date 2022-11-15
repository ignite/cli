package dry

import (
	"embed"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/typed"
)

//go:embed files/component/* files/component/**/*
var fsComponent embed.FS

// NewGenerator returns the generator to scaffold a basic type in  module.
func NewGenerator(opts *typed.Options) (*genny.Generator, error) {
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
