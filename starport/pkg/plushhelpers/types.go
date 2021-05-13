package plushhelpers

import (
	"fmt"

	"github.com/gobuffalo/plush"
)

func castTo(typ string) string {
	switch typ {
	case "string":
		return "cast.ToStringE"
	case "uint":
		return "cast.ToUintE"
	case "int":
		return "cast.ToIntE"
	default:
		panic(fmt.Sprintf("unknown type %s", typ))
	}
}

// CastArgs to type using github.com/spf13/cast.
// Don't forget to import github.com/spf13/cast in templates.
func CastArgs(typ string, i int) string {
	return fmt.Sprintf("%s(args[%d])", castTo(typ), i)
}

func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("castArgs", CastArgs)
}
