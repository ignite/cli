package plushhelpers

import "github.com/gobuffalo/plush"

const (
	datatypeString = "string"
	datatypeUint   = "uint"
	datatypeInt    = "int"
	datatypeBool   = "bool"
	valueFalse     = "false"
)

// ExtendPlushContext sets available helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("castArg", castArg)
	ctx.Set("castToBytes", CastToBytes)
	ctx.Set("castToString", CastToString)
	ctx.Set("genValidArg", GenerateValidArg)
	ctx.Set("genUniqueArg", GenerateUniqueArg)
	ctx.Set("genValidIndex", GenerateValidIndex)
	ctx.Set("genNotFoundIndex", GenerateNotFoundIndex)
}
