// Package multiformatname provides names automatically converted into multiple naming convention
package multiformatname

import (
	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xstrcase"
)

// Name represents a name with multiple naming convention representations.
// Supported naming convention are: camel, pascal, and kebab cases.
type Name struct {
	Original   string
	LowerCamel string
	UpperCamel string
	PascalCase string
	LowerCase  string
	UpperCase  string
	Kebab      string
	Snake      string
}

type Checker func(name string) error

// MustNewName returns a new multi-format name from a name.
func MustNewName(name string, additionalChecks ...Checker) Name {
	n, err := NewName(name, additionalChecks...)
	if err != nil {
		panic(err)
	}
	return n
}

// NewName returns a new multi-format name from a name.
func NewName(name string, additionalChecks ...Checker) (Name, error) {
	checks := append([]Checker{basicCheckName}, additionalChecks...)

	for _, check := range checks {
		if err := check(name); err != nil {
			return Name{}, err
		}
	}

	return Name{
		Original:   name,
		LowerCamel: strcase.ToLowerCamel(name),
		UpperCamel: xstrcase.UpperCamel(name),
		PascalCase: strcase.ToCamel(name),
		LowerCase:  xstrcase.Lowercase(name),
		UpperCase:  xstrcase.Uppercase(name),
		Kebab:      strcase.ToKebab(name),
		Snake:      strcase.ToSnake(name),
	}, nil
}

// NoNumber prevents using number in a name.
func NoNumber(name string) error {
	for _, c := range name {
		if '0' <= c && c <= '9' {
			return errors.New("name cannot contain number")
		}
	}

	return nil
}

// basicCheckName performs basic checks common for all names.
func basicCheckName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	// check  characters
	c := name[0]
	authorized := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
	if !authorized {
		return errors.Errorf("name cannot contain %v as first character", string(c))
	}

	for _, c := range name[1:] {
		// A name can contain letter, hyphen or underscore
		authorized := ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-' || c == '_'
		if !authorized {
			return errors.Errorf("name cannot contain %v", string(c))
		}
	}

	return nil
}
