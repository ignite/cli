package testutil

import (
	"embed"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/xgenny"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS
)

// Register testutil template using existing generator.
// Register is meant to be used by modules that depend on this module.
func Register(gen *genny.Generator, appPath string) error {
	return xgenny.Box(gen, xgenny.NewEmbedWalker(fsStargate, "stargate/", appPath))
}
