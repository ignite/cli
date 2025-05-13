package testutil

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

//go:embed files/* files/**/*
var files embed.FS

// Register testutil template using existing generator.
// Register is meant to be used by modules that depend on this module.
func Register(gen *genny.Generator) error {
	subFs, err := fs.Sub(files, "files")
	if err != nil {
		return errors.Errorf("fail to generate sub: %w", err)
	}

	return gen.OnlyFS(subFs, nil, nil)
}
