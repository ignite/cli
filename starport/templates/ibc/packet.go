package ibc

import (
	"embed"
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/module"
)

var (
	//go:embed packet/component/* packet/component/**/*
	fsPacketComponent embed.FS

	//go:embed packet/messages/* packet/messages/**/*
	fsPacketMessages embed.FS

	// ibcTemplateComponent is the template to scaffold a new packet in an IBC module
	ibcTemplateComponent = xgenny.NewEmbedWalker(fsPacketComponent, "packet/component/")

	// ibcTemplateMessages is the template to scaffold send message for a packet
	ibcTemplateMessages = xgenny.NewEmbedWalker(fsPacketMessages, "packet/messages/")
)

// PacketOptions are options to scaffold a packet in a IBC module
type PacketOptions struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	PacketName multiformatname.Name
	Fields     []field.Field
	AckFields  []field.Field
	NoMessage  bool
}

// NewPacket returns the generator to scaffold a packet in an IBC module
func NewPacket(replacer placeholder.Replacer, opts *PacketOptions) (*genny.Generator, error) {
	g := genny.New()

	// Add the component
	g.RunFn(moduleModify(replacer, opts))
	g.RunFn(protoModify(replacer, opts))
	g.RunFn(eventModify(replacer, opts))
	if err := g.Box(ibcTemplateComponent); err != nil {
		return g, err
	}

	// Add the send message
	if !opts.NoMessage {
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(handlerTxModify(replacer, opts))
		g.RunFn(clientCliTxModify(replacer, opts))
		g.RunFn(codecModify(replacer, opts))
		if err := g.Box(ibcTemplateMessages); err != nil {
			return g, err
		}
	}

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("packetName", opts.PacketName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("fields", opts.Fields)
	ctx.Set("ackFields", opts.AckFields)
	ctx.Set("title", strings.Title)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{packetName}}", opts.PacketName.Snake))
	return g, nil
}

func moduleModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `%[1]v
case *types.%[2]vPacketData_%[3]vPacket:
	packetAck, err := am.keeper.OnRecv%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err.Error())
	} else {
		// Encode packet acknowledgment
		packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
		if err != nil {
			return nil, []byte{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		ack = channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType%[3]vPacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%%t", err != nil)),
		),
	)`
		replacementRecv := fmt.Sprintf(
			templateRecv,
			PlaceholderIBCPacketModuleRecv,
			strings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderIBCPacketModuleRecv, replacementRecv)

		// Ack packet dispatch
		templateAck := `%[1]v
case *types.%[2]vPacketData_%[3]vPacket:
	err := am.keeper.OnAcknowledgement%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket, ack)
	if err != nil {
		return nil, err
	}
	eventType = types.EventType%[3]vPacket`
		replacementAck := fmt.Sprintf(
			templateAck,
			PlaceholderIBCPacketModuleAck,
			strings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketModuleAck, replacementAck)

		// Timeout packet dispatch
		templateTimeout := `%[1]v
case *types.%[2]vPacketData_%[3]vPacket:
	err := am.keeper.OnTimeout%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		return nil, err
	}`
		replacementTimeout := fmt.Sprintf(
			templateTimeout,
			PlaceholderIBCPacketModuleTimeout,
			strings.Title(opts.ModuleName),
			opts.PacketName.UpperCamel,
		)
		content = replacer.Replace(content, PlaceholderIBCPacketModuleTimeout, replacementTimeout)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/packet.proto", opts.ModuleName)
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
			packetFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+1)
		}

		var ackFields string
		for i, field := range opts.AckFields {
			ackFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+1)
		}

		templateMessage := `%[1]v
// %[2]vPacketData defines a struct for the packet payload
message %[2]vPacketData {
	%[3]v}

// %[2]vPacketAck defines a struct for the packet acknowledgment
message %[2]vPacketAck {
	%[4]v}
`
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
		path := fmt.Sprintf("x/%s/types/events_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v
EventType%[2]vPacket       = "%[3]v_packet"
`
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
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// RPC
		templateRPC := `%[1]v
  rpc Send%[2]v(MsgSend%[2]v) returns (MsgSend%[2]vResponse);`
		replacementRPC := fmt.Sprintf(
			templateRPC,
			PlaceholderProtoTxRPC,
			opts.PacketName.UpperCamel,
		)
		content := replacer.Replace(f.String(), PlaceholderProtoTxRPC, replacementRPC)

		var sendFields string
		for i, field := range opts.Fields {
			sendFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name.LowerCamel, i+5)
		}

		// Message
		// TODO: Include timestamp height
		// This addition would include using the type ibc.core.client.v1.Height
		// Ex: https://github.com/cosmos/cosmos-sdk/blob/816306b85addae6350bd380997f2f4bf9dce9471/proto/ibc/applications/transfer/v1/tx.proto
		templateMessage := `%[1]v
message MsgSend%[2]v {
  string sender = 1;
  string port = 2;
  string channelID = 3;
  uint64 timeoutTimestamp = 4;
%[3]v}

message MsgSend%[2]vResponse {
}
`
		replacementMessage := fmt.Sprintf(
			templateMessage,
			PlaceholderProtoTxMessage,
			opts.PacketName.UpperCamel,
			sendFields,
		)
		content = replacer.Replace(content, PlaceholderProtoTxMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerTxModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
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
		case *types.MsgSend%[2]v:
					res, err := msgServer.Send%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
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
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(CmdSend%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder, opts.PacketName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), module.Placeholder, replacement)

		// Register the module packet
		templateRegistry := `%[1]v
cdc.RegisterConcrete(&MsgSend%[2]v{}, "%[3]v/Send%[2]v", nil)
`
		replacementRegistry := fmt.Sprintf(
			templateRegistry,
			module.Placeholder2,
			opts.PacketName.UpperCamel,
			opts.ModuleName,
		)
		content = replacer.Replace(content, module.Placeholder2, replacementRegistry)

		// Register the module packet interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgSend%[2]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, module.Placeholder3, opts.PacketName.UpperCamel)
		content = replacer.Replace(content, module.Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
