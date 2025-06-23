package ignitecmd

import (
	"os"
	"path"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdmodel "github.com/ignite/cli/v29/ignite/cmd/bubblemodel"
	igcfg "github.com/ignite/cli/v29/ignite/config"
	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	flagNodeDirPrefix = "node-dir-prefix"
)

func NewTestnetMultiNode() *cobra.Command {
	c := &cobra.Command{
		Use:   "multi-node",
		Short: "Initialize and provide multi-node on/off functionality",
		Long: `Initialize the test network with the number of nodes and bonded from the config.yml file::
			...
                  validators:
                        - name: alice
                        bonded: 100000000stake
                        - name: validator1
                        bonded: 100000000stake
                        - name: validator2
                        bonded: 200000000stake
                        - name: validator3
                        bonded: 300000000stake


			The "multi-node" command allows developers to easily set up, initialize, and manage multiple nodes for a 
			testnet environment. This command provides full flexibility in enabling or disabling each node as desired, 
			making it a powerful tool for simulating a multi-node blockchain network during development.

			Usage:
					ignite testnet multi-node [flags]

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

	return c
}

func testnetMultiNodeHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
		cliui.WithoutUserInteraction(getYes(cmd)),
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
	nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)

	outputDir, err := xfilepath.Join(igcfg.DirPath, xfilepath.Path(path.Join("local-chains", c.Name(), "testnet")))()
	if err != nil {
		return err
	}

	ports, err := availableport.Find(uint(numVal)) //nolint:gosec,nolintlint // conversion is fine
	if err != nil {
		return err
	}

	args := chain.MultiNodeArgs{
		OutputDir:             outputDir,
		NumValidator:          strconv.Itoa(numVal),
		ValidatorsStakeAmount: amountDetails,
		NodeDirPrefix:         nodeDirPrefix,
		ListPorts:             ports,
	}

	resetOnce, _ := cmd.Flags().GetBool(flagResetOnce)
	if resetOnce {
		// If resetOnce is true, the app state will be reset by deleting the output directory.
		if err := os.RemoveAll(outputDir); err != nil {
			return err
		}
	}

	if err = c.TestnetMultiNode(cmd.Context(), args); err != nil {
		return err
	}

	model, err := cmdmodel.NewModel(cmd.Context(), c.Name(), args)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(model, tea.WithInput(cmd.InOrStdin())).Run()
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
