package starportcmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin"
)

var (
	pluginHome   = fmt.Sprintf("%s/.starport/plugins", os.Getenv("HOME"))
	pluginLoader plugin.Loader
)

// NewScaffoldPlugins returns the command to plugin.
func NewScaffoldPlugins(chainID string, pluginConfigs []chainconfig.Plugin) []*cobra.Command {
	var err error

	if pluginLoader == nil {
		pluginLoader, err = plugin.NewLoader(chainID)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	cmds := make([]*cobra.Command, 0)

	for i, cfg := range pluginConfigs {
		i := i // Fix scopelint

		if pluginLoader.IsInstalled(cfg) {
			plugin, err := pluginLoader.LoadPlugin(pluginConfigs[i], pluginHome)
			if err != nil {
				log.Println(err)
				continue
			}

			pluginCmd := &cobra.Command{
				Use:   cfg.Name,
				Short: cfg.Description,
			}

			funcList := plugin.List()

			for _, f := range funcList {
				f := f // Fix scope

				pluginCmd.AddCommand(&cobra.Command{
					Use:   f.Name,
					Short: plugin.Help(f.Name),
					Long:  plugin.Help(f.Name),
					Args:  cobra.ExactArgs(len(f.ParamTypes)),
					RunE: func(cmd *cobra.Command, args []string) error {
						return plugin.Execute(f.Name, args)
					},
				})
			}

			cmds = append(cmds, pluginCmd)
		}
	}

	return cmds
}
