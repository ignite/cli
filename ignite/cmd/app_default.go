package ignitecmd

import (
	"github.com/spf13/cobra"

	appsconfig "github.com/ignite/cli/ignite/config/apps"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/plugin"
)

type defaultApp struct {
	use     string
	short   string
	aliases []string
	path    string
}

const (
	AppNetworkVersion = "v0.1.1"
	AppNetworkPath    = "github.com/ignite/cli-plugin-network@" + AppNetworkVersion
)

// defaultApps holds the app that are considered trustable and for which
// a command will added if the app is not already installed.
// When the user executes that command, the app is automatically installed.
var defaultApps = []defaultApp{
	{
		use:     "network",
		short:   "Launch a blockchain in production",
		aliases: []string{"n"},
		path:    AppNetworkPath,
	},
}

// ensureDefaultApps ensures that all defaultApps are wether registered
// in cfg OR have an install command added to rootCmd.
func ensureDefaultApps(rootCmd *cobra.Command, cfg *appsconfig.Config) {
	for _, app := range defaultApps {
		// Check if app is declared in global config
		if cfg.HasApp(app.path) {
			// app found nothing to do
			continue
		}
		// app not found in config, add a proxy install command
		rootCmd.AddCommand(newAppInstallCmd(app))
	}
}

// newAppInstallCmd mimics the app command but acts as proxy to first:
// - register the config in the global config
// - load the app
// - execute the command thanks to the loaded app.
func newAppInstallCmd(app defaultApp) *cobra.Command {
	return &cobra.Command{
		Use:                app.use,
		Short:              app.short,
		Aliases:            app.aliases,
		DisableFlagParsing: true, // Avoid -h to skip command run
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := parseGlobalApps()
			if err != nil {
				return err
			}

			// add app to config
			appCfg := appsconfig.App{Path: app.path}
			cfg.Apps = append(cfg.Apps, appCfg)
			if err := cfg.Save(); err != nil {
				return err
			}

			session := cliui.New()
			defer session.End()

			// load and link the app
			plugins, err := plugin.Load(
				cmd.Context(),
				[]appsconfig.App{appCfg},
				plugin.CollectEvents(session.EventBus()),
			)
			if err != nil {
				return err
			}
			defer plugins[0].KillClient()

			// Keep reference of the root command before removal
			rootCmd := cmd.Root()
			// Remove this command before call to linkApps because an app is
			// usually not allowed to override an existing command.
			rootCmd.RemoveCommand(cmd)
			if err := linkApps(rootCmd, plugins); err != nil {
				return err
			}
			// Execute the command
			return rootCmd.Execute()
		},
	}
}
