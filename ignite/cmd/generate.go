package ignitecmd

import "github.com/spf13/cobra"

// NewGenerate returns a command that groups code generation related sub commands.
func NewGenerate() *cobra.Command {
	c := &cobra.Command{
		Use:   "generate [command]",
		Short: "Generate clients, API docs from source code",
		Long: `Generate clients, API docs from source code.

Such as compiling protocol buffer files into Go or implement particular functionality, for example, generating an OpenAPI spec.

Produced source code can be regenerated by running a command again and is not meant to be edited by hand.`,
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.AddCommand(NewGenerateGo())
	c.AddCommand(NewGenerateVuex())
	c.AddCommand(NewGenerateDart())
	c.AddCommand(NewGenerateOpenAPI())

	return c
}
