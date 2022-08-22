package singleton

import (
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/typed"
)

var (
	//go:embed stargate/component/* stargate/component/**/*
	fsStargateComponent embed.FS

	//go:embed stargate/messages/* stargate/messages/**/*
	fsStargateMessages embed.FS

	//go:embed stargate/simapp/* stargate/simapp/**/*
	fsStargateSimapp embed.FS
)

// NewStargate returns the generator to scaffold a new indexed type in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	var (
		g = genny.New()

		messagesTemplate = xgenny.NewEmbedWalker(
			fsStargateMessages,
			"stargate/messages/",
			opts.AppPath,
		)
		componentTemplate = xgenny.NewEmbedWalker(
			fsStargateComponent,
			"stargate/component/",
			opts.AppPath,
		)
		simappTemplate = xgenny.NewEmbedWalker(
			fsStargateSimapp,
			"stargate/simapp/",
			opts.AppPath,
		)
	)

	g.RunFn(typesKeyModify(opts))
	g.RunFn(protoRPCModify(replacer, opts))
	g.RunFn(moduleGRPCGatewayModify(replacer, opts))
	g.RunFn(clientCliQueryModify(replacer, opts))
	g.RunFn(genesisProtoModify(replacer, opts))
	g.RunFn(genesisTypesModify(replacer, opts))
	g.RunFn(genesisModuleModify(replacer, opts))
	g.RunFn(genesisTestsModify(replacer, opts))
	g.RunFn(genesisTypesTestsModify(replacer, opts))

	// Modifications for new messages
	if !opts.NoMessage {
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(handlerModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))
		g.RunFn(typesCodecModify(replacer, opts))

		if !opts.NoSimulation {
			g.RunFn(moduleSimulationModify(replacer, opts))
			if err := typed.Box(simappTemplate, opts, g); err != nil {
				return nil, err
			}
		}

		if err := typed.Box(messagesTemplate, opts, g); err != nil {
			return nil, err
		}
	}

	return g, typed.Box(componentTemplate, opts, g)
}

func typesKeyModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/keys.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%[1]vKey= "%[1]v-value-"
)
`, opts.TypeName.UpperCamel)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoRPCModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `import "%s/%s.proto";
%s`
		replacementImport := fmt.Sprintf(templateImport,
			opts.ModuleName,
			opts.TypeName.Snake,
			typed.Placeholder,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacementImport)

		// Add gogo.proto
		replacementGogoImport := typed.EnsureGogoProtoImported(path, typed.Placeholder)
		content = replacer.Replace(content, typed.Placeholder, replacementGogoImport)

		// Add the service
		templateService := `// Queries a %[2]v by index.
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[3]v/%[4]v/%[5]v";
	}
%[1]v`
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		replacementService := fmt.Sprintf(templateService, typed.Placeholder2,
			opts.TypeName.UpperCamel,
			appModulePath,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementService)

		// Add the service messages
		templateMessage := `message QueryGet%[2]vRequest {}

message QueryGet%[2]vResponse {
	%[2]v %[2]v = 1 [(gogoproto.nullable) = false];
}
%[1]v`
		replacementMessage := fmt.Sprintf(templateMessage, typed.Placeholder3,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
		)
		content = replacer.Replace(content, typed.Placeholder3, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func moduleGRPCGatewayModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module.go")
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

func clientCliQueryModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/query.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `cmd.AddCommand(CmdShow%[2]v())
%[1]v`
		replacement := fmt.Sprintf(template, typed.Placeholder,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisProtoModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `import "%[2]v/%[3]v.proto";
