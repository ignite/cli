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
	g.RunFn(typesKeyModify(opts))
	g.RunFn((typesCodecModify(opts)))
	g.RunFn((clientCliTxModify(opts)))
	g.RunFn((clientCliQueryModify(opts)))
	g.RunFn((typesQuerierModify(opts)))
	g.RunFn((keeperQuerierModify(opts)))
	g.RunFn((clientRestRestModify(opts)))
	g.RunFn((frontendSrcStoreAppModify(opts)))
	if err := g.Box(packr.New("typed/templates", "./templates")); err != nil {
		return g, err
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
			if field.Datatype != "string" {
				strconv = true
			}
		}
		return strconv
	})
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{typeName}}", opts.TypeName))
	g.Transformer(genny.Replace("{{TypeName}}", strings.Title(opts.TypeName)))
	return g, nil
}

const placeholder = "// this line is used by starport scaffolding"
const placeholder2 = "// this line is used by starport scaffolding # 2"
const placeholder4 = "<!-- this line is used by starport scaffolding # 4 -->"

func handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		case types.MsgCreate%[2]v:
			return handleMsgCreate%[2]v(ctx, k, msg)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
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

func frontendSrcStoreAppModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "frontend/src/store/app.js"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		fields := ""
		for _, field := range opts.Fields {
			fields += fmt.Sprintf(`"%[1]v", `, field.Name)
		}
		replacement := fmt.Sprintf(`%[1]v
		{ type: "%[2]v", fields: [%[3]v] },`, placeholder, opts.TypeName, fields)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
