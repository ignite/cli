package scaffolder

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/datatype"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
	"github.com/ignite/cli/ignite/templates/typed"
	"github.com/ignite/cli/ignite/templates/typed/dry"
	"github.com/ignite/cli/ignite/templates/typed/list"
	maptype "github.com/ignite/cli/ignite/templates/typed/map"
	"github.com/ignite/cli/ignite/templates/typed/singleton"
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

	withoutMessage    bool
	withoutSimulation bool
	signer            string
}

// newAddTypeOptions returns a addTypeOptions with default options.
func newAddTypeOptions(moduleName string) addTypeOptions {
	return addTypeOptions{
		moduleName: moduleName,
		signer:     "creator",
	}
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
func TypeWithoutMessage() AddTypeOption {
	return func(o *addTypeOptions) {
		o.withoutMessage = true
	}
}

// TypeWithoutSimulation disables generating messages simulation.
func TypeWithoutSimulation() AddTypeOption {
	return func(o *addTypeOptions) {
		o.withoutSimulation = true
	}
}

// TypeWithSigner provides a custom signer name for the message.
func TypeWithSigner(signer string) AddTypeOption {
	return func(o *addTypeOptions) {
		o.signer = signer
	}
}

// AddType adds a new type to a scaffolded app.
// if none of the list, map or singleton given, a dry type without anything extra (like a storage layer, models, CLI etc.)
// will be scaffolded.
// if no module is given, the type will be scaffolded inside the app's default module.
func (s Scaffolder) AddType(
	ctx context.Context,
	cacheStorage cache.Storage,
	typeName string,
	tracer *placeholder.Tracer,
	kind AddTypeKind,
	options ...AddTypeOption,
) (sm xgenny.SourceModification, err error) {
	// apply options.
	o := newAddTypeOptions(s.modpath.Package)
	for _, apply := range append(options, AddTypeOption(kind)) {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(o.moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName := mfName.LowerCase

	name, err := multiformatname.NewName(typeName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name, o.withoutMessage); err != nil {
		return sm, err
	}

	// Check and parse provided fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, o.fields); err != nil {
		return sm, err
	}
	tFields, err := parseTypeFields(o)
	if err != nil {
		return sm, err
	}

	mfSigner, err := multiformatname.NewName(o.signer)
	if err != nil {
		return sm, err
	}

	isIBC, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return sm, err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:      s.modpath.Package,
			AppPath:      s.path,
			ModulePath:   s.modpath.RawPath,
			ModuleName:   moduleName,
			TypeName:     name,
			Fields:       tFields,
			NoMessage:    o.withoutMessage,
			NoSimulation: o.withoutSimulation,
			MsgSigner:    mfSigner,
			IsIBC:        isIBC,
		}
		gens []*genny.Generator
	)
	// Check and support MsgServer convention
	gens, err = supportMsgServer(
		gens,
		tracer,
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			AppPath:    opts.AppPath,
		},
	)
	if err != nil {
		return sm, err
	}

	gens, err = supportGenesisTests(
		gens,
		opts.AppPath,
		opts.AppName,
		opts.ModulePath,
		opts.ModuleName,
	)
	if err != nil {
		return sm, err
	}

	gens, err = supportSimulation(
		gens,
		opts.AppPath,
		opts.ModulePath,
		opts.ModuleName,
	)
	if err != nil {
		return sm, err
	}

	// create the type generator depending on the model
	switch {
	case o.isList:
		g, err = list.NewGenerator(tracer, opts)
	case o.isMap:
		g, err = mapGenerator(tracer, opts, o.indexes)
	case o.isSingleton:
		g, err = singleton.NewGenerator(tracer, opts)
	default:
		g, err = dry.NewGenerator(opts)
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

	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// checkForbiddenTypeIndex returns true if the name is forbidden as a index name.
func checkForbiddenTypeIndex(index string) error {
	indexSplit := strings.Split(index, datatype.Separator)
	if len(indexSplit) > 1 {
		index = indexSplit[0]
		indexType := datatype.Name(indexSplit[1])
		if f, ok := datatype.IsSupportedType(indexType); !ok || f.NonIndex {
			return fmt.Errorf("invalid index type %s", indexType)
		}
	}
	return checkForbiddenTypeField(index)
}

// checkForbiddenTypeField returns true if the name is forbidden as a field name.
func checkForbiddenTypeField(field string) error {
	mfName, err := multiformatname.NewName(field)
	if err != nil {
		return err
	}

	switch mfName.LowerCase {
	case
		"id",
		"params",
		"appendedvalue",
		datatype.TypeCustom:
		return fmt.Errorf("%s is used by type scaffolder", field)
	}

	return checkGoReservedWord(field)
}

// parseTypeFields validates the fields and returns an error if the validation fails.
func parseTypeFields(opts addTypeOptions) (field.Fields, error) {
	signer := ""
	if opts.isList || opts.isMap || opts.isSingleton {
		if !opts.withoutMessage {
			signer = opts.signer
		}
		return field.ParseFields(opts.fields, checkForbiddenTypeField, signer)
	}
	// For simple types, only check if it's a reserved keyword and don't pass a signer.
	return field.ParseFields(opts.fields, checkGoReservedWord, signer)
}

// mapGenerator returns the template generator for a map.
func mapGenerator(replacer placeholder.Replacer, opts *typed.Options, indexes []string) (*genny.Generator, error) {
	// Parse indexes with the associated type
	parsedIndexes, err := field.ParseFields(indexes, checkForbiddenTypeIndex)
	if err != nil {
		return nil, err
	}

	// Indexes and type fields must be disjoint
	exists := make(map[string]struct{})
	for _, name := range opts.Fields {
		exists[name.Name.LowerCamel] = struct{}{}
	}
	for _, index := range parsedIndexes {
		if _, ok := exists[index.Name.LowerCamel]; ok {
			return nil, fmt.Errorf("%s cannot simultaneously be an index and a field", index.Name.Original)
		}
	}

	opts.Indexes = parsedIndexes
	return maptype.NewGenerator(replacer, opts)
}
