package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewApp creates new command named `app` to create Cosmos scaffolds customized
// by the user given options.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [github.com/org/repo]",
		Short: "Generates an empty application",
		Args:  cobra.ExactArgs(1),
		RunE:  appHandler,
	}
	c.Flags().String("address-prefix", "cosmos", "Address prefix")
	return c
}

func appHandler(cmd *cobra.Command, args []string) error {
	name := args[0]
	addressPrefix, _ := cmd.Flags().GetString("address-prefix")
	sc := scaffolder.New("")
	path, err := sc.Init(name, scaffolder.AddressPrefix(addressPrefix))
	if err != nil {
		return err
	}
	message := `
⭐️ Successfully created a Cosmos app '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

NOTE: add --verbose flag for verbose (detailed) output.
`
	fmt.Printf(message, path)
	return nil
}
