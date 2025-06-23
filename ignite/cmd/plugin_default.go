package ignitecmd

import (
	"github.com/spf13/cobra"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type defaultPlugin struct {
	use     string
	short   string
	aliases []string
	path    string
}

const (
	PluginRelayerVersion     = "hermes/v0.3.0"
	PluginRelayerPath        = "github.com/ignite/apps/hermes@" + PluginRelayerVersion
	PluginAppRegistryVersion = "appregistry/v0.1.3"
	PluginAppRegistryPath    = "github.com/ignite/apps/appregistry@" + PluginAppRegistryVersion
)

// defaultPlugins holds the plugin that are considered trustable and for which
// a command will added if the plugin is not already installed.
// When the user executes that command, the plugin is automatically installed.
var defaultPlugins = []defaultPlugin{
	{
		use:     "relayer",
		short:   "Connect blockchains with an IBC relayer",
		aliases: []string{"r"},
		path:    PluginRelayerPath,
	},
	{
		use:     "appregistry",
		short:   "Browse the Ignite App Registry App",
		aliases: []string{"mp"},
		path:    PluginAppRegistryPath,
	},
}

// ensureDefaultPlugins ensures that all defaultPlugins are whether registered
// in cfg OR have an install command added to rootCmd.
func ensureDefaultPlugins(rootCmd *cobra.Command, cfg *pluginsconfig.Config) {
	for _, dp := range defaultPlugins {
		// Check if plugin is declared in global config
		if cfg.HasPlugin(dp.path) {
			// plugin found nothing to do
			continue
		}
		// plugin not found in config, add a proxy install command
		rootCmd.AddCommand(newPluginInstallCmd(dp))
	}
}

// newPluginInstallCmd mimics the plugin command but acts as proxy to first:
// - register the config in the global config
// - load the plugin
// - execute the command thanks to the loaded plugin.
func newPluginInstallCmd(dp defaultPlugin) *cobra.Command {
	return &cobra.Command{
		Use:                dp.use,
		Short:              dp.short,
		Aliases:            dp.aliases,
		DisableFlagParsing: true, // Avoid -h to skip command run
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := parseGlobalPlugins()
			if err != nil {
				return err
			}

			// add plugin to config
			pluginCfg := pluginsconfig.Plugin{
				Path: dp.path,
			}
			cfg.Apps = append(cfg.Apps, pluginCfg)
			if err := cfg.Save(); err != nil {
				return err
			}

			session := cliui.New(cliui.WithoutUserInteraction(getYes(cmd)))
			defer session.End()

			// load and link the plugin
			plugins, err := plugin.Load(
				cmd.Context(),
				[]pluginsconfig.Plugin{pluginCfg},
				plugin.CollectEvents(session.EventBus()),
			)
			if err != nil {
				return err
			}
			defer plugins[0].KillClient()

			// Keep reference of the root command before removal
			rootCmd := cmd.Root()
			// Remove this command before call to linkPlugins because a plugin is
			// usually not allowed to override an existing command.
			rootCmd.RemoveCommand(cmd)
			if err := linkPlugins(cmd.Context(), rootCmd, plugins); err != nil {
				return err
			}
			// Execute the command
			return rootCmd.Execute()
		},
	}
}
