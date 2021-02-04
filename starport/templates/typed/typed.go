package typed

import (
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var (
	launchpadTemplate      = packr.New("typed/templates/launchpad", "./launchpad")
	stargateTemplate       = packr.New("typed/templates/stargate", "./stargate")
	stargateLegacyTemplate = packr.New("typed/templates/stargate", "./stargate_legacy")
)

func box(box *packr.Box, opts *Options, g *genny.Generator) error {
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
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName))
	g.Transformer(genny.Replace("{{TypeName}}", strings.Title(opts.TypeName)))
	return nil
}
