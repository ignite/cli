package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/templates/app"
	"github.com/ignite/cli/v29/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
	"github.com/ignite/cli/v29/ignite/templates/testutil"
)

// Init initializes a new app with name and given options.
func Init(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	root, name, addressPrefix string,
	noDefaultModule, skipGit, skipProto, minimal, isConsumerChain bool,
	params, moduleConfigs []string,
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
		minimal,
		isConsumerChain,
		params,
		moduleConfigs,
	); err != nil {
		return "", err
	}

	if err = finish(ctx, cacheStorage, path, pathInfo.RawPath, skipProto); err != nil {
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
	noDefaultModule, minimal, isConsumerChain bool,
	params, moduleConfigs []string,
) error {
	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	// Parse configs with the associated type
	configsFields, err := field.ParseFields(moduleConfigs, checkForbiddenTypeIndex)
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
		IsChainMinimal:   minimal,
		IsConsumerChain:  isConsumerChain,
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

	if err := cosmosgen.InstallDepTools(ctx, absRoot); err != nil {
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
			Configs:    configsFields,
			IsIBC:      false,
		}
		// Check if the module name is valid
		if err := checkModuleName(opts.AppPath, opts.ModuleName); err != nil {
			return err
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

	return nil
}
