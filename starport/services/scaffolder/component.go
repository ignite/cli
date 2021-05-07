package scaffolder

import (
	"context"
	"go/token"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// isComponentCreated checks if the component has been already created with Starport in the project
func isComponentCreated(appPath, moduleName, compName string) (bool, error) {
	// Check for type, packet or message creation
	created, err := isTypeCreated(appPath, moduleName, compName)
	if err != nil {
		return false, err
	}
	if created {
		return created, nil
	}
	created, err = isPacketCreated(appPath, moduleName, compName)
	if err != nil {
		return false, err
	}
	if created {
		return created, nil
	}
	created, err = isQueryCreated(appPath, moduleName, compName)
	if err != nil {
		return false, err
	}
	if created {
		return created, nil
	}
	return isMsgCreated(appPath, moduleName, compName)
}

// supportMsgServer checks if the module supports the MsgServer convention
// if not, the module codebase is modified to support it
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func supportMsgServer(
	appPath string,
	opts *modulecreate.MsgServerOptions,
) error {
	// Check if convention used
	msgServerDefined, err := isMsgServerDefined(appPath, opts.ModuleName)
	if err != nil {
		return err
	}
	if !msgServerDefined {
		// Patch the module to support the convention
		g, err := modulecreate.AddMsgServerConventionToLegacyModule(opts)
		if err != nil {
			return err
		}
		run := genny.WetRunner(context.Background())
		run.With(g)
		return run.Run()
	}
	return nil
}

// isMsgServerDefined checks if the module uses the MsgServer convention for transactions
// this is checked by verifying the existence of the tx.proto file
func isMsgServerDefined(appPath, moduleName string) (bool, error) {
	txProto, err := filepath.Abs(filepath.Join(appPath, "proto", moduleName, "tx.proto"))
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(txProto); os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// isForbiddenComponentName returns true if the name is forbidden as a component name
func isForbiddenComponentName(name string) bool {
	switch name {
	case
		"logger",
		"keeper",
		"query",
		"genesis",
		"types",
		"tx":
		return true
	}

	return isGoReservedWord(name)
}

// isGoReservedWord checks if the name can't be used because it is a go reserved keyword
func isGoReservedWord(name string) bool {
	// Check keyword or literal
	if token.Lookup(name).IsKeyword() {
		return true
	}

	// Check with builtin identifier
	switch name {
	case
		"panic",
		"recover",
		"append",
		"bool",
		"byte",
		"cap",
		"close",
		"complex",
		"complex64",
		"complex128",
		"uint16",
		"copy",
		"false",
		"float32",
		"float64",
		"imag",
		"int",
		"int8",
		"int16",
		"uint32",
		"int32",
		"int64",
		"iota",
		"len",
		"make",
		"new",
		"nil",
		"uint64",
		"print",
		"println",
		"real",
		"string",
		"true",
		"uint",
		"uint8",
		"uintptr":
		return true
	}
	return false
}
