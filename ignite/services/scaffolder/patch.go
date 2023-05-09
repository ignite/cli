package scaffolder

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

// supportSimulation checks if module_simulation.go exists,
// appends the generator to create the file if it doesn't.
func supportSimulation(
	gens []*genny.Generator,
	appPath,
	modulePath,
	moduleName string,
) ([]*genny.Generator, error) {
	simulation, err := modulecreate.AddSimulation(
		appPath,
		modulePath,
		moduleName,
	)
	if err != nil {
		return gens, err
	}
	gens = append(gens, simulation)
	return gens, nil
}

// supportGenesisTests checks if types/genesis_test.go exists
// appends the generator to create the file if it doesn't.
func supportGenesisTests(
	gens []*genny.Generator,
	appPath,
	appName,
	modulePath,
	moduleName string,
	isIBC bool,
) ([]*genny.Generator, error) {
	genesisTest, err := modulecreate.AddGenesisTest(
		appPath,
		appName,
		modulePath,
		moduleName,
		isIBC,
	)
	if err != nil {
		return gens, err
	}
	gens = append(gens, genesisTest)
	return gens, nil
}

// supportMsgServer checks if the module supports the MsgServer convention
// appends the generator to support it if it doesn't
// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-031-msg-service.md
func supportMsgServer(
	gens []*genny.Generator,
	replacer placeholder.Replacer,
	appPath string,
	opts *modulecreate.MsgServerOptions,
) ([]*genny.Generator, error) {
	// Check if convention used
	msgServerDefined, err := isMsgServerDefined(appPath, opts.AppName, opts.ModuleName)
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
// this is checked by verifying the existence of the tx.proto file.
func isMsgServerDefined(appPath, appName, moduleName string) (bool, error) {
	txProto, err := filepath.Abs(filepath.Join(appPath, "proto", appName, moduleName, "tx.proto"))
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(txProto); os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
