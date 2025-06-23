package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

// NewScaffoldChainRegistry returns the command to scaffold the chain registry chain.json and assets.json files.
func NewScaffoldChainRegistry() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain-registry",
		Short: "Configs for the chain registry",
		Long: `Scaffold the chain registry chain.json and assets.json files.

The chain registry is a GitHub repo, hosted at https://github.com/cosmos/chain-registry, that
contains the chain.json and assets.json files of most of chains in the Cosmos ecosystem.
It is good practices, when creating a new chain, and about to launch a testnet or mainnet, to
publish the chain's metadata in the chain registry.

Read more about the chain.json at https://github.com/cosmos/chain-registry?tab=readme-ov-file#chainjson
Read more about the assets.json at https://github.com/cosmos/chain-registry?tab=readme-ov-file#assetlists`,
		Args:    cobra.NoArgs,
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldChainRegistryFiles,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func scaffoldChainRegistryFiles(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	c, err := chain.NewWithHomeFlags(cmd)
	if err != nil {
		return err
	}

	appPath := flagGetPath(cmd)
	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	if err = sc.CreateChainRegistryFiles(c, cfg); err != nil {
		return err
	}

	// no need for post scaffolding, as we are just creating two files
	// that are not part of the build process

	session.Printf("🎉 chain-registry files successfully scaffolded\n")

	return nil
}
