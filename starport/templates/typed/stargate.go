package typed

import (
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
)

type typedStargate struct {
}

// NewStargate returns the generator to scaffold a new type in a Stargate module
func NewStargate(opts *Options) (*genny.Generator, error) {
	t := typedStargate{}
	g := genny.New()

	g.RunFn(t.protoQueryModify(opts))
	g.RunFn(t.moduleGRPCGatewayModify(opts))
	g.RunFn(t.typesKeyModify(opts))
	g.RunFn(t.clientCliQueryModify(opts))
	g.RunFn(t.typesQueryModify(opts))

	// Genesis modifications
	t.genesisModify(opts, g)

	if !opts.NoMessage {
		// Modifications for new messages
		g.RunFn(t.handlerModify(opts))
		g.RunFn(t.protoTxModify(opts))
		g.RunFn(t.typesCodecModify(opts))
		g.RunFn(t.clientCliTxModify(opts))

		// Messages template
		if err := Box(stargateMessagesTemplate, opts, g); err != nil {
			return nil, err
		}
	}

	g.RunFn(t.frontendSrcStoreAppModify(opts))

	return g, Box(stargateComponentTemplate, opts, g)
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

func (t *typedStargate) protoTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, PlaceholderProtoTxImport,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxImport, replacementImport, 1)

		// RPC service
		templateRPC := `%[1]v
  rpc Create%[2]v(MsgCreate%[2]v) returns (MsgCreate%[2]vResponse);
  rpc Update%[2]v(MsgUpdate%[2]v) returns (MsgUpdate%[2]vResponse);
  rpc Delete%[2]v(MsgDelete%[2]v) returns (MsgDelete%[2]vResponse);`
		replacementRPC := fmt.Sprintf(templateRPC, PlaceholderProtoTxRPC,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, PlaceholderProtoTxRPC, replacementRPC, 1)

		// Messages
		var createFields string
		for i, field := range opts.Fields {
			createFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+2)
		}
		var updateFields string
		for i, field := range opts.Fields {
			updateFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+3)
		}

		templateMessages := `%[1]v
message MsgCreate%[2]v {
  string creator = 1;
%[3]v}

message MsgCreate%[2]vResponse {
  uint64 id = 1;
}

message MsgUpdate%[2]v {
  string creator = 1;
  uint64 id = 2;
%[4]v}

message MsgUpdate%[2]vResponse { }

message MsgDelete%[2]v {
  string creator = 1;
  uint64 id = 2;
}

message MsgDelete%[2]vResponse { }
`
		replacementMessages := fmt.Sprintf(templateMessages, PlaceholderProtoTxMessage,
			strings.Title(opts.TypeName),
			createFields,
			updateFields,
		)
		content = strings.Replace(content, PlaceholderProtoTxMessage, replacementMessages, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func (t *typedStargate) protoQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, Placeholder,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), Placeholder, replacementImport, 1)

		// RPC service
		templateRPC := `%[1]v
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/{id}";
	}
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementRPC := fmt.Sprintf(templateRPC, Placeholder2,
			strings.Title(opts.TypeName),
			opts.TypeName,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content = strings.Replace(content, Placeholder2, replacementRPC, 1)

		// Messages
		templateMessages := `%[1]v
message QueryGet%[2]vRequest {
	uint64 id = 1;
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
		replacementMessages := fmt.Sprintf(templateMessages, Placeholder3,
			strings.Title(opts.TypeName),
			opts.TypeName,
		)
		content = strings.Replace(content, Placeholder3, replacementMessages, 1)

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

func (t *typedStargate) typesCodecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := strings.Replace(f.String(), Placeholder, replacementImport, 1)

		// Concrete
		templateConcrete := `%[1]v
cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
`
		replacementConcrete := fmt.Sprintf(templateConcrete, Placeholder2, strings.Title(opts.TypeName), opts.ModuleName)
		content = strings.Replace(content, Placeholder2, replacementConcrete, 1)

		// Interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, Placeholder3, strings.Title(opts.TypeName))
		content = strings.Replace(content, Placeholder3, replacementInterface, 1)

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

func (t *typedStargate) frontendSrcStoreAppModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "vue/src/views/Types.vue"
		f, err := r.Disk.Find(path)
		if os.IsNotExist(err) {
			// Skip modification if the app doesn't contain front-end
			return nil
		}
		if err != nil {
			return err
		}
		replacement := fmt.Sprintf(`%[1]v
		<SpType modulePath="%[2]v.%[3]v.%[4]v" moduleType="%[5]v"  />`,
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
