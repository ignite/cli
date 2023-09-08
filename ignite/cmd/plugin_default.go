package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/plugin"
)

// ensureDefaultPlugins ensures that all defaultPlugins are wether registered
// in cfg OR have an install command added to rootCmd.
func ensureDefaultPlugins(rootCmd *cobra.Command, cfg *pluginsconfig.Config) error {
	defaultPlugins, err := plugin.GetDefaultPlugins()
	if err != nil {
		return err
	}

	for _, dp := range defaultPlugins {
		// Check if plugin is declared in global config
		if cfg.HasPlugin(dp.Path) {
			// plugin found nothing to do
			continue
		}
		// plugin not found in config, add a proxy install command
		rootCmd.AddCommand(newPluginInstallCmd(dp))
	}

	return nil
}

// newPluginInstallCmd mimics the plugin command but acts as proxy to first:
// - register the config in the global config
// - load the plugin
// - execute the command thanks to the loaded plugin.
func newPluginInstallCmd(dp plugin.DefaultPlugin) *cobra.Command {
	return &cobra.Command{
		Use:                dp.Use,
		Short:              dp.Short,
		Aliases:            dp.Aliases,
		DisableFlagParsing: true, // Avoid -h to skip command run
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := parseGlobalPlugins()
			if err != nil {
				return err
			}
			if cfg.HasPlugin(dp.Path) {
				// plugin already declared in global plugins, this shouldn't happen
				// because this is actually why this command has been added, so let's
				// break violently
				panic(fmt.Sprintf("plugin %q unexpected in global config", dp.Path))
			}

			// add plugin to config
			pluginCfg := pluginsconfig.Plugin{
				Path: dp.Path,
			}
			cfg.Plugins = append(cfg.Plugins, pluginCfg)
			if err := cfg.Save(); err != nil {
				return err
			}

			session := cliui.New()
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
			if err := linkPlugins(rootCmd, plugins); err != nil {
				return err
			}
			// Execute the command
			return rootCmd.Execute()
		},
	}
}
