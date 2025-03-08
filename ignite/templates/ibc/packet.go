package ibc

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/testutil"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

var (
	//go:embed files/packet/component/* files/packet/component/**/*
	fsPacketComponent embed.FS

	//go:embed files/packet/messages/* files/packet/messages/**/*
	fsPacketMessages embed.FS
)

// PacketOptions are options to scaffold a packet in a IBC module.
type PacketOptions struct {
	AppName    string
	AppPath    string
	ProtoDir   string
	ProtoVer   string
	ModuleName string
	ModulePath string
	PacketName multiformatname.Name
	MsgSigner  multiformatname.Name
	Fields     field.Fields
	AckFields  field.Fields
	NoMessage  bool
}

// ProtoFile returns the path to the proto folder.
func (opts *PacketOptions) ProtoFile(fname string) string {
	return filepath.Join(opts.AppPath, opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}

// NewPacket returns the generator to scaffold a packet in an IBC module.
func NewPacket(replacer placeholder.Replacer, opts *PacketOptions) (*genny.Generator, error) {
	var (
		g = genny.New()

		componentTemplate = xgenny.NewEmbedWalker(
			fsPacketComponent,
			"files/packet/component/",
			opts.AppPath,
		)
		messagesTemplate = xgenny.NewEmbedWalker(
			fsPacketMessages,
			"files/packet/messages/",
			opts.AppPath,
		)
	)

	// Add the component
	g.RunFn(moduleModify(replacer, opts))
	g.RunFn(protoModify(opts))
	g.RunFn(eventModify(replacer, opts))
	if err := g.Box(componentTemplate); err != nil {
		return g, err
	}

	// Add the send message
	if !opts.NoMessage {
		g.RunFn(protoTxModify(opts))
		g.RunFn(clientCliTxModify(opts))
		g.RunFn(codecModify(opts))
		if err := g.Box(messagesTemplate); err != nil {
			return g, err
		}
	}

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("protoVer", opts.ProtoVer)
	ctx.Set("packetName", opts.PacketName)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("fields", opts.Fields)
	ctx.Set("ackFields", opts.AckFields)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))
	g.Transformer(genny.Replace("{{packetName}}", opts.PacketName.Snake))

	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, opts.AppPath); err != nil {
		return g, err
	}

	return g, nil
}

