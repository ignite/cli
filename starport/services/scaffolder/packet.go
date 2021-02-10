package scaffolder

import (
	"fmt"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"os"
	"path/filepath"
)

const (
	ibcModuleImplementation = "module_ibc.go"
	keeperDirectory = "keeper"
)


// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddPacket(moduleName string, packetName string, fields ...string) error {
	version, err := s.version()
	if err != nil {
		return err
	}
	_ = version.Major()
	_, err = gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// Module must exist
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Module must implement IBC
	ok, err = isIBCModule(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	// Check packet doesn't exist
	ok, err = isPacketCreated(s.path, moduleName, packetName)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("the packet %s already exist", packetName)
	}

	// Parse provided field
	_, err = parseFields(fields)
	if err != nil {
		return err
	}

	return nil
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

// isPacketCreated returns true if the provided packet already exists in the module
// we naively check the existence of keeper/<packetName>.go for this check
func isPacketCreated(appPath, moduleName, packetName string) (isCreated bool, err error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, moduleDir, moduleName, keeperDirectory, packetName+".go"))
	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		// Packet doesn't exist
		return false, nil
	}

	return true, err
}