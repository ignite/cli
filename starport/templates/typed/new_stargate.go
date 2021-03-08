package typed

import (
	"fmt"
	"os"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/genny"
)

type typedStargate struct {
}

// New ...
func NewStargate(opts *Options) (*genny.Generator, error) {
	t := typedStargate{}
	g := genny.New()

	if opts.Legacy {
		g.RunFn(t.handlerModifyLegacy(opts))
	} else {
		g.RunFn(t.handlerModify(opts))
		g.RunFn(t.protoTxImportModify(opts))
		g.RunFn(t.protoTxRPCModify(opts))
		g.RunFn(t.protoTxMessageModify(opts))
	}

	g.RunFn(t.typesKeyModify(opts))
	g.RunFn(t.typesCodecModify(opts))
	g.RunFn(t.typesCodecImportModify(opts))
	g.RunFn(t.typesCodecInterfaceModify(opts))
	g.RunFn(t.protoRPCImportModify(opts))
	g.RunFn(t.protoRPCModify(opts))
	g.RunFn(t.protoRPCMessageModify(opts))
	g.RunFn(t.moduleGRPCGatewayModify(opts))
	g.RunFn(t.clientCliTxModify(opts))
	g.RunFn(t.clientCliQueryModify(opts))
	g.RunFn(t.typesQueryModify(opts))
	g.RunFn(t.keeperQueryModify(opts))
	g.RunFn(t.clientRestRestModify(opts))
	g.RunFn(t.frontendSrcStoreAppModify(opts))
	t.genesisModify(opts, g)

	if opts.Legacy {
		return g, Box(stargateLegacyTemplate, opts, g)
	}
	return g, Box(stargateTemplate, opts, g)
}

func (t *typedStargate) handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := strings.Replace(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer, 1)

		templateHandlers := `%[1]v
		case *types.MsgCreate%[2]v:
					res, err := msgServer.Create%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUpdate%[2]v:
					res, err := msgServer.Update%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDelete%[2]v:
					res, err := msgServer.Delete%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			Placeholder,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, Placeholder, replacementHandlers, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoTxImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
import "%s/%s.proto";`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxImport,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxImport, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoTxRPCModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
  rpc Create%[2]v(MsgCreate%[2]v) returns (MsgCreate%[2]vResponse);
  rpc Update%[2]v(MsgUpdate%[2]v) returns (MsgUpdate%[2]vResponse);
  rpc Delete%[2]v(MsgDelete%[2]v) returns (MsgDelete%[2]vResponse);`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxRPC,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxRPC, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoTxMessageModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var createFields string
		for i, field := range opts.Fields {
			createFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+2)
		}
		var updateFields string
		for i, field := range opts.Fields {
			updateFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+3)
		}

		template := `%[1]v
message MsgCreate%[2]v {
  string creator = 1;
%[3]v}

message MsgCreate%[2]vResponse {
  string id = 1;
}

message MsgUpdate%[2]v {
  string creator = 1;
  string id = 2;
%[4]v}

message MsgUpdate%[2]vResponse { }

message MsgDelete%[2]v {
  string creator = 1;
  string id = 2;
}

message MsgDelete%[2]vResponse { }
`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxMessage,
			strings.Title(opts.TypeName),
			createFields,
			updateFields,
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxMessage, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%s
import "%s/%s.proto";`
		replacement := fmt.Sprintf(template, Placeholder,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/{id}";
	}
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacement := fmt.Sprintf(template, Placeholder2,
			strings.Title(opts.TypeName),
			opts.TypeName,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content := strings.Replace(f.String(), Placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoRPCMessageModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
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
		replacement := fmt.Sprintf(template, Placeholder3,
			strings.Title(opts.TypeName),
			opts.TypeName,
		)
		content := strings.Replace(f.String(), Placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) moduleGRPCGatewayModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacement := `"context"`
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		replacement = `types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))`
		content = strings.Replace(content, Placeholder2, replacement, 1)
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
	%[1]vKey= "%[1]v-value-"
	%[1]vCountKey= "%[1]v-count-"
)
`, strings.Title(opts.TypeName))
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
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
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
cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
`
		replacement := fmt.Sprintf(template, Placeholder2, strings.Title(opts.TypeName), opts.ModuleName)
		content := strings.Replace(f.String(), Placeholder2, replacement, 1)
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
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)`
		replacement := fmt.Sprintf(template, Placeholder3, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), Placeholder3, replacement, 1)
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
	cmd.AddCommand(CmdUpdate%[2]v())
	cmd.AddCommand(CmdDelete%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
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
	cmd.AddCommand(CmdShow%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) typesQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/query.go", opts.ModuleName)
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

func (t *typedStargate) keeperQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/keeper/query.go", opts.ModuleName)
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
		replacement2 := fmt.Sprintf(template2, Placeholder, opts.ModulePath, opts.ModuleName)
		replacement3 := fmt.Sprintf(template3, Placeholder2, strings.Title(opts.TypeName))
		content := f.String()
		content = strings.Replace(content, replacement, "", 1)
		content = strings.Replace(content, Placeholder, replacement2, 1)
		content = strings.Replace(content, Placeholder2, replacement3, 1)
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
		replacement := fmt.Sprintf(template, Placeholder2)
		content := strings.Replace(f.String(), Placeholder2, replacement, 1)

		template = `%[1]v
    r.HandleFunc("/%[2]v/%[3]v/{id}", get%[4]vHandler(clientCtx)).Methods("GET")
    r.HandleFunc("/%[2]v/%[3]v", list%[4]vHandler(clientCtx)).Methods("GET")
`
		replacement = fmt.Sprintf(template, Placeholder3, opts.ModuleName, pluralize.NewClient().Plural(opts.TypeName), strings.Title(opts.TypeName))
		content = strings.Replace(content, Placeholder3, replacement, 1)

		template = `%[1]v
    r.HandleFunc("/%[2]v/%[3]v", create%[4]vHandler(clientCtx)).Methods("POST")
    r.HandleFunc("/%[2]v/%[3]v/{id}", update%[4]vHandler(clientCtx)).Methods("POST")
    r.HandleFunc("/%[2]v/%[3]v/{id}", delete%[4]vHandler(clientCtx)).Methods("POST")
`
		replacement = fmt.Sprintf(template, Placeholder44, opts.ModuleName, pluralize.NewClient().Plural(opts.TypeName), strings.Title(opts.TypeName))
		content = strings.Replace(content, Placeholder44, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) frontendSrcStoreAppModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "vue/src/views/Index.vue"
		f, err := r.Disk.Find(path)
		if os.IsNotExist(err) {
			// Skip modification if the app doesn't contain front-end
			return nil
		}
		if err != nil {
			return err
		}
		replacement := fmt.Sprintf(`%[1]v
		<SpType modulePath="%[2]v/%[3]v/%[2]v.%[3]v.%[4]v" moduleType="%[5]v"  />`,
			Placeholder4,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), Placeholder4, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// --- Legacy Stargate Types

func (t *typedStargate) handlerModifyLegacy(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
case *types.MsgCreate%[2]v:
	return handleMsgCreate%[2]v(ctx, k, msg)
case *types.MsgUpdate%[2]v:
	return handleMsgUpdate%[2]v(ctx, k, msg)
case *types.MsgDelete%[2]v:
	return handleMsgDelete%[2]v(ctx, k, msg)
`
		replacement := fmt.Sprintf(template, Placeholder, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
