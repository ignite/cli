package modulecreate

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

// AddSimulation returns the generator to generate module_simulation.go file
func AddSimulation(appPath, modulePath, moduleName string, params ...field.Field) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(fsSimapp, "simapp/", appPath)
	)

	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("params", params)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))

	if err := xgenny.Box(g, template); err != nil {
		return nil, err
	}

	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}
