package typed

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
)

// New ...
func NewStargate(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(typesCodecModify(opts))
	g.RunFn(typesCodecImportModify(opts))
	g.RunFn(typesCodecInterfaceModify(opts))
	g.RunFn(clientCliTxModify(opts))

	return g, Box(stargateTemplate, opts, g)
}

func handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := strings.Replace(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer, 1)

		templateHandlers := `%[1]v
		case *types.Msg%[2]v:
					res, err := msgServer.%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			Placeholder,
			strings.Title(opts.MsgName),
		)
		content = strings.Replace(content, Placeholder, replacementHandlers, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxRPCModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
  rpc %[2]v(Msg%[2]v) returns (Msg%[2]vResponse);`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxRPC,
			strings.Title(opts.MsgName),
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxRPC, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxMessageModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		var msgFields string
		for i, field := range opts.Fields {
			msgFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+2)
		}
		var resFields string
		for i, field := range opts.ResFields {
			resFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+1)
		}

		template := `%[1]v
message Msg%[2]v {
  string creator = 1;
%[3]v}

message Msg%[2]vResponse {
%[4]v}
`
		replacement := fmt.Sprintf(template, PlaceholderProtoTxMessage,
			strings.Title(opts.MsgName),
			msgFields,
			resFields,
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxMessage, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecImportModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
cdc.RegisterConcrete(&Msg%[2]v{}, "%[3]v/%[2]v", nil)
`
		replacement := fmt.Sprintf(template, Placeholder2, strings.Title(opts.MsgName), opts.ModuleName)
		content := strings.Replace(f.String(), Placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typesCodecInterfaceModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]v{},
)`
		replacement := fmt.Sprintf(template, Placeholder3, strings.Title(opts.MsgName))
		content := strings.Replace(f.String(), Placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(Cmd%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder, strings.Title(opts.MsgName))
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

