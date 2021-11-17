package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// NewScaffoldPluginConfig Read Config module information from yml file
func NewScaffoldPluginConfig() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugins --plugins [pluginName]",
		Short: "Scaffold function that call up plugins and check its availability",
		Long:  "Scaffold function that call up plugins and check its available plugins info & its validity",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("must read config for configs")
			//TODO Read config files to get
			return nil
		},
	}

	flagSetPath(c)
	c.Flags().String(flagModule, "", "IBC Module to add the packet into")
	c.Flags().String(flagSigner, "", "Label for the message signer (default: creator)")

	return c
}

func isInstalled(cmd *cobra.Command, args []string) error {

	fmt.Printf("all plugins list")

	return nil
}

func install(cmd *cobra.Command, args []string) error {

	return nil
}
