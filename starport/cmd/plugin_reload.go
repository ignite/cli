package starportcmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/pluginsrpc"
)

// NewPluginReload creates a new reload command to manually refresh chain plugins.
func NewPluginReload() *cobra.Command {
	c := &cobra.Command{
		Use:   "reload",
		Short: "Reload plugins specified in config file.",
		RunE:  pluginReloadHandler,
	}

	flagSetPath(c)
	c.Flags().StringP(flagConfig, "c", "", "Starport config file (default: ./config.yml)")

	return c
}

func pluginReloadHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Reloading plugins...")
	defer s.Stop()

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
	}

	config, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	if config != "" {
		chainOption = append(chainOption, chain.ConfigFile(config))
	}

	// create the chain
	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	chainConfig, err := c.Config()
	if err != nil {
		return err
	}

	chainId, err := c.ID()
	if err != nil {
		return err
	}

	pluginManager := pluginsrpc.NewManager(chainId, chainConfig)
	if err := pluginManager.PullBuild(cmd.Context()); err != nil {
		return err
	}

	fmt.Println("🔄  Reloaded plugins.")
	return nil
}

// Add support for custom config files
func getCommandConfig(cmd *cobra.Command) (chainconfig.Config, error) {
	configPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return chainconfig.Config{}, err
	}

	if configPath != "" {
		return chainconfig.ParseFile(configPath)
	}

	// Check if custom home is provided
	appPath := flagGetPath(cmd)
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return chainconfig.Config{}, err
	}

	_, gappPath, err := gomodulepath.Find(absPath)
	if err != nil {
		return chainconfig.Config{}, err
	}

	configPath, err = chainconfig.LocateDefault(gappPath)
	if err != nil {
		return chainconfig.Config{}, err
	}

	return chainconfig.ParseFile(configPath)
}
