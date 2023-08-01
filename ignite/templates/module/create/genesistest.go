package modulecreate

import (
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

// AddGenesisTest returns the generator to generate genesis_test.go files.
func AddGenesisTest(appPath, appName, modulePath, moduleName string, isIBC bool) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(fsGenesisTest, "files/genesistest/", appPath)
	)

	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("appName", appName)
	ctx.Set("isIBC", isIBC)
	ctx.Set("title", xstrings.Title)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", appName))
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))

	if err := xgenny.Box(g, template); err != nil {
		return nil, err
	}

	return g, nil
}
