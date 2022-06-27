package message

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/typed"
)

// NewStargate returns the generator to scaffold a empty message in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(handlerModify(replacer, opts))
	g.RunFn(protoTxRPCModify(replacer, opts))
	g.RunFn(protoTxMessageModify(replacer, opts))
	g.RunFn(typesCodecModify(replacer, opts))
	g.RunFn(clientCliTxModify(replacer, opts))

	template := xgenny.NewEmbedWalker(
		fsStargateMessage,
		"stargate/message",
		opts.AppPath,
	)

	if !opts.NoSimulation {
		g.RunFn(moduleSimulationModify(replacer, opts))
		simappTemplate := xgenny.NewEmbedWalker(
			fsStargateSimapp,
			"stargate/simapp",
			opts.AppPath,
		)
		if err := Box(simappTemplate, opts, g); err != nil {
			return nil, err
		}
	}
	return g, Box(template, opts, g)
}

func handlerModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "handler.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `case *types.Msg%[2]v:
					res, err := msgServer.%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
%[1]v`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			Placeholder,
			opts.MsgName.UpperCamel,
		)
		content = replacer.Replace(content, Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxRPCModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `  rpc %[2]v(Msg%[2]v) returns (Msg%[2]vResponse);
%[1]v`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxRPC,
			opts.MsgName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderProtoTxRPC, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxMessageModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var msgFields string
		for i, field := range opts.Fields {
			msgFields += fmt.Sprintf("  %s;\n", field.ProtoType(i+2))
		}
		var resFields string
		for i, field := range opts.ResFields {
			resFields += fmt.Sprintf("  %s;\n", field.ProtoType(i+1))
		}

		template := `message Msg%[2]v {
  string %[5]v = 1;
%[3]v}

message Msg%[2]vResponse {
%[4]v}

%[1]v`
		replacement := fmt.Sprintf(template,
			PlaceholderProtoTxMessage,
			opts.MsgName.UpperCamel,
			msgFields,
			resFields,
			opts.MsgSigner.LowerCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderProtoTxMessage, replacement)

		// Ensure custom types are imported
		protoImports := append(opts.ResFields.ProtoImports(), opts.Fields.ProtoImports()...)
		customFields := append(opts.ResFields.Custom(), opts.Fields.Custom()...)
		for _, f := range customFields {
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

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), Placeholder, replacementImport)

		templateRegisterConcrete := `cdc.RegisterConcrete(&Msg%[2]v{}, "%[3]v/%[2]v", nil)
%[1]v`
		replacementRegisterConcrete := fmt.Sprintf(
			templateRegisterConcrete,
			Placeholder2,
			opts.MsgName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, Placeholder2, replacementRegisterConcrete)

		templateRegisterImplementations := `registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]v{},
)
%[1]v`
		replacementRegisterImplementations := fmt.Sprintf(
			templateRegisterImplementations,
			Placeholder3,
			opts.MsgName.UpperCamel,
		)
		content = replacer.Replace(content, Placeholder3, replacementRegisterImplementations)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/tx.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `cmd.AddCommand(Cmd%[2]v())
%[1]v`
		replacement := fmt.Sprintf(template, Placeholder, opts.MsgName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func moduleSimulationModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module_simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := typed.ModuleSimulationMsgModify(
			replacer,
			f.String(),
			opts.ModuleName,
			opts.MsgName,
		)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
