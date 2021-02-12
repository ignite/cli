package ibc

import (
	"fmt"

	"github.com/tendermint/starport/starport/templates/module"

	"github.com/tendermint/starport/starport/templates/typed"

	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

var (
	ibcTemplate = packr.New("ibc/templates/packet", "./packet")
)

// Options ...
type PacketOptions struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	PacketName string
	Fields     []typed.Field
}

// New ...
func NewIBC(opts *PacketOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleModify(opts))
	g.RunFn(protoModify(opts))
	g.RunFn(typeModify(opts))
	g.RunFn(eventModify(opts))

	// Modification for the new tx
	g.RunFn(protoTxModify(opts))
	g.RunFn(handlerTxModify(opts))
	g.RunFn(clientCliTxModify(opts))
	g.RunFn(codecModify(opts))

	if err := g.Box(ibcTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("packetName", opts.PacketName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("fields", opts.Fields)
	ctx.Set("title", strings.Title)

	ctx.Set("nodash", func(s string) string {
		return strings.ReplaceAll(s, "-", "")
	})

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{packetName}}", opts.PacketName))
	return g, nil
}

func moduleModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `%[1]v
case *types.%[2]vPacketData_%[3]vPacket:
	err := am.keeper.OnRecv%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())
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
			strings.Title(opts.PacketName),
		)
		content := strings.Replace(f.String(), PlaceholderIBCPacketModuleRecv, replacementRecv, 1)

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
			strings.Title(opts.PacketName),
		)
		content = strings.Replace(content, PlaceholderIBCPacketModuleAck, replacementAck, 1)

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
			strings.Title(opts.PacketName),
		)
		content = strings.Replace(content, PlaceholderIBCPacketModuleTimeout, replacementTimeout, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoModify(opts *PacketOptions) genny.RunFn {
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
			strings.Title(opts.PacketName),
			opts.PacketName,
			fieldCount+2,
			PlaceholderIBCPacketProtoFieldNumber,
		)
		content = strings.Replace(content, PlaceholderIBCPacketProtoField, replacementField, 1)

		// Add the message definition
		var messageFields string
		for i, field := range opts.Fields {
			messageFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+1)
		}
		templateMessage := `%[1]v
// %[2]vPacketData defines a struct for the packet payload
message %[2]vPacketData {
	%[3]v}
`
		replacementMessage := fmt.Sprintf(
			templateMessage,
			PlaceholderIBCPacketProtoMessage,
			strings.Title(opts.PacketName),
			messageFields,
		)
		content = strings.Replace(content, PlaceholderIBCPacketProtoMessage, replacementMessage, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func typeModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/packet.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v
// ValidateBasic is used for validating the packet
func (p %[3]vPacketData) ValidateBasic() error {
	
	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p %[3]vPacketData) GetBytes() ([]byte, error) {
	var modulePacket %[2]vPacketData

	modulePacket.Packet = &%[2]vPacketData_%[3]vPacket{&p}

	return ModuleCdc.MarshalBinaryBare(&modulePacket)
}`
		replacement := fmt.Sprintf(
			template,
			PlaceholderIBCPacketType,
			strings.Title(opts.ModuleName),
			strings.Title(opts.PacketName),
		)
		content := strings.Replace(f.String(), PlaceholderIBCPacketType, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func eventModify(opts *PacketOptions) genny.RunFn {
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
			strings.Title(opts.PacketName),
			opts.PacketName,
		)
		content := strings.Replace(f.String(), PlaceholderIBCPacketEvent, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Imports
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(
			templateImport,
			PlaceholderProtoTxImport,
			opts.ModuleName,
			opts.PacketName,
		)
		content := strings.Replace(f.String(), PlaceholderProtoTxImport, replacementImport, 1)

		// RPC
		templateRPC := `%[1]v
  rpc Send%[2]v(MsgSend%[2]v) returns (MsgSend%[2]vResponse);`
		replacementRPC := fmt.Sprintf(
			templateRPC,
			PlaceholderProtoTxRPC,
			strings.Title(opts.PacketName),
		)
		content = strings.Replace(content, PlaceholderProtoTxRPC, replacementRPC, 1)

		var sendFields string
		for i, field := range opts.Fields {
			sendFields += fmt.Sprintf("  %s %s = %d;\n", field.Datatype, field.Name, i+5)
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
			strings.Title(opts.PacketName),
			sendFields,
		)
		content = strings.Replace(content, PlaceholderProtoTxMessage, replacementMessage, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerTxModify(opts *PacketOptions) genny.RunFn {
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
		case *types.MsgSend%[2]v:
					res, err := msgServer.Send%[2]v(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
		replacementHandlers := fmt.Sprintf(templateHandlers,
			Placeholder,
			strings.Title(opts.PacketName),
		)
		content = strings.Replace(content, Placeholder, replacementHandlers, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(CmdSend%[2]v())
`
		replacement := fmt.Sprintf(template, Placeholder, strings.Title(opts.PacketName))
		content := strings.Replace(f.String(), Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := strings.Replace(f.String(), module.Placeholder, replacement, 1)

		// Register the module packet
		templateRegistry := `%[1]v
cdc.RegisterConcrete(&MsgSend%[2]v{}, "%[3]v/Send%[2]v", nil)
`
		replacementRegistry := fmt.Sprintf(templateRegistry, module.Placeholder2, strings.Title(opts.PacketName), opts.ModuleName)
		content = strings.Replace(content, module.Placeholder2, replacementRegistry, 1)

		// Register the module packet interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgSend%[2]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, module.Placeholder3, strings.Title(opts.PacketName))
		content = strings.Replace(content, module.Placeholder3, replacementInterface, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
