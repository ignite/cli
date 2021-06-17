// Package multiformatname provides names automatically converted into multiple naming convention
package multiformatname

import (
	"errors"
	"fmt"

	"github.com/iancoleman/strcase"
)

// MultiFormatName represents a name with multiple naming convention representations
// Supported naming convention are: camel, pascal, and kebab cases
type Name struct {
	Original   string
	LowerCamel string
	UpperCamel string
	Kebab      string
	Snake      string
}

// NewMultiFormatName returns a new multi-format name from a name
func NewName(name string) (Name, error) {
	if err := CheckName(name); err != nil {
		return Name{}, err
	}

	return Name{
		Original:   name,
		LowerCamel: strcase.ToLowerCamel(name),
		UpperCamel: strcase.ToCamel(name),
		Kebab:      strcase.ToKebab(name),
		Snake:      strcase.ToSnake(name),
	}, nil
}

// CheckName checks name validity
func CheckName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	// check  characters
	c := name[0]
	authorized := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
	if !authorized {
		return fmt.Errorf("name cannot contain %v as first character", string(c))
	}

	for _, c := range name[1:] {
		// A name can contains letter, hyphen or underscore
		authorized := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-' || c == '_'
		if !authorized {
			return fmt.Errorf("name cannot contain %v", string(c))
		}
	}

	return nil
}
