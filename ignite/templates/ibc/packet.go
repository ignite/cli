package ibc

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/testutil"
)

var (
	//go:embed packet/component/* packet/component/**/*
	fsPacketComponent embed.FS

	//go:embed packet/messages/* packet/messages/**/*
	fsPacketMessages embed.FS
)

// PacketOptions are options to scaffold a packet in a IBC module
type PacketOptions struct {
	AppName    string
	AppPath    string
	ModuleName string
	ModulePath string
	PacketName multiformatname.Name
	MsgSigner  multiformatname.Name
	Fields     field.Fields
	AckFields  field.Fields
	NoMessage  bool
}

// NewPacket returns the generator to scaffold a packet in an IBC module
func NewPacket(replacer placeholder.Replacer, opts *PacketOptions) (*genny.Generator, error) {
	var (
		g = genny.New()

		messagesTemplate = xgenny.NewEmbedWalker(
			fsPacketMessages,
			"packet/messages/",
			opts.AppPath,
		)
		componentTemplate = xgenny.NewEmbedWalker(
			fsPacketComponent,
			"packet/component/",
			opts.AppPath,
		)
	)

	// Add the component
	g.RunFn(moduleModify(replacer, opts))
	g.RunFn(protoModify(replacer, opts))
	g.RunFn(eventModify(replacer, opts))
	if err := g.Box(componentTemplate); err != nil {
		return g, err
	}

	// Add the send message
	if !opts.NoMessage {
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(handlerTxModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))
		g.RunFn(codecModify(replacer, opts))
		if err := g.Box(messagesTemplate); err != nil {
			return g, err
		}
	}

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("packetName", opts.PacketName)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("fields", opts.Fields)
	ctx.Set("ackFields", opts.AckFields)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{packetName}}", opts.PacketName.Snake))

	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, opts.AppPath); err != nil {
		return g, err
	}

	return g, nil
}

func moduleModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module_ibc.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `case *types.%[2]vPacketData_%[3]vPacket:
	packetAck, err := am.keeper.OnRecv%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err.Error())
	} else {
		// Encode packet acknowledgment
		packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error()).Error())
		}
		ack = channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType%[3]vPacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%%t", err != nil)),
		),
	)
%[1]v`
		replacementRecv := fmt.Sprintf(
			templateRecv,
			PlaceholderIBCPacketModuleRecv,
			xstrings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderIBCPacketModuleRecv, replacementRecv)

		// Ack packet dispatch
		templateAck := `case *types.%[2]vPacketData_%[3]vPacket:
	err := am.keeper.OnAcknowledgement%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket, ack)
	if err != nil {
		return err
	}
	eventType = types.EventType%[3]vPacket
%[1]v`
		replacementAck := fmt.Sprintf(
			templateAck,
			PlaceholderIBCPacketModuleAck,
			xstrings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketModuleAck, replacementAck)

		// Timeout packet dispatch
		templateTimeout := `case *types.%[2]vPacketData_%[3]vPacket:
	err := am.keeper.OnTimeout%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		return err
	}
%[1]v`
		replacementTimeout := fmt.Sprintf(
			templateTimeout,
			PlaceholderIBCPacketModuleTimeout,
			xstrings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketModuleTimeout, replacementTimeout)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "packet.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := f.String()

		// Add the field in the module packet
		fieldCount := strings.Count(content, PlaceholderIBCPacketProtoFieldNumber)
		templateField := `%[1]v
				%[2]vPacketData %[3]vPacket = %[4]v; %[5]v`
		replacementField := fmt.Sprintf(
			templateField,
			PlaceholderIBCPacketProtoField,
			opts.PacketName.UpperCamel,
			opts.PacketName.LowerCamel,
			fieldCount+2,
			PlaceholderIBCPacketProtoFieldNumber,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketProtoField, replacementField)

		// Add the message definition for packet and acknowledgment
		var packetFields string
		for i, field := range opts.Fields {
			packetFields += fmt.Sprintf("  %s;\n", field.ProtoType(i+1))
		}

		var ackFields string
		for i, field := range opts.AckFields {
			ackFields += fmt.Sprintf("  %s;\n", field.ProtoType(i+1))
		}

		// Ensure custom types are imported
		protoImports := append(opts.Fields.ProtoImports(), opts.AckFields.ProtoImports()...)
		customFields := append(opts.Fields.Custom(), opts.AckFields.Custom()...)
		for _, f := range customFields {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", PlaceholderProtoPacketImport, importModule)
			content = replacer.Replace(content, PlaceholderProtoPacketImport, replacementImport)
		}

		templateMessage := `// %[2]vPacketData defines a struct for the packet payload
message %[2]vPacketData {
%[3]v}

// %[2]vPacketAck defines a struct for the packet acknowledgment
message %[2]vPacketAck {
	%[4]v}
%[1]v`
		replacementMessage := fmt.Sprintf(
			templateMessage,
			PlaceholderIBCPacketProtoMessage,
			opts.PacketName.UpperCamel,
			packetFields,
			ackFields,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketProtoMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func eventModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/events_ibc.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `EventType%[2]vPacket       = "%[3]v_packet"
%[1]v`
		replacement := fmt.Sprintf(
			template,
			PlaceholderIBCPacketEvent,
			opts.PacketName.UpperCamel,
			opts.PacketName.LowerCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderIBCPacketEvent, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// RPC
		templateRPC := `  rpc Send%[2]v(MsgSend%[2]v) returns (MsgSend%[2]vResponse);
%[1]v`
		replacementRPC := fmt.Sprintf(
			templateRPC,
			PlaceholderProtoTxRPC,
			opts.PacketName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderProtoTxRPC, replacementRPC)

		var sendFields string
		for i, field := range opts.Fields {
			sendFields += fmt.Sprintf("  %s;\n", field.ProtoType(i+5))
		}

		// Ensure custom types are imported
		protoImports := opts.Fields.ProtoImports()
		for _, f := range opts.Fields.Custom() {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", PlaceholderProtoTxImport, importModule)
			content = replacer.Replace(content, PlaceholderProtoTxImport, replacementImport)
		}

		// Message
		// TODO: Include timestamp height
		// This addition would include using the type ibc.core.client.v1.Height
		// Ex: https://github.com/cosmos/cosmos-sdk/blob/816306b85addae6350bd380997f2f4bf9dce9471/proto/ibc/applications/transfer/v1/tx.proto
		templateMessage := `message MsgSend%[2]v {
  string %[3]v = 1;
  string port = 2;
  string channelID = 3;
  uint64 timeoutTimestamp = 4;
%[4]v}

message MsgSend%[2]vResponse {
}
%[1]v`
		replacementMessage := fmt.Sprintf(
			templateMessage,
			PlaceholderProtoTxMessage,
			opts.PacketName.UpperCamel,
			opts.MsgSigner.LowerCamel,
			sendFields,
		)
		content = replacer.Replace(content, PlaceholderProtoTxMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerTxModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "handler.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `case *types.MsgSend%[2]v:
					res, err := msgServer.Send%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
%[1]v`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			Placeholder,
			opts.PacketName.UpperCamel,
		)
		content = replacer.Replace(content, Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/tx.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `cmd.AddCommand(CmdSend%[2]v())
%[1]v`
		replacement := fmt.Sprintf(template, Placeholder, opts.PacketName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), module.Placeholder, replacement)

		// Register the module packet
		templateRegistry := `cdc.RegisterConcrete(&MsgSend%[2]v{}, "%[3]v/Send%[2]v", nil)
%[1]v`
		replacementRegistry := fmt.Sprintf(
			templateRegistry,
			module.Placeholder2,
			opts.PacketName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, module.Placeholder2, replacementRegistry)

		// Register the module packet interface
		templateInterface := `registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgSend%[2]v{},
)
%[1]v`
		replacementInterface := fmt.Sprintf(templateInterface, module.Placeholder3, opts.PacketName.UpperCamel)
		content = replacer.Replace(content, module.Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
