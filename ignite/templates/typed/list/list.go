package list

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

var (
	//go:embed files/component/* files/component/**/*
	fsComponent embed.FS

	//go:embed files/messages/* files/messages/**/*
	fsMessages embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsSimapp embed.FS
)

// NewGenerator returns the generator to scaffold a new type in a module.
func NewGenerator(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	subMessages, err := fs.Sub(fsMessages, "files/messages")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}
	subComponent, err := fs.Sub(fsComponent, "files/component")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}
	subSimapp, err := fs.Sub(fsSimapp, "files/simapp")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	g.RunFn(protoQueryModify(opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn(keeperModify(opts))
	g.RunFn(clientCliQueryModify(replacer, opts))

	// Genesis modifications
	genesisModify(opts, g)

	if !opts.NoMessage {
		// Modifications for new messages
		g.RunFn(protoTxModify(opts))
		g.RunFn(typesCodecModify(opts))
		g.RunFn(clientCliTxModify(replacer, opts))

		if !opts.NoSimulation {
			g.RunFn(moduleSimulationModify(opts))
			if err := typed.Box(subSimapp, opts, g); err != nil {
				return nil, err
			}
		}

		// Messages template
		if err := typed.Box(subMessages, opts, g); err != nil {
			return nil, err
		}
	}

	g.RunFn(frontendSrcStoreAppModify(replacer, opts))

	return g, typed.Box(subComponent, opts, g)
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

		// RPC service
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		// Create, update, delete rpcs. Better to append them altogether, single traversal.
		typenamePascal := opts.TypeName.PascalCase
		protoutil.Append(serviceMsg,
			protoutil.NewRPC(
				fmt.Sprintf("Create%s", typenamePascal),
				fmt.Sprintf("MsgCreate%s", typenamePascal),
				fmt.Sprintf("MsgCreate%sResponse", typenamePascal),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Update%s", typenamePascal),
				fmt.Sprintf("MsgUpdate%s", typenamePascal),
				fmt.Sprintf("MsgUpdate%sResponse", typenamePascal),
			),
			protoutil.NewRPC(
				fmt.Sprintf("Delete%s", typenamePascal),
				fmt.Sprintf("MsgDelete%s", typenamePascal),
				fmt.Sprintf("MsgDelete%sResponse", typenamePascal),
			),
		)

		// - Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(protoFile, true, protoImports...); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Messages
		creator := protoutil.NewField(opts.MsgSigner.Snake, "string", 1)
		creator.Options = append(creator.Options, protoutil.NewOption("cosmos_proto.scalar", "cosmos.AddressString", protoutil.Custom())) // set the scalar annotation
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.Snake)
		createFields := []*proto.NormalField{creator}
		for i, field := range opts.Fields {
			createFields = append(createFields, field.ToProtoField(i+2))
		}
		udfields := []*proto.NormalField{creator, protoutil.NewField("id", "uint64", 2)}
		updateFields := udfields
		for i, field := range opts.Fields {
			updateFields = append(updateFields, field.ToProtoField(i+3))
		}

		msgCreate := protoutil.NewMessage(
			fmt.Sprintf("MsgCreate%s", typenamePascal),
			protoutil.WithFields(createFields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgCreateResponse := protoutil.NewMessage(
			fmt.Sprintf("MsgCreate%sResponse", typenamePascal),
			protoutil.WithFields(protoutil.NewField("id", "uint64", 1)),
		)
		msgUpdate := protoutil.NewMessage(
			fmt.Sprintf("MsgUpdate%s", typenamePascal),
			protoutil.WithFields(updateFields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgUpdateResponse := protoutil.NewMessage(fmt.Sprintf("MsgUpdate%sResponse", typenamePascal))
		msgDelete := protoutil.NewMessage(
			fmt.Sprintf("MsgDelete%s", typenamePascal),
			protoutil.WithFields(udfields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgDeleteResponse := protoutil.NewMessage(fmt.Sprintf("MsgDelete%sResponse", typenamePascal))
		protoutil.Append(
			protoFile,
			msgCreate,
			msgCreateResponse,
			msgUpdate,
			msgUpdateResponse,
			msgDelete,
			msgDeleteResponse,
		)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

// Modifies query.proto to add the required RPCs and Messages.
//
// What it depends on:
//   - Existence of a service with name "Query". Adds the rpc's there.
func protoQueryModify(opts *typed.Options) genny.RunFn {
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
		// Imports for the new type and gogoImport.
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}

		// Add to Query:
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		typenamePascal, typenameSnake := opts.TypeName.PascalCase, opts.TypeName.Snake
		rpcQueryGet := protoutil.NewRPC(
			fmt.Sprintf("Get%s", typenamePascal),
			fmt.Sprintf("QueryGet%sRequest", typenamePascal),
			fmt.Sprintf("QueryGet%sResponse", typenamePascal),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s/{id}",
						appModulePath, opts.ModuleName, opts.ProtoVer, opts.TypeName.Snake,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("Get%[1]v Queries a %[1]v by id.", typenamePascal))

		rpcQueryAll := protoutil.NewRPC(
			fmt.Sprintf("List%s", typenamePascal),
			fmt.Sprintf("QueryAll%sRequest", typenamePascal),
			fmt.Sprintf("QueryAll%sResponse", typenamePascal),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s",
						appModulePath, opts.ModuleName, opts.ProtoVer, opts.TypeName.Snake,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcQueryGet, fmt.Sprintf("List%[1]v Queries a list of %[1]v items.", typenamePascal))
		protoutil.Append(serviceQuery, rpcQueryGet, rpcQueryAll)

		// Add messages
		paginationType, paginationName := "cosmos.base.query.v1beta1.Page", "pagination"
		gogoOption := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())

		queryGetRequest := protoutil.NewMessage(
			fmt.Sprintf("QueryGet%sRequest", typenamePascal),
			protoutil.WithFields(protoutil.NewField("id", "uint64", 1)),
		)
		field := protoutil.NewField(typenameSnake, typenamePascal, 1, protoutil.WithFieldOptions(gogoOption))
		queryGetResponse := protoutil.NewMessage(
			fmt.Sprintf("QueryGet%sResponse", typenamePascal),
			protoutil.WithFields(field))

		queryAllRequest := protoutil.NewMessage(
			fmt.Sprintf("QueryAll%sRequest", typenamePascal),
			protoutil.WithFields(protoutil.NewField(paginationName, paginationType+"Request", 1)),
		)
		field = protoutil.NewField(typenameSnake, typenamePascal, 1, protoutil.Repeated(), protoutil.WithFieldOptions(gogoOption))
		queryAllResponse := protoutil.NewMessage(
			fmt.Sprintf("QueryAll%sResponse", typenamePascal),
			protoutil.WithFields(field, protoutil.NewField(paginationName, paginationType+"Response", 2)),
		)
		protoutil.Append(protoFile, queryGetRequest, queryGetResponse, queryAllRequest, queryAllResponse)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

// typesKeyModify modifies the keys.go file to add a new collection prefix.
func typesKeyModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/keys.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
var (
	%[1]vKey= collections.NewPrefix("%[2]v/value/")
	%[1]vCountKey= collections.NewPrefix("%[2]v/count/")
)
`,
			opts.TypeName.PascalCase,
			opts.TypeName.LowerCase,
		)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// keeperModify modifies the keeper to add a new collections item type.
func keeperModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "keeper/keeper.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := xast.ModifyStruct(
			f.String(),
			"Keeper",
			xast.AppendStructValue(
				fmt.Sprintf("%[1]vSeq", opts.TypeName.UpperCamel),
				"collections.Sequence",
			),
			xast.AppendStructValue(
				opts.TypeName.UpperCamel,
				fmt.Sprintf("collections.Map[uint64, types.%[1]v]", opts.TypeName.PascalCase),
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
				fmt.Sprintf(`collections.NewMap(sb, types.%[1]vKey, "%[2]v", collections.Uint64Key, codec.CollValue[types.%[1]v](cdc))`,
					opts.TypeName.PascalCase,
					opts.TypeName.LowerCamel,
				),
			),
			xast.AppendFuncStruct(
				"Keeper",
				fmt.Sprintf("%[1]vSeq", opts.TypeName.UpperCamel),
				fmt.Sprintf(`collections.NewSequence(sb, types.%[2]vCountKey, "%[3]vSequence")`,
					opts.TypeName.UpperCamel,
					opts.TypeName.PascalCase,
					opts.TypeName.LowerCamel,
				),
			),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(f.String(), xast.WithNamedImport("sdk", "github.com/cosmos/cosmos-sdk/types"))
		if err != nil {
			return err
		}

		// Interface
		templateInterface := `registrar.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[1]v{},
	&MsgUpdate%[1]v{},
	&MsgDelete%[1]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, opts.TypeName.PascalCase)
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

func clientCliTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "module/autocli.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `{
			RpcMethod: "Create%[2]v",
			Use: "create-%[3]v %[6]s",
			Short: "Create %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[5]s},
		},
		{
			RpcMethod: "Update%[2]v",
			Use: "update-%[3]v [id] %[6]s",
			Short: "Update %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, %[5]s},
		},
		{
			RpcMethod: "Delete%[2]v",
			Use: "delete-%[3]v [id]",
			Short: "Delete %[4]v",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
		},
		%[1]v`

		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLITx,
			opts.TypeName.PascalCase,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
			opts.Fields.ProtoFieldNameAutoCLI(),
			opts.Fields.CLIUsage(),
		)

		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLITx, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "module/autocli.go")
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
			Short: "Gets a %[4]v by id",
			Alias: []string{"show-%[3]v"},
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
		},
		%[1]v`
		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLIQuery,
			opts.TypeName.PascalCase,
			opts.TypeName.Kebab,
			opts.TypeName.Original,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLIQuery, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func frontendSrcStoreAppModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := chain.DefaultVueTypesPath
		f, err := r.Disk.Find(path)
		if os.IsNotExist(err) {
			// Skip modification if the app doesn't contain front-end
			return nil
		}
		if err != nil {
			return err
		}
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		replacement := fmt.Sprintf(`%[1]v
		<SpType modulePath="%[2]v.%[3]v" moduleType="%[4]v"  />`,
			typed.Placeholder4,
			strings.ReplaceAll(appModulePath, "/", "."),
			opts.ModuleName,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.Placeholder4, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
