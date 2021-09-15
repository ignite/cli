// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
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
	isForbiddenField func(string) error,
) (Fields, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var parsedFields Fields
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

		// Check if is a static type
		if datatype, ok := StaticDataTypes[datatypeName]; ok {
			parsedFields = append(parsedFields, Field{
				Name:         name,
				Datatype:     datatype,
				DatatypeName: datatypeName,
			})
			continue
		}

		parsedFields = append(parsedFields, Field{
			Name:         name,
			Datatype:     datatypeName,
			DatatypeName: TypeCustom,
		})
	}
	return parsedFields, nil
}
