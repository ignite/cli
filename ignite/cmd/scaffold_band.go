package ignitecmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldBandchain creates a new BandChain oracle in the module
func NewScaffoldBandchain() *cobra.Command {
	c := &cobra.Command{
		Use:   "band [queryName] --module [moduleName]",
		Short: "Scaffold an IBC BandChain query oracle to request real-time data",
		Long:  "Scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createBandchainHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().String(flagModule, "", "IBC Module to add the packet into")
	c.Flags().String(flagSigner, "", "Label for the message signer (default: creator)")

	return c
}

func createBandchainHandler(cmd *cobra.Command, args []string) error {
	var (
		oracle  = args[0]
		appPath = flagGetPath(cmd)
		signer  = flagGetSigner(cmd)
	)

	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}
	if module == "" {
		return errors.New("please specify a module to create the BandChain oracle into: --module <module_name>")
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	var options []scaffolder.OracleOption
	if signer != "" {
		options = append(options, scaffolder.OracleWithSigner(signer))
	}

	sc, err := newApp(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.AddOracle(cacheStorage, placeholder.New(), module, oracle, options...)
	if err != nil {
		return err
	}

	s.Stop()

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	fmt.Println(modificationsStr)

	fmt.Printf(`
🎉 Created a Band oracle query "%[1]v".

Note: BandChain module uses version "bandchain-1".
Make sure to update the keys.go file accordingly.

// x/%[2]v/types/keys.go
const Version = "bandchain-1"

`, oracle, module)

	return nil
}
