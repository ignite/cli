package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagNoDefaultModule = "no-module"
)

// NewScaffoldChain creates new command to scaffold a Comos-SDK based blockchain.
func NewScaffoldChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain [github.com/org/repo]",
		Short: "Fully-featured Cosmos SDK blockchain",
		Long:  "Scaffold a new Cosmos SDK blockchain with a default directory structure",
		Args:  cobra.ExactArgs(1),
		RunE:  scaffoldChainHandler,
	}

	c.Flags().StringP(flagPath, "p", ".", "path to scaffold the chain")
	c.Flags().String(flagAddressPrefix, "cosmos", "Address prefix")
	c.Flags().Bool(flagNoDefaultModule, false, "Prevent scaffolding a default module in the app")

	return c
}

func scaffoldChainHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		name               = args[0]
		addressPrefix, _   = cmd.Flags().GetString(flagAddressPrefix)
		noDefaultModule, _ = cmd.Flags().GetBool(flagNoDefaultModule)
		appPath            = flagGetPath(cmd)
	)

	appdir, err := scaffolder.Init(placeholder.New(), appPath, name, addressPrefix, noDefaultModule)
	if err != nil {
		return err
	}

	s.Stop()

	path, err := relativePath(appdir)
	if err != nil {
		return err
	}

	message := `
⭐️ Successfully created a new blockchain '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport chain serve

Documentation: https://docs.starport.network
`
	fmt.Printf(message, path)

	return nil
}
