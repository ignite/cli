package doc

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/*
var files embed.FS

// Options represents the options to scaffold a migration document.
type Options struct {
	Path     string
	FileName string
	Config   string
}

// NewGenerator returns the generator to scaffold a migration doc.
func NewGenerator(opts Options) (*genny.Generator, error) {
	subFs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("Config", opts.Config)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{Name}}", opts.FileName))

	return g, nil
}
