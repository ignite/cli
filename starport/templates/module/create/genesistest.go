package modulecreate

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// AddGenesisTest returns the generator to generate genesis_test.go
func AddGenesisTest(appName, modulePath, moduleName string) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(genesisTestTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("appName", appName)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))
	return g, nil
}
