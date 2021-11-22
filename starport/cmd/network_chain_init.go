package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagValidatorAccount = "validator-account"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [launch-id]",
		Short: "Initialize a chain from a published chain ID",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInitHandler,
	}

	c.Flags().String(flagValidatorAccount, cosmosaccount.DefaultAccount, "Account for the chain validator")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// check if the provided account for the validator exists
	validatorAccount, _ := cmd.Flags().GetString(flagValidatorAccount)
	_, err = nb.AccountRegistry().GetByName(validatorAccount)
	if err != nil {
		return err
	}

	// if a chain has already been initialized with this launch ID, we ask for confirmation before erasing the directory
	chainHome, exist, err := cosmosutil.IsChainHomeExist(launchID)
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
	initOptions := initOptionWithHomeFlag(cmd, network.MustNotInitializedBefore())
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if err := blockchain.Init(cmd.Context()); err != nil {
		return err
	}

	// ask validator information
	v, err := askValidatorInfo(validatorAccount)
	if err != nil {
		return err
	}

	gentxPath, err := blockchain.InitAccount(cmd.Context(), v, validatorAccount)
	if err != nil {
		return err
	}
	fmt.Printf("%s Gentx generated: %s\n", clispinner.Bullet, gentxPath)

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
