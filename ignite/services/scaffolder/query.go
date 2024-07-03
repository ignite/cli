package scaffolder

import (
	"context"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field"
	"github.com/ignite/cli/v28/ignite/templates/query"
)

// AddQuery adds a new query to scaffolded app.
func (s Scaffolder) AddQuery(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName,
	queryName,
	description string,
	reqFields,
	resFields []string,
	paginated bool,
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

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.appPath, moduleName, name, true); err != nil {
		return err
	}

	// Check and parse provided request fields
	if ok := containsCustomTypes(reqFields); ok {
		return sm, errors.New("query request params can't contain custom type")
	}
	parsedReqFields, err := field.ParseFields(reqFields, checkGoReservedWord)
	if err != nil {
		return sm, err
	}

	// Check and parse provided response fields
	if err := checkCustomTypes(ctx, s.appPath, s.modpath.Package, s.protoDir, moduleName, resFields); err != nil {
		return err
	}
	parsedResFields, err := field.ParseFields(resFields, checkGoReservedWord)
	if err != nil {
		return sm, err
	}

	var (
		g    *genny.Generator
		opts = &query.Options{
			AppName:     s.modpath.Package,
			AppPath:     s.appPath,
			ProtoDir:    s.protoDir,
			ModulePath:  s.modpath.RawPath,
			ModuleName:  moduleName,
			QueryName:   name,
			ReqFields:   parsedReqFields,
			ResFields:   parsedResFields,
			Description: description,
			Paginated:   paginated,
		}
	)

	// Scaffold
	g, err = query.NewGenerator(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath, false)
}
