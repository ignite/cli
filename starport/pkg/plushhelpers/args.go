package plushhelpers

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/plush"
)

// CastArgs to type using github.com/spf13/cast.
// Don't forget to import github.com/spf13/cast in templates.
func CastArgs(actual string, i int) string {
	return fmt.Sprintf("cast.To%sE(args[%d])", strings.Title(actual), i)
}

// GenerateValidArg will produce a valid value for the specified type.
// This function doesn't guarantee to produce unique values.
// Note that return value needs to be wrapped into a string.
func GenerateValidArg(typ string) string {
	switch typ {
	case "string":
		return "xyz"
	case "uint":
		return "111"
	case "int":
		return "111"
	default:
		panic(fmt.Sprintf("unknown type %s", typ))
	}
}

// ExtendPlushContext sets available helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("castArgs", CastArgs)
	ctx.Set("genValidArg", GenerateValidArg)
}
