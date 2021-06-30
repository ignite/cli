package starportcmd

import "github.com/spf13/cobra"

// NewScaffold returns a command that groups scaffolding related sub commands.
func NewScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:     "scaffold [command]",
		Short:   "Scaffold a new blockchain or scaffold components inside an existing one",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewScaffoldChain())
	c.AddCommand(NewScaffoldQuery())
	c.AddCommand(NewScaffoldModule())
	c.AddCommand(NewScaffoldWasm())
	c.AddCommand(NewScaffoldPacket())
	c.AddCommand(NewScaffoldMessage())

	return c
}
