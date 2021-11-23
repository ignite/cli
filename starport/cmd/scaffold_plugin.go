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
	pluginHome   = fmt.Sprintf("%s/.starport/plugins", os.Getenv("HOME")) // TODO:
	pluginLoader plugin.Loader
)

func init() {
	var err error

	pluginLoader, err = plugin.NewLoader()
	if err != nil {
		log.Panic(err)
	}
}

// NewScaffoldPlugins returns the command to plugin.
func NewScaffoldPlugins(pluginConfigs []chainconfig.Plugin) []*cobra.Command {
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
					Short: f.Name, // TODO: Any alternatives?
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
