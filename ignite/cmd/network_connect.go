package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/services/network"
)

// NewNetworkConnect connects the monitoring modules of launched chains with SPN
func NewNetworkConnect() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkConnectHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().String(flagSourceGasPrice, "0.0000025", "Gas price used for transactions on source chain")
	c.Flags().String(flagTargetGasPrice, "0.0000025", "Gas price used for transactions on target chain")
	c.Flags().Int64(flagSourceGasLimit, 300000, "Gas limit used for transactions on source chain")
	c.Flags().Int64(flagTargetGasLimit, 300000, "Gas limit used for transactions on target chain")
	return c
}

// ignite network --local connect [launch-id] [target-rpc]

// Flag
// --target-faucet

// Flag values with defaults
// --source-gaslimit
// --target-gaslimit
// --source-gasprice
// --target-gasprice
// --source-account (from)
// --target-account (from)

// Hardcoded flags
// --ordered
// --source-rpc "http://0.0.0.0:26657"
// --source-faucet "http://0.0.0.0:4500"
// --source-port "monitoringc"
// --target-port "monitoringp"
// --source-version "monitoring-1"
// --target-version "monitoring-1"
// --source-prefix "spn"
// --target-prefix "cosmos" (fetch from genesis)

func networkConnectHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	nodeAPI := args[1]

	clientID, err := clientCreate(cmd, launchID, nodeAPI)
	if err != nil {
		return err
	}

	session.StopSpinner()
	return nil
}
