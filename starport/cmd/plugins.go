package starportcmd

import (
	"github.com/spf13/cobra"
)

func NewPlugins() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugins",
		Short: "plugins for starport",
	}

	c.AddCommand(NewPluginsInstall())
	c.AddCommand(NewPluginsUse())

	return c
}
