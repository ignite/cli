package ibc

import (
	"fmt"

	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

var (
	ibcTemplate = packr.New("ibc/templates/packet", "./packet")
)

// Field ...
type Field struct {
	Name         string
	Datatype     string
	DatatypeName string
}

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	TypeName   string
	Fields     []Field
	Legacy     bool
}

// New ...
func NewIBC(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleModify(opts))
	g.RunFn(protoModify(opts))
	g.RunFn(typeModify(opts))
	g.RunFn(eventModify(opts))

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

func moduleModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketModuleRecv
		_ = `
case ...:
	err := am.keeper.OnRecv<Foo>Packet(ctx, packet, data)
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
case ...:
	err := am.keeper.OnAcknowledgement<Foo>Packet(ctx, packet, data, ack)
	if err != nil {
		return nil, err
	}
	eventType = types.EventType<Foo>Packet
}
`

		// PlaceholderIBCPacketModuleTimeout
		_ = `
case ...:
	err := am.keeper.OnTimeoutPacket<Foo>Packet(ctx, packet, data)
	if err != nil {
		return nil, err
	}
}
`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func protoModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketProto
		_ = `
// <%= title(moduleName) %>PacketData defines a struct for the packet payload
message <%= title(moduleName) %>PacketData {
}
`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func typeModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		_, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// PlaceholderIBCPacketType
		_ = `
// ValidateBasic is used for validating the packet
func (p <%= title(moduleName) %>PacketData) ValidateBasic() error {
	
	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p <%= title(moduleName) %>PacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&p))
}`

		// newFile := genny.NewFileS(path, content)
		return nil // return r.File(newFile)
	}
}

func eventModify(opts *Options) genny.RunFn {
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
