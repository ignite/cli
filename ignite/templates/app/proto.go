package app

import (
	"embed"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/xgenny"
)

//go:embed files/proto/* files/buf.work.yaml
var fsProto embed.FS

// NewBufGenerator returns the generator to buf build files.
func NewBufGenerator(appPath string) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsProto,
			"files",
			appPath,
		)
	)
	return g, xgenny.Box(g, template)
}
