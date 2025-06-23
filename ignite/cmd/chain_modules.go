package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewChainModules returns the modules command.
func NewChainModules() *cobra.Command {
	c := &cobra.Command{
		Use:   "modules",
		Short: "Manage modules",
		Long:  "The modules command allows you to manage modules in the codebase.",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(
		NewChainModulesList(),
	)

	return c
}
