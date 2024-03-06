package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/app"
	"github.com/ignite/cli/v28/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v28/ignite/templates/module/create"
	"github.com/ignite/cli/v28/ignite/templates/testutil"
)

// Init initializes a new app with name and given options.
func Init(
	ctx context.Context,
	cacheStorage cache.Storage,
	runner *xgenny.Runner,
	root, name, addressPrefix string,
	noDefaultModule, skipGit, minimal bool,
	params, moduleConfigs []string,
) (path string, gomodule string, err error) {
	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", "", err
	}

	// Create a new folder named as the blockchain when a custom path is not specified
	var appFolder string
	if root == "" {
		appFolder = pathInfo.Root
	}

	if root, err = filepath.Abs(root); err != nil {
		return path, gomodule, err
	}
	path = filepath.Join(root, appFolder)
	gomodule = pathInfo.RawPath

	// create the project
	if _, err = generate(
		runner,
		pathInfo,
		addressPrefix,
		path,
		noDefaultModule,
		minimal,
		params,
		moduleConfigs,
	); err != nil {
		return path, gomodule, err
	}
	return path, pathInfo.RawPath, nil
}

//nolint:interfacer
func generate(
	runner *xgenny.Runner,
	pathInfo gomodulepath.Path,
	addressPrefix,
	absRoot string,
	noDefaultModule, minimal bool,
	params, moduleConfigs []string,
) (xgenny.SourceModification, error) {
	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return xgenny.SourceModification{}, err
	}

	// Parse configs with the associated type
	configsFields, err := field.ParseFields(moduleConfigs, checkForbiddenTypeIndex)
	if err != nil {
		return xgenny.SourceModification{}, err
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
	})
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, absRoot); err != nil {
		return xgenny.SourceModification{}, err
	}

	// generate module template
	smc, err := runner.RunAndApply(g)
	if err != nil {
		return smc, err
	}

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
			return smc, err
		}

		g, err = modulecreate.NewGenerator(runner.Tracer(), opts)
		if err != nil {
			return smc, err
		}

		runner.Path = runner.TempPath
		// generate module template
		smm, err := runner.RunAndApply(g)
		if err != nil {
			return smc, err
		}
		smc.Merge(smm)
	}

	return smc, err
}
