package starportcmd

import "github.com/spf13/cobra"

func NewModule() *cobra.Command {
	c := &cobra.Command{
		Use:   "module",
		Short: "Manage cosmos modules for your app",
	}
	c.AddCommand(
		NewModuleImport(),
		NewModuleCreate(),
	)
	return c
}
