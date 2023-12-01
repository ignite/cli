package plugin

import (
	"context"
	"embed"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/pkg/errors"

	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
)

//go:embed template/*
var fsPluginSource embed.FS

// Scaffold generates a plugin structure under dir/path.Base(moduleName).
func Scaffold(ctx context.Context, dir, moduleName string, sharedHost bool) (string, error) {
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
		return "", errors.Errorf("directory %q already exists, abort scaffolding", finalDir)
	}

	if err := g.Box(template); err != nil {
		return "", errors.WithStack(err)
	}

	pctx := plush.NewContextWithContext(ctx)
	pctx.Set("ModuleName", moduleName)
	pctx.Set("Name", name)
	pctx.Set("SharedHost", sharedHost)

	g.Transformer(xgenny.Transformer(pctx))
	r := genny.WetRunner(ctx)
	err := r.With(g)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err := r.Run(); err != nil {
		return "", errors.WithStack(err)
	}

	return finalDir, nil
}
