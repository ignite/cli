// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

// Fields represents a Field slice
type Fields []Field

// GoCLIImports return all go CLI imports
func (f Fields) GoCLIImports() []string {
	allImports := make([]string, 0)
	for _, field := range f {
		allImports = append(allImports, field.GoCLIImports()...)
	}
	return allImports
}

// ProtoImports return all proto imports
func (f Fields) ProtoImports() []string {
	allImports := make([]string, 0)
	for _, field := range f {
		allImports = append(allImports, field.ProtoImports()...)
	}
	return allImports
}

// String return all inline fields args for command usage
func (f Fields) String() string {
	args := ""
	for _, field := range f {
		args += fmt.Sprintf(" [%s]", field.Name.Kebab)
	}
	return args
}

// Custom return a list of custom fields
func (f Fields) Custom() []string {
	fields := make([]string, 0)
	for _, field := range f {
		if field.DatatypeName == TypeCustom {
			dataType, err := multiformatname.NewName(field.Datatype)
			if err != nil {
				panic(err)
			}
			fields = append(fields, dataType.Snake)
		}
	}
	return fields
}
