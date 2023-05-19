package app

import (
	"embed"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

//go:embed files/proto/* files/buf.work.yaml
var fsProto embed.FS

// NewBufGenerator returns the generator to buf build files.
func NewBufGenerator(appPath string) (*genny.Generator, error) {
	g := genny.New()
	ctx := plush.NewContext()
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	protoTemplate := xgenny.NewEmbedWalker(
		fsProto,
		"files",
		appPath,
	)
	return g, xgenny.Box(g, protoTemplate)
}