%[1]v`
		replacementProtoImport := fmt.Sprintf(
			templateProtoImport,
			typed.PlaceholderGenesisProtoImport,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisProtoImport, replacementProtoImport)

		// Parse proto file to determine the field numbers
		highestNumber, err := typed.GenesisStateHighestFieldNumber(path)
		if err != nil {
			return err
		}

		templateProtoState := `%[2]v %[3]v = %[4]v;
  %[1]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			typed.PlaceholderGenesisProtoState,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			highestNumber+1,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisProtoState, replacementProtoState)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := typed.PatchGenesisTypeImport(replacer, f.String())

		templateTypesDefault := `%[2]v: nil,
%[1]v`
		replacementTypesDefault := fmt.Sprintf(
			templateTypesDefault,
			typed.PlaceholderGenesisTypesDefault,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesDefault, replacementTypesDefault)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a fields
		sampleFields := ""
		for _, field := range opts.Fields {
			sampleFields += field.GenesisArgs(rand.Intn(100) + 1)
		}

		templateState := `%[2]v: &types.%[2]v{
		%[3]v},
		%[1]v`
		replacementState := fmt.Sprintf(
			templateState,
			module.PlaceholderGenesisTestState,
			opts.TypeName.UpperCamel,
			sampleFields,
		)
		content := replacer.Replace(f.String(), module.PlaceholderGenesisTestState, replacementState)

		templateAssert := `require.Equal(t, genesisState.%[2]v, got.%[2]v)
%[1]v`
		replacementTests := fmt.Sprintf(
			templateAssert,
			module.PlaceholderGenesisTestAssert,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, module.PlaceholderGenesisTestAssert, replacementTests)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a fields
		sampleFields := ""
		for _, field := range opts.Fields {
			sampleFields += field.GenesisArgs(rand.Intn(100) + 1)
		}

		templateValid := `%[2]v: &types.%[2]v{
		%[3]v},
%[1]v`
		replacementValid := fmt.Sprintf(
			templateValid,
			module.PlaceholderTypesGenesisValidField,
			opts.TypeName.UpperCamel,
			sampleFields,
		)
		content := replacer.Replace(f.String(), module.PlaceholderTypesGenesisValidField, replacementValid)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set if defined
if genState.%[3]v != nil {
	k.Set%[3]v(ctx, *genState.%[3]v)
}
%[1]v`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `// Get all %[2]v
%[2]v, found := k.Get%[3]v(ctx)
if found {
	genesis.%[3]v = &%[2]v
}
%[1]v`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			typed.PlaceholderGenesisModuleExport,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisModuleExport, replacementModuleExport)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `import "%s/%s.proto";
%s`
		replacementImport := fmt.Sprintf(templateImport,
			opts.ModuleName,
			opts.TypeName.Snake,
			typed.PlaceholderProtoTxImport,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderProtoTxImport, replacementImport)

		// RPC service
		templateRPC := `  rpc Create%[2]v(MsgCreate%[2]v) returns (MsgCreate%[2]vResponse);
  rpc Update%[2]v(MsgUpdate%[2]v) returns (MsgUpdate%[2]vResponse);
  rpc Delete%[2]v(MsgDelete%[2]v) returns (MsgDelete%[2]vResponse);
%[1]v`
		replacementRPC := fmt.Sprintf(templateRPC, typed.PlaceholderProtoTxRPC,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxRPC, replacementRPC)

		// Messages
		var fields string
		for i, field := range opts.Fields {
			fields += fmt.Sprintf("  %s;\n", field.ProtoType(i+3))
		}

		// Ensure custom types are imported
		protoImports := opts.Fields.ProtoImports()
		for _, f := range opts.Fields.Custom() {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", typed.PlaceholderProtoTxImport, importModule)
			content = replacer.Replace(content, typed.PlaceholderProtoTxImport, replacementImport)
		}

		templateMessages := `message MsgCreate%[2]v {
  string %[3]v = 1;
%[4]v}
message MsgCreate%[2]vResponse {}

message MsgUpdate%[2]v {
  string %[3]v = 1;
%[4]v}
message MsgUpdate%[2]vResponse {}

message MsgDelete%[2]v {
  string %[3]v = 1;
}
message MsgDelete%[2]vResponse {}
%[1]v`
		replacementMessages := fmt.Sprintf(templateMessages, typed.PlaceholderProtoTxMessage,
			opts.TypeName.UpperCamel,
			opts.MsgSigner.LowerCamel,
			fields,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxMessage, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "handler.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), typed.PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `case *types.MsgCreate%[2]v:
					res, err := msgServer.Create%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdate%[2]v:
					res, err := msgServer.Update%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDelete%[2]v:
					res, err := msgServer.Delete%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
%[1]v`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			typed.Placeholder,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/tx.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `cmd.AddCommand(CmdCreate%[2]v())
	cmd.AddCommand(CmdUpdate%[2]v())
	cmd.AddCommand(CmdDelete%[2]v())
%[1]v`
		replacement := fmt.Sprintf(template, typed.Placeholder, opts.TypeName.UpperCamel)
		content := replacer.Replace(f.String(), typed.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := f.String()

		// Import
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content = replacer.ReplaceOnce(content, typed.Placeholder, replacementImport)

		// Concrete
		templateConcrete := `cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
%[1]v`
		replacementConcrete := fmt.Sprintf(
			templateConcrete,
			typed.Placeholder2,
			opts.TypeName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementConcrete)

		// Interface
		templateInterface := `registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)
%[1]v`
		replacementInterface := fmt.Sprintf(
			templateInterface,
			typed.Placeholder3,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
