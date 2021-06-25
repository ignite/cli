package starportcmd

import "github.com/spf13/cobra"

// NewGenerate returns a command that groups code generation related sub commands.
func NewGenerate() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate [command]",
		Short:   "Generate code",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewGenerateGo())
	c.AddCommand(NewGenerateVuex())
	c.AddCommand(NewGenerateOpenAPI())

	return c
}
