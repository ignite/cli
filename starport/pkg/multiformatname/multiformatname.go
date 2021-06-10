// Package multiformatname provides names automatically converted into multiple naming convention
package multiformatname

import (
	"errors"
	"fmt"

	"github.com/iancoleman/strcase"
)

// MultiFormatName represents a name with multiple naming convention representations
// Supported naming convention are: camel, pascal, and kebab cases
type MultiFormatName struct {
	Original   string
	LowerCamel string
	UpperCamel string
	Kebab      string
}

// NewMultiFormatName returns a new multi-format name from a name
func NewMultiFormatName(name string) (MultiFormatName, error) {
	if err := CheckName(name); err != nil {
		return MultiFormatName{}, err
	}

	return MultiFormatName{
		Original: 	name,
		LowerCamel: strcase.ToLowerCamel(name),
		UpperCamel: strcase.ToCamel(name),
		Kebab:      strcase.ToKebab(name),
	}, nil
}

// CheckName checks name validity
func CheckName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	// check  characters
	for _, c := range name {
		// A name can contains letter, hyphen or underscore
		authorized := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '-' || c == '_'
		if !authorized {
			return fmt.Errorf("name cannot contain %v", string(c))
		}
	}

	return nil
}
