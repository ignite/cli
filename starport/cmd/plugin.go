package starportcmd

import (
	"log"

	"github.com/spf13/cobra"
)

// TODO: Log issues.
// What is common method to log on Starport?

var (
	pluginHandler PluginCmdHandler = &pluginCmdHandler{}
)

// PluginCmdHandler provides interfaces to handle subcommands of plugin.
type PluginCmdHandler interface {
	HandleInstall(cmd *cobra.Command, args []string) error
	HandleList(cmd *cobra.Command, args []string) error
}

type pluginCmdHandler struct {
}

func (p *pluginCmdHandler) HandleInstall(cmd *cobra.Command, args []string) error {
	log.Println("HandleInstall", args)

	// TODO:

	return nil
}

func (p *pluginCmdHandler) HandleList(cmd *cobra.Command, args []string) error {
	log.Println("HandleList")

	// TODO:

	return nil
}

// NewPlugin creates a new plugin command to manage plugin.
func NewPlugin() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin list and install.",
		Args:  cobra.ExactArgs(1),
	}

	c.AddCommand(pluginListCmd())
	c.AddCommand(pluginInstallCmd())

	return c
}

func pluginListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List plugins cofigured",
		RunE:  pluginHandler.HandleList,
	}

	return c
}

func pluginInstallCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "install",
		Short: "Install new plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pluginHandler.HandleInstall(cmd, args)
		},
	}

	return c
}
