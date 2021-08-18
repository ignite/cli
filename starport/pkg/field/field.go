// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
	TypeUint64 = "uint64"
)

// Field represents a field inside a structure for a component
// it can a field contained in a type or inside the response of a query, etc...
type Field struct {
	Name         multiformatname.Name
	Datatype     string
	DatatypeName string
}

// ParseFields parses the provided fields, analyses the types and checks there is no duplicated field
func ParseFields(fields []string, isForbiddenField func(string) error) ([]Field, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var parsedFields []Field
	for _, field := range fields {
		fieldSplit := strings.Split(field, ":")
		if len(fieldSplit) > 2 {
			return parsedFields, fmt.Errorf("invalid field format: %s, should be 'name' or 'name:type'", field)
		}

		name, err := multiformatname.NewName(fieldSplit[0])
		if err != nil {
			return parsedFields, err
		}

		// Ensure the field name is not a Go reserved name, it would generate an incorrect code
		if err := isForbiddenField(name.LowerCamel); err != nil {
			return parsedFields, fmt.Errorf("%s can't be used as a field name: %s", name, err.Error())
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name.LowerCamel]; exists {
			return parsedFields, fmt.Errorf("the field %s is duplicated", name.Original)
		}
		existingFields[name.LowerCamel] = true

		// Parse the type if it is provided, otherwise string is used by defaut
		datatypeName, datatype := TypeString, TypeString
		isTypeSpecified := len(fieldSplit) == 2
		if isTypeSpecified {
			acceptedTypes := map[string]string{
				"string": TypeString,
				"bool":   TypeBool,
				"int":    TypeInt32,
				"uint":   TypeUint64,
			}

			if t, ok := acceptedTypes[fieldSplit[1]]; ok {
				datatype = t
				datatypeName = fieldSplit[1]
			} else {
				return parsedFields, fmt.Errorf("the field type %s doesn't exist", fieldSplit[1])
			}
		}

		parsedFields = append(parsedFields, Field{
			Name:         name,
			Datatype:     datatype,
			DatatypeName: datatypeName,
		})
	}

	return parsedFields, nil
}
