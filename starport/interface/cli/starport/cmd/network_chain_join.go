package starportcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [chain-id]",
		Short: "Propose to join to a network as a validator",
		RunE:  networkChainJoinHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	var (
		chainID = args[0]
	)

	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	nb, err := newNetworkBuilder(networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	// init the blockchain.
	blockchain, err := nb.InitBlockchainFromChainID(cmd.Context(), chainID, false)

	if err == context.Canceled {
		s.Stop()
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	// get blockchain's info.
	info, err := blockchain.Info()
	if err != nil {
		return err
	}

	// hold default values and user inputs for target chain to later use these to join to the chain.
	var (
		proposal      networkbuilder.Proposal
		account       networkbuilder.Account
		denom         string
		publicAddress string
	)
	if info.Config.Validator.Staked != "" {
		if c, err := types.ParseCoin(info.Config.Validator.Staked); err == nil {
			denom = c.Denom
		}
	}
	acc, _ := info.Config.AccountByName(info.Config.Validator.Name)

	// ask to create an account on target blockchain.
	printSection(fmt.Sprintf("Account on the blockchain %s", chainID))
	account.Name, err = createAccount(nb, fmt.Sprintf("%s blockchain", chainID))
	if err != nil {
		return err
	}

	// ask to create an account proposal,
	printSection("Account proposal")

	if err := cliquiz.Ask(
		cliquiz.NewQuestion("Account coins",
			&account.Coins,
			cliquiz.DefaultAnswer(strings.Join(acc.Coins, ",")),
		),
	); err != nil {
		return err
	}

	// ask to create a validator proposal.
	fmt.Println()
	printSection("Validator proposal")

	questions := []cliquiz.Question{
		cliquiz.NewQuestion("Staking amount", &proposal.Validator.StakingAmount, cliquiz.DefaultAnswer("95000000stake")),
		cliquiz.NewQuestion("Moniker", &proposal.Validator.Moniker, cliquiz.DefaultAnswer("mynode")),
		cliquiz.NewQuestion("Commission rate", &proposal.Validator.CommissionRate, cliquiz.DefaultAnswer("0.10")),
		cliquiz.NewQuestion("Commission max rate", &proposal.Validator.CommissionMaxRate, cliquiz.DefaultAnswer("0.20")),
		cliquiz.NewQuestion("Commission max change rate", &proposal.Validator.CommissionMaxChangeRate, cliquiz.DefaultAnswer("0.01")),
		cliquiz.NewQuestion("Min self delegation", &proposal.Validator.MinSelfDelegation, cliquiz.DefaultAnswer("1")),
		cliquiz.NewQuestion("Gas prices", &proposal.Validator.GasPrices, cliquiz.DefaultAnswer("0.025"+denom)),
		cliquiz.NewQuestion("Website", &proposal.Meta.Website),
		cliquiz.NewQuestion("Identity", &proposal.Meta.Identity),
		cliquiz.NewQuestion("Details", &proposal.Meta.Details),
	}

	if !xchisel.IsEnabled() {
		opts := []cliquiz.Option{
			cliquiz.Required(),
		}
		ip, err := ipify.GetIp()
		if err == nil {
			opts = append(opts, cliquiz.DefaultAnswer(fmt.Sprintf("%s:26656", ip)))
		}
		questions = append(questions, cliquiz.NewQuestion("Peer's address", &publicAddress, opts...))
	}

	s.Stop()

	if err := cliquiz.Ask(questions...); err != nil {
		return err
	}
	gentx, a, err := blockchain.IssueGentx(cmd.Context(), account, proposal)
	if err != nil {
		return err
	}

	prettyGentx, err := gentx.Pretty()
	if err != nil {
		return err
	}

	// ask to confirm to join to network.
	fmt.Printf("\nYour validator's details (Gentx): \n\n%s\n\n", prettyGentx)
	fmt.Printf("Confirm joining to %s as a validator with the Gentx above:\n", chainID)

	var shouldJoin string
	var answerYes, answerNo = "yes", "no"

	for {
		if err := cliquiz.Ask(cliquiz.NewQuestion(fmt.Sprintf("Please type %q or %q", answerYes, answerNo),
			&shouldJoin,
			cliquiz.Required(),
		)); err != nil {
			return err
		}
		if shouldJoin == answerNo {
			s.Stop()
			fmt.Println("said no")
			return nil
		}
		if shouldJoin == answerYes {
			break
		}
	}

	// propose to join to the network.
	coins, err := types.ParseCoins(account.Coins) // parse the coins of the account
	if err != nil {
		return err
	}

	selfDelegation, err := types.ParseCoin(proposal.Validator.StakingAmount) // parse the self delegation of this account for the validator
	if err != nil {
		return err
	}

	s.SetText("Proposing...")
	s.Start()

	if err := blockchain.Join(cmd.Context(), a.Address, publicAddress, coins, gentx, selfDelegation); err != nil {
		return err
	}
	s.Stop()

	if a.Mnemonic != "" {
		fmt.Printf("\n*** IMPORTANT - Save your mnemonic in a secret place:\n%s\n", a.Mnemonic)
	}

	fmt.Println("\n📜 Proposal to join as a validator has been submitted successfully!")
	return nil
}
