package typed

import (
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/testutil"
)

func Box(box packd.Walker, opts *Options, g *genny.Generator) error {
	if err := g.Box(box); err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("Indexes", opts.Indexes)
	ctx.Set("NoMessage", opts.NoMessage)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName))
	ctx.Set("strconv", func() bool {
		strconv := false
		for _, field := range opts.Fields {
			if field.DatatypeName != "string" {
				strconv = true
			}
		}
		return strconv
	})

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName.Snake))

	// Create the 'testutil' package with the test helpers
	return testutil.Register(g, opts.AppPath)
}
