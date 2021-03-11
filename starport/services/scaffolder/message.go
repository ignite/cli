package scaffolder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/message"
)

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddMessage(moduleName string, msgName string, msgDesc string, fields []string, resField []string) error {
	version, err := s.version()
	if err != nil {
		return err
	}
	majorVersion := version.Major()
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the msg name is not a Go reserved name, it would generate an incorrect code
	if isGoReservedWord(msgName) {
		return fmt.Errorf("%s can't be used as a type name", msgName)
	}

	// Check msg is not already created
	ok, err = isMsgCreated(s.path, moduleName, msgName)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s message is already added", msgName)
	}

	// Parse provided fields
	parsedMsgFields, err := parseFields(fields, isForbiddenMessageField)
	if err != nil {
		return err
	}
	parsedResFields, err := parseFields(resField, isGoReservedWord)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &message.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			MsgName:    msgName,
			Fields:     parsedMsgFields,
			ResFields:  parsedResFields,
			MsgDesc:    msgDesc,
		}
	)
	// generate depending on the version
	if majorVersion == cosmosver.Launchpad {
		return errors.New("message scaffolding not supported on Launchpad")
	}
	// check if the msgServer convention is used
	var msgServerDefined bool
	msgServerDefined, err = isMsgServerDefined(s.path, moduleName)
	if err != nil {
		return err
	}
	if !msgServerDefined {
		// TODO: Determine if we want to support blockchains not using MsgServer convention
		return errors.New("the blockchain must use MsgServer convention")
	}

	// Scaffold
	g, err = message.NewStargate(opts)
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
	if err := s.protoc(pwd, path.RawPath, majorVersion); err != nil {
		return err
	}
	return fmtProject(pwd)
}

// isMsgCreated checks if the message is already scaffolded
func isMsgCreated(appPath, moduleName, msgName string) (isCreated bool, err error) {
	absPath, err := filepath.Abs(filepath.Join(
		appPath,
		moduleDir,
		moduleName,
		typesDirectory,
		"message_"+msgName+".go",
	))
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

// isForbiddenTypeField returns true if the name is forbidden as a message name
func isForbiddenMessageField(name string) bool {
	if name == "creator" {
		return true
	}

	return isGoReservedWord(name)
}
