package list

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/typed"
)

var (
	//go:embed files/component/* files/component/**/*
	fsComponent embed.FS

	//go:embed files/messages/* files/messages/**/*
	fsMessages embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsSimapp embed.FS
)

// NewGenerator returns the generator to scaffold a new type in a module
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
			fsSimapp,
			"files/simapp/",
			opts.AppPath,
		)
	)

	g.RunFn(protoQueryModify(opts))
	g.RunFn(typesKeyModify(opts))
	g.RunFn(clientCliQueryModify(replacer, opts))

	// Genesis modifications
	genesisModify(replacer, opts, g)

	if !opts.NoMessage {
		// Modifications for new messages
		g.RunFn(protoTxModify(opts))
		g.RunFn(typesCodecModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))

		if !opts.NoSimulation {
			g.RunFn(moduleSimulationModify(replacer, opts))
			if err := typed.Box(simappTemplate, opts, g); err != nil {
				return nil, err
			}
		}

		// Messages template
		if err := typed.Box(messagesTemplate, opts, g); err != nil {
			return nil, err
		}
	}

	g.RunFn(frontendSrcStoreAppModify(replacer, opts))

	return g, typed.Box(componentTemplate, opts, g)
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
		// Import
		if err = protoutil.AddImports(pf, true, opts.TypeImport()); err != nil {
			return fmt.Errorf("failed while adding imports to %s: %w", path, err)
		}

		// RPC service
		s, err := protoutil.GetServiceByName(pf, "Msg")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		// better to append them altogether, single traversal.
		name := opts.TypeName.UpperCamel
		create := protoutil.NewRPC("Create"+name, "MsgCreate"+name, "MsgCreate"+name+"Response")
		update := protoutil.NewRPC("Update"+name, "MsgUpdate"+name, "MsgUpdate"+name+"Response")
		delete := protoutil.NewRPC("Delete"+name, "MsgDelete"+name, "MsgDelete"+name+"Response")
		protoutil.Append(s, create, update, delete)

		// - Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v.proto", opts.AppName, opts.ModuleName, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		// we already know an import exists, pass false for fallback.
		if err = protoutil.AddImports(pf, true, protoImports...); err != nil {
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Messages
		creator := protoutil.NewField("string", opts.MsgSigner.LowerCamel, 1)
		createFields := []*proto.NormalField{creator}
		for i, field := range opts.Fields {
			createFields = append(createFields, field.ToProtoField(i+2))
		}
		udfields := []*proto.NormalField{creator, protoutil.NewField("uint64", "id", 2)}
		updateFields := udfields
		for i, field := range opts.Fields {
			updateFields = append(updateFields, field.ToProtoField(i+3))
		}

		msgCreate := protoutil.NewMessage("MsgCreate"+name, protoutil.WithFields(createFields...))
		msgCreateResp := protoutil.NewMessage(
			"MsgCreate"+name+"Response",
			protoutil.WithFields(protoutil.NewField("uint64", "id", 1)),
		)
		msgUpdate := protoutil.NewMessage("MsgUpdate"+name, protoutil.WithFields(updateFields...))
		msgUpdateResp := protoutil.NewMessage("MsgUpdate" + name + "Response")
		msgDelete := protoutil.NewMessage("MsgDelete"+name, protoutil.WithFields(udfields...))
		msgDeleteResp := protoutil.NewMessage("MsgDelete" + name + "Response")
		protoutil.Append(pf,
			msgCreate, msgCreateResp, msgUpdate, msgUpdateResp, msgDelete, msgDeleteResp,
		)
		newFile := genny.NewFileS(path, protoutil.Printer(pf))

		return r.File(newFile)
	}
}

// Modifies query.proto to add the required RPCs and Messages.
//
// What it depends on:
//   - Existence of a service with name "Query". Adds the rpc's there.
func protoQueryModify(opts *typed.Options) genny.RunFn {
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
		// Imports for the new type and gogoproto.
		gogoproto := protoutil.NewImport("gogoproto/gogo.proto")
		if err = protoutil.AddImports(pf, true, gogoproto, opts.TypeImport()); err != nil {
			return fmt.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Add to Query:
		srv, err := protoutil.GetServiceByName(pf, "Query")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}
		appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)
		typ := opts.TypeName.UpperCamel
		single := protoutil.NewRPC(typ, "QueryGet"+typ+"Request", "QueryGet"+typ+"Response",
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/{id}",
						appModulePath, opts.ModuleName, opts.TypeName.Snake,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		all := protoutil.NewRPC(typ+"All", "QueryAll"+typ+"Request", "QueryAll"+typ+"Response",
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
		protoutil.Append(srv, single, all)

		// Add messages
		pagT, pagN := "cosmos.base.query.v1beta1.Page", "pagination"
		gogoOpt := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())

		queryReq := protoutil.NewMessage(
			"QueryGet"+typ+"Request",
			protoutil.WithFields(protoutil.NewField("uint64", "id", 1)),
		)
		field := protoutil.NewField(typ, typ, 1, protoutil.WithFieldOptions(gogoOpt))
		queryResp := protoutil.NewMessage("QueryGet"+typ+"Response", protoutil.WithFields(field))

		queryAllReq := protoutil.NewMessage(
			"QueryAll"+typ+"Request",
			protoutil.WithFields(protoutil.NewField(pagT+"Request", pagN, 1)),
		)
		field = protoutil.NewField(typ, typ, 1, protoutil.Repeated(), protoutil.WithFieldOptions(gogoOpt))
		queryAllResp := protoutil.NewMessage(
			"QueryAll"+typ+"Response",
			protoutil.WithFields(field, protoutil.NewField(pagT+"Response", pagN, 2)),
		)

		protoutil.Append(pf, queryReq, queryResp, queryAllReq, queryAllResp)

		newFile := genny.NewFileS(path, protoutil.Printer(pf))
		return r.File(newFile)
	}
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
	%[1]vCountKey= "%[1]v/count/"
)
`, opts.TypeName.UpperCamel)
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

		// Import
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), typed.Placeholder, replacementImport)

		// Concrete
		templateConcrete := `cdc.RegisterConcrete(&MsgCreate%[2]v{}, "%[3]v/Create%[2]v", nil)
cdc.RegisterConcrete(&MsgUpdate%[2]v{}, "%[3]v/Update%[2]v", nil)
cdc.RegisterConcrete(&MsgDelete%[2]v{}, "%[3]v/Delete%[2]v", nil)
%[1]v`
		replacementConcrete := fmt.Sprintf(templateConcrete, typed.Placeholder2, opts.TypeName.UpperCamel, opts.ModuleName)
		content = replacer.Replace(content, typed.Placeholder2, replacementConcrete)

		// Interface
		templateInterface := `registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgCreate%[2]v{},
	&MsgUpdate%[2]v{},
	&MsgDelete%[2]v{},
)
%[1]v`
		replacementInterface := fmt.Sprintf(templateInterface, typed.Placeholder3, opts.TypeName.UpperCamel)
		content = replacer.Replace(content, typed.Placeholder3, replacementInterface)

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

func frontendSrcStoreAppModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "vue/src/views/Types.vue")
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
