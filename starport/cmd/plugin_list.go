package starportcmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/pluginsrpc"
)

const (
	flagState = "state"
)

// NewPluginList creates a new list command to retrieve plugins.
func NewPluginList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List built plugins.",
		Long:  "List plugins specified in the config file that are also loaded.",
		RunE:  pluginListHandler,
	}

	c.Flags().StringP(flagConfig, "c", "", "Starport config file (default: ./config.yml)")
	c.Flags().StringP(flagState, "s", "configured", "Plugin state (configured, downloaded, built)")

	return c
}

func pluginListHandler(cmd *cobra.Command, args []string) error {
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

	pluginState, err := cmd.Flags().GetString("state")
	if err != nil {
		return err
	}
	log.Println(pluginState)

	pluginManager := pluginsrpc.NewManager(chainId, chainConfig)
	plugins, err := pluginManager.List(cmd.Context(), pluginsrpc.PluginStateFromString(pluginState))
	if err != nil {
		return err
	}

	var output string
	if len(plugins) > 0 {
		for _, plugin := range plugins {
			output += fmt.Sprintf("%s\n", plugin)
		}
	} else {
		output = "None"
	}

	log.Println("Plugins: \n", output)
	return nil
}
