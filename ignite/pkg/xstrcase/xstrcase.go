package xstrcase

import (
	"strings"

	protogenerator "github.com/cosmos/gogoproto/protoc-gen-gogo/generator"
	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
)

// UpperCamel returns the name with upper camel and no special character.
func UpperCamel(name string) string {
	return protogenerator.CamelCase(strcase.ToSnake(name))
}

// Lowercase returns the name with lower case and no special character.
func Lowercase(name string) string {
	return strings.ToLower(
		strings.ReplaceAll(
			xstrings.NoDash(name),
			"_",
			"",
		),
	)
}

// Uppercase returns the name with upper case and no special character.
func Uppercase(name string) string {
	return strings.ToUpper(
		strings.ReplaceAll(
			xstrings.NoDash(name),
			"_",
			"",
		),
	)
}
