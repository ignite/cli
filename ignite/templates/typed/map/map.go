package maptype

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
	"github.com/ignite/cli/v29/ignite/templates/typed"
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
	if opts.Index.DatatypeName != datatype.Bool {
		generateTest = true
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
	g.RunFn(keeperModify(opts))
	g.RunFn(clientCliQueryModify(replacer, opts))
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisModuleModify(opts))
	g.RunFn(genesisTestsModify(opts))
	g.RunFn(genesisTypesTestsModify(opts))

	// Modifications for new messages
	if !opts.NoMessage {
		g.RunFn(protoTxModify(opts))
		g.RunFn(clientCliTxModify(replacer, opts))
		g.RunFn(typesCodecModify(opts))

		if !opts.NoSimulation {
			g.RunFn(moduleSimulationModify(opts))
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

// keeperModify modifies the keeper to add a new collections map type.
func keeperModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/keeper.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := xast.ModifyStruct(
			f.String(),
			"Keeper",
			xast.AppendStructValue(
				opts.TypeName.UpperCamel,
				fmt.Sprintf("collections.Map[%[1]v, types.%[2]v]", opts.Index.DataType(), opts.TypeName.UpperCamel),
			),
		)
		if err != nil {
			return err
		}

		// add parameter to the struct into the new keeper method.
		content, err = xast.ModifyFunction(
			content,
			"NewKeeper",
			xast.AppendFuncStruct(
				"Keeper",
				opts.TypeName.UpperCamel,
				fmt.Sprintf(`collections.NewMap(sb, types.%[1]vKey, "%[2]v", %[3]v, codec.CollValue[types.%[1]v](cdc))`,
					opts.TypeName.UpperCamel,
					opts.TypeName.LowerCamel,
					opts.Index.CollectionsKeyValueType(),
				),
				-1,
			),
		)
		if err != nil {
			return err
		}

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
		// Add initial import for the new type
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		protoIndex := fmt.Sprintf("{%s}", opts.Index.ProtoFieldName())
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		typenameUpper, typenameSnake, typenameLower := opts.TypeName.UpperCamel, opts.TypeName.Snake, opts.TypeName.LowerCamel
		rpcQueryGet := protoutil.NewRPC(
			fmt.Sprintf("Get%s", typenameUpper),
			fmt.Sprintf("QueryGet%sRequest", typenameUpper),
			fmt.Sprintf("QueryGet%sResponse", typenameUpper),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s",
						appModulePath, opts.ModuleName, typenameSnake, protoIndex,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("Queries a %v by index.", typenameUpper))

		rpcQueryAll := protoutil.NewRPC(
			fmt.Sprintf("List%s", typenameUpper),
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
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			// shouldn't really occur.
			return errors.Errorf("failed to add imports to %s: %w", path, err)
		}

		// Add the messages.
		paginationType, paginationName := "cosmos.base.query.v1beta1.Page", "pagination"
		queryGetRequest := protoutil.NewMessage(
			fmt.Sprintf("QueryGet%sRequest", typenameUpper),
			protoutil.WithFields(opts.Index.ToProtoField(1)),
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

		template := `{
			RpcMethod: "List%[2]v",
			Use: "list-%[3]v",
			Short: "List all %[4]v",
		},
		{
			RpcMethod: "Get%[2]v",
			Use: "get-%[3]v [id]",
			Short: "Gets a %[4]v",
			Alias: []string{"show-%[3]v"},
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField:"%[5]s"}},
		},
		%[1]v`
		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLIQuery,
			opts.TypeName.UpperCamel,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
			opts.Index.ProtoFieldName(),
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

func genesisTypesModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(f.String(), xast.WithLastImport("fmt"))
		if err != nil {
			return err
		}

		content, err = xast.ModifyFunction(content, "DefaultGenesis", xast.AppendFuncStruct(
			"GenesisState",
			fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
			fmt.Sprintf("[]%[1]v{}", opts.TypeName.UpperCamel),
			-1,
		))
		if err != nil {
			return err
		}

		// lines of code to call the key function with the indexes of the element
		keyCall := fmt.Sprintf(`fmt.Sprint(elem.%s)`, opts.Index.Name.UpperCamel)
		templateTypesValidate := `// Check for duplicated index in %[1]v
%[1]vIndexMap := make(map[string]struct{})

for _, elem := range gs.%[2]vList {
	index := %[3]v
	if _, ok := %[1]vIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for %[1]v")
	}
	%[1]vIndexMap[index] = struct{}{}
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			keyCall,
		)
		content, err = xast.ModifyFunction(
			content,
			"Validate",
			xast.AppendFuncCode(replacementTypesValidate),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set all the %[1]v
for _, elem := range genState.%[2]vList {
	if err := k.%[2]v.Set(ctx, elem.%[3]v, elem); err != nil {
		return err
	}
}`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			opts.Index.Name.UpperCamel,
		)
		content, err := xast.ModifyFunction(
			f.String(),
			"InitGenesis",
			xast.AppendFuncCode(replacementModuleInit),
		)
		if err != nil {
			return err
		}

		templateModuleExport := `if err := k.%[1]v.Walk(ctx, nil, func(_ %[2]v, val types.%[1]v) (stop bool, err error) {
		genesis.%[1]vList = append(genesis.%[1]vList, val)
		return false, nil
	}); err != nil {
		return nil, err
	}`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			opts.TypeName.UpperCamel,
			opts.Index.DataType(),
		)
		content, err = xast.ModifyFunction(
			content,
			"ExportGenesis",
			xast.AppendFuncCode(replacementModuleExport),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTestsModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a list of two different indexes to use as sample
		sampleIndexes := make([]string, 2)
		for i := 0; i < 2; i++ {
			sampleIndexes[i] = opts.Index.GenesisArgs(i)
		}

		// add parameter to the struct into the new method.
		content, err := xast.ModifyFunction(
			f.String(),
			"TestGenesis",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf(
					"[]types.%[1]v{{ %[2]v }, { %[3]v }}",
					opts.TypeName.UpperCamel,
					sampleIndexes[0],
					sampleIndexes[1],
				),
				-1,
			),
			xast.AppendFuncCode(fmt.Sprintf("require.ElementsMatch(t, genesisState.%[1]vList, got.%[1]vList)", opts.TypeName.UpperCamel)),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesTestsModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a list of two different indexes to use as sample
		sampleIndexes := make([]string, 2)
		for i := 0; i < 2; i++ {
			sampleIndexes[i] = opts.Index.GenesisArgs(i)
		}

		templateDuplicated := `{
	desc:     "duplicated %[1]v",
	genState: &types.GenesisState{
		%[2]vList: []types.%[2]v{
			{
				%[3]v},
			{
				%[3]v},
		},
	},
	valid:    false,
}`
		replacementDuplicated := fmt.Sprintf(
			templateDuplicated,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			sampleIndexes[0],
		)

		// add parameter to the struct into the new method.
		content, err := xast.ModifyFunction(
			f.String(),
			"TestGenesisState_Validate",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf(
					"[]types.%[1]v{{ %[2]v }, { %[3]v }}",
					opts.TypeName.UpperCamel,
					sampleIndexes[0],
					sampleIndexes[1],
				),
				-1,
			),
			xast.AppendFuncTestCase(replacementDuplicated),
		)
		if err != nil {
			return err
		}

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
		index := opts.Index.ToProtoField(2)
		var fields []*proto.NormalField
		for i, f := range opts.Fields {
			fields = append(fields, f.ToProtoField(i+3)) // +3 because of the index
		}

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range append(opts.Fields.ProtoImports(), opts.Index.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(protoFile, false, protoImports...); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		creator := protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1)
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.LowerCamel)
		commonFields := []*proto.NormalField{creator}
		commonFields = append(commonFields, index)

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

		index := fmt.Sprintf(`{ProtoField: "%s"}, `, opts.Index.ProtoFieldName())
		indexStr := fmt.Sprintf("[%s] ", opts.Index.ProtoFieldName())
		var positionalArgs, positionalArgsStr string
		for _, field := range opts.Fields {
			positionalArgs += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
			positionalArgsStr += fmt.Sprintf("[%s] ", field.ProtoFieldName())
		}

		positionalArgs = index + positionalArgs
		positionalArgsStr = indexStr + positionalArgsStr

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
			strings.TrimSpace(index),
			strings.TrimSpace(indexStr),
		)

		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLITx, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(f.String(), xast.WithLastNamedImport("sdk", "github.com/cosmos/cosmos-sdk/types"))
		if err != nil {
			return err
		}

		// Interface
		templateInterface := `registrar.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[1]v{},
	&MsgUpdate%[1]v{},
	&MsgDelete%[1]v{},
)`
		replacementInterface := fmt.Sprintf(
			templateInterface,
			opts.TypeName.UpperCamel,
		)
		content, err = xast.ModifyFunction(
			content,
			"RegisterInterfaces",
			xast.AppendFuncAtLine(replacementInterface, 0),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
