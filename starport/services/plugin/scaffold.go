package plugin

import (
	"embed"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

var (
	//go:embed template/*
	fsPluginSource embed.FS
)

// Scaffold generates a plugin structure under dir/path.Base(moduleName).
func Scaffold(dir, moduleName string) error {
	var (
		name     = filepath.Base(moduleName)
		finalDir = path.Join(dir, name)
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsPluginSource,
			"template",
			finalDir,
		)
	)
	if _, err := os.Stat(finalDir); err == nil {
		// finalDir already exists, don't overwrite stuff
		return errors.Errorf("dir %q already exists, abort scaffolding", finalDir)
	}
	if err := g.Box(template); err != nil {
		return errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", moduleName)
	ctx.Set("Name", name)
	// plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	r := genny.WetRunner(ctx)
	err := r.With(g)
	if err != nil {
		return errors.WithStack(err)
	}
	err = r.Run()
	return errors.WithStack(err)
}
