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
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v28/ignite/templates/module/create"
)

// CreateConfigs creates a new configs in the scaffolded module.
func (s Scaffolder) CreateConfigs(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName string,
	configs ...string,
) (sm xgenny.SourceModification, err error) {
	appName := s.modpath.Package
	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = s.modpath.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	// Check if the module already exist
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if !ok {
		return sm, errors.Errorf("the module %v not exist", moduleName)
	}

	if err := checkConfigCreated(s.path, appName, moduleName, configs); err != nil {
		return sm, err
	}

	// Parse config with the associated type
	configsFields, err := field.ParseFields(configs, checkForbiddenTypeIndex)
	if err != nil {
		return sm, err
	}

	opts := modulecreate.ConfigsOptions{
		ModuleName: moduleName,
		Configs:    configsFields,
		AppName:    s.modpath.Package,
		AppPath:    s.path,
	}

	g, err := modulecreate.NewModuleConfigs(opts)
	if err != nil {
		return sm, err
	}
	gens := []*genny.Generator{g}

	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}

	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// checkConfigCreated checks if the config has been already created.
func checkConfigCreated(appPath, appName, moduleName string, configs []string) (err error) {
	path := filepath.Join(appPath, "api", appName, moduleName, "module")
	ok, err := goanalysis.HasAnyStructFieldsInPkg(path, "Module", configs)
	if err != nil {
		return err
	}

	if ok {
		return errors.Errorf(
			"duplicated configs (%s) module %s",
			strings.Join(configs, " "),
			moduleName,
		)
	}
	return nil
}
