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
	"github.com/tendermint/starport/starport/templates/typed/dry"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
	"github.com/tendermint/starport/starport/templates/typed/singleton"
)

// AddTypeOption configures options for AddType.
type AddTypeOption func(*addTypeOptions)

// AddTypeKind configures the type kind option for AddType.
type AddTypeKind func(*addTypeOptions)

type addTypeOptions struct {
	moduleName string
	fields     []string

	isList      bool
	isMap       bool
	isSingleton bool

	indexes []string

	withoutMessage bool
}

// ListType makes the type stored in a list convention in the storage.
func ListType() AddTypeKind {
	return func(o *addTypeOptions) {
		o.isList = true
	}
}

// MapType makes the type stored in a key-value convention in the storage with a custom
// index option.
func MapType(indexes ...string) AddTypeKind {
	return func(o *addTypeOptions) {
		o.isMap = true
		o.indexes = indexes
	}
}

// SingletonType makes the type stored in a fixed place as a single entry in the storage.
func SingletonType() AddTypeKind {
	return func(o *addTypeOptions) {
		o.isSingleton = true
	}
}

// DryType only creates a type with a basic definition.
func DryType() AddTypeKind {
	return func(o *addTypeOptions) {}
}

// TypeWithModule module to scaffold type into.
func TypeWithModule(name string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.moduleName = name
	}
}

// TypeWithFields adds fields to the type to be scaffolded.
func TypeWithFields(fields ...string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.fields = fields
	}
}

// TypeWithoutMessage disables generating sdk compatible messages and tx related APIs.
func TypeWithoutMessage(fields ...string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.withoutMessage = true
	}
}

// AddType adds a new type to a scaffolded app.
// if non of the list, map or singleton given, a dry type without anything extra (like a storage layer, models, CLI etc.)
// will be scaffolded.
// if no module is given, the type will be scaffolded inside the app's default module.
func (s *Scaffolder) AddType(
	typeName string,
	tracer *placeholder.Tracer,
	kind AddTypeKind,
	options ...AddTypeOption,
) (sm xgenny.SourceModification, err error) {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	// apply options.
	o := addTypeOptions{
		moduleName: path.Package, // app's default module.
	}

	for _, apply := range append(options, AddTypeOption(kind)) {
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
		g, err = mapGenerator(tracer, opts, o.indexes)
	case o.isSingleton:
		g, err = singleton.NewStargate(tracer, opts)
	default:
		g, err = dry.NewStargate(opts)
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
		"appendedValue",
		"creator":
		return fmt.Errorf("%s is used by type scaffolder", name)
	}

	return checkGoReservedWord(name)
}

// mapGenerator returns the template generator for a map
func mapGenerator(replacer placeholder.Replacer, opts *typed.Options, indexes []string) (*genny.Generator, error) {
	// Parse indexes with the associated type
	parsedIndexes, err := field.ParseFields(indexes, checkForbiddenTypeField)
	if err != nil {
		return nil, err
	}

	// Indexes and type fields must be disjoint
	exists := make(map[string]struct{})
	for _, field := range opts.Fields {
		exists[field.Name.LowerCamel] = struct{}{}
	}
	for _, index := range parsedIndexes {
		if _, ok := exists[index.Name.LowerCamel]; ok {
			return nil, fmt.Errorf("%s cannot simultaneously be an index and a field", index.Name.Original)
		}
	}

	opts.Indexes = parsedIndexes
	return indexed.NewStargate(replacer, opts)
}
