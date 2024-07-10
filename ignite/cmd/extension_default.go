package ignitecmd

import (
	"github.com/spf13/cobra"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type defaultExtension struct {
	use     string
	short   string
	aliases []string
	path    string
}

const (
	ExtensionNetworkVersion = "v0.2.2"
	ExtensionNetworkPath    = "github.com/ignite/cli-plugin-network@" + ExtensionNetworkVersion
	ExtensionRelayerVersion = "hermes/v0.2.2"
	ExtensionRelayerPath    = "github.com/ignite/apps/hermes@" + ExtensionRelayerVersion
)

// defaultExtension holds the extensions that are considered trustable and for which
// a command will added if the extensions is not already installed.
// When the user executes that command, the extensions is automatically installed.
var defaultExtensions = []defaultExtension{
	{
		use:     "network",
		short:   "Launch a blockchain in production",
		aliases: []string{"n"},
		path:    ExtensionNetworkPath,
	},
	{
		use:     "relayer",
		short:   "Connect blockchains with an IBC relayer",
		aliases: []string{"r"},
		path:    ExtensionRelayerPath,
	},

	// TODO(@julienrbrt) eventually add Ignite Connect
}

// ensureDefaultExtensions ensures that all defaultExtensions are whether registered
// in cfg OR have an install command added to rootCmd.
func ensureDefaultExtensions(rootCmd *cobra.Command, cfg *pluginsconfig.Config) {
	for _, dp := range defaultExtensions {
		// Check if plugin is declared in global config
		if cfg.HasPlugin(dp.path) {
			// plugin found nothing to do
			continue
		}
		// plugin not found in config, add a proxy install command
		rootCmd.AddCommand(newPExtensionInstallCmd(dp))
	}
}

// newPExtensionInstallCmd mimics the plugin command but acts as proxy to first:
// - register the config in the global config
// - load the plugin
// - execute the command thanks to the loaded plugin.
func newPExtensionInstallCmd(dp defaultExtension) *cobra.Command {
	return &cobra.Command{
		Use:                dp.use,
		Short:              dp.short,
		Aliases:            dp.aliases,
		DisableFlagParsing: true, // Avoid -h to skip command run
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := parseGlobalExtensions()
			if err != nil {
				return err
			}

			// add plugin to config
			pluginCfg := pluginsconfig.Plugin{
				Path: dp.path,
			}
			cfg.Extensions = append(cfg.Extensions, pluginCfg)
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
			if err := linkExtensions(cmd.Context(), rootCmd, plugins); err != nil {
				return err
			}
			// Execute the command
			return rootCmd.Execute()
		},
	}
}
