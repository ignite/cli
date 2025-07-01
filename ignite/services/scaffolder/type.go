package scaffolder

import (
	"context"
	"strings"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
	"github.com/ignite/cli/v29/ignite/templates/typed"
	"github.com/ignite/cli/v29/ignite/templates/typed/dry"
	"github.com/ignite/cli/v29/ignite/templates/typed/list"
	maptype "github.com/ignite/cli/v29/ignite/templates/typed/map"
	"github.com/ignite/cli/v29/ignite/templates/typed/singleton"
)

const maxLength = 64

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

	index string

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

// MapType makes the type stored in a key-value convention in the storage with an index option.
func MapType(index string) AddTypeKind {
	return func(o *addTypeOptions) {
		o.isMap = true
		o.index = index
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
	return func(*addTypeOptions) {}
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
	typeName string,
	kind AddTypeKind,
	options ...AddTypeOption,
) error {
	// apply options.
	o := newAddTypeOptions(s.modpath.Package)
	for _, apply := range append(options, AddTypeOption(kind)) {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(o.moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName := mfName.LowerCase

	name, err := multiformatname.NewName(typeName)
	if err != nil {
		return err
	}

	if err := checkComponentValidity(s.appPath, moduleName, name, o.withoutMessage); err != nil {
		return err
	}

	// Check and parse provided fields
	if err := checkCustomTypes(ctx, s.appPath, s.modpath.Package, s.protoDir, moduleName, o.fields); err != nil {
		return err
	}
	tFields, err := parseTypeFields(o)
	if err != nil {
		return err
	}

	mfSigner, err := multiformatname.NewName(o.signer)
	if err != nil {
		return err
	}

	isIBC, err := isIBCModule(s.appPath, moduleName)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:      s.modpath.Package,
			ProtoDir:     s.protoDir,
			ProtoVer:     "v1", // TODO(@julienrbrt): possibly in the future add flag to specify custom proto version.
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

	// create the type generator depending on the model
	switch {
	case o.isList:
		g, err = list.NewGenerator(s.Tracer(), opts)
	case o.isMap:
		g, err = mapGenerator(s.Tracer(), opts, o.index)
	case o.isSingleton:
		g, err = singleton.NewGenerator(s.Tracer(), opts)
	default:
		g, err = dry.NewGenerator(opts)
	}
	if err != nil {
		return err
	}

	// run the generation
	return s.Run(append(gens, g)...)
}

// checkMaxLength checks if the index length exceeds the maximum allowed length.
func checkMaxLength(name string) error {
	if len(name) > maxLength {
		return errors.Errorf("index exceeds maximum allowed length of %d characters", maxLength)
	}
	return nil
}

// checkForbiddenTypeIndex returns true if the name is forbidden as an index name.
func checkForbiddenTypeIndex(index string) error {
	indexSplit := strings.Split(index, datatype.Separator)
	if len(indexSplit) > 1 {
		index = indexSplit[0]
		indexType := datatype.Name(indexSplit[1])
		if f, ok := datatype.IsSupportedType(indexType); !ok || f.NonIndex {
			return errors.Errorf("invalid index type %s", indexType)
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
		return errors.Errorf("%s is used by type scaffolder", field)
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
func mapGenerator(replacer placeholder.Replacer, opts *typed.Options, index string) (*genny.Generator, error) {
	// Parse indexes with the associated type
	if strings.Contains(index, ",") {
		return nil, errors.Errorf("multi-index map isn't supported")
	}

	parsedIndexes, err := field.ParseFields([]string{index}, checkForbiddenTypeIndex)
	if err != nil {
		return nil, err
	}

	if len(parsedIndexes) == 0 {
		return nil, errors.Errorf("no index found, a valid map index must be provided")
	}

	// Indexes and type fields must be disjoint
	exists := make(map[string]struct{})
	for _, name := range opts.Fields {
		exists[name.Name.LowerCamel] = struct{}{}
	}

	if _, ok := exists[parsedIndexes[0].Name.LowerCamel]; ok {
		return nil, errors.Errorf("%s cannot simultaneously be an index and a field", parsedIndexes[0].Name.Original)
	}

	opts.Index = parsedIndexes[0]
	return maptype.NewGenerator(replacer, opts)
}
