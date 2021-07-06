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
	"github.com/tendermint/starport/starport/templates/typed/basic"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
	"github.com/tendermint/starport/starport/templates/typed/singleton"
)

// AddTypeOption configures options for AddType.
type AddTypeOption func(*addTypeOptions)

type addTypeOptions struct {
	moduleName string
	fields     []string

	isList      bool
	isMap       bool
	isSingleton bool

	index string

	withoutMessage bool
}

// ListType enables adding a list storage.
func ListType() AddTypeOption {
	return func(o *addTypeOptions) {
		o.isList = true
	}
}

// MapType enables adding a map storage.
func MapType(index string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.isMap = true
		o.index = index
	}
}

// SingletonType enables adding a singleton storage.
func SingletonType() AddTypeOption {
	return func(o *addTypeOptions) {
		o.isSingleton = true
	}
}

// TypeWithModule specifies module to scaffold type inside.
func TypeWithModule(name string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.moduleName = name
	}
}

// TypeWithFields adds fields to the type to be scaffold.
func TypeWithFields(fields ...string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.fields = fields
	}
}

// TypeWithoutMessage disables CRUD for type.
func TypeWithoutMessage(fields ...string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.withoutMessage = true
	}
}

// AddType adds a new type stype to scaffolded app.
// if non of the list, map or singleton given, a basic type without anything extra will be scaffold.
// if no module given, type will be scaffold inside the app's default module.
func (s *Scaffolder) AddType(
	typeName string,
	tracer *placeholder.Tracer,
	options ...AddTypeOption,
) (sm xgenny.SourceModification, err error) {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	// apply options.
	o := addTypeOptions{
		moduleName: path.Package,
	}

	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(o.moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName := mfName.Lowercase

	name, err := multiformatname.NewName(typeName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name); err != nil {
		return sm, err
	}

	tFields, err := field.ParseFields(o.fields, checkForbiddenTypeField)
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
			NoMessage:  o.withoutMessage,
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

	// create the type generator depending on the model
	// TODO: rename the template packages to make it consistent with the type new naming
	switch {
	case o.isList:
		g, err = typed.NewStargate(tracer, opts)
	case o.isMap:
		g, err = indexed.NewStargate(tracer, opts)
	case o.isSingleton:
		g, err = singleton.NewStargate(tracer, opts)
	default:
		g, err = basic.NewStargate(opts)
	}
	if err != nil {
		return sm, err
	}

	// run the generation
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
