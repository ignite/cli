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

const (
	flagRecover  = "recover"
	flagMnemonic = "mnemonic"
	flagKeyName  = "key-name"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [launch-id]",
		Short: "Initialize a chain from a published chain ID",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInitHandler,
	}

	c.Flags().Bool(flagRecover, false, "Recover chain account from a mnemonic")
	c.Flags().String(flagMnemonic, "", "Mnemonic for recovered account")
	c.Flags().String(flagKeyName, "", "key name for the chain account")

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	var (
		recover, _  = cmd.Flags().GetBool(flagRecover)
		mnemonic, _ = cmd.Flags().GetString(flagMnemonic)
		keyName, _  = cmd.Flags().GetString(flagKeyName)
	)
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
	v, err := askValidatorInfo()
	if err != nil {
		return err
	}

	// ask user for mnemonic if not provided and recover option defined
	if mnemonic == "" && recover {
		cliquiz.Ask(cliquiz.NewQuestion("Enter your account mnemonic",
			&mnemonic,
			cliquiz.DefaultAnswer(""),
			cliquiz.Required(),
		))
	}

	acc, gentxPath, err := blockchain.InitAccount(cmd.Context(), v, keyName, mnemonic)
	if err != nil {
		return err
	}
	fmt.Printf(`Account created
Address: %s
Mnemonic: %s
`, acc.Address, acc.Mnemonic)
	fmt.Printf("Gentx generated: %s\n", gentxPath)

	return nil
}

// askValidatorInfo prompts to the user questions to query validator information
func askValidatorInfo() (v chain.Validator, err error) {
	questions := append([]cliquiz.Question{},
		cliquiz.NewQuestion("Staking amount",
			&v.StakingAmount,
			cliquiz.DefaultAnswer("95000000stake"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Moniker",
			&v.Moniker,
			cliquiz.DefaultAnswer("mynode"),
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
		cliquiz.NewQuestion("Min self delegation",
			&v.MinSelfDelegation,
			cliquiz.DefaultAnswer("1"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Gas prices",
			&v.GasPrices,
			cliquiz.DefaultAnswer("0.025stake"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Details", &v.Details),
		cliquiz.NewQuestion("Identity", &v.Identity),
		cliquiz.NewQuestion("Website", &v.Website),
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
