package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/ignite/cli/ignite/templates/app"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

// Init initializes a new app with name and given options.
func Init(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	root, name, addressPrefix string,
	noDefaultModule, skipGit bool,
) (path string, err error) {
	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", err
	}

	// Create a new folder named as the blockchain when a custom path is not specified
	var appFolder string
	if root == "" {
		appFolder = pathInfo.Root
	}

	if root, err = filepath.Abs(root); err != nil {
		return "", err
	}

	path = filepath.Join(root, appFolder)

	// create the project
	if err := generate(ctx, tracer, pathInfo, addressPrefix, path, noDefaultModule); err != nil {
		return "", err
	}

	if err := finish(ctx, cacheStorage, path, pathInfo.RawPath); err != nil {
		return "", err
	}

	if !skipGit {
		// Initialize git repository and perform the first commit
		if err := xgit.InitAndCommit(path); err != nil {
			return "", err
		}
	}

	return path, nil
}

//nolint:interfacer
func generate(
	ctx context.Context,
	tracer *placeholder.Tracer,
	pathInfo gomodulepath.Path,
	addressPrefix,
	absRoot string,
	noDefaultModule bool,
) error {
	githubPath := gomodulepath.ExtractAppPath(pathInfo.RawPath)
	if !strings.Contains(githubPath, "/") {
		// A username must be added when the app module path has a single element
		githubPath = fmt.Sprintf("username/%s", githubPath)
	}

	g, err := app.New(&app.Options{
		// generate application template
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		AppPath:          absRoot,
		GitHubPath:       githubPath,
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    addressPrefix,
	})
	if err != nil {
		return err
	}

	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		runner.Root = absRoot
		return runner.Run()
	}
	if err := run(genny.WetRunner(ctx), g); err != nil {
		return err
	}

	// generate module template
	if !noDefaultModule {
		opts := &modulecreate.CreateOptions{
			ModuleName: pathInfo.Package, // App name
			ModulePath: pathInfo.RawPath,
			AppName:    pathInfo.Package,
			AppPath:    absRoot,
			IsIBC:      false,
		}
		g, err = modulecreate.NewGenerator(opts)
		if err != nil {
			return err
		}
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}
		g = modulecreate.NewAppModify(tracer, opts)
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}

	}

	// FIXME(tb) untagged version of ignite/cli triggers a 404 not found when go
	// mod tidy requests the sumdb, until we understand why, we disable sumdb.
	// related issue:  https://github.com/golang/go/issues/56174
	opt := exec.StepOption(step.Env("GOSUMDB=off"))

	return gocmd.ModTidy(ctx, absRoot, opt)
}
