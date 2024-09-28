package ignitecmd

import (
	"math/rand"
	"strconv"

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

	return testnetMultiNode(cmd, session)
}

func testnetMultiNode(cmd *cobra.Command, session *cliui.Session) error {
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

	numVal, amountDetails, err := getValidatorAmountStake(cfg.MultiNode)
	if err != nil {
		return err
	}
	args := chain.MultiNodeArgs{
		ChainID:               cfg.MultiNode.ChainID,
		ValidatorsStakeAmount: amountDetails,
		OutputDir:             cfg.MultiNode.OutputDir,
		NumValidator:          strconv.Itoa(numVal),
	}

	return c.TestnetMultiNode(cmd.Context(), args)
}

// getValidatorAmountStake returns the number of validators and the amountStakes arg from config.MultiNode
func getValidatorAmountStake(cfg base.MultiNode) (int, string, error) {
	var amounts string
	count := 0

	if len(cfg.Validators) == 0 {
		numVal := cfg.RandomValidators.Count
		minStake, err := sdk.ParseCoinNormalized(cfg.RandomValidators.MinStake)
		if err != nil {
			return count, amounts, err
		}
		maxStake, err := sdk.ParseCoinNormalized(cfg.RandomValidators.MaxStake)
		if err != nil {
			return count, amounts, err
		}
		minS := minStake.Amount.Uint64()
		maxS := maxStake.Amount.Uint64()
		for i := 0; i < numVal; i++ {
			stakeAmount := minS + rand.Uint64()%(maxS-minS+1)
			if amounts == "" {
				amounts = math.NewIntFromUint64(stakeAmount).String()
				count += 1
			} else {
				amounts = amounts + "," + math.NewIntFromUint64(stakeAmount).String()
				count += 1
			}
		}
	} else {
		for _, v := range cfg.Validators {
			stakeAmount, err := sdk.ParseCoinNormalized(v.Stake)
			if err != nil {
				return count, amounts, err
			}
			if amounts == "" {
				amounts = stakeAmount.Amount.String()
				count += 1
			} else {
				amounts = amounts + "," + stakeAmount.Amount.String()
				count += 1
			}
		}
	}

	return count, amounts, nil
}
