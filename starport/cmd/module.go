package starportcmd

import "github.com/spf13/cobra"

// NewModule creates a new command that holds some other sub commands
// related to scaffolding sdk modules.
func NewModule() *cobra.Command {
	c := &cobra.Command{
		Use:   "module",
		Short: "Manage Cosmos SDK modules for your blockchain",
	}
	c.AddCommand(
		NewModuleImport(),
		NewModuleCreate(),
	)
	return c
}
