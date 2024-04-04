package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/app"
	"github.com/ignite/cli/v29/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
	"github.com/ignite/cli/v29/ignite/templates/testutil"
)

// Init initializes a new app with name and given options.
func Init(
	ctx context.Context,
	runner *xgenny.Runner,
	root, name, addressPrefix string,
	noDefaultModule, minimal, isConsumerChain bool,
	params, moduleConfigs []string,
) (string, string, error) {
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
		return "", "", err
	}

	var (
		path     = filepath.Join(root, appFolder)
		gomodule = pathInfo.RawPath
	)
	// create the project
	_, err = generate(
		ctx,
		runner,
		pathInfo,
		addressPrefix,
		path,
		noDefaultModule,
		minimal,
		isConsumerChain,
		params,
		moduleConfigs,
	)
	return path, gomodule, err
}

//nolint:interfacer
func generate(
	ctx context.Context,
	runner *xgenny.Runner,
	pathInfo gomodulepath.Path,
	addressPrefix,
	absRoot string,
	noDefaultModule, minimal, isConsumerChain bool,
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
		// A username must be added when the app module appPath has a single element
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
		return xgenny.SourceModification{}, err
	}
	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, absRoot); err != nil {
		return xgenny.SourceModification{}, err
	}

	// generate module template
	runner.Root = absRoot
	smc, err := runner.RunAndApply(g)
	if err != nil {
		return smc, err
	}

	if err := cosmosgen.InstallDepTools(ctx, absRoot); err != nil {
		return smc, err
	}

	if !noDefaultModule {
		opts := &modulecreate.CreateOptions{
			ModuleName: pathInfo.Package, // App name
			ModulePath: pathInfo.RawPath,
			AppName:    pathInfo.Package,
			AppPath:    absRoot,
			ProtoPath:  defaults.ProtoPath,
			Params:     paramsFields,
			Configs:    configsFields,
			IsIBC:      false,
		}
		// Check if the module name is valid
		if err := checkModuleName(opts.AppPath, opts.ModuleName); err != nil {
			return smc, err
		}

		g, err = modulecreate.NewGenerator(opts)
		if err != nil {
			return smc, err
		}

		// generate module template
		smm, err := runner.RunAndApply(g, modulecreate.NewAppModify(runner.Tracer(), opts))
		if err != nil {
			return smc, err
		}
		smc.Merge(smm)
	}

	return smc, err
}
