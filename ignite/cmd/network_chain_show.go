package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/network"
)

const flagOut = "out"

// NewNetworkChainShow creates a new chain show
// command to show a chain details on SPN.
func NewNetworkChainShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "Show details of a chain",
	}
	c.AddCommand(
		newNetworkChainShowInfo(),
		newNetworkChainShowGenesis(),
		newNetworkChainShowAccounts(),
		newNetworkChainShowValidators(),
		newNetworkChainShowPeers(),
	)
	return c
}

func networkChainLaunch(cmd *cobra.Command, args []string, session cliui.Session) (NetworkBuilder, uint64, error) {
	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return nb, 0, err
	}
	// parse launch ID.
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return nb, launchID, err
	}
	return nb, launchID, err
}
