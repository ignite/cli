package scaffolder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/ibc"
)

const (
	ibcModuleImplementation = "module_ibc.go"
)

// AddPacket adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddPacket(
	moduleName,
	packetName string,
	packetFields,
	ackFields []string,
	noMessage bool,
) error {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	if err := checkComponentValidity(s.path, moduleName, packetName); err != nil {
		return err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	// Parse packet fields
	parsedPacketFields, err := parseFields(packetFields, isForbiddenPacketField)
	if err != nil {
		return err
	}

	// Parse acknowledgment fields
	parsedAcksFields, err := parseFields(ackFields, isGoReservedWord)
	if err != nil {
		return err
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.PacketOptions{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			PacketName: packetName,
			Fields:     parsedPacketFields,
			AckFields:  parsedAcksFields,
			NoMessage:  noMessage,
		}
	)
	g, err = ibc.NewPacket(opts)
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := s.protoc(pwd, path.RawPath); err != nil {
		return err
	}
	return fmtProject(pwd)
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

// isForbiddenPacketField returns true if the name is forbidden as a packet name
func isForbiddenPacketField(name string) bool {
	switch name {
	case
		"sender",
		"port",
		"channelID":
		return true
	}

	return isGoReservedWord(name)
}
