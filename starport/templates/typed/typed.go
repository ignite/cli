package typed

import (
	"embed"
	"strings"

	templatesutils "github.com/tendermint/starport/starport/templates"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	//go:embed stargate_legacy/* stargate_legacy/**/*
	fsStargateLegacy embed.FS

	//go:embed launchpad/* launchpad/**/*
	fsLaunchpad embed.FS

	stargateTemplate       = xgenny.NewEmbedWalker(fsStargate, "stargate/")
	stargateLegacyTemplate = xgenny.NewEmbedWalker(fsStargateLegacy, "stargate_legacy/")
	launchpadTemplate      = xgenny.NewEmbedWalker(fsLaunchpad, "launchpad/")
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
	ctx.Set("nodash", func(s string) string {
		return strings.ReplaceAll(s, "-", "")
	})

	// Used for proto package name
	ctx.Set("noNumberPrefix", templatesutils.NoNumberPrefix)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName))
	g.Transformer(genny.Replace("{{TypeName}}", strings.Title(opts.TypeName)))
	return nil
}
