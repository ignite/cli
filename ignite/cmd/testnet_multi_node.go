package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

func NewTestnetMultiNode() *cobra.Command {
	c := &cobra.Command{
		Use:   "multi-node",
		Short: "Create a network test multi node",
		Long: `Create a test network with the number of nodes from the config.yml file:
			...
			multi-node:
				validators:
					- name: validator1
					stake: 100000000stake
					- name: validator2
					stake: 200000000stake
					- name: validator3
					stake: 200000000stake
					- name: validator4
					stake: 200000000stake

		`,
		Args: cobra.NoArgs,
		RunE: testnetMultiNodeHandler,
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

func testnetMultiNodeHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
	)
	defer session.End()

	// Otherwise run the serve command directly
	return testnetInplace1(cmd, session)
}

func testnetInplace1(cmd *cobra.Command, session *cliui.Session) error {
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
	keyringBackend, err := c.KeyringBackend()
	if err != nil {
		return err
	}
	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringBackend(keyringBackend)),
		cosmosaccount.WithHome(home),
	)
	if err != nil {
		return err
	}

	var (
		operatorAddress sdk.ValAddress
		accounts        string
		accErr          *cosmosaccount.AccountDoesNotExistError
	)

	fmt.Println(cfg.MultiNode.Validators)

	for _, acc := range cfg.Accounts {
		sdkAcc, err := ca.GetByName(acc.Name)
		if errors.As(err, &accErr) {
			sdkAcc, _, err = ca.Create(acc.Name)
		}
		if err != nil {
			return err
		}

		sdkAddr, err := sdkAcc.Address(getAddressPrefix(cmd))
		if err != nil {
			return err
		}
		if len(cfg.Validators) == 0 {
			return errors.Errorf("no validators found for account %s", sdkAcc.Name)
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

	args := chain.InPlaceArgs{
		NewChainID:         chainID,
		NewOperatorAddress: operatorAddress.String(),
		AccountsToFund:     accounts,
	}
	return c.TestnetInPlace(cmd.Context(), args)
}
