package scaffolder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/datatype"
	"github.com/ignite/cli/ignite/templates/ibc"
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
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName,
	packetName string,
	packetFields,
	ackFields []string,
	options ...PacketOption,
) (sm xgenny.SourceModification, err error) {
	// apply options.
	o := newPacketOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(packetName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name, o.withoutMessage); err != nil {
		return sm, err
	}

	mfSigner, err := multiformatname.NewName(o.signer)
	if err != nil {
		return sm, err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if !ok {
		return sm, fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	signer := ""
	if !o.withoutMessage {
		signer = o.signer
	}

	// Check and parse packet fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, packetFields); err != nil {
		return sm, err
	}
	parsedPacketFields, err := field.ParseFields(packetFields, checkForbiddenPacketField, signer)
	if err != nil {
		return sm, err
	}

	// check and parse acknowledgment fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, ackFields); err != nil {
		return sm, err
	}
	parsedAcksFields, err := field.ParseFields(ackFields, checkGoReservedWord, signer)
	if err != nil {
		return sm, err
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.PacketOptions{
			AppName:    s.modpath.Package,
			AppPath:    s.path,
			ModulePath: s.modpath.RawPath,
			ModuleName: moduleName,
			PacketName: name,
			Fields:     parsedPacketFields,
			AckFields:  parsedAcksFields,
			NoMessage:  o.withoutMessage,
			MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewPacket(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// isIBCModule returns true if the provided module implements the IBC module interface
// we naively check the existence of module_ibc.go for this check.
func isIBCModule(appPath string, moduleName string) (bool, error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, moduleDir, moduleName, ibcModuleImplementation))
	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
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
		return fmt.Errorf("%s is used by the packet scaffolder", name)
	}

	return checkGoReservedWord(name)
}
