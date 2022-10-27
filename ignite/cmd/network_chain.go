package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkChain creates a new chain command that holds some other
// sub commands related to launching a network for a chain.
func NewNetworkChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain",
		Short: "Build networks",
	}

	c.AddCommand(
		NewNetworkChainList(),
		NewNetworkChainPublish(),
		NewNetworkChainInit(),
		NewNetworkChainInstall(),
		NewNetworkChainJoin(),
		NewNetworkChainPrepare(),
		NewNetworkChainShow(),
		NewNetworkChainLaunch(),
		NewNetworkChainRevertLaunch(),
	)

	return c
}
