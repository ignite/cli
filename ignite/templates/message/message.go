package message

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/testutil"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

var (
	//go:embed files/message/* files/message/**/*
	fsMessage embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsSimapp embed.FS
)

func Box(box packd.Walker, opts *Options, g *genny.Generator) error {
	if err := g.Box(box); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("ProtoVer", opts.ProtoVer)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("MsgName", opts.MsgName)
	ctx.Set("MsgDesc", opts.MsgDesc)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("ResFields", opts.ResFields)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))
	g.Transformer(genny.Replace("{{msgName}}", opts.MsgName.Snake))

	// Create the 'testutil' package with the test helpers
	return testutil.Register(g, opts.AppPath)
}

// NewGenerator returns the generator to scaffold a empty message in a module.
func NewGenerator(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(protoTxRPCModify(opts))
	g.RunFn(protoTxMessageModify(opts))
	g.RunFn(typesCodecModify(opts))
	g.RunFn(clientCliTxModify(replacer, opts))

	template := xgenny.NewEmbedWalker(
		fsMessage,
		"files/message",
		opts.AppPath,
	)

	if !opts.NoSimulation {
		g.RunFn(moduleSimulationModify(opts))
		simappTemplate := xgenny.NewEmbedWalker(
			fsSimapp,
			"files/simapp",
			opts.AppPath,
		)
		if err := Box(simappTemplate, opts, g); err != nil {
			return nil, err
		}
	}
	return g, Box(template, opts, g)
}

// protoTxRPCModify modifies the tx.proto file to add the required RPCs and messages.
//
// What it expects:
//   - A service named "Msg" to exist in the proto file, it appends the RPCs inside it.
func protoTxRPCModify(opts *Options) genny.RunFn {
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
		// = Add new rpc to Msg.
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		typenameUpper := opts.MsgName.UpperCamel
		protoutil.Append(
			serviceMsg,
			protoutil.NewRPC(
				typenameUpper,
				fmt.Sprintf("Msg%s", typenameUpper),
				fmt.Sprintf("Msg%sResponse", typenameUpper),
			),
		)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func protoTxMessageModify(opts *Options) genny.RunFn {
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
		// Prepare the fields and create the messages.
		creator := protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1)
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.LowerCamel)
		msgFields := []*proto.NormalField{creator}
		for i, field := range opts.Fields {
			msgFields = append(msgFields, field.ToProtoField(i+2))
		}
		var resFields []*proto.NormalField
		for i, field := range opts.ResFields {
			resFields = append(resFields, field.ToProtoField(i+1))
		}

		typenameUpper := opts.MsgName.UpperCamel
		msg := protoutil.NewMessage(
			"Msg"+typenameUpper,
			protoutil.WithFields(msgFields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgResp := protoutil.NewMessage("Msg"+typenameUpper+"Response", protoutil.WithFields(resFields...))
		protoutil.Append(protoFile, msg, msgResp)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range append(opts.ResFields.ProtoImports(), opts.Fields.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range append(opts.ResFields.Custom(), opts.Fields.Custom()...) {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		if err = protoutil.AddImports(protoFile, true, protoImports...); err != nil {
			return errors.Errorf("failed to add imports to %s: %w", path, err)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func typesCodecModify(opts *Options) genny.RunFn {
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

		templateRegisterImplementations := `registrar.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[1]v{},
)`
		replacementRegisterImplementations := fmt.Sprintf(
			templateRegisterImplementations,
			opts.MsgName.UpperCamel,
		)

		content, err = xast.ModifyFunction(
			content,
			"RegisterInterfaces",
			xast.AppendFuncAtLine(replacementRegisterImplementations, 0),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/autocli.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var positionalArgs string
		for _, field := range opts.Fields {
			positionalArgs += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
		}

		template := `{
			RpcMethod: "%[2]v",
			Use: "%[3]v",
			Short: "Send a %[4]v tx",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[5]s},
		},
		%[1]v`

		replacement := fmt.Sprintf(
			template,
			typed.PlaceholderAutoCLITx,
			opts.MsgName.UpperCamel,
			strings.TrimSpace(fmt.Sprintf("%s%s", opts.MsgName.Kebab, opts.Fields.String())),
			opts.MsgName.Original,
			strings.TrimSpace(positionalArgs),
		)

		content := replacer.Replace(f.String(), typed.PlaceholderAutoCLITx, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func moduleSimulationModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := typed.ModuleSimulationMsgModify(
			f.String(),
			opts.ModulePath,
			opts.ModuleName,
			opts.MsgName,
			opts.MsgSigner,
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
