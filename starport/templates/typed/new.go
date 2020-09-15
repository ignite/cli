package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

const placeholder = "// this line is used by starport scaffolding"
const placeholder2 = "// this line is used by starport scaffolding # 2"
const placeholder3 = "// this line is used by starport scaffolding # 3"
const placeholder4 = "<!-- this line is used by starport scaffolding # 4 -->"

// New ...
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(handlerModify(opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn(typesCodecModify(opts))
	g.RunFn(typesCodecImportModify(opts))
	g.RunFn(typesCodecInterfaceModify(opts))
	g.RunFn(protoRPCImportModify(opts))
	g.RunFn(protoRPCModify(opts))
	g.RunFn(protoRPCMessageModify(opts))
	g.RunFn(clientCliTxModify(opts))
	g.RunFn(clientCliQueryModify(opts))
	g.RunFn((typesQuerierModify(opts)))
	g.RunFn(keeperQuerierModify(opts))
	g.RunFn(clientRestRestModify(opts))
	g.RunFn(frontendSrcStoreAppModify(opts))
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

func protoRPCImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
import "%s/v1beta/%s.proto";`
		replacement := fmt.Sprintf(template, placeholder,
			opts.AppName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoRPCModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
	rpc All%s(QueryAll%sRequest) returns (QueryAll%sResponse);`
		replacement := fmt.Sprintf(template, placeholder2,
			strings.Title(opts.TypeName),
			strings.Title(opts.TypeName),
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoRPCMessageModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
message QueryAll%sRequest {
	cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAll%sResponse {
	repeated Msg%s %s = 1;
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}`
		replacement := fmt.Sprintf(template, placeholder3,
			strings.Title(opts.TypeName),
			strings.Title(opts.TypeName),
			strings.Title(opts.TypeName),
			opts.TypeName,
		)
		content := strings.Replace(f.String(), placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
case *types.Msg%[2]v:
return handleMsgCreate%[2]v(ctx, k, msg)`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesKeyModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/keys.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%sKey= "%s"
)
`, strings.Title(opts.TypeName), strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
sdk "github.com/cosmos/cosmos-sdk/types"`
		replacement := fmt.Sprintf(template, placeholder)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
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
cdc.RegisterConcrete(Msg%[2]v{}, "%[3]v/Create%[2]v", nil)`
		replacement := fmt.Sprintf(template, placeholder2, strings.Title(opts.TypeName), opts.AppName)
		content := strings.Replace(f.String(), placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecInterfaceModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]v{},
)`
		replacement := fmt.Sprintf(template, placeholder3, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder3, replacement, 1)
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
		template := `%s

	cmd.AddCommand(CmdCreate%s())`
		replacement := fmt.Sprintf(template, placeholder,
			strings.Title(opts.TypeName),
		)
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
		template := `%s

	cmd.AddCommand(CmdList%s())`
		replacement := fmt.Sprintf(template, placeholder,
			strings.Title(opts.TypeName),
		)
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
		return list%[2]v(ctx, k, legacyQuerierCdc)`
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
		template := `%s
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
`
		replacement := fmt.Sprintf(template, placeholder)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func frontendSrcStoreAppModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "vue/src/store/app.js"
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
