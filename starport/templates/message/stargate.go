package message

import (
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

// NewStargate returns the generator to scaffold a empty message in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(handlerModify(replacer, opts))
	g.RunFn(protoTxRPCModify(replacer, opts))
	g.RunFn(protoTxMessageModify(replacer, opts))
	g.RunFn(typesCodecModify(replacer, opts))
	g.RunFn(clientCliTxModify(replacer, opts))

	return g, Box(stargateTemplate, opts, g)
}

func handlerModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `%[1]v
		case *types.Msg%[2]v:
					res, err := msgServer.%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
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
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
  rpc %[2]v(Msg%[2]v) returns (Msg%[2]vResponse);`
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
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var msgFields string
		for i, field := range opts.Fields {
			msgFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+2)
		}
		var resFields string
		for i, field := range opts.ResFields {
			resFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+1)
		}

		template := `%[1]v
message Msg%[2]v {
  string creator = 1;
%[3]v}

message Msg%[2]vResponse {
%[4]v}
`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxMessage,
			opts.MsgName.UpperCamel,
			msgFields,
			resFields,
		)
		content := replacer.Replace(f.String(), PlaceholderProtoTxMessage, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacementImport := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), Placeholder, replacementImport)

		templateRegisterConcrete := `%[1]v
cdc.RegisterConcrete(&Msg%[2]v{}, "%[3]v/%[2]v", nil)
`
		replacementRegisterConcrete := fmt.Sprintf(
			templateRegisterConcrete,
			Placeholder2,
			opts.MsgName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, Placeholder2, replacementRegisterConcrete)

		templateRegisterImplementations := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]v{},
)`
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
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(Cmd%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder, opts.MsgName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
