package scaffolder

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// supportGenesisTests checks if types/genesis_test.go exists
// appends the generator to create the file if it doesn't
func supportGenesisTests(
	gens []*genny.Generator,
	appPath,
	appName,
	modulePath,
	moduleName string,
) ([]*genny.Generator, error) {
	modulePath, err := filepath.Abs(filepath.Join(appPath, "x", moduleName))
	if err != nil {
		return nil, err
	}

	gmPath := filepath.Join(modulePath, "genesis_test.go")
	if _, err := os.Stat(gmPath); os.IsNotExist(err) {
		g, err := modulecreate.AddGenesisModuleTest(appName, modulePath, moduleName)
		if err != nil {
			return nil, err
		}
		gens = append(gens, g)
	} else if err != nil {
		return nil, err
	}

	gtPath := filepath.Join(modulePath, "types/genesis_test.go")
	if _, err := os.Stat(gtPath); os.IsNotExist(err) {
		g, err := modulecreate.AddGenesisTypesTest(appName, modulePath, moduleName)
		if err != nil {
			return nil, err
		}
		gens = append(gens, g)
	} else if err != nil {
		return nil, err
	}
	return gens, nil
}

// supportMsgServer checks if the module supports the MsgServer convention
// appends the generator to support it if it doesn't
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func supportMsgServer(
	gens []*genny.Generator,
	replacer placeholder.Replacer,
	appPath string,
	opts *modulecreate.MsgServerOptions,
) ([]*genny.Generator, error) {
	// Check if convention used
	msgServerDefined, err := isMsgServerDefined(appPath, opts.ModuleName)
	if err != nil {
		return nil, err
	}
	if !msgServerDefined {
		// Patch the module to support the convention
		g, err := modulecreate.AddMsgServerConventionToLegacyModule(replacer, opts)
		if err != nil {
			return nil, err
		}
		gens = append(gens, g)
	}
	return gens, nil
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
