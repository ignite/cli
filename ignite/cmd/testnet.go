package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewTestNet returns a command that groups scaffolding related sub commands.
func NewTestNet() *cobra.Command {
	c := &cobra.Command{
		Use:   "testnet [command]",
		Short: "Start a testnet local",
		Long: `Start a testnet local

`,
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(
		NewTestNetInPlace(),
	)

	return c
}
