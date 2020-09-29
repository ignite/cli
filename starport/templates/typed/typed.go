package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

const placeholder = "// this line is used by starport scaffolding # 1"
const placeholder2 = "// this line is used by starport scaffolding # 2"
const placeholder3 = "// this line is used by starport scaffolding # 3"
const placeholder4 = "<!-- this line is used by starport scaffolding # 4 -->"
const placeholder44 = "// this line is used by starport scaffolding # 4"

func box(sdkVersion string, opts *Options, g *genny.Generator) error {
	if err := g.Box(packr.New("typed/templates", "./"+sdkVersion)); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("title", strings.Title)
	ctx.Set("strconv", func() bool {
		strconv := false
		for _, field := range opts.Fields {
			if field.DatatypeName != "string" {
				strconv = true
			}
		}
		return strconv
	})
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName))
	g.Transformer(genny.Replace("{{TypeName}}", strings.Title(opts.TypeName)))
	return nil
}

func frontendSrcStoreAppModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "vue/src/views/Index.vue"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		fields := ""
		for _, field := range opts.Fields {
			fields += fmt.Sprintf(`'%[1]v', `, field.Name)
		}
		replacement := fmt.Sprintf(`%[1]v
		<sp-type-form type="%[2]v" :fields="[%[3]v]" />`, placeholder4, opts.TypeName, fields)
		content := strings.Replace(f.String(), placeholder4, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
