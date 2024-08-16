package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/spf13/cobra"
)

func NewTestNetInPlace() *cobra.Command {
	c := &cobra.Command{
		Use:   "in-place",
		Short: "Create and start a testnet from current local state",
		Long: `Testnet in-place command is used to create and start a testnet from current local state.
		After utilizing this command the network will start. We can create testnet from mainnet state and mint more coins for accounts from config.yml file.
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
	home, err := c.Home()
	if err != nil {
		return err
	}
	keyringbankend, err := c.KeyringBackend()
	if err != nil {
		return err
	}
	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringBackend(keyringbankend)),
		cosmosaccount.WithHome(home),
	)
	if err != nil {
		return err
	}

	var operatorAddress sdk.ValAddress
	var accounts string
	for _, acc := range cfg.Accounts {
		var sdkAcc cosmosaccount.Account
		if sdkAcc, err = ca.GetByName(acc.Name); err != nil {
			sdkAcc, _, err = ca.Create(acc.Name)
			if err != nil {
				return err
			}
		}
		sdkAddr, err := sdkAcc.Address(getAddressPrefix(cmd))
		if err != nil {
			return err
		}
		if cfg.Validators[0].Name == acc.Name {
			accAddr, err := sdk.AccAddressFromBech32(sdkAddr)
			if err != nil {
				return err
			}
			operatorAddress = sdk.ValAddress(accAddr)
		}
		accounts = accounts + "," + sdkAddr
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
