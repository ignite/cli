package dry

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

//go:embed files/component/* files/component/**/*
var fsComponent embed.FS

// NewGenerator returns the generator to scaffold a basic type in  module.
func NewGenerator(opts *typed.Options) (*genny.Generator, error) {
	subFs, err := fs.Sub(fsComponent, "files/component")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}
	g := genny.New()
	return g, typed.Box(subFs, opts, g)
}
