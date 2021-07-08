package indexed

import (
	"embed"
	"fmt"
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

	//go:embed stargate/tests/component/* stargate/tests/component/**/*
	fsStargateTestsComponent embed.FS

	//go:embed stargate/tests/messages/* stargate/tests/messages/**/*
	fsStargateTestsMessages embed.FS

	// stargateMapComponentTemplate allows to scaffold a new map component in a Stargate module
	stargateMapComponentTemplate = xgenny.NewEmbedWalker(fsStargateComponent, "stargate/component/")

	// stargateMapMessagesTemplate allows to scaffold map CRUD messages in a Stargate module
	stargateMapMessagesTemplate = xgenny.NewEmbedWalker(fsStargateMessages, "stargate/messages/")

	// stargateMapTestsComponentTemplate allows to scaffold tests for map component in a Stargate module
	stargateMapTestsComponentTemplate = xgenny.NewEmbedWalker(fsStargateTestsComponent, "stargate/tests/component/")

	// stargateMapTestsMessagesTemplate allows to scaffold tests for map CRUD messages in a Stargate module
	stargateMapTestsMessagesTemplate = xgenny.NewEmbedWalker(fsStargateTestsMessages, "stargate/tests/messages/")
)

// NewStargate returns the generator to scaffold a new indexed type in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	// Tests are currently not generated for map with a custom index
	// TODO: handle test generation for map with custom index
	generateTest := len(opts.Indexes) == 1 && opts.Indexes[0].Name.UpperCamel == "Index"

	g := genny.New()

	g.RunFn(protoRPCModify(replacer, opts))
	g.RunFn(moduleGRPCGatewayModify(replacer, opts))
	g.RunFn(clientCliQueryModify(replacer, opts))
	g.RunFn(genesisProtoModify(replacer, opts))
	g.RunFn(genesisTypesModify(replacer, opts))
	g.RunFn(genesisModuleModify(replacer, opts))

	// Modifications for new messages
	if !opts.NoMessage {
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(handlerModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))
		g.RunFn(typesCodecModify(replacer, opts))

		if err := typed.Box(stargateMapMessagesTemplate, opts, g); err != nil {
			return nil, err
		}
		if generateTest {
			if err := typed.Box(stargateMapTestsMessagesTemplate, opts, g); err != nil {
				return nil, err
			}
		}
	}

	if generateTest {
		if err := typed.Box(stargateMapTestsComponentTemplate, opts, g); err != nil {
			return nil, err
		}
	}
	return g, typed.Box(stargateMapComponentTemplate, opts, g)
}

func protoRPCModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, typed.Placeholder,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacementImport)

		var lowercaseIndexes []string
		for _, index := range opts.Indexes {
			lowercaseIndexes = append(lowercaseIndexes, fmt.Sprintf("{%s}", index.Name.Lowercase))
		}
		indexPath := strings.Join(lowercaseIndexes, "/")

		// Add the service
		templateService := `%[1]v

	// Queries a %[3]v by index.
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/%[7]v";
	}

	// Queries a list of %[3]v items.
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementService := fmt.Sprintf(templateService, typed.Placeholder2,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
			indexPath,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementService)

		// Add the service messages
		var queryIndexFields string
		for i, index := range opts.Indexes {
			queryIndexFields += fmt.Sprintf(
				"  %s %s = %d;\n",
				index.Datatype,
				index.Name.LowerCamel,
				i+1,
			)
		}

		templateMessage := `%[1]v
message QueryGet%[2]vRequest {
	%[4]v
}

message QueryGet%[2]vResponse {
	%[2]v %[3]v = 1;
}

