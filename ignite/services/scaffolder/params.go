package scaffolder

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v28/ignite/templates/module/create"
)

// CreateParams creates a new params in the scaffolded module.
func (s Scaffolder) CreateParams(
	ctx context.Context,
	cacheStorage cache.Storage,
	runner *xgenny.Runner,
	moduleName string,
	params ...string,
) error {
	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = s.modpath.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfName.LowerCase

	// Check if the module already exist
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %v not exist", moduleName)
	}

	if err := checkParamCreated(s.path, moduleName, params); err != nil {
		return err
	}

	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	opts := modulecreate.ParamsOptions{
		ModuleName: moduleName,
		Params:     paramsFields,
		AppName:    s.modpath.Package,
		AppPath:    s.path,
	}

	g, err := modulecreate.NewModuleParam(opts)
	if err != nil {
		return err
	}
	gens := []*genny.Generator{g}

	if err := runner.Run(gens...); err != nil {
		return err
	}
	return finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// checkParamCreated checks if the parameter has been already created.
func checkParamCreated(appPath, moduleName string, params []string) error {
	path := filepath.Join(appPath, "x", moduleName, "types")
	ok, err := goanalysis.HasAnyStructFieldsInPkg(path, "Params", params)
	if err != nil {
		return err
	}

	if ok {
		return errors.Errorf(
			"duplicated params (%s) module %s",
			strings.Join(params, " "),
			moduleName,
		)
	}
	return nil
}
