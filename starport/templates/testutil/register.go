package testutil

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

const (
	modulePathKey = "ModulePath"
	testUtilDir   = "testutil"
	sampleDir     = "sample"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS
	//go:embed stargate/testutil/sample/*
	fsSample embed.FS

	stargateTemplate = xgenny.NewEmbedWalker(fsStargate, "stargate/")
	sampleTemplate   = xgenny.NewEmbedWalker(fsSample, "stargate/")
)

// Register testutil template using existing generator.
// Register is meant to be used by modules that depend on this module.
//nolint:interfacer
func Register(ctx *plush.Context, gen *genny.Generator, appPath string) error {
	if !ctx.Has(modulePathKey) {
		return fmt.Errorf("ctx is missing value for the key %s", modulePathKey)
	}
	path := filepath.Join(appPath, testUtilDir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return gen.Box(stargateTemplate)
	} else if err != nil {
		return err
	}

	path = filepath.Join(path, sampleDir)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return gen.Box(sampleTemplate)
	}
	return err
}
