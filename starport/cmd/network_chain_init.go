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
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

const (
	flagValidatorAccount         = "validator-account"
	flagValidatorWebsite         = "validator-website"
	flagValidatorDetails         = "validator-details"
	flagValidatorSecurityContact = "validator-security-contact"
	flagValidatorMoniker         = "validator-moniker"
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
	c.Flags().String(flagValidatorWebsite, "", "Add validator website")
	c.Flags().String(flagValidatorDetails, "", "Add validator description")
	c.Flags().String(flagValidatorSecurityContact, "", "Add validator Security Contact")
	c.Flags().String(flagValidatorMoniker, "", "Add validator moniker")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// check if the provided account for the validator exists.
	validatorAccount, _ := cmd.Flags().GetString(flagValidatorAccount)
	if _, err = nb.AccountRegistry.GetByName(validatorAccount); err != nil {
		return err
	}

	// if a chain has already been initialized with this launch ID, we ask for confirmation
	// before erasing the directory.
	chainHome, exist, err := networkchain.IsChainHomeExist(launchID)
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
		nb.Spinner.Stop()
		if _, err := prompt.Run(); err != nil {
			fmt.Println("said no")
			return nil
		}
		nb.Spinner.Start()
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	launchInfo, err := n.LaunchInfo(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(launchInfo))
	if err != nil {
		return err
	}

	if err := c.Init(cmd.Context()); err != nil {
		return err
	}

	// ask validator information.
	v, err := askValidatorInfo(cmd)
	if err != nil {
		return err
	}

	gentxPath, err := c.InitAccount(cmd.Context(), v, validatorAccount)
	if err != nil {
		return err
	}
	fmt.Printf("%s Gentx generated: %s\n", clispinner.Bullet, gentxPath)

	return nil
}

// askValidatorInfo prompts to the user questions to query validator information
func askValidatorInfo(cmd *cobra.Command) (chain.Validator, error) {
	var (
		account, _         = cmd.Flags().GetString(flagValidatorAccount)
		website, _         = cmd.Flags().GetString(flagValidatorWebsite)
		details, _         = cmd.Flags().GetString(flagValidatorDetails)
		securityContact, _ = cmd.Flags().GetString(flagValidatorSecurityContact)
		moniker, _         = cmd.Flags().GetString(flagValidatorMoniker)
	)

	v := chain.Validator{
		Name:              account,
		Website:           website,
		Details:           details,
		Moniker:           moniker,
		SecurityContact:   securityContact,
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
