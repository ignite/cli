package ignitecmd

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdmodel "github.com/ignite/cli/v29/ignite/cmd/model"
	"github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	flagNodeDirPrefix = "node-dir-prefix"
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
				output-dir: ./.testchain-testnet/
				chain-id: testchain-test-1
				node-dir-prefix: validator

	or random amount stake
			....
			multi-node:
				random_validators:
					count: 4
					min_stake: 50000000stake
					max_stake: 150000000stake
				output-dir: ./.testchain-testnet/
				chain-id: testchain-test-1
				node-dir-prefix: validator


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
	c.Flags().BoolP(flagResetOnce, "r", false, "reset the app state once on init")
	c.Flags().String(flagNodeDirPrefix, "validator", "prefix of dir node")

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

	numVal, amountDetails, err := getValidatorAmountStake(cfg.Validators)
	if err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)

	outputDir := filepath.Join(homeDir, ".ignite/local-chains/"+c.Name()+"d/testnet/")
	args := chain.MultiNodeArgs{
		OutputDir:             outputDir,
		NumValidator:          strconv.Itoa(numVal),
		ValidatorsStakeAmount: amountDetails,
		NodeDirPrefix:         nodeDirPrefix,
	}

	resetOnce, _ := cmd.Flags().GetBool(flagResetOnce)
	if resetOnce {
		// If resetOnce is true, the app state will be reset by deleting the output directory.
		err := os.RemoveAll(outputDir)
		if err != nil {
			return err
		}
	}

	err = c.TestnetMultiNode(cmd.Context(), args)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	m := cmdmodel.NewModel(cmd.Context(), c.Name(), args)
	_, err = tea.NewProgram(m).Run()
	return err
}

// getValidatorAmountStake returns the number of validators and the amountStakes arg from config.MultiNode.
func getValidatorAmountStake(validators []v1.Validator) (int, string, error) {
	numVal := len(validators)
	var amounts string

	for _, v := range validators {
		stakeAmount, err := sdk.ParseCoinNormalized(v.Bonded)
		if err != nil {
			return numVal, amounts, err
		}
		if amounts == "" {
			amounts = stakeAmount.Amount.String()
		} else {
			amounts = amounts + "," + stakeAmount.Amount.String()
		}
	}

	return numVal, amounts, nil
}
