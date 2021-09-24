package testutil

import (
	"embed"
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

const (
	modulePathKey = "ModulePath"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS
)

// Register testutil template using existing generator.
// Register is meant to be used by modules that depend on this module.
//nolint:interfacer
func Register(ctx *plush.Context, gen *genny.Generator, appPath string) error {
	if !ctx.Has(modulePathKey) {
		return fmt.Errorf("ctx is missing value for the key %s", modulePathKey)
	}
	// Check if the testutil folder already exists
	return xgenny.Box(gen, xgenny.NewEmbedWalker(fsStargate, "stargate/", appPath))
}
