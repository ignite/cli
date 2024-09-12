package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewTestnet returns a command that groups scaffolding related sub commands.
func NewTestnet() *cobra.Command {
	c := &cobra.Command{
		Use:   "testnet [command]",
		Short: "Start a testnet local",
		Long: `Start a testnet local

`,
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(
		NewTestnetInPlace(),
	)

	return c
}
