package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/pluginsrpc"
)

// NewPluginReload creates a new reload command to manually refresh chain plugins.
func NewPluginPull() *cobra.Command {
	c := &cobra.Command{
		Use:   "pull",
		Short: "Pull plugins specified in config file.",
		RunE:  pluginPullHandler,
	}

	c.Flags().StringP(flagConfig, "c", "", "Starport config file (default: ./config.yml)")

	return c
}

func pluginPullHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Pulling plugins...")
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

	chainConfig, err := getCommandConfig(cmd)
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

	if err := pluginManager.Pull(cmd.Context()); err != nil {
		return err
	}

	fmt.Println("🔄  Pulled plugins.")
	return nil
}
