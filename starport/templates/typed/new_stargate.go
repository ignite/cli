package typed

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

type typedStargate struct {
}

// New ...
func NewStargate(opts *Options) (*genny.Generator, error) {
	t := typedStargate{}
	g := genny.New()
	g.RunFn(t.handlerModify(opts))
	g.RunFn(t.typesKeyModify(opts))
	g.RunFn(t.typesCodecModify(opts))
	g.RunFn(t.typesCodecImportModify(opts))
	g.RunFn(t.typesCodecInterfaceModify(opts))
	g.RunFn(t.protoRPCImportModify(opts))
	g.RunFn(t.protoRPCModify(opts))
	g.RunFn(t.protoRPCMessageModify(opts))
	g.RunFn(t.clientCliTxModify(opts))
	g.RunFn(t.clientCliQueryModify(opts))
	g.RunFn(t.typesQuerierModify(opts))
	g.RunFn(t.keeperQuerierModify(opts))
	g.RunFn(t.clientRestRestModify(opts))
	return g, box(cosmosver.Stargate, opts, g)
}

func (t *typedStargate) handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
case *types.MsgCreate%[2]v:
	return handleMsgCreate%[2]v(ctx, k, msg)

case *types.MsgSet%[2]v:
	return handleMsgSet%[2]v(ctx, k, msg)

case *types.MsgDelete%[2]v:
	return handleMsgDelete%[2]v(ctx, k, msg)
`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
import "%s/v1beta/%s.proto";`
		replacement := fmt.Sprintf(template, placeholder,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	rpc One%[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse);
	rpc All%[2]v(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse);
`
		replacement := fmt.Sprintf(template, placeholder2, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCMessageModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/v1beta/querier.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
message QueryGet%[2]vRequest {
	string id = 1;
}

message QueryGet%[2]vResponse {
	%[2]v %[2]v = 1;
}

message QueryAll%[2]vRequest {
	cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAll%[2]vResponse {
	repeated %[2]v %[2]v = 1;
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}`
		replacement := fmt.Sprintf(template, placeholder3,
			strings.Title(opts.TypeName),
			opts.TypeName,
		)
		content := strings.Replace(f.String(), placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) typesKeyModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/keys.go", opts.ModuleName)
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

func (t *typedStargate) typesCodecImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) typesCodecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
cdc.RegisterConcrete(MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(MsgSet%[2]v{}, "%[3]v/Set%[2]v", nil)
cdc.RegisterConcrete(MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
`
		replacement := fmt.Sprintf(template, placeholder2, strings.Title(opts.TypeName), opts.ModuleName)
		content := strings.Replace(f.String(), placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) typesCodecInterfaceModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgSet%[2]v{},
	&MsgDelete%[2]v{},
)`
		replacement := fmt.Sprintf(template, placeholder3, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) clientCliTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

	cmd.AddCommand(CmdCreate%[2]v())
	cmd.AddCommand(CmdSet%[2]v())
	cmd.AddCommand(CmdDelete%[2]v())
`
		replacement := fmt.Sprintf(template, placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) clientCliQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

	cmd.AddCommand(CmdList%[2]v())
	cmd.AddCommand(CmdGet%[2]v())
`
		replacement := fmt.Sprintf(template, placeholder,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) typesQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/querier.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `
const (
	QueryGet%[2]v  = "get-%[1]v"
	QueryList%[2]v = "list-%[1]v"
)
`
		content := f.String() + fmt.Sprintf(template, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) keeperQuerierModify(opts *Options) genny.RunFn {
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
	case types.QueryGet%[2]v:
		return get%[2]v(ctx, path[1], k, legacyQuerierCdc)

	case types.QueryList%[2]v:
		return list%[2]v(ctx, k, legacyQuerierCdc)
`
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

func (t *typedStargate) clientRestRestModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/rest/rest.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%s
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
`
		replacement := fmt.Sprintf(template, placeholder2)
		content := strings.Replace(f.String(), placeholder2, replacement, 1)

		template = `%[1]v
    r.HandleFunc("/%[2]v/%[3]v/{id}", get%[4]vHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/%[2]v/%[3]v", list%[4]vHandler(clientCtx)).Methods("GET")
`
		replacement = fmt.Sprintf(template, placeholder3, opts.ModuleName, pluralize.NewClient().Plural(opts.TypeName), strings.Title(opts.TypeName))
		content = strings.Replace(content, placeholder3, replacement, 1)

		template = `%[1]v
    r.HandleFunc("/%[2]v/%[3]v", create%[4]vHandler(clientCtx)).Methods("POST")
    r.HandleFunc("/%[2]v/%[3]v/{id}", set%[4]vHandler(clientCtx)).Methods("PUT")
    r.HandleFunc("/%[2]v/%[3]v/{id}", delete%[4]vHandler(clientCtx)).Methods("DELETE")
`
		replacement = fmt.Sprintf(template, placeholder44, opts.ModuleName, pluralize.NewClient().Plural(opts.TypeName), strings.Title(opts.TypeName))
		content = strings.Replace(content, placeholder44, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
