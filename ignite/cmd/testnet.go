package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewTestnet returns a command that groups scaffolding related sub commands.
func NewTestnet() *cobra.Command {
	c := &cobra.Command{
		Use:     "testnet [command]",
		Short:   "Simulate and manage test networks",
		Long:    `Comprehensive toolset for managing and simulating blockchain test networks. It allows users to either run a test network in place using mainnet data or set up a multi-node environment for more complex testing scenarios. Additionally, it includes a subcommand for simulating the chain, which is useful for fuzz testing and other testing-related tasks.`,
		Aliases: []string{"t"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(
		NewTestnetInPlace(),
		NewTestnetMultiNode(),
		NewChainSimulate(), // While this is not per se a testnet command, it is related to testing.
	)

	return c
}
