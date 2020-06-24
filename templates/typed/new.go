package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// New ...
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(handlerModify(opts))
	g.RunFn(aliasModify(opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn((typesCodecModify(opts)))
	g.RunFn((clientCliTxModify(opts)))
	g.RunFn((clientCliQueryModify(opts)))
	g.RunFn((typesQuerierModify(opts)))
	g.RunFn((keeperQuerierModify(opts)))
	g.RunFn((clientRestRestModify(opts)))
	g.RunFn((uiIndexModify(opts)))
	g.RunFn((uiScriptModify(opts)))
	if err := g.Box(packr.New("typed/templates", "./templates")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("title", func(s string) string {
		return strings.Title(s)
	})
	ctx.Set("strconv", func() bool {
		strconv := false
		for _, field := range opts.Fields {
			if field.Datatype != "string" {
				strconv = true
			}
		}
		return strconv
	})
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", fmt.Sprintf("%s", opts.AppName)))
	g.Transformer(genny.Replace("{{typeName}}", fmt.Sprintf("%s", opts.TypeName)))
	g.Transformer(genny.Replace("{{TypeName}}", fmt.Sprintf("%s", strings.Title(opts.TypeName))))
	return g, nil
}

const placeholder = "// this line is used by startport scaffolding"
const placeholder2 = "// this line is used by startport scaffolding # 2"
const placeholder3 = "// this line is used by startport scaffolding # 3"
const placeholder4 = "<!-- this line is used by startport scaffolding # 4 -->"

func handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		case MsgCreate%[2]v:
			return handleMsgCreate%[2]v(ctx, k, msg)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func aliasModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/alias.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
var (
	NewMsgCreate%[1]v = types.NewMsgCreate%[1]v
)

type (
	MsgCreate%[1]v = types.MsgCreate%[1]v
)
		`, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesKeyModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/key.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%[2]vPrefix = "%[1]v-"
)
		`, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		cdc.RegisterConcrete(MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName), opts.AppName)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		GetCmdCreate%[2]v(cdc),`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
			GetCmdList%[2]v(queryRoute, cdc),`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/querier.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `
const (QueryList%[2]v = "list-%[1]v")
		`
		content := f.String() + fmt.Sprintf(template, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keeperQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/keeper/querier.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `"%[1]v/x/%[2]v/types"`
		template2 := `%[1]v
	"%[2]v/x/%[3]v/types"
		`
		template3 := `%[1]v
		case types.QueryList%[2]v:
			return list%[2]v(ctx, k)`
		replacement := fmt.Sprintf(template, opts.ModulePath, opts.AppName)
		replacement2 := fmt.Sprintf(template2, placeholder, opts.ModulePath, opts.AppName)
		replacement3 := fmt.Sprintf(template3, placeholder2, strings.Title(opts.TypeName))
		content := f.String()
		content = strings.Replace(content, replacement, "", 1)
		content = strings.Replace(content, placeholder, replacement2, 1)
		content = strings.Replace(content, placeholder2, replacement3, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientRestRestModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/rest/rest.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	r.HandleFunc("/%[2]v/%[4]v", list%[3]vHandler(cliCtx, "%[2]v")).Methods("GET")
	r.HandleFunc("/%[2]v/%[4]v", create%[3]vHandler(cliCtx)).Methods("POST")`
		replacement := fmt.Sprintf(template, placeholder, opts.AppName, strings.Title(opts.TypeName), opts.TypeName)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func uiIndexModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "ui/index.html"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := fmt.Sprintf(`
			<h2>List of "%[1]v" items</h2>
			<div class="type-%[1]v-list-%[1]v"></div>
      <h3>Create a new %[1]v:</h3>`, opts.TypeName)
		for _, field := range opts.Fields {
			template = template + fmt.Sprintf(`
			<input placeholder="%[1]v" class="type-%[2]v-field-%[1]v" type="text" />`, field.Name, opts.TypeName)
		}
		template = template + fmt.Sprintf(`
			<button class="type-%[1]v-create">Create %[1]v</button>
		`, opts.TypeName) + "  " + placeholder4
		content := strings.Replace(f.String(), placeholder4, template, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func uiScriptModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "ui/script.js"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		fields := ""
		for _, field := range opts.Fields {
			fields = fields + fmt.Sprintf("\"%[1]v\", ", field.Name)
		}
		template := `%[1]v
	["%[2]v", [%[3]v]],`
		replacement := fmt.Sprintf(template, placeholder, opts.TypeName, fields)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
