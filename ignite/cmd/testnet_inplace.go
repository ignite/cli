package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/spf13/cobra"
)

func NewTestNetInPlace() *cobra.Command {
	c := &cobra.Command{
		Use:   "in-place",
		Short: "Create and start a testnet from current local state",
		Long: `Testnet in-place command is used to create and start a testnet from current local state.
		After utilizing this command the network will start. We can create testnet from mainnet state and mint more coins for accounts from config.yml file.

		In the config.yml file, there should be at least the address account to fund, operator address, home of the local state node.

		For example:

			Configuration acounts to fund:
				accounts: 
					- name: alice
					address: "cosmos1wa3u4knw74r598quvzydvca42qsmk6jrzmgy07"
					- name: bob
					address: "cosmos10uls38gddhhlywla0sjlvqg8pjvcffx4lu25c4"

			Configuration validators:
				validators:
					- name: alice
					home: "$HOME/.testchaind/validator1"
		`,
		Args: cobra.NoArgs,
		RunE: testnetInPlaceHandler,
	}
	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().AddFlagSet(flagSetVerbose())

	c.Flags().Bool(flagQuitOnFail, false, "quit program if the app fails to start")
	return c
}

func testnetInPlaceHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
	)
	defer session.End()

	// Otherwise run the serve command directly
	return testnetInplace(cmd, session)
}

func testnetInplace(cmd *cobra.Command, session *cliui.Session) error {
	chainOption := []chain.Option{
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.CheckCosmosSDKVersion(),
	}

	if flagGetCheckDependencies(cmd) {
		chainOption = append(chainOption, chain.CheckDependencies())
	}

	// check if custom config is defined
	config, _ := cmd.Flags().GetString(flagConfig)
	if config != "" {
		chainOption = append(chainOption, chain.ConfigFile(config))
	}

	c, err := chain.NewWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cfg, err := c.Config()
	if err != nil {
		return err
	}
	var operatorAddress sdk.ValAddress
	var accounts string
	for _, acc := range cfg.Accounts {
		if cfg.Validators[0].Name == acc.Name {
			accAddr, err := sdk.AccAddressFromBech32(acc.Address)
			if err != nil {
				return err
			}
			operatorAddress = sdk.ValAddress(accAddr)
		}
		accounts = accounts + "," + acc.Address
	}

	chainID, err := c.ID()
	if err != nil {
		return err
	}

	args := chain.InplaceArgs{
		NewChainID:         chainID,
		NewOperatorAddress: operatorAddress.String(),
		AcountsToFund:      accounts,
	}
	return c.TestNetInPlace(cmd.Context(), args)
}
