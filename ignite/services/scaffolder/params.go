package scaffolder

import (
	"context"
	"fmt"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

// CreateParams creates a new params in the scaffolded module.
func (s Scaffolder) CreateParams(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName string,
	params ...string,
) (sm xgenny.SourceModification, err error) {
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
		return sm, fmt.Errorf("the module %v not exist", moduleName)
	}

	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return sm, err
	}

	opts := modulecreate.ParamsOptions{
		ModuleName: moduleName,
		Params:     paramsFields,
		AppName:    s.modpath.Package,
		AppPath:    s.path,
	}

	g, err := modulecreate.NewModuleParam(tracer, opts)
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
