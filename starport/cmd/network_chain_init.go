package starportcmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ignite-hq/cli/starport/pkg/cliquiz"
	"github.com/ignite-hq/cli/starport/pkg/clispinner"
	"github.com/ignite-hq/cli/starport/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/starport/services/chain"
	"github.com/ignite-hq/cli/starport/services/network"
	"github.com/ignite-hq/cli/starport/services/network/networkchain"
)

const (
	flagValidatorAccount         = "validator-account"
	flagValidatorWebsite         = "validator-website"
	flagValidatorDetails         = "validator-details"
	flagValidatorSecurityContact = "validator-security-contact"
	flagValidatorMoniker         = "validator-moniker"
	flagValidatorIdentity        = "validator-identity"
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
	c.Flags().String(flagValidatorWebsite, "", "Associate a website to the validator")
	c.Flags().String(flagValidatorDetails, "", "Provide details about the validator")
	c.Flags().String(flagValidatorSecurityContact, "", "Provide a validator security contact email")
	c.Flags().String(flagValidatorMoniker, "", "Provide a custom validator moniker")
	c.Flags().String(flagValidatorIdentity, "", "Provide a validator identity signature (ex. UPort or Keybase)")
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
	launchID, err := network.ParseLaunchID(args[0])
	if err != nil {
		return err
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

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
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
		identity, _        = cmd.Flags().GetString(flagValidatorIdentity)
	)

	v := chain.Validator{
		Name:              account,
		Website:           website,
		Details:           details,
		Moniker:           moniker,
		Identity:          identity,
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
