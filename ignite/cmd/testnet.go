package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewTestnet returns a command that groups scaffolding related sub commands.
func NewTestnet() *cobra.Command {
	c := &cobra.Command{
		Use:     "testnet [command]",
		Short:   "Simulate (Fuzz) the chain or start a testnet, either in place (using mainnet data) or with multiple nodes.",
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
