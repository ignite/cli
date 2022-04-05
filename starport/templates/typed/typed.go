package typed

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/ignite-hq/cli/starport/pkg/xstrings"
	"github.com/ignite-hq/cli/starport/templates/field/plushhelpers"
	"github.com/ignite-hq/cli/starport/templates/testutil"
)

func Box(box packd.Walker, opts *Options, g *genny.Generator) error {
	if err := g.Box(box); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("OwnerName", opts.OwnerName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("Indexes", opts.Indexes)
	ctx.Set("NoMessage", opts.NoMessage)
	ctx.Set("strconv", func() bool {
		strconv := false
		for _, field := range opts.Fields {
			if field.DatatypeName != "string" {
				strconv = true
			}
		}
		return strconv
	})

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName.Snake))

	// Create the 'testutil' package with the test helpers
	return testutil.Register(g, opts.AppPath)
}
