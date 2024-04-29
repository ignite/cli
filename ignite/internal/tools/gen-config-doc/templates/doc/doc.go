package doc

import (
	"embed"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
)

//go:embed files/*
var fsFiles embed.FS

// Options represents the options to scaffold a migration document.
type Options struct {
	Path   string
	Config string
}

// NewGenerator returns the generator to scaffold a migration doc.
func NewGenerator(opts Options) (*genny.Generator, error) {
	var (
		g           = genny.New()
		docTemplate = xgenny.NewEmbedWalker(
			fsFiles,
			"files/",
			opts.Path,
		)
	)

	if err := g.Box(docTemplate); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("Config", opts.Config)
	g.Transformer(xgenny.Transformer(ctx))

	return g, nil
}
