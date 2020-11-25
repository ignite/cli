package starportcmd

import "github.com/spf13/cobra"

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
	return c
}
