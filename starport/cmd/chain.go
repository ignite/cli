package starportcmd

import "github.com/spf13/cobra"

// NewChain returns a command that groups sub commands related to compiling, serving
// blockchains and so on.
func NewChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Short:   "Build, initialize, and start a blockchain in development",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewChainServe())
	c.AddCommand(NewChainBuild())
	c.AddCommand(NewChainInit())
	c.AddCommand(NewChainFaucet())

	return c
}
