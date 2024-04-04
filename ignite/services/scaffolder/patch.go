package scaffolder

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

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
	msgServerDefined, err := isMsgServerDefined(appPath, opts.AppName, opts.ProtoPath, opts.ModuleName)
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
func isMsgServerDefined(appPath, appName, protoPath, moduleName string) (bool, error) {
	txProto, err := filepath.Abs(filepath.Join(appPath, protoPath, appName, moduleName, "tx.proto"))
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(txProto); os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
