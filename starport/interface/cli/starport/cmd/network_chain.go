package starportcmd

import "github.com/spf13/cobra"

// NewNetworkChain creates a new chain command that holds some other
// sub commands related to launching a network for a chain.
func NewNetworkChain() *cobra.Command {
	c := &cobra.Command{
		Use:               "chain",
		Short:             "Build networks",
		PersistentPreRunE: ensureSPNAccountHook,
	}
	c.AddCommand(NewNetworkChainCreate())
	c.AddCommand(NewNetworkChainJoin())
	c.AddCommand(NewNetworkChainStart())
	c.AddCommand(NewNetworkChainShow())
	c.AddCommand(NewNetworkChainList())
	return c
}
