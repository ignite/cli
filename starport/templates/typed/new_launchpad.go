package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

type typedLaunchpad struct {
}

// New ...
func NewLaunchpad(opts *Options) (*genny.Generator, error) {
	t := typedLaunchpad{}
	g := genny.New()
	g.RunFn(t.handlerModify(opts))
	g.RunFn(t.typesKeyModify(opts))
	g.RunFn(t.typesCodecModify(opts))
	g.RunFn(t.clientCliTxModify(opts))
	g.RunFn(t.clientCliQueryModify(opts))
	g.RunFn(t.typesQuerierModify(opts))
	g.RunFn(t.keeperQuerierModify(opts))
	g.RunFn(t.clientRestRestModify(opts))
	g.RunFn(frontendSrcStoreAppModify(opts))
	return g, box(cosmosver.Launchpad, opts, g)
}

func (t *typedLaunchpad) handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		case types.MsgCreate%[2]v:
			return handleMsgCreate%[2]v(ctx, k, msg)
		case types.MsgSet%[2]v:
			return handleMsgSet%[2]v(ctx, k, msg)
		case types.MsgDelete%[2]v:
			return handleMsgDelete%[2]v(ctx, k, msg)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) typesKeyModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/key.go", opts.ModuleName)
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

func (t *typedLaunchpad) typesCodecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		cdc.RegisterConcrete(MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
		cdc.RegisterConcrete(MsgSet%[2]v{}, "%[3]v/Set%[2]v", nil)
		cdc.RegisterConcrete(MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName), opts.ModuleName)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) clientCliTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		GetCmdCreate%[2]v(cdc),
		GetCmdSet%[2]v(cdc),
		GetCmdDelete%[2]v(cdc),`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) clientCliQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
			GetCmdList%[2]v(queryRoute, cdc),
			GetCmdGet%[2]v(queryRoute, cdc),`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) typesQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/querier.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `
		const QueryList%[2]v = "list-%[1]v"
		const QueryGet%[2]v = "get-%[1]v"
		`
		content := f.String() + fmt.Sprintf(template, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) keeperQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/keeper/querier.go", opts.ModuleName)
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
			return list%[2]v(ctx, k)
		case types.QueryGet%[2]v:
			return get%[2]v(ctx, path[1:], k)`
		replacement := fmt.Sprintf(template, opts.ModulePath, opts.ModuleName)
		replacement2 := fmt.Sprintf(template2, placeholder, opts.ModulePath, opts.ModuleName)
		replacement3 := fmt.Sprintf(template3, placeholder2, strings.Title(opts.TypeName))
		content := f.String()
		content = strings.Replace(content, replacement, "", 1)
		content = strings.Replace(content, placeholder, replacement2, 1)
		content = strings.Replace(content, placeholder2, replacement3, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedLaunchpad) clientRestRestModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/rest/rest.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
		r.HandleFunc("/%[2]v/%[4]v", create%[3]vHandler(cliCtx)).Methods("POST")
		r.HandleFunc("/%[2]v/%[4]v", list%[3]vHandler(cliCtx, "%[2]v")).Methods("GET")
		r.HandleFunc("/%[2]v/%[4]v/{key}", get%[3]vHandler(cliCtx, "%[2]v")).Methods("GET")
		r.HandleFunc("/%[2]v/%[4]v", set%[3]vHandler(cliCtx)).Methods("PUT")
		r.HandleFunc("/%[2]v/%[4]v", delete%[3]vHandler(cliCtx)).Methods("DELETE")

		`
		replacement := fmt.Sprintf(template, placeholder, opts.ModuleName, strings.Title(opts.TypeName), opts.TypeName)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
