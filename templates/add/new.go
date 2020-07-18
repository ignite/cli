package add

import (
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// New ...
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(packr.New("wasm", "./wasm")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	return g, nil
}

const placeholder = "// this line is used by starport scaffolding"
const placeholder2 = "// this line is used by starport scaffolding # 2"
const placeholder4 = "<!-- this line is used by starport scaffolding # 4 -->"
