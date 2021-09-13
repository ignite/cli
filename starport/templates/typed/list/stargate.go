package list

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/typed"
)

var (
	//go:embed stargate/component/* stargate/component/**/*
	fsStargateComponent embed.FS

	//go:embed stargate/messages/* stargate/messages/**/*
	fsStargateMessages embed.FS

	// stargateComponentTemplate is the template for a Stargate module type component
	stargateComponentTemplate = xgenny.NewEmbedWalker(fsStargateComponent, "stargate/component/")

	// stargateMessagesTemplate is the template for a Stargate module type interaction messages
	stargateMessagesTemplate = xgenny.NewEmbedWalker(fsStargateMessages, "stargate/messages/")
)

// NewStargate returns the generator to scaffold a new type in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(protoQueryModify(replacer, opts))
	g.RunFn(moduleGRPCGatewayModify(replacer, opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn(clientCliQueryModify(replacer, opts))

	// Genesis modifications
	genesisModify(replacer, opts, g)

	if !opts.NoMessage {
		// Modifications for new messages
		g.RunFn(handlerModify(replacer, opts))
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(typesCodecModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))

		// Messages template
		if err := typed.Box(stargateMessagesTemplate, opts, g); err != nil {
			return nil, err
		}
	}

	g.RunFn(frontendSrcStoreAppModify(replacer, opts))

	return g, typed.Box(stargateComponentTemplate, opts, g)
}

func handlerModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), typed.PlaceholderHandlerMsgServer, replacementMsgServer)

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
			typed.Placeholder,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, typed.PlaceholderProtoTxImport,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderProtoTxImport, replacementImport)

		// RPC service
		templateRPC := `%[1]v
  rpc Create%[2]v(MsgCreate%[2]v) returns (MsgCreate%[2]vResponse);
  rpc Update%[2]v(MsgUpdate%[2]v) returns (MsgUpdate%[2]vResponse);
  rpc Delete%[2]v(MsgDelete%[2]v) returns (MsgDelete%[2]vResponse);`
		replacementRPC := fmt.Sprintf(templateRPC, typed.PlaceholderProtoTxRPC,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxRPC, replacementRPC)

		// Messages
		var createFields string
		for i, field := range opts.Fields {
			createFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+2)
		}
		var updateFields string
		for i, field := range opts.Fields {
			updateFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+3)
		}

		// Ensure custom types are imported
		for _, f := range opts.Fields.Custom() {
			importModule := fmt.Sprintf(`
import "%[1]v/%[2]v.proto";`, opts.ModuleName, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", typed.PlaceholderProtoTxImport, importModule)
			content = replacer.Replace(content, typed.PlaceholderProtoTxImport, replacementImport)
		}

		templateMessages := `%[1]v
message MsgCreate%[2]v {
  string %[3]v = 1;
%[4]v}

message MsgCreate%[2]vResponse {
  uint64 id = 1;
}

message MsgUpdate%[2]v {
  string %[3]v = 1;
  uint64 id = 2;
%[5]v}

message MsgUpdate%[2]vResponse {}

message MsgDelete%[2]v {
  string %[3]v = 1;
  uint64 id = 2;
}

message MsgDelete%[2]vResponse {}
`
		replacementMessages := fmt.Sprintf(templateMessages, typed.PlaceholderProtoTxMessage,
			opts.TypeName.UpperCamel,
			opts.MsgSigner.LowerCamel,
			createFields,
			updateFields,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxMessage, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoQueryModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, typed.Placeholder,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacementImport)

		// Add gogo.proto
		replacementGogoImport := typed.EnsureGogoProtoImported(path, typed.Placeholder)
		content = replacer.Replace(content, typed.Placeholder, replacementGogoImport)

		// RPC service
		templateRPC := `%[1]v

	// Queries a %[3]v by id.
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/{id}";
	}

	// Queries a list of %[3]v items.
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementRPC := fmt.Sprintf(templateRPC, typed.Placeholder2,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementRPC)

		// Messages
		templateMessages := `%[1]v
message QueryGet%[2]vRequest {
	uint64 id = 1;
}

message QueryGet%[2]vResponse {
	%[2]v %[2]v = 1 [(gogoproto.nullable) = false];
}

message QueryAll%[2]vRequest {
	cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAll%[2]vResponse {
	repeated %[2]v %[2]v = 1 [(gogoproto.nullable) = false];
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}`
		replacementMessages := fmt.Sprintf(templateMessages, typed.Placeholder3,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
		)
		content = replacer.Replace(content, typed.Placeholder3, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func moduleGRPCGatewayModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacement := `"context"`
		content := replacer.ReplaceOnce(f.String(), typed.Placeholder, replacement)
		replacement = `types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))`
		content = replacer.ReplaceOnce(content, typed.Placeholder2, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesKeyModify(opts *typed.Options) genny.RunFn {
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
`, opts.TypeName.UpperCamel)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), typed.Placeholder, replacementImport)

		// Concrete
		templateConcrete := `%[1]v
cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
`
		replacementConcrete := fmt.Sprintf(templateConcrete, typed.Placeholder2, opts.TypeName.UpperCamel, opts.ModuleName)
		content = replacer.Replace(content, typed.Placeholder2, replacementConcrete)

		// Interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, typed.Placeholder3, opts.TypeName.UpperCamel)
		content = replacer.Replace(content, typed.Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
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
		replacement := fmt.Sprintf(template, typed.Placeholder, opts.TypeName.UpperCamel)
		content := replacer.Replace(f.String(), typed.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
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
		replacement := fmt.Sprintf(template, typed.Placeholder,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func frontendSrcStoreAppModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
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
			typed.Placeholder4,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.Placeholder4, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
