package modulecreate

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

// genesisTestCtx returns the generator to generate genesis_test.go
func genesisTestCtx(appPath, appName, modulePath, moduleName string) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsGenesisTest,
			"genesistest/",
			appPath,
		)
	)
	if err := g.Box(template); err != nil {
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

// AddGenesisModuleTest returns the generator to generate genesis_test.go
func AddGenesisModuleTest(appPath, appName, modulePath, moduleName string) (*genny.Generator, error) {
	g, err := genesisTestCtx(appPath, appName, modulePath, moduleName)
	if err != nil {
		return g, err
	}
	return g, g.Box(genesisModuleTestTemplate)
}

// AddGenesisTypesTest returns the generator to generate types/genesis_test.go
func AddGenesisTypesTest(appPath, appName, modulePath, moduleName string) (*genny.Generator, error) {
	g, err := genesisTestCtx(appPath, appName, modulePath, moduleName)
	if err != nil {
		return g, err
	}
	return g, g.Box(genesisTypesTestTemplate)
}
