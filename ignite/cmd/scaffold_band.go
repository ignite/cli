package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/services/scaffolder"
)

const tplScaffoldBandSuccess = `
🎉 Created a Band oracle query "%[1]v".
Note: BandChain module uses version "bandchain-1".
Make sure to update the keys.go file accordingly.
// x/%[2]v/types/keys.go
const Version = "bandchain-1"
`

// NewScaffoldBandchain creates a new BandChain oracle in the module.
func NewScaffoldBandchain() *cobra.Command {
	c := &cobra.Command{
		Use:     "band [queryName] --module [moduleName]",
		Short:   "Scaffold an IBC BandChain query oracle to request real-time data",
		Long:    "Scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    createBandchainHandler,
		Hidden:  true,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().String(flagModule, "", "IBC module to add the packet into")
	c.Flags().String(flagSigner, "", "label for the message signer (default: creator)")

	return c
}

func createBandchainHandler(cmd *cobra.Command, args []string) error {
	var (
		oracle  = args[0]
		appPath = flagGetPath(cmd)
		signer  = flagGetSigner(cmd)
	)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

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
		options = append(options, scaffolder.OracleWithSigner(signer)) //nolint: staticcheck
	}

	sc, err := scaffolder.New(cmd.Context(), appPath)
	if err != nil {
		return err
	}

	//nolint: staticcheck
	err = sc.AddOracle(cmd.Context(), cacheStorage, placeholder.New(), module, oracle, options...)
	if err != nil {
		return err
	}

	modificationsStr, err := sc.ApplyModifications()
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cacheStorage, false); err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf(tplScaffoldBandSuccess, oracle, module)

	return nil
}
