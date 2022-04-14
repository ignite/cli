package ignitecmd

import "github.com/spf13/cobra"

// NewChain returns a command that groups sub commands related to compiling, serving
// blockchains and so on.
func NewChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Short:   "Build, initialize and start a blockchain node or perform other actions on the blockchain",
		Long:    `Build, initialize and start a blockchain node or perform other actions on the blockchain.`,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(
		NewChainServe(),
		NewChainBuild(),
		NewChainInit(),
		NewChainFaucet(),
		NewChainSimulate(),
	)

	return c
}
