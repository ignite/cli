package starportcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldBandchain creates a new Bandchain oracle in the module
func NewScaffoldBandchain() *cobra.Command {
	c := &cobra.Command{
		Use:   "band [queryName] --module [moduleName]",
		Short: "Scaffold an IBC Bandchain query oracle to request real-time data",
		Long:  "Scaffold an IBC Bandchain query oracle to request real-time data from Bandchain scripts in a specific IBC-enabled Cosmos SDK module",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createBandchainHandler,
	}

	c.Flags().String(flagModule, "", "IBC Module to add the packet into")

	return c
}

func createBandchainHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	oracle := args[0]
	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}
	if module == "" {
		return errors.New("please specify a module to create the packet into: --module <module_name>")
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddOracle(placeholder.New(), module, oracle)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf(`
🎉 Created a Band oracle query "%[1]v".

Note: BandChain module uses version "bandchain-1".
Make sure to update the keys.go file accordingly.

// x/bandmodule/types/keys.go
const Version = "bandchain-1"

`, args[0])
	return nil
}
