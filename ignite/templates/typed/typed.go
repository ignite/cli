package typed

import (
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

func Box(box fs.FS, opts *Options, g *genny.Generator) error {
	if err := g.OnlyFS(box, nil, nil); err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("ProtoVer", opts.ProtoVer)
	ctx.Set("IsIBC", opts.IsIBC)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("Index", opts.Index)
	ctx.Set("NoMessage", opts.NoMessage)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName, opts.ProtoVer))
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
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName.Snake))

	return nil
}
