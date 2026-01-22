package scaffolder

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
	"github.com/ignite/cli/v29/ignite/templates/ibc"
)

const (
	ibcModuleImplementation = "module_ibc.go"
)

// packetOptions represents configuration for the packet scaffolding.
type packetOptions struct {
	withoutMessage bool
	signer         string
}

// newPacketOptions returns a packetOptions with default options.
func newPacketOptions() packetOptions {
	return packetOptions{
		signer: "creator",
	}
}

// PacketOption configures the packet scaffolding.
type PacketOption func(*packetOptions)

// PacketWithoutMessage disables generating sdk compatible messages and tx related APIs.
func PacketWithoutMessage() PacketOption {
	return func(o *packetOptions) {
		o.withoutMessage = true
	}
}

// PacketWithSigner provides a custom signer name for the packet.
func PacketWithSigner(signer string) PacketOption {
	return func(m *packetOptions) {
		m.signer = signer
	}
}

// AddPacket adds a new type stype to scaffolded app by using optional type fields.
func (s Scaffolder) AddPacket(
	ctx context.Context,
	moduleName,
	packetName string,
	packetFields,
	ackFields []string,
	options ...PacketOption,
) error {
	// apply options.
	o := newPacketOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(packetName)
	if err != nil {
		return err
	}

	if err := checkComponentValidity(s.appPath, moduleName, name, o.withoutMessage); err != nil {
		return err
	}

	mfSigner, err := multiformatname.NewName(o.signer)
	if err != nil {
		return err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	signer := ""
	if !o.withoutMessage {
		signer = o.signer
	}

	// Check and parse packet fields
	if err := checkCustomTypes(ctx, s.appPath, s.modpath.Package, s.protoDir, moduleName, packetFields); err != nil {
		return err
	}
	parsedPacketFields, err := field.ParseFields(packetFields, checkForbiddenPacketField, signer)
	if err != nil {
		return err
	}

	// check and parse acknowledgment fields
	if err := checkCustomTypes(ctx, s.appPath, s.modpath.Package, s.protoDir, moduleName, ackFields); err != nil {
		return err
	}
	parsedAcksFields, err := field.ParseFields(ackFields, checkGoReservedWord, signer)
	if err != nil {
		return err
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.PacketOptions{
			AppName:    s.modpath.Package,
			ProtoDir:   s.protoDir,
			ProtoVer:   "v1", // TODO(@julienrbrt): possibly in the future add flag to specify custom proto version.
			ModulePath: s.modpath.RawPath,
			ModuleName: moduleName,
			PacketName: name,
			Fields:     parsedPacketFields,
			AckFields:  parsedAcksFields,
			NoMessage:  o.withoutMessage,
			MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewPacket(opts)
	if err != nil {
		return err
	}
	return s.Run(g)
}

// isIBCModule returns true if the provided module implements the IBC module interface
// we naively check the existence of module_ibc.go for this check.
func isIBCModule(appPath string, moduleName string) (bool, error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, moduleDir, moduleName, modulePkg, ibcModuleImplementation))
	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	} else if err == nil {
		// Is an IBC module
		return true, err
	}

	// check the legacy Path
	absPathLegacy, err := filepath.Abs(filepath.Join(appPath, moduleDir, moduleName, ibcModuleImplementation))
	if err != nil {
		return false, err
	}
	_, err = os.Stat(absPathLegacy)
	if os.IsNotExist(err) {
		// Not an IBC module
		return false, nil
	}

	return true, err
}

// checkForbiddenPacketField returns true if the name is forbidden as a packet name.
func checkForbiddenPacketField(name string) error {
	mfName, err := multiformatname.NewName(name)
	if err != nil {
		return err
	}

	switch mfName.LowerCase {
	case
		"sender",
		"port",
		"channelid",
		datatype.TypeCustom:
		return errors.Errorf("%s is used by the packet scaffolder", name)
	}

	return checkGoReservedWord(name)
}
