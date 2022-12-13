package singleton

import (
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/emicklei/proto"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/typed"
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

	g.RunFn(typesKeyModify(opts))
	g.RunFn(protoRPCModify(opts))
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

func typesKeyModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/keys.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%[1]vKey= "%[1]v/value/"
)
`, opts.TypeName.UpperCamel)
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
		path := opts.ProtoPath("query.proto")
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
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Find service.
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		typenameUpper := opts.TypeName.UpperCamel
		rpcQueryGet := protoutil.NewRPC(typenameUpper, "QueryGet"+typenameUpper+"Request", "QueryGet"+typenameUpper+"Response",
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
		queryGetResponse := protoutil.NewMessage("QueryGet"+typenameUpper+"Response", protoutil.WithFields(field))
		protoutil.Append(protoFile, queryGetRequest, queryGetResponse)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
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

// Modifies the genesis.proto file to add a new field.
//
// What it depends on:
//   - Existence of a message with name "GenesisState". Adds the field there.
func genesisProtoModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoPath("genesis.proto")
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
			return fmt.Errorf("failed to add imports to %s: %w", path, err)
		}

		// Add field to GenesisState message.
		genesisState, err := protoutil.GetMessageByName(protoFile, typed.ProtoGenesisStateMessage)
		if err != nil {
			return fmt.Errorf("failed while looking up message '%s' in %s: %w", typed.ProtoGenesisStateMessage, path, err)
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

// protoTxModify modifies the tx.proto file to add the required RPCs and messages.
//
// What it expects:
//   - A service named "Msg" to exist in the proto file, it appends the RPCs inside it.
func protoTxModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoPath("tx.proto")
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
			return fmt.Errorf("failed while adding imports to %s: %w", path, err)
		}
		// Add the RPC service.
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return fmt.Errorf("failed while looking up a message 'Msg' in %s: %w", path, err)
		}
		// Append create, update, delete rpcs. Better to append them altogether, single traversal.
		name := opts.TypeName.UpperCamel
		protoutil.Append(serviceMsg,
			protoutil.NewRPC("Create"+name, "MsgCreate"+name, "MsgCreate"+name+"Response"),
			protoutil.NewRPC("Update"+name, "MsgUpdate"+name, "MsgUpdate"+name+"Response"),
			protoutil.NewRPC("Delete"+name, "MsgDelete"+name, "MsgDelete"+name+"Response"),
		)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v.proto", opts.AppName, opts.ModuleName, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			// shouldn't really occur.
			return fmt.Errorf("failed while adding imports to %s: %w", path, err)
		}

		// Add the messages
		creator := protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1)
		fields := []*proto.NormalField{creator}
		for i, field := range opts.Fields {
			fields = append(fields, field.ToProtoField(i+3))
		}
		msgCreate := protoutil.NewMessage(
			"MsgCreate"+name, protoutil.WithFields(fields...),
		)
		msgCreateResponse := protoutil.NewMessage("MsgCreate" + name + "Response")
		msgUpdate := protoutil.NewMessage(
			"MsgUpdate"+name, protoutil.WithFields(fields...),
		)
		msgUpdateResponse := protoutil.NewMessage("MsgUpdate" + name + "Response")
		msgDelete := protoutil.NewMessage(
			"MsgDelete"+name, protoutil.WithFields(creator),
		)
		msgDeleteResponse := protoutil.NewMessage("MsgDelete" + name + "Response")
		protoutil.Append(protoFile,
			msgCreate, msgCreateResponse, msgUpdate, msgUpdateResponse, msgDelete, msgDeleteResponse,
		)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
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
