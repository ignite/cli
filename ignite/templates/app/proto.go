package app

import (
	"embed"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/xembed"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
)

//go:embed files/proto/* files/buf.work.yaml.plush
var fsProto embed.FS

// NewBufGenerator returns the generator to buf build files.
func NewBufGenerator(appPath, protoPath string) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsProto,
			"files",
			appPath,
		)
	)
	if err := xgenny.Box(g, template); err != nil {
		return nil, err
	}

	ctx := plush.NewContext()
	ctx.Set("ProtoPath", protoPath)
	g.Transformer(xgenny.Transformer(ctx))

	return g, nil
}

// BufFiles returns a list of Buf.Build files.
func BufFiles() ([]string, error) {
	files, err := xembed.FileList(fsProto, "files")
	if err != nil {
		return nil, err
	}
	// remove all .plush extensions.
	for i, file := range files {
		files[i] = strings.TrimSuffix(file, ".plush")
	}
	return files, nil
}
