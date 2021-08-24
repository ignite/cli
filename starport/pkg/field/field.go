// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	TypeCustom = "custom"
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
	TypeUint64 = "uint64"

	FolderX     = "x"
	FolderTypes = "types"

	TypeSeparator = ":"
)

var (
	staticTypes = map[string]string{
		"string": TypeString,
		"bool":   TypeBool,
		"int":    TypeInt32,
		"uint":   TypeUint64,
	}
)

type (
	// Field represents a field inside a structure for a component
	// it can be a field contained in a type or inside the response of a query, etc...
	Field struct {
		Name         multiformatname.Name
		Datatype     string
		DatatypeName string
		Nested       Fields
	}

	// Fields represents a Field slice
	Fields []Field
)

// validateField validate the field name and type, and run the forbidden method
func validateField(field string, isForbiddenField func(string) error) (multiformatname.Name, string, error) {
	fieldSplit := strings.Split(field, TypeSeparator)
	if len(fieldSplit) > 2 {
		return multiformatname.Name{}, "", fmt.Errorf("invalid field format: %s, should be 'name' or 'name:type'", field)
	}

	name, err := multiformatname.NewName(fieldSplit[0])
	if err != nil {
		return name, "", err

	}

	// Ensure the field name is not a Go reserved name, it would generate an incorrect code
	if err := isForbiddenField(name.LowerCamel); err != nil {
		return name, "", fmt.Errorf("%s can't be used as a field name: %s", name, err.Error())
	}

	// Check if the object has an explicit type. The default is a string
	dataTypeName := TypeString
	isTypeSpecified := len(fieldSplit) == 2
	if isTypeSpecified {
		dataTypeName = fieldSplit[1]
	}
	return name, dataTypeName, nil
}

// ParseFields parses the provided fields, analyses the types
// and checks there is no duplicated field
func ParseFields(
	fields []string,
	module string,
	isForbiddenField func(string) error,
) (parsedFields Fields, err error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	for _, field := range fields {
		name, datatypeName, err := validateField(field, isForbiddenField)
		if err != nil {
			return parsedFields, err
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name.LowerCamel]; exists {
			return parsedFields, fmt.Errorf("the field %s is duplicated", name.Original)
		}
		existingFields[name.LowerCamel] = true

		// Parse the type if it is provided, otherwise string is used by default
		if datatype, ok := staticTypes[datatypeName]; ok {
			parsedFields = append(parsedFields, Field{
				Name:         name,
				Datatype:     datatype,
				DatatypeName: datatypeName,
			})
			continue
		}

		// Check if the custom type is valid and fetch the fields
		path := filepath.Join(FolderX, module, FolderTypes)
		structFields, err := cosmosanalysis.FindStructFields(path, datatypeName)
		if err != nil {
			return parsedFields, err
		}
		nestedFields, err := ParseFields(structFields, module, isForbiddenField)
		if err != nil {
			return parsedFields, err
		}

		parsedFields = append(parsedFields, Field{
			Name:         name,
			Datatype:     datatypeName,
			DatatypeName: TypeCustom,
			Nested:       nestedFields,
		})
	}
	return
}

// NeedCast return true if the field slice needs
// external cast library
func (f Fields) NeedCast() bool {
	for _, field := range f {
		if field.DatatypeName != TypeString {
			return true
		}
	}
	return false
}
