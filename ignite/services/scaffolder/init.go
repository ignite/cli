package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/ignite/cli/ignite/templates/app"
	"github.com/ignite/cli/ignite/templates/field"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
	"github.com/ignite/cli/ignite/templates/testutil"
)

// Init initializes a new app with name and given options.
func Init(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	root, name, addressPrefix string,
	noDefaultModule, skipGit bool,
	params []string,
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
	if err := generate(
		ctx,
		tracer,
		pathInfo,
		addressPrefix,
		path,
		noDefaultModule,
		params,
	); err != nil {
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
	params []string,
) error {
	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	githubPath := gomodulepath.ExtractAppPath(pathInfo.RawPath)
	if !strings.Contains(githubPath, "/") {
		// A username must be added when the app module path has a single element
		githubPath = fmt.Sprintf("username/%s", githubPath)
	}

	g, err := app.NewGenerator(&app.Options{
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
	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, absRoot); err != nil {
		return err
	}

	run := func(runner *genny.Runner, gen *genny.Generator) error {
		if err := runner.With(gen); err != nil {
			return err
		}
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
			Params:     paramsFields,
			IsIBC:      false,
		}
		g, err = modulecreate.NewGenerator(opts)
		if err != nil {
			return err
		}
		if err := run(genny.WetRunner(ctx), g); err != nil {
			return err
		}
		g = modulecreate.NewAppModify(tracer, opts)
		if err := run(genny.WetRunner(ctx), g); err != nil {
			return err
		}

	}
	return gocmd.ModTidy(ctx, absRoot)
}
