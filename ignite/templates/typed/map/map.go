package maptype

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field/datatype"
	"github.com/ignite/cli/v28/ignite/templates/module"
	"github.com/ignite/cli/v28/ignite/templates/typed"
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

// NewGenerator returns the generator to scaffold a new map type in a module.
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

		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		var protoIndexes []string
		for _, index := range opts.Indexes {
			protoIndexes = append(protoIndexes, fmt.Sprintf("{%s}", index.ProtoFieldName()))
		}
		indexPath := strings.Join(protoIndexes, "/")
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		typenameUpper, typenameSnake, typenameLower := opts.TypeName.UpperCamel, opts.TypeName.Snake, opts.TypeName.LowerCamel
		rpcQueryGet := protoutil.NewRPC(
			typenameUpper,
			fmt.Sprintf("QueryGet%sRequest", typenameUpper),
			fmt.Sprintf("QueryGet%sResponse", typenameUpper),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s",
						appModulePath, opts.ModuleName, typenameSnake, indexPath,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("Queries a %v by index.", typenameUpper))

		rpcQueryAll := protoutil.NewRPC(
			fmt.Sprintf("%sAll", typenameUpper),
			fmt.Sprintf("QueryAll%sRequest", typenameUpper),
			fmt.Sprintf("QueryAll%sResponse", typenameUpper),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s",
						appModulePath, opts.ModuleName, typenameSnake,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("Queries a list of %v items.", typenameUpper))
		protoutil.Append(serviceQuery, rpcQueryGet, rpcQueryAll)

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
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			// shouldn't really occur.
			return errors.Errorf("failed to add imports to %s: %w", path, err)
		}

		// Add the messages.
		var queryIndexFields []*proto.NormalField
		for i, index := range opts.Indexes {
			queryIndexFields = append(queryIndexFields, index.ToProtoField(i+1))
		}
		paginationType, paginationName := "cosmos.base.query.v1beta1.Page", "pagination"
		queryGetRequest := protoutil.NewMessage(
			fmt.Sprintf("QueryGet%sRequest", typenameUpper),
			protoutil.WithFields(queryIndexFields...),
		)
		gogoOption := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		queryGetResponse := protoutil.NewMessage(
			fmt.Sprintf("QueryGet%sResponse", typenameUpper),
			protoutil.WithFields(protoutil.NewField(typenameLower, typenameUpper, 1, protoutil.WithFieldOptions(gogoOption))),
		)
		queryAllRequest := protoutil.NewMessage(
			fmt.Sprintf("QueryAll%sRequest", typenameUpper),
			protoutil.WithFields(protoutil.NewField(paginationName, paginationType+"Request", 1)),
		)
		queryAllResponse := protoutil.NewMessage(
			fmt.Sprintf("QueryAll%sResponse", typenameUpper),
			protoutil.WithFields(
				protoutil.NewField(
					typenameLower,
					typenameUpper,
					1,
					protoutil.Repeated(),
					protoutil.WithFieldOptions(gogoOption),
				),
				protoutil.NewField(paginationName, fmt.Sprintf("%sResponse", paginationType), 2),
			),
		)
		protoutil.Append(protoFile, queryGetRequest, queryGetResponse, queryAllRequest, queryAllResponse)

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

		var positionalArgs string
		for _, field := range opts.Indexes {
			positionalArgs += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
		}

		template := `{
			RpcMethod: "%[2]vAll",
			Use: "list-%[3]v",
			Short: "List all %[4]v",
		},
		{
			RpcMethod: "%[2]v",
			Use: "show-%[3]v [id]",
			Short: "Shows a %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%s},
		},
		%[1]v`
		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLIQuery,
			opts.TypeName.UpperCamel,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
			strings.TrimSpace(positionalArgs),
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
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Get next available sequence number from GenesisState.
		genesisState, err := protoutil.GetMessageByName(protoFile, typed.ProtoGenesisStateMessage)
		if err != nil {
			return errors.Errorf("failed while looking up message '%s' in %s: %w", typed.ProtoGenesisStateMessage, path, err)
		}
		seqNumber := protoutil.NextUniqueID(genesisState)

		// Create new option and append to GenesisState message.
		typenameLower, typenameUpper := opts.TypeName.LowerCamel, opts.TypeName.UpperCamel
		gogoOption := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		typeListField := protoutil.NewField(
			typenameLower+"List", typenameUpper, seqNumber, protoutil.Repeated(), protoutil.WithFieldOptions(gogoOption),
		)
		protoutil.Append(genesisState, typeListField)

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
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/genesis.go")
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
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/genesis_test.go")
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

		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		if err = protoutil.AddImports(protoFile, true, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		// RPC service
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		// better to append them altogether, single traversal.
		typenameUpper := opts.TypeName.UpperCamel
		protoutil.Append(serviceMsg,
			protoutil.NewRPC(
				fmt.Sprintf("Create%s", typenameUpper),
				fmt.Sprintf("MsgCreate%s", typenameUpper),
				fmt.Sprintf("MsgCreate%sResponse", typenameUpper),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Update%s", typenameUpper),
				fmt.Sprintf("MsgUpdate%s", typenameUpper),
				fmt.Sprintf("MsgUpdate%sResponse", typenameUpper),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Delete%s", typenameUpper),
				fmt.Sprintf("MsgDelete%s", typenameUpper),
				fmt.Sprintf("MsgDelete%sResponse", typenameUpper),
			),
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
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		creator := protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1)
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.LowerCamel)
		commonFields := []*proto.NormalField{creator}
		commonFields = append(commonFields, indexes...)

		msgCreate := protoutil.NewMessage(
			"MsgCreate"+typenameUpper,
			protoutil.WithFields(append(commonFields, fields...)...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgCreateResponse := protoutil.NewMessage(fmt.Sprintf("MsgCreate%sResponse", typenameUpper))

		msgUpdate := protoutil.NewMessage(
			"MsgUpdate"+typenameUpper,
			protoutil.WithFields(append(commonFields, fields...)...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgUpdateResponse := protoutil.NewMessage(fmt.Sprintf("MsgUpdate%sResponse", typenameUpper))

		msgDelete := protoutil.NewMessage(
			"MsgDelete"+typenameUpper,
			protoutil.WithFields(commonFields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgDeleteResponse := protoutil.NewMessage(fmt.Sprintf("MsgDelete%sResponse", typenameUpper))
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

<<<<<<< HEAD
		var positionalArgs, positionalArgsStr string
		var indexes, indexesStr string
		for _, field := range opts.Fields {
			positionalArgs += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
			positionalArgsStr += fmt.Sprintf("[%s] ", field.ProtoFieldName())
		}
		for _, field := range opts.Indexes {
			indexes += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
			indexesStr += fmt.Sprintf("[%s] ", field.ProtoFieldName())
		}
		positionalArgs = indexes + positionalArgs
		positionalArgsStr = indexesStr + positionalArgsStr
=======
		index := fmt.Sprintf(`{ProtoField: "%s"}, `, opts.Index.ProtoFieldName())
		indexStr := fmt.Sprintf("[%s] ", opts.Index.ProtoFieldName())
		positionalArgs := index + opts.Fields.ProtoFieldNameAutoCLI()
		positionalArgsStr := indexStr + opts.Fields.CLIUsage()
>>>>>>> 7c532030 (feat(templates): set varargs in autocli config when needed (#4624))

		template := `{
			RpcMethod: "Create%[2]v",
			Use: "create-%[3]v %[6]s",
			Short: "Create a new %[4]v",
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
			Use: "delete-%[3]v %[8]s",
			Short: "Delete %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[7]s},
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
			strings.TrimSpace(indexes),
			strings.TrimSpace(indexesStr),
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
