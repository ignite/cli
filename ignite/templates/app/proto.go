package app

import (
	"embed"
	"fmt"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
=======
	"github.com/ignite/cli/v29/ignite/pkg/xembed"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
>>>>>>> 6364ecbf (feat: support custom proto path (#4071))
)

//go:embed files/{{protoDir}}/* files/buf.work.yaml.plush
var fsProto embed.FS

// NewBufGenerator returns the generator to buf build files.
func NewBufGenerator(appPath, protoDir string) (*genny.Generator, error) {
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
	ctx.Set("ProtoDir", protoDir)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", protoDir))

	return g, nil
}

func CutTemplatePrefix(name string) (string, bool) {
	return strings.CutPrefix(name, fmt.Sprintf("%s/", "{{protoDir}}"))
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
