package starportcmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [launch-id] [validator-account]",
		Short: "Initialize a chain from a published chain ID",
		Args:  cobra.ExactArgs(2),
		RunE:  networkChainInitHandler,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	nb, s, endRoutine, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer endRoutine()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}

	// check if the provided account for the validator exists
	validatorName := args[1]
	if err := checkAccountExist(cmd, validatorName); err != nil {
		return err
	}

	// if a chain has already been initialized with this launch ID, we ask for confirmation before erasing the directory
	chainHome, exist, err := checkChainHomeExist(launchID)
	if err != nil {
		return err
	}
	if !getYes(cmd) && exist {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("The chain has already been initialized under: %s. Would you like to overwrite the home directory",
				chainHome,
			),
			IsConfirm: true,
		}
		s.Stop()
		if _, err := prompt.Run(); err != nil {
			fmt.Println("said no")
			return nil
		}
		s.Start()
	}

	// initialize the blockchain from the launch ID
	initOptions := initOptionWithHomeFlag(cmd, []network.InitOption{network.MustNotInitializedBefore()})
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if err := blockchain.Init(cmd.Context()); err != nil {
		return err
	}

	// ask validator information
	v, err := askValidatorInfo(validatorName)
	if err != nil {
		return err
	}

	gentxPath, err := blockchain.InitAccount(cmd.Context(), v, validatorName)
	if err != nil {
		return err
	}
	fmt.Printf("%s Gentx generated: %s\n", clicpinner.Bullet, gentxPath)

	return nil
}

// askValidatorInfo prompts to the user questions to query validator information
func askValidatorInfo(validatorName string) (chain.Validator, error) {
	// TODO: allowing more customization for the validator
	v := chain.Validator{
		Name:              validatorName,
		Moniker:           validatorName,
		GasPrices:         "0stake",
		MinSelfDelegation: "1",
	}

	questions := append([]cliquiz.Question{},
		cliquiz.NewQuestion("Staking amount",
			&v.StakingAmount,
			cliquiz.DefaultAnswer("95000000stake"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission rate",
			&v.CommissionRate,
			cliquiz.DefaultAnswer("0.10"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max rate",
			&v.CommissionMaxRate,
			cliquiz.DefaultAnswer("0.20"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max change rate",
			&v.CommissionMaxChangeRate,
			cliquiz.DefaultAnswer("0.01"),
			cliquiz.Required(),
		),
	)
	return v, cliquiz.Ask(questions...)
}

// checkChainHomeExist checks if a home with the provided launchID already exist
func checkChainHomeExist(launchID uint64) (string, bool, error) {
	home, err := network.ChainHome(launchID)
	if err != nil {
		return home, false, err
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		return home, false, nil
	}
	return home, true, err
}
