package scaffolder

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	"github.com/tendermint/starport/starport/templates/typed"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
)

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
	TypeUint64 = "uint64"
)

type AddTypeOption struct {
	Indexed   bool
	NoMessage bool
}

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddType(
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
	tFields, err := parseFields(fields, checkForbiddenTypeField)
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
	)
	// Check and support MsgServer convention
	if err := supportMsgServer(
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			OwnerName:  opts.OwnerName,
		},
	); err != nil {
		return err
	}

	// Check if indexed type
	if addTypeOptions.Indexed {
		g, err = indexed.NewStargate(opts)
	} else {
		// Scaffolding a type with ID
		g, err = typed.NewStargate(opts)
	}
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
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

// parseFields parses the provided fields, analyses the types and checks there is no duplicated field
func parseFields(fields []string, isForbiddenField func(string) error) ([]typed.Field, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var tFields []typed.Field
	for _, f := range fields {
		fs := strings.Split(f, ":")
		name := fs[0]

		// Ensure the field name is not a Go reserved name, it would generate an incorrect code
		if err := isForbiddenField(name); err != nil {
			return tFields, fmt.Errorf("%s can't be used as a field name: %s", name, err.Error())
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name]; exists {
			return tFields, fmt.Errorf("the field %s is duplicated", name)
		}
		existingFields[name] = true

		datatypeName, datatype := TypeString, TypeString
		acceptedTypes := map[string]string{
			"string": TypeString,
			"bool":   TypeBool,
			"int":    TypeInt32,
			"uint":   TypeUint64,
		}
		isTypeSpecified := len(fs) == 2
		if isTypeSpecified {
			if t, ok := acceptedTypes[fs[1]]; ok {
				datatype = t
				datatypeName = fs[1]
			} else {
				return tFields, fmt.Errorf("the field type %s doesn't exist", fs[1])
			}
		}
		tFields = append(tFields, typed.Field{
			Name:         name,
			Datatype:     datatype,
			DatatypeName: datatypeName,
		})
	}

	return tFields, nil
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