func moduleModify(replacer placeholder.Replacer, opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module/module_ibc.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `case *types.%[2]vPacketData_%[3]vPacket:
	packetAck, err := im.keeper.OnRecv%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err)
	} else {
		// Encode packet acknowledgment
		packetAckBytes, err := im.cdc.MarshalJSON(&packetAck)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error()))
		}
		ack = channeltypes.NewResultAcknowledgement(packetAckBytes)
	}


	sdkCtx := sdk.UnwrapSDKContext(ctx)
    sdkCtx.EventManager().EmitEvent(
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
	err := im.keeper.OnAcknowledgement%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket, ack)
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
	err := im.keeper.OnTimeout%[3]vPacket(ctx, modulePacket, *packet.%[3]vPacket)
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

// Modifies packet.proto to add a field on the oneof element of the message created and
// add a couple of messages.
//
// What it depends on:
//   - Existence of a Oneof field named 'packet'.
func protoModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("packet.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		name := fmt.Sprintf("%sPacketData", xstrings.Title(opts.ModuleName))
		message, err := protoutil.GetMessageByName(protoFile, name)
		if err != nil {
			return errors.Errorf("failed while looking up '%s' message in %s: %w", name, path, err)
		}
		// Use a directly Apply call here, modifying oneofs isn't common enough to warrant a separate function.
		var packet *proto.Oneof
		protoutil.Apply(message, nil, func(c *protoutil.Cursor) bool {
			if o, ok := c.Node().(*proto.Oneof); ok {
				if o.Name == "packet" {
					packet = o
					return false
				}
			}
			// continue traversing.
			return true
		})
		if packet == nil {
			return errors.Errorf("could not find 'oneof packet' in message '%s' of file %s", name, path)
		}
		// Count fields of oneof:
		maximum := 1
		protoutil.Apply(packet, nil, func(c *protoutil.Cursor) bool {
			if o, ok := c.Node().(*proto.OneOfField); ok {
				if o.Sequence > maximum {
					maximum = o.Sequence
				}
			}
			return true
		})
		// Add it to Oneof.
		typenameUpper, typenameLower := opts.PacketName.UpperCamel, opts.PacketName.LowerCamel
		packetField := protoutil.NewOneofField(typenameLower+"Packet", typenameUpper+"PacketData", maximum+1)
		protoutil.Append(packet, packetField)

		// Add the message definition for packet and acknowledgment
		var packetFields []*proto.NormalField
		for i, field := range opts.Fields {
			packetFields = append(packetFields, field.ToProtoField(i+1))
		}
		packetData := protoutil.NewMessage(typenameUpper+"PacketData", protoutil.WithFields(packetFields...))
		protoutil.AttachComment(packetData, typenameUpper+"PacketData defines a struct for the packet payload")
		var ackFields []*proto.NormalField
		for i, field := range opts.AckFields {
			ackFields = append(ackFields, field.ToProtoField(i+1))
		}
		packetAck := protoutil.NewMessage(typenameUpper+"PacketAck", protoutil.WithFields(ackFields...))
		protoutil.AttachComment(packetAck, typenameUpper+"PacketAck defines a struct for the packet acknowledgment")
		protoutil.Append(protoFile, packetData, packetAck)

		// Add any custom imports.
		var protoImports []*proto.Import
		for _, imp := range append(opts.Fields.ProtoImports(), opts.AckFields.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range append(opts.Fields.Custom(), opts.AckFields.Custom()...) {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		if err := protoutil.AddImports(protoFile, true, protoImports...); err != nil {
			return errors.Errorf("failed while adding imports to %s: %w", path, err)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
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

// Modifies tx.proto to add a new RPC and the required messages.
//
// What it depends on:
//   - Existence of a service named 'Msg'. The other elements don't depend on already existing
//     elements in the file.
func protoTxModify(opts *PacketOptions) genny.RunFn {
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

		// Add RPC to service Msg.
		serviceMsg, err := protoutil.GetServiceByName(protoFile, "Msg")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Msg' in %s: %w", path, err)
		}
		typenameUpper := opts.PacketName.UpperCamel
		send := protoutil.NewRPC(
			fmt.Sprintf("Send%s", typenameUpper),
			fmt.Sprintf("MsgSend%s", typenameUpper),
			fmt.Sprintf("MsgSend%sResponse", typenameUpper),
		)
		protoutil.Append(serviceMsg, send)

		// Create fields for MsgSend.
		var sendFields []*proto.NormalField
		for i, field := range opts.Fields {
			sendFields = append(sendFields, field.ToProtoField(i+5))
		}
		sendFields = append(sendFields,
			protoutil.NewField(opts.MsgSigner.LowerCamel, "string", 1),
			protoutil.NewField("port", "string", 2),
			protoutil.NewField("channelID", "string", 3),
			protoutil.NewField("timeoutTimestamp", "uint64", 4),
		)
		creatorOpt := protoutil.NewOption(typed.MsgSignerOption, opts.MsgSigner.LowerCamel)

		// Create MsgSend, MsgSendResponse and add to file.
		msgSend := protoutil.NewMessage(
			"MsgSend"+typenameUpper,
			protoutil.WithFields(sendFields...),
			protoutil.WithMessageOptions(creatorOpt),
		)
		msgSendResponse := protoutil.NewMessage("MsgSend" + typenameUpper + "Response")
		protoutil.Append(protoFile, msgSend, msgSendResponse)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range opts.Fields.ProtoImports() {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range opts.Fields.Custom() {
			protopath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protopath))
		}
		if err := protoutil.AddImports(protoFile, true, protoImports...); err != nil {
			return errors.Errorf("error while processing %s: %w", path, err)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

// clientCliTxModify does not use AutoCLI here, because it as a better UX as it is.
func clientCliTxModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		filePath := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/tx.go")
		f, err := r.Disk.Find(filePath)
		if err != nil {
			return err
		}
		replacement := fmt.Sprintf("cmd.AddCommand(CmdSend%[1]v())", opts.PacketName.UpperCamel)
		content, err := xast.ModifyFunction(
			f.String(),
			"GetTxCmd",
			xast.AppendFuncCode(replacement),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(filePath, content)
		return r.File(newFile)
	}
}

func codecModify(opts *PacketOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		content, err := xast.AppendImports(f.String(), xast.WithLastNamedImport("sdk", "github.com/cosmos/cosmos-sdk/types"))
		if err != nil {
			return err
		}

		// Register the module packet interface
		templateInterface := `registrar.RegisterImplementations((*sdk.Msg)(nil),
	&MsgSend%[1]v{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, opts.PacketName.UpperCamel)
		content, err = xast.ModifyFunction(
			content,
			"RegisterInterfaces",
			xast.AppendFuncAtLine(replacementInterface, 0),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
