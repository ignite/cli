package starportcmd

import (
	"fmt"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/services/chain"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagRecover  = "recover"
	flagMnemonic = "mnemomic"
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

	// TODO create flagset
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	var (
		//	recover, _     = cmd.Flags().GetBool(flagRecover)
			mnemonic, _       = cmd.Flags().GetString(flagMnemonic)
		keyName, _ = cmd.Flags().GetString(flagKeyName)
	)

	s := clispinner.New()
	defer s.Stop()

	var (
		wg sync.WaitGroup
		ev = events.NewBus()
	)
	wg.Add(1)

	defer wg.Wait()
	defer ev.Shutdown()

	go printEvents(&wg, ev, s)

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing launchID: %s", err.Error())
	}

	nb, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
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

	acc, gentxPath, err := blockchain.InitAccount(cmd.Context(), v, keyName, mnemonic)
	if err != nil {
		return err
	}
	fmt.Printf("account created: %s\n%s\n", acc.Address, acc.Mnemonic)
	fmt.Printf("gentx generated: %s\n", gentxPath)

	return nil
}

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
