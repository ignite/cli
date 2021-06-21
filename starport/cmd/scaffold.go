package starportcmd

import "github.com/spf13/cobra"

// NewScaffold returns a command that groups scaffolding related sub commands.
func NewScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:     "scaffold [command]",
		Short:   "Scaffold a blockchain or add features to it",
		Aliases: []string{"sc"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewScaffoldChain())

	return c
}
