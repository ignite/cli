package scaffolder

import (
	"os"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/query"
)

// AddQuery adds a new query to scaffolded app
func (s *Scaffolder) AddQuery(
	moduleName,
	queryName,
	description string,
	reqFields,
	resFields []string,
	paginated bool,
) error {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	if err := checkComponentValidity(s.path, moduleName, queryName); err != nil {
		return err
	}

	// Parse provided fields
	parsedReqFields, err := parseFields(reqFields, checkGoReservedWord)
	if err != nil {
		return err
	}
	parsedResFields, err := parseFields(resFields, checkGoReservedWord)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &query.Options{
			AppName:     path.Package,
			ModulePath:  path.RawPath,
			ModuleName:  moduleName,
			OwnerName:   owner(path.RawPath),
			QueryName:   queryName,
			ReqFields:   parsedReqFields,
			ResFields:   parsedResFields,
			Description: description,
			Paginated:   paginated,
		}
	)

	// Scaffold
	g, err = query.NewStargate(s.tracer, opts)
	if err != nil {
		return err
	}
	if err := xgenny.RunWithValidation(s.tracer, g); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := s.protoc(pwd, path.RawPath); err != nil {
		return err
	}
	return fmtProject(pwd)
}
