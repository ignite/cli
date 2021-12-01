package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/pluginsrpc"
)

// NewPluginReload creates a new reload command to manually refresh chain plugins.
func NewPluginBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Build plugins specified in config file.",
		RunE:  pluginBuildHandler,
	}

	flagSetPath(c)
	c.Flags().StringP(flagConfig, "c", "", "Starport config file (default: ./config.yml)")

	return c
}

func pluginBuildHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Building plugins...")
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

	pluginManager, err := pluginsrpc.NewManager(chainId, chainConfig)
	if err != nil {
		return err
	}

	if err := pluginManager.Build(cmd.Context()); err != nil {
		return err
	}

	fmt.Println("🔄  Built plugins.")
	return nil
}
