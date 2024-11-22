package singleton

import (
	"crypto/rand"
	"embed"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/module"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

var (
	//go:embed files/messages/* files/messages/**/*
	fsMessages embed.FS

	//go:embed files/component/* files/component/**/*
	fsComponent embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsimapp embed.FS
)

// NewGenerator returns the generator to scaffold a new indexed type in a module.
func NewGenerator(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	var (
		g = genny.New()

		messagesTemplate = xgenny.NewEmbedWalker(
			fsMessages,
			"files/messages/",
			opts.AppPath,
		)
		componentTemplate = xgenny.NewEmbedWalker(
			fsComponent,
			"files/component/",
			opts.AppPath,
		)
		simappTemplate = xgenny.NewEmbedWalker(
			fsimapp,
			"files/simapp/",
			opts.AppPath,
		)
	)

	g.RunFn(protoRPCModify(opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn(keeperModify(replacer, opts))
	g.RunFn(clientCliQueryModify(replacer, opts))
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(genesisTypesModify(replacer, opts))
	g.RunFn(genesisModuleModify(replacer, opts))
	g.RunFn(genesisTestsModify(replacer, opts))
	g.RunFn(genesisTypesTestsModify(replacer, opts))

	// Modifications for new messages
	if !opts.NoMessage {
		g.RunFn(protoTxModify(opts))
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

// typesKeyModify modifies the keys.go file to add a new collection prefix.
func typesKeyModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/keys.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
var (
	%[1]vKey= collections.NewPrefix("%[2]v/value/")
)
`,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
		)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// keeperModify modifies the keeper to add a new collections item type.
func keeperModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/keeper.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateKeeperType := `%[2]v collections.Item[types.%[2]v]
	%[1]v`
		replacementModuleType := fmt.Sprintf(
			templateKeeperType,
			typed.PlaceholderCollectionType,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderCollectionType, replacementModuleType)

		templateKeeperInstantiate := `%[2]v: collections.NewItem(sb, types.%[2]vKey, "%[3]v", codec.CollValue[types.%[2]v](cdc)),
	%[1]v`
		replacementInstantiate := fmt.Sprintf(
			templateKeeperInstantiate,
			typed.PlaceholderCollectionInstantiate,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderCollectionInstantiate, replacementInstantiate)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Modifies query.proto to add the required RPCs and Messages.
//
// What it depends on:
//   - Existence of a service with name "Query". Adds the rpc's there.
func protoRPCModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// Import the type and gogoImport.
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Find service.
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		typenameUpper := opts.TypeName.UpperCamel
		rpcQueryGet := protoutil.NewRPC(
			fmt.Sprintf("Get%s", typenameUpper),
			fmt.Sprintf("QueryGet%sRequest", typenameUpper),
			fmt.Sprintf("QueryGet%sResponse", typenameUpper),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s",
						appModulePath, opts.ModuleName, opts.TypeName.Snake,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("Queries a %v by index.", typenameUpper))
		protoutil.Append(serviceQuery, rpcQueryGet)

		// Add the service messages
		queryGetRequest := protoutil.NewMessage("QueryGet" + typenameUpper + "Request")
		field := protoutil.NewField(typenameUpper, typenameUpper, 1,
			protoutil.WithFieldOptions(protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())),
		)
		queryGetResponse := protoutil.NewMessage(fmt.Sprintf("QueryGet%sResponse", typenameUpper), protoutil.WithFields(field))
		protoutil.Append(protoFile, queryGetRequest, queryGetResponse)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func clientCliQueryModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/autocli.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `{
			RpcMethod: "Get%[2]v",
			Use: "get-%[3]v",
			Short: "Gets a %[4]v",
			Alias: []string{"show-%[3]v"},
		},
		%[1]v`
		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLIQuery,
			opts.TypeName.UpperCamel,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLIQuery, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Modifies the genesis.proto file to add a new field.
//
// What it depends on:
//   - Existence of a message with name "GenesisState". Adds the field there.
func genesisProtoModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		if err = protoutil.AddImports(protoFile, true, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed to add imports to %s: %w", path, err)
		}

		// Add field to GenesisState message.
		genesisState, err := protoutil.GetMessageByName(protoFile, typed.ProtoGenesisStateMessage)
		if err != nil {
			return errors.Errorf("failed while looking up message '%s' in %s: %w", typed.ProtoGenesisStateMessage, path, err)
		}
		seqNumber := protoutil.NextUniqueID(genesisState)
		field := protoutil.NewField(
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			seqNumber,
		)
		protoutil.Append(genesisState, field)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
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

		templateTypesDefault := `%[2]v: nil,
%[1]v`
		replacementTypesDefault := fmt.Sprintf(
			templateTypesDefault,
			typed.PlaceholderGenesisTypesDefault,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisTypesDefault, replacementTypesDefault)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a fields
		sampleFields := ""
		for _, field := range opts.Fields {
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				return err
			}
			sampleFields += field.GenesisArgs(int(n.Int64()) + 1)
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
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				return err
			}
			sampleFields += field.GenesisArgs(int(n.Int64()) + 1)
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
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set if defined
if genState.%[3]v != nil {
	if err := k.%[3]v.Set(ctx, *genState.%[3]v); err != nil {
		return err
	}
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
%[2]v, err := k.%[3]v.Get(ctx)
if err != nil {
	return nil, err
}
genesis.%[3]v = &%[2]v
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

// protoTxModify modifies the tx.proto file to add the required RPCs and messages.
//
// What it expects:
//   - A service named "Msg" to exist in the proto file, it appends the RPCs inside it.
func protoTxModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// Add initial import for the new type:
		if err = protoutil.AddImports(protoFile, true, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports to %s: %w", path, err)
		}
		// Add the RPC service.
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return errors.Errorf("failed while looking up a message 'Msg' in %s: %w", path, err)
		}
		// Append create, update, delete rpcs. Better to append them altogether, single traversal.
		name := opts.TypeName.UpperCamel
		protoutil.Append(serviceMsg,
			protoutil.NewRPC(
				fmt.Sprintf("Create%s", name),
				fmt.Sprintf("MsgCreate%s", name),
				fmt.Sprintf("MsgCreate%sResponse", name),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Update%s", name),
				fmt.Sprintf("MsgUpdate%s", name),
				fmt.Sprintf("MsgUpdate%sResponse", name),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Delete%s", name),
				fmt.Sprintf("MsgDelete%s", name),
				fmt.Sprintf("MsgDelete%sResponse", name),
			),
		)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			// shouldn't really occur.
			return errors.Errorf("failed while adding imports to %s: %w", path, err)
		}

		// Add the messages
		creator := protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1)
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.LowerCamel)
		fields := []*proto.NormalField{creator}
		for i, field := range opts.Fields {
			fields = append(fields, field.ToProtoField(i+3))
		}
		msgCreate := protoutil.NewMessage(
			"MsgCreate"+name,
			protoutil.WithFields(fields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgCreateResponse := protoutil.NewMessage(fmt.Sprintf("MsgCreate%sResponse", name))
		msgUpdate := protoutil.NewMessage(
			"MsgUpdate"+name,
			protoutil.WithFields(fields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgUpdateResponse := protoutil.NewMessage(fmt.Sprintf("MsgUpdate%sResponse", name))
		msgDelete := protoutil.NewMessage(
			"MsgDelete"+name,
			protoutil.WithFields(creator),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgDeleteResponse := protoutil.NewMessage(fmt.Sprintf("MsgDelete%sResponse", name))
		protoutil.Append(protoFile,
			msgCreate, msgCreateResponse, msgUpdate, msgUpdateResponse, msgDelete, msgDeleteResponse,
		)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/autocli.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var positionalArgs, positionalArgsStr string
		for _, field := range opts.Fields {
			positionalArgs += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
			positionalArgsStr += fmt.Sprintf("[%s] ", field.ProtoFieldName())
		}

		template := `{
			RpcMethod: "Create%[2]v",
			Use: "create-%[3]v %[6]s",
			Short: "Create %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[5]s},
		},
		{
			RpcMethod: "Update%[2]v",
			Use: "update-%[3]v %[6]s",
			Short: "Update %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[5]s},
		},
		{
			RpcMethod: "Delete%[2]v",
			Use: "delete-%[3]v",
			Short: "Delete %[4]v",
		},
		%[1]v`

		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLITx,
			opts.TypeName.UpperCamel,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
			strings.TrimSpace(positionalArgs),
			strings.TrimSpace(positionalArgsStr),
		)

		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLITx, replacement)
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

		// Interface
		templateInterface := `registrar.RegisterImplementations((*sdk.Msg)(nil),
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
