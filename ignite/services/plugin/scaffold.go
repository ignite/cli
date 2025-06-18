package plugin

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
)

//go:embed template/*
var fsPluginSource embed.FS

// Scaffold generates a plugin structure under dir/path.Base(appName).
func Scaffold(ctx context.Context, session *cliui.Session, dir, appName string, sharedHost bool) (string, error) {
	subFs, err := fs.Sub(fsPluginSource, "template")
	if err != nil {
		return "", errors.WithStack(err)
	}

	var (
		name     = filepath.Base(appName)
		title    = toTitle(name)
		finalDir = path.Join(dir, name)
	)
	if _, err := os.Stat(finalDir); err == nil {
		// finalDir already exists, don't overwrite stuff
		return "", errors.Errorf("directory %q already exists, abort scaffolding", finalDir)
	}

	g := genny.New()
	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return "", errors.WithStack(err)
	}

	pctx := plush.NewContextWithContext(ctx)
	pctx.Set("AppName", appName)
	pctx.Set("Name", name)
	pctx.Set("Title", title)
	pctx.Set("SharedHost", sharedHost)

	g.Transformer(xgenny.Transformer(pctx))
	r := xgenny.NewRunner(ctx, finalDir)
	_, err = r.RunAndApply(g, xgenny.ApplyPreRun(func(_, _, duplicated []string) error {
		if len(duplicated) == 0 {
			return nil
		}
		question := fmt.Sprintf("Do you want to overwrite the existing files? \n%s", strings.Join(duplicated, "\n"))
		return session.AskConfirm(question)
	}))
	if err != nil {
		return "", err
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
