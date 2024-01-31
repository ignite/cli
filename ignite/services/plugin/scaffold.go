package plugin

import (
	"context"
	"embed"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
)

//go:embed template/*
var fsPluginSource embed.FS

// Scaffold generates a plugin structure under dir/path.Base(moduleName).
func Scaffold(ctx context.Context, dir, moduleName string, sharedHost bool) (string, error) {
	var (
		name     = filepath.Base(moduleName)
		title    = toTitle(name)
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
	pctx.Set("Title", title)
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

	if err := gocmd.ModTidy(ctx, finalDir); err != nil {
		return "", errors.WithStack(err)
	}
	if err := gocmd.Fmt(ctx, finalDir); err != nil {
		return "", errors.WithStack(err)
	}

	return finalDir, nil
}

func toTitle(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(cases.Title(language.English).String(s), "_", ""), "-", "")
}
