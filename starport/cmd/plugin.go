package starportcmd

import "github.com/spf13/cobra"

// NewPlugin returns a command that groups sub commands related to chain plugins.
func NewPlugin() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugin [command]",
		Short: "Manage plugins specified in config file.",
		Long:  `Manage plugins specified in config file.`,
		Args:  cobra.ExactArgs(1),
	}

	flagSetPath(c)
	c.AddCommand(NewPluginReload())
	c.AddCommand(NewPluginPull())
	c.AddCommand(NewPluginBuild())

	return c
}
