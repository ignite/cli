package field

import (
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

// validateField validates the field Name and type, and checks the name is not forbidden by Ignite CLI.
func validateField(field string, isForbiddenField func(string) error) (multiformatname.Name, datatype.Name, error) {
	name, dataTypeName, err := parseField(field)
	if err != nil {
		return name, "", err
	}

	// Ensure the field Name is not a Go reserved Name, it would generate an incorrect code
	if err := isForbiddenField(name.LowerCamel); err != nil {
		return name, "", errors.Errorf("%s can't be used as a field Name: %w", name, err)
	}

	return name, dataTypeName, nil
}

// parseField parses the field string and returns the multiformat name and datatype name.
func parseField(field string) (multiformatname.Name, datatype.Name, error) {
	fieldSplit := strings.Split(field, datatype.Separator)
	if len(fieldSplit) > 2 {
		return multiformatname.Name{}, "", errors.Errorf("invalid field format: %s, should be 'Name' or 'Name:type'", field)
	}

	name, err := multiformatname.NewName(fieldSplit[0])
	if err != nil {
		return name, "", err
	}

	// Check if the object has an explicit type. The default is a string
	dataTypeName := datatype.String
	isTypeSpecified := len(fieldSplit) == 2
	if isTypeSpecified {
		dataTypeName = datatype.Name(fieldSplit[1])
	}
	return name, dataTypeName, nil
}

// MultipleCoins checks if the provided fields contain more than one coin type.
func MultipleCoins(fields []string) (bool, error) {
	coinsCount := 0
	for _, field := range fields {
		_, datatypeName, err := parseField(field)
		if err != nil {
			return false, err
		}
		if datatypeName == datatype.Coins || datatypeName == datatype.DecCoins ||
			datatypeName == datatype.CoinSliceAlias || datatypeName == datatype.DecCoinSliceAlias {
			coinsCount++
		}
	}
	return coinsCount > 1, nil
}

// ParseFields parses the provided fields, analyses the types
// and checks there is no duplicated field.
func ParseFields(
	fields []string,
	isForbiddenField func(string) error,
	forbiddenFieldNames ...string,
) (Fields, error) {
	// Used to check duplicated field
	existingFields := make(map[string]struct{})
	for _, name := range forbiddenFieldNames {
		if name != "" {
			existingFields[name] = struct{}{}
		}
	}

	var parsedFields Fields
	for _, field := range fields {
		name, datatypeName, err := validateField(field, isForbiddenField)
		if err != nil {
			return parsedFields, err
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name.LowerCamel]; exists {
			return parsedFields, errors.Errorf("the field %s is duplicated", name.Original)
		}
		existingFields[name.LowerCamel] = struct{}{}

		// Check if is a static type
		if _, ok := datatype.IsSupportedType(datatypeName); ok {
			parsedFields = append(parsedFields, Field{
				Name:         name,
				DatatypeName: datatypeName,
			})
			continue
		}

		parsedFields = append(parsedFields, Field{
			Name:         name,
			Datatype:     string(datatypeName),
			DatatypeName: datatype.TypeCustom,
		})
	}
	return parsedFields, nil
}
