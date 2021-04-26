package typed

import (
	"embed"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/pkg/xstrings"
)

var (
	//go:embed stargate/component/* stargate/component/**/*
	fsStargateComponent embed.FS

	//go:embed stargate/messages/* stargate/messages/**/*
	fsStargateMessages embed.FS

	// stargateComponentTemplate is the template for a Stargate module type component
	stargateComponentTemplate = xgenny.NewEmbedWalker(fsStargateComponent, "stargate/component/")

	// stargateMessagesTemplate is the template for a Stargate module type interaction messages
	stargateMessagesTemplate = xgenny.NewEmbedWalker(fsStargateMessages, "stargate/messages/")
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
	ctx.Set("Fields", opts.Fields)
	ctx.Set("title", strings.Title)
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

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName))
	g.Transformer(genny.Replace("{{TypeName}}", strings.Title(opts.TypeName)))
	return nil
}
