package ignitecmd

import (
	"fmt"
	"math/rand"

	"cosmossdk.io/math"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
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
	or
			....
			multi-node:
				random_validators:
				count: 4
				min_stake: 50000000stake
				max_stake: 150000000stake

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

	validatorDetails, err := getValidatorAmountStake(cfg.MultiNode)
	if err != nil {
		return err
	}
	fmt.Println(validatorDetails)
	fmt.Println(cfg.MultiNode.OutputDir)

	return nil
}

func getValidatorAmountStake(cfg base.MultiNode) ([]math.Int, error) {
	var amounts []math.Int

	if len(cfg.Validators) == 0 {
		numVal := cfg.RandomValidators.Count
		minStake, err := sdk.ParseCoinNormalized(cfg.RandomValidators.MinStake)
		if err != nil {
			return amounts, err
		}
		maxStake, err := sdk.ParseCoinNormalized(cfg.RandomValidators.MaxStake)
		if err != nil {
			return amounts, err
		}
		minS := minStake.Amount.Uint64()
		maxS := maxStake.Amount.Uint64()
		for i := 0; i < numVal; i++ {
			stakeAmount := minS + rand.Uint64()%(maxS-minS+1)
			amounts = append(amounts, math.NewIntFromUint64(stakeAmount))
		}
	} else {
		for _, v := range cfg.Validators {
			stakeAmount, err := sdk.ParseCoinNormalized(v.Stake)
			if err != nil {
				return amounts, err
			}
			amounts = append(amounts, stakeAmount.Amount)
		}
	}

	return amounts, nil
}
