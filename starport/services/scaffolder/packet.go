package scaffolder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/ibc"
)

const (
	ibcModuleImplementation = "module_ibc.go"
)

// packetOptions represents configuration for the packet scaffolding
type packetOptions struct {
	withoutMessage bool
	signer         string
}

// newPacketOptions returns a packetOptions with default options
func newPacketOptions() packetOptions {
	return packetOptions{
		signer: "creator",
	}
}

// PacketOption configures the packet scaffolding
type PacketOption func(*packetOptions)

// PacketWithoutMessage disables generating sdk compatible messages and tx related APIs.
func PacketWithoutMessage() PacketOption {
	return func(o *packetOptions) {
		o.withoutMessage = true
	}
}

// PacketWithSigner provides a custom signer name for the packet
func PacketWithSigner(signer string) PacketOption {
	return func(m *packetOptions) {
		m.signer = signer
	}
}

// AddPacket adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddPacket(
	tracer *placeholder.Tracer,
	moduleName,
	packetName string,
	packetFields,
	ackFields []string,
	options ...PacketOption,
) (sm xgenny.SourceModification, err error) {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	// apply options.
	o := newPacketOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.Lowercase

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

	// Parse packet fields
	parsedPacketFields, err := field.ParseFields(packetFields, checkForbiddenPacketField)
	if err != nil {
		return sm, err
	}

	// Parse acknowledgment fields
	parsedAcksFields, err := field.ParseFields(ackFields, checkGoReservedWord)
	if err != nil {
		return sm, err
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.PacketOptions{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
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
	pwd, err := os.Getwd()
	if err != nil {
		return sm, err
	}
	return sm, s.finish(pwd, path.RawPath)
}

// isIBCModule returns true if the provided module implements the IBC module interface
// we naively check the existence of module_ibc.go for this check
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

// checkForbiddenPacketField returns true if the name is forbidden as a packet name
func checkForbiddenPacketField(name string) error {
	switch name {
	case
		"sender",
		"port",
		"channelID":
		return fmt.Errorf("%s is used by the packet scaffolder", name)
	}

	return checkGoReservedWord(name)
}
