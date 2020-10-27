package starportcmd

import "github.com/spf13/cobra"

func NewNetworkChain() *cobra.Command {
	c := &cobra.Command{
		Use:  "chain",
		Args: cobra.ExactArgs(1),
	}
	c.AddCommand(NewNetworkChainCreate())
	c.AddCommand(NewNetworkChainJoin())
	return c
}
