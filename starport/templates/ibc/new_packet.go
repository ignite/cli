package ibc

import (
	"fmt"

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

	// CODEC!!!

	if err := g.Box(ibcTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("title", strings.Title)

	ctx.Set("nodash", func(s string) string {
		return strings.ReplaceAll(s, "-", "")
	})

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func moduleModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketModuleRecv
		_ = `
case *types.<ModuleName>PacketData_<PacketName>Packet:
	err := am.keeper.OnRecv<Foo>Packet(ctx, modulePacket, packet.<PacketName>Packet)
	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType<Foo>Packet,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%t", err != nil)),
		),
	)
}
`

		// PlaceholderIBCPacketModuleAck
		_ = `
case *types.<ModuleName>PacketData_<PacketName>Packet:
	err := am.keeper.OnAcknowledgement<Foo>Packet(ctx, modulePacket, packet.<PacketName>Packet, ack)
	if err != nil {
		return nil, err
	}
	eventType = types.EventType<Foo>Packet
}
`

		// PlaceholderIBCPacketModuleTimeout
		_ = `
case *types.<ModuleName>PacketData_<PacketName>Packet:
	err := am.keeper.OnTimeoutPacket<Foo>Packet(ctx, modulePacket, packet.<PacketName>Packet)
	if err != nil {
		return nil, err
	}
}
`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func protoModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketProtoField
		_ = strings.Count("", PlaceholderIBCPacketProtoFieldNumber)
		_ = `
		PacketData packet = count; // placeholder
`

		// PlaceholderIBCPacketProtoMessage
		_ = `
// <%= title(moduleName) %>PacketData defines a struct for the packet payload
message <%= title(packetName) %>PacketData {
}
`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func typeModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketType
		_ = `
// ValidateBasic is used for validating the packet
func (p <%= title(packetName) %>PacketData) ValidateBasic() error {
	
	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p <%= title(packetName) %>PacketData) GetBytes() []byte {
	var modulePacket <%= title(packetName) %>PacketData

	modulePacket.Packet = &<ModuleName>PacketData_<PacketName>Packet{p}

	return ModuleCdc.MustMarshalBinaryBare(&p)
}`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func eventModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketEvent
		_ = `
EventTypePacket       = "<%= moduleName %>_packet"
`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}
