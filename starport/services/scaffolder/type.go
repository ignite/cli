package scaffolder

import (
	"fmt"
	"os"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	"github.com/tendermint/starport/starport/templates/typed"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
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
) (sm xgenny.SourceModification, err error) {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.Lowercase

	name, err := multiformatname.NewName(typeName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name); err != nil {
		return sm, err
	}

	// Parse provided field
	tFields, err := field.ParseFields(fields, checkForbiddenTypeField)
	if err != nil {
		return sm, err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			TypeName:   name,
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
		return sm, err
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
		return sm, err
	}
	gens = append(gens, g)
	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return sm, err
	}
	return sm, s.finish(pwd, path.RawPath)
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
