package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewApp creates new command named `app` to create Cosmos scaffolds customized
// by the user given options.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [github.com/org/repo]",
		Short: "Scaffold a new blockchain",
		Long: "Scaffold a new Cosmos SDK blockchain with a default directory structure",
		Args:  cobra.ExactArgs(1),
		RunE:  appHandler,
	}
	c.Flags().String("address-prefix", "cosmos", "Address prefix")
	return c
}

func appHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		name             = args[0]
		addressPrefix, _ = cmd.Flags().GetString("address-prefix")
	)

	sc, err := scaffolder.New("",
		scaffolder.AddressPrefix(addressPrefix),
	)
	if err != nil {
		return err
	}

	appdir, err := sc.Init(name)
	if err != nil {
		return err
	}

	s.Stop()

	message := `
⭐️ Successfully created a new blockchain '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

Documentation: https://docs.starport.network
`
	fmt.Printf(message, appdir)

	return nil
}
