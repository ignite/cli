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

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/xgenny"
)

//go:embed template/*
var fsPluginSource embed.FS

// Scaffold generates a plugin structure under dir/path.Base(moduleName).
func Scaffold(dir, moduleName string, sharedHost bool) (string, error) {
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
		return "", errors.Errorf("dir %q already exists, abort scaffolding", finalDir)
	}
	if err := g.Box(template); err != nil {
		return "", errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", moduleName)
	ctx.Set("Name", name)
	ctx.Set("SharedHost", sharedHost)

	g.Transformer(xgenny.Transformer(ctx))
	r := genny.WetRunner(ctx)
	err := r.With(g)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if err := r.Run(); err != nil {
		return "", errors.WithStack(err)
	}
	// FIXME(tb) we need to disable sumdb to get the branch version of CLI
	// because our git history is too fat.
	opt := exec.StepOption(step.Env("GOSUMDB=off"))
	if err := gocmd.ModTidy(context.TODO(), finalDir, opt); err != nil {
		return "", errors.WithStack(err)
	}
	return finalDir, nil
}
