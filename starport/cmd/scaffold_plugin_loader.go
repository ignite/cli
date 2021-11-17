package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewScaffoldPluginLoader Read Config module information from yml file
func NewScaffoldPluginLoader() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugins --name [Config]",
		Short: "Scaffold an plugins command to call up configs",
		Long:  "Scaffold an plugins command to call up configs regarding plugins info & plugins infos",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("must read config for configs")
			return nil
		},
	}

	flagSetPath(c)
	c.Flags().String(flagModule, "", "IBC Module to add the packet into")
	c.Flags().String(flagSigner, "", "Label for the message signer (default: creator)")

	return c
}

func LoadConfig(cmd *cobra.Command, args []string) error {
	fmt.Printf("all plugins list")

	return nil
}

func LoadPlugin(cmd *cobra.Command, args []string) error {
	fmt.Printf("all plugins list")

	return nil
}
