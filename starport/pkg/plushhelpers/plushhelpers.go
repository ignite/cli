package plushhelpers

import "github.com/gobuffalo/plush"

// ExtendPlushContext sets available helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("castArgs", CastArgs)
	ctx.Set("genValidArg", GenerateValidArg)
	ctx.Set("castToBytes", CastToBytes)
}
