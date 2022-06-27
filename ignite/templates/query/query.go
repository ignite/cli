package query

import (
	"embed"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"

	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS
)

func Box(box packd.Walker, opts *Options, g *genny.Generator) error {
	if err := g.Box(box); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("QueryName", opts.QueryName)
	ctx.Set("Description", opts.Description)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("ReqFields", opts.ReqFields)
	ctx.Set("ResFields", opts.ResFields)
	ctx.Set("Paginated", opts.Paginated)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{queryName}}", opts.QueryName.Snake))
	return nil
}
