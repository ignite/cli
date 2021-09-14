package modulecreate

import (
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

// genesisTestCtx returns the generator to generate genesis_test.go
func genesisTestCtx(appName, modulePath, moduleName string) *genny.Generator {
	g := genny.New()
	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("appName", appName)
	ctx.Set("title", strings.Title)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))
	return g
}

// AddGenesisModuleTest returns the generator to generate genesis_test.go
func AddGenesisModuleTest(appPath, appName, modulePath, moduleName string) (*genny.Generator, error) {
	g := genesisTestCtx(appName, modulePath, moduleName)
	return g, g.Box(xgenny.NewEmbedWalker(
		fsGenesisModuleTest,
		"genesistest/module/",
		appPath,
		true,
	))
}

// AddGenesisTypesTest returns the generator to generate types/genesis_test.go
func AddGenesisTypesTest(appPath, appName, modulePath, moduleName string) (*genny.Generator, error) {
	g := genesisTestCtx(appName, modulePath, moduleName)
	return g, g.Box(xgenny.NewEmbedWalker(
		fsGenesisTypesTest,
		"genesistest/types/",
		appPath,
		true,
	))
}
