package testutil

import (
	"embed"
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

const modulePathKey = "ModulePath"

var (
	//go:embed stargate/* stargate/**/*
	fs embed.FS

	testutilTemplate = xgenny.NewEmbedWalker(fs, "stargate/")
)

// Register testutil template using existing generator.
// Register is meant to be used by modules that depend on this module.
//nolint:interfacer
func Register(ctx *plush.Context, gen *genny.Generator) error {
	if !ctx.Has(modulePathKey) {
		return fmt.Errorf("ctx is missing value for the ket %s", modulePathKey)
	}
	return gen.Box(testutilTemplate)
}
