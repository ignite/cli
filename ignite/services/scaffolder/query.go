package scaffolder

import (
	"context"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field"
	"github.com/ignite/cli/v28/ignite/templates/query"
)

// AddQuery adds a new query to scaffolded app.
func (s Scaffolder) AddQuery(
	ctx context.Context,
	runner *xgenny.Runner,
	moduleName,
	queryName,
	description string,
	reqFields,
	resFields []string,
	paginated bool,
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

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return err
	}

	if err := checkComponentValidity(s.path, moduleName, name, true); err != nil {
		return err
	}

	// Check and parse provided request fields
	if ok := containsCustomTypes(reqFields); ok {
		return errors.New("query request params can't contain custom type")
	}
	parsedReqFields, err := field.ParseFields(reqFields, checkGoReservedWord)
	if err != nil {
		return err
	}

	// Check and parse provided response fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, resFields); err != nil {
		return err
	}
	parsedResFields, err := field.ParseFields(resFields, checkGoReservedWord)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &query.Options{
			AppName:     s.modpath.Package,
			AppPath:     s.path,
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
	g, err = query.NewGenerator(runner.Tracer(), opts)
	if err != nil {
		return err
	}

	return runner.Run(g)
}