message QueryAll%[2]vRequest {
	cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAll%[2]vResponse {
	repeated %[2]v %[3]v = 1;
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}`
		replacementMessage := fmt.Sprintf(templateMessage,
			typed.Placeholder3,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			queryIndexFields,
		)
		content = replacer.Replace(content, typed.Placeholder3, replacementMessage)

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

func genesisProtoModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/genesis.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `%[1]v
import "%[2]v/%[3]v.proto";`
		replacementProtoImport := fmt.Sprintf(
			templateProtoImport,
			typed.PlaceholderGenesisProtoImport,
			opts.ModuleName,
			opts.TypeName.Snake,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisProtoImport, replacementProtoImport)

		// Determine the new field number
		fieldNumber := strings.Count(content, typed.PlaceholderGenesisProtoStateField) + 1

		templateProtoState := `%[1]v
		repeated %[2]v %[3]vList = %[4]v; %[5]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			typed.PlaceholderGenesisProtoState,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			fieldNumber,
			typed.PlaceholderGenesisProtoStateField,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisProtoState, replacementProtoState)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := typed.PatchGenesisTypeImport(replacer, f.String())

		templateTypesImport := `"fmt"`
		content = replacer.ReplaceOnce(content, typed.PlaceholderGenesisTypesImport, templateTypesImport)

		templateTypesDefault := `%[1]v
%[2]vList: []*%[2]v{},`
		replacementTypesDefault := fmt.Sprintf(
			templateTypesDefault,
			typed.PlaceholderGenesisTypesDefault,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesDefault, replacementTypesDefault)

		// lines of code to call the key function with the indexes of the element
		var indexArgs []string
		for _, index := range opts.Indexes {
			indexArgs = append(indexArgs, "elem."+index.Name.UpperCamel)
		}
		keyCall := fmt.Sprintf("%sKey(%s)", opts.TypeName.UpperCamel, strings.Join(indexArgs, ","))

		templateTypesValidate := `%[1]v
// Check for duplicated index in %[2]v
%[2]vIndexMap := make(map[string]struct{})

for _, elem := range gs.%[3]vList {
	index := %[4]v
	if _, ok := %[2]vIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for %[2]v")
	}
	%[2]vIndexMap[index] = struct{}{}
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			typed.PlaceholderGenesisTypesValidate,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			fmt.Sprintf("string(%s)", keyCall),
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesValidate, replacementTypesValidate)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `%[1]v
// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, *elem)
}

`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `%[1]v
// Get all %[2]v
%[2]vList := k.GetAll%[3]v(ctx)
for _, elem := range %[2]vList {
	elem := elem
	genesis.%[3]vList = append(genesis.%[3]vList, &elem)
}
`
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
		var indexes string
		for i, index := range opts.Indexes {
			indexes += fmt.Sprintf(
				"  %s %s = %d;\n",
				index.Datatype,
				index.Name.LowerCamel,
				i+2,
			)
		}

		var fields string
		for i, field := range opts.Fields {
			fields += fmt.Sprintf(
				"  %s %s = %d;\n",
				field.Datatype,
				field.Name.LowerCamel,
				i+2+len(opts.Indexes),
			)
		}

		templateMessages := `%[1]v
message MsgCreate%[2]v {
  string creator = 1;
%[3]v
%[4]v}
message MsgCreate%[2]vResponse { }

message MsgUpdate%[2]v {
  string creator = 1;
%[3]v
%[4]v}
message MsgUpdate%[2]vResponse { }

message MsgDelete%[2]v {
  string creator = 1;
%[3]v
}
message MsgDelete%[2]vResponse { }
`
		replacementMessages := fmt.Sprintf(templateMessages, typed.PlaceholderProtoTxMessage,
			opts.TypeName.UpperCamel,
			indexes,
			fields,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxMessage, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
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

func typesCodecModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := f.String()

		// Import
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content = replacer.ReplaceOnce(content, typed.Placeholder, replacementImport)

		// Concrete
		templateConcrete := `%[1]v
cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
`
		replacementConcrete := fmt.Sprintf(
			templateConcrete,
			typed.Placeholder2,
			opts.TypeName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementConcrete)

		// Interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)`
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
