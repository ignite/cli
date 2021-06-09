package scaffolder

import (
	"fmt"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	"github.com/tendermint/starport/starport/templates/typed"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
	"os"
)

type AddTypeOption struct {
	Indexed   bool
	NoMessage bool
}

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddType(
	tracer *placeholder.Tracer,
	addTypeOptions AddTypeOption,
	moduleName,
	typeName string,
	fields ...string,
) error {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	if err := checkComponentValidity(s.path, moduleName, typeName); err != nil {
		return err
	}

	// Parse provided field
	tFields, err := field.ParseFields(fields, checkForbiddenTypeField)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			TypeName:   typeName,
			Fields:     tFields,
			NoMessage:  addTypeOptions.NoMessage,
		}
		gens []*genny.Generator
	)
	// Check and support MsgServer convention
	g, err = supportMsgServer(
		tracer,
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			OwnerName:  opts.OwnerName,
		},
	)
	if err != nil {
		return err
	}
	if g != nil {
		gens = append(gens, g)
	}

	// Check if indexed type
	if addTypeOptions.Indexed {
		g, err = indexed.NewStargate(tracer, opts)
	} else {
		// Scaffolding a type with ID
		g, err = typed.NewStargate(tracer, opts)
	}
	if err != nil {
		return err
	}
	gens = append(gens, g)
	if err := xgenny.RunWithValidation(tracer, gens...); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return s.finish(pwd, path.RawPath)
}

// checkForbiddenTypeField returns true if the name is forbidden as a field name
func checkForbiddenTypeField(name string) error {
	switch name {
	case
		"id",
		"index",
		"appendedValue",
		"creator":
		return fmt.Errorf("%s is used by type scaffolder", name)
	}

	return checkGoReservedWord(name)
}
