package starportcmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin"
)

var (
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

	for _, p := range pluginConfigs {
		if pluginLoader.IsInstalled(p) {
			cmds = append(cmds, &cobra.Command{
				Use:   p.Name,
				Short: p.Description,
				RunE: func(cmd *cobra.Command, args []string) error {
					// TODO: Run plugin here.
					fmt.Println("Run plugin ", cmd.Use)
					return nil
				},
			})
		}
	}

	return cmds
}
