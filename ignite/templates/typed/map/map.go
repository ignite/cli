package maptype

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/datatype"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/typed"
)

var (
	//go:embed files/messages/* files/messages/**/*
	fsMessages embed.FS

	//go:embed files/tests/messages/* files/tests/messages/**/*
	fsTestsMessages embed.FS

	//go:embed files/component/* files/component/**/*
	fsComponent embed.FS

	//go:embed files/tests/component/* files/tests/component/**/*
	fsTestsComponent embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsSimapp embed.FS
)

// NewGenerator returns the generator to scaffold a new map type in a module
func NewGenerator(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	// Tests are not generated for map with a custom index that contains only booleans
	// because we can't generate reliable tests for this type
	var generateTest bool
	for _, index := range opts.Indexes {
		if index.DatatypeName != datatype.Bool {
			generateTest = true
		}
	}

	var (
		g = genny.New()

		messagesTemplate = xgenny.NewEmbedWalker(
			fsMessages,
			"files/messages/",
			opts.AppPath,
		)
		testsMessagesTemplate = xgenny.NewEmbedWalker(
			fsTestsMessages,
			"files/tests/messages/",
			opts.AppPath,
		)
		componentTemplate = xgenny.NewEmbedWalker(
			fsComponent,
			"files/component/",
			opts.AppPath,
		)
		testsComponentTemplate = xgenny.NewEmbedWalker(
			fsTestsComponent,
			"files/tests/component/",
			opts.AppPath,
		)
		simappTemplate = xgenny.NewEmbedWalker(
			fsSimapp,
			"files/simapp/",
			opts.AppPath,
		)
	)

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
		if generateTest {
			if err := typed.Box(testsMessagesTemplate, opts, g); err != nil {
				return nil, err
			}
		}
	}

	if generateTest {
		if err := typed.Box(testsComponentTemplate, opts, g); err != nil {
			return nil, err
		}
	}
	return g, typed.Box(componentTemplate, opts, g)
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

		pf, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		gogoImp := protoutil.NewImport("gogoproto/gogo.proto")
		if err = protoutil.AddImports(pf, true, gogoImp, opts.TypeImport()); err != nil {
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}

		var protoIndexes []string
		for _, index := range opts.Indexes {
			protoIndexes = append(protoIndexes, fmt.Sprintf("{%s}", index.ProtoFieldName()))
		}
		indexPath := strings.Join(protoIndexes, "/")
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		srv, err := protoutil.GetServiceByName(pf, "Query")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		typU, typS, typL := opts.TypeName.UpperCamel, opts.TypeName.Snake, opts.TypeName.LowerCamel
		single := protoutil.NewRPC(typU, "QueryGet"+typU+"Request", "QueryGet"+typU+"Response",
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s",
						appModulePath, opts.ModuleName, typS, indexPath,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		all := protoutil.NewRPC(typU+"All", "QueryAll"+typU+"Request", "QueryAll"+typU+"Response",
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s",
						appModulePath, opts.ModuleName, typS,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.Append(srv, single, all)

		//  Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v.proto", opts.AppName, opts.ModuleName, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(pf, false, protoImports...); err != nil {
			// shouldn't really occur.
			return fmt.Errorf("failed to add imports to %s: %w", path, err)
		}

		// Add the messages.
		var queryIndexFields []*proto.NormalField
		for i, index := range opts.Indexes {
			queryIndexFields = append(queryIndexFields, index.ToProtoField(i+1))
		}
		pagT, pagN := "cosmos.base.query.v1beta1.Page", "pagination"
		msgGetReq := protoutil.NewMessage(
			"QueryGet"+typU+"Request",
			protoutil.WithFields(queryIndexFields...),
		)
		gogoproto := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		msgGetResp := protoutil.NewMessage(
			"QueryGet"+typU+"Response",
			protoutil.WithFields(protoutil.NewField(typU, typL, 1, protoutil.WithFieldOptions(gogoproto))),
		)
		msgAllReq := protoutil.NewMessage(
			"QueryAll"+typU+"Request",
			protoutil.WithFields(protoutil.NewField(pagT+"Request", pagN, 1)),
		)
		msgAllResp := protoutil.NewMessage(
			"QueryAll"+typU+"Response",
			protoutil.WithFields(
				protoutil.NewField(
					typU,
					typL,
					1,
					protoutil.Repeated(),
					protoutil.WithFieldOptions(gogoproto),
				),
				protoutil.NewField(pagT+"Response", pagN, 2),
			),
		)
		protoutil.Append(pf, msgGetReq, msgGetResp, msgAllReq, msgAllResp)

		newFile := genny.NewFileS(path, protoutil.Printer(pf))
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
		template := `cmd.AddCommand(CmdList%[2]v())
	cmd.AddCommand(CmdShow%[2]v())
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
		pf, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		gogoproto := protoutil.NewImport("gogoproto/gogo.proto")
		if err = protoutil.AddImports(pf, true, gogoproto, opts.TypeImport()); err != nil {
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Get next available sequence number from GenesisState.
		m, err := protoutil.GetMessageByName(pf, typed.ProtoGenesisStateMessage)
		if err != nil {
			return fmt.Errorf("failed while looking up message '%s' in %s: %w", typed.ProtoGenesisStateMessage, path, err)
		}
		seqNumber := protoutil.NextUniqueID(m)

		// Create new option and append to GenesisState message.
		typL, typU := opts.TypeName.LowerCamel, opts.TypeName.UpperCamel
		opt := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		typeList := protoutil.NewField(
			typU, typL+"List", seqNumber, protoutil.Repeated(), protoutil.WithFieldOptions(opt),
		)
		protoutil.Append(m, typeList)

		newFile := genny.NewFileS(path, protoutil.Printer(pf))
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

		templateTypesImport := `"fmt"`
		content = replacer.ReplaceOnce(content, typed.PlaceholderGenesisTypesImport, templateTypesImport)

		templateTypesDefault := `%[2]vList: []%[2]v{},
%[1]v`
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

		templateTypesValidate := `// Check for duplicated index in %[2]v
%[2]vIndexMap := make(map[string]struct{})

for _, elem := range gs.%[3]vList {
	index := %[4]v
	if _, ok := %[2]vIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for %[2]v")
	}
	%[2]vIndexMap[index] = struct{}{}
}
%[1]v`
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
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, elem)
}
%[1]v`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `genesis.%[3]vList = k.GetAll%[3]v(ctx)
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

func genesisTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a list of two different indexes to use as sample
		sampleIndexes := make([]string, 2)
		for i := 0; i < 2; i++ {
			for _, index := range opts.Indexes {
				sampleIndexes[i] += index.GenesisArgs(i)
			}
		}

		templateState := `%[2]vList: []types.%[2]v{
		{
			%[3]v},
		{
			%[4]v},
	},
	%[1]v`
		replacementState := fmt.Sprintf(
			templateState,
			module.PlaceholderGenesisTestState,
			opts.TypeName.UpperCamel,
			sampleIndexes[0],
			sampleIndexes[1],
		)
		content := replacer.Replace(f.String(), module.PlaceholderGenesisTestState, replacementState)

		templateAssert := `require.ElementsMatch(t, genesisState.%[2]vList, got.%[2]vList)
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

		// Create a list of two different indexes to use as sample
		sampleIndexes := make([]string, 2)
		for i := 0; i < 2; i++ {
			for _, index := range opts.Indexes {
				sampleIndexes[i] += index.GenesisArgs(i)
			}
		}

		templateValid := `%[2]vList: []types.%[2]v{
	{
		%[3]v},
	{
		%[4]v},
},
%[1]v`
		replacementValid := fmt.Sprintf(
			templateValid,
			module.PlaceholderTypesGenesisValidField,
			opts.TypeName.UpperCamel,
			sampleIndexes[0],
			sampleIndexes[1],
		)
		content := replacer.Replace(f.String(), module.PlaceholderTypesGenesisValidField, replacementValid)

		templateDuplicated := `{
	desc:     "duplicated %[2]v",
	genState: &types.GenesisState{
		%[3]vList: []types.%[3]v{
			{
				%[4]v},
			{
				%[4]v},
		},
	},
	valid:    false,
},
%[1]v`
		replacementDuplicated := fmt.Sprintf(
			templateDuplicated,
			module.PlaceholderTypesGenesisTestcase,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			sampleIndexes[0],
		)
		content = replacer.Replace(content, module.PlaceholderTypesGenesisTestcase, replacementDuplicated)

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

		pf, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		if err = protoutil.AddImports(pf, true, opts.TypeImport()); err != nil {
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}

		// RPC service
		s, err := protoutil.GetServiceByName(pf, "Msg")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		// better to append them altogether, single traversal.
		typU := opts.TypeName.UpperCamel
		protoutil.Append(s,
			protoutil.NewRPC("Create"+typU, "MsgCreate"+typU, "MsgCreate"+typU+"Response"),
			protoutil.NewRPC("Update"+typU, "MsgUpdate"+typU, "MsgUpdate"+typU+"Response"),
			protoutil.NewRPC("Delete"+typU, "MsgDelete"+typU, "MsgDelete"+typU+"Response"),
		)

		// Messages
		var indexes []*proto.NormalField
		for i, index := range opts.Indexes {
			indexes = append(indexes, index.ToProtoField(i+2))
		}

		var fields []*proto.NormalField
		for i, f := range opts.Fields {
			fields = append(fields, f.ToProtoField(i+2+len(opts.Indexes)))
		}

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range append(opts.Fields.ProtoImports(), opts.Indexes.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v.proto", opts.AppName, opts.ModuleName, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(pf, false, protoImports...); err != nil {
			// shouldn't really occur.
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}
		commonFields := []*proto.NormalField{protoutil.NewField("string", opts.MsgSigner.LowerCamel, 1)}
		commonFields = append(commonFields, indexes...)

		msgCreate := protoutil.NewMessage(
			"MsgCreate"+typU,
			protoutil.WithFields(append(commonFields, fields...)...),
		)
		msgCreateResp := protoutil.NewMessage("MsgCreate" + typU + "Response")

		msgUpdate := protoutil.NewMessage(
			"MsgUpdate"+typU,
			protoutil.WithFields(append(commonFields, fields...)...),
		)
		msgUpdateResp := protoutil.NewMessage("MsgUpdate" + typU + "Response")

		msgDelete := protoutil.NewMessage("MsgDelete"+typU, protoutil.WithFields(commonFields...))
		msgDeleteResp := protoutil.NewMessage("MsgDelete" + typU + "Response")
		protoutil.Append(pf,
			msgCreate, msgCreateResp, msgUpdate, msgUpdateResp, msgDelete, msgDeleteResp,
		)

		newFile := genny.NewFileS(path, protoutil.Printer(pf))
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
