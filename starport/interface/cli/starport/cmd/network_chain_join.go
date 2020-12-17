package starportcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
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
	blockchain, err := nb.Init(cmd.Context(), chainID, networkbuilder.SourceChainID())

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
		account       chain.Account
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
	account, err = createChainAccount(cmd.Context(), blockchain, fmt.Sprintf("%s blockchain", chainID), acc.Name)
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
		cliquiz.NewQuestion("Staking amount", &proposal.Validator.StakingAmount,
			cliquiz.DefaultAnswer("95000000stake"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Moniker",
			&proposal.Validator.Moniker,
			cliquiz.DefaultAnswer("mynode"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission rate",
			&proposal.Validator.CommissionRate,
			cliquiz.DefaultAnswer("0.10"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max rate",
			&proposal.Validator.CommissionMaxRate,
			cliquiz.DefaultAnswer("0.20"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Commission max change rate",
			&proposal.Validator.CommissionMaxChangeRate,
			cliquiz.DefaultAnswer("0.01"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Min self delegation",
			&proposal.Validator.MinSelfDelegation,
			cliquiz.DefaultAnswer("1"),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Gas prices",
			&proposal.Validator.GasPrices,
			cliquiz.DefaultAnswer("0.025"+denom),
			cliquiz.Required(),
		),
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
	gentx, err := blockchain.IssueGentx(cmd.Context(), account, proposal)
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
		if err := cliquiz.Ask(cliquiz.NewQuestion(fmt.Sprintf("Please enter %q or %q", answerYes, answerNo),
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

	if err := blockchain.Join(cmd.Context(), account.Address, publicAddress, coins, gentx, selfDelegation); err != nil {
		return err
	}
	s.Stop()

	fmt.Println("\n📜  Proposal about joining as a validator has been successfully submitted!")
	return nil
}

func createChainAccount(ctx context.Context, blockchain *networkbuilder.Blockchain, title, defaultAccountName string) (account chain.Account, err error) {
	var (
		createAccount = "Create a new account"
		importAccount = "Import an account from mnemonic"
	)
	var (
		qs = []*survey.Question{
			{
				Name: "account",
				Prompt: &survey.Select{
					Message: "Choose an account:",
					Options: []string{createAccount, importAccount},
				},
			},
		}
		answers = struct {
			Account string `survey:"account"`
		}{}
	)
	if err := survey.Ask(qs, &answers); err != nil {
		if err == terminal.InterruptErr {
			return account, context.Canceled
		}
		return account, err
	}

	switch answers.Account {
	case createAccount:
		var name string
		if err := cliquiz.Ask(cliquiz.NewQuestion("Account name", &name, cliquiz.DefaultAnswer(defaultAccountName), cliquiz.Required())); err != nil {
			return account, err
		}

		if account, err = blockchain.CreateAccount(ctx, chain.Account{Name: name}); err != nil {
			return account, err
		}

		fmt.Printf("\n%s account has been created successfully!\nAccount address: %s \nMnemonic: %s\n\n",
			title,
			account.Address,
			account.Mnemonic,
		)

	case importAccount:
		var name string
		var mnemonic string
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Account name", &name, cliquiz.DefaultAnswer(defaultAccountName), cliquiz.Required()),
			cliquiz.NewQuestion("Mnemonic", &mnemonic, cliquiz.Required()),
		); err != nil {
			return account, err
		}

		if account, err = blockchain.CreateAccount(ctx, chain.Account{
			Name:     name,
			Mnemonic: mnemonic,
		}); err != nil {
			return account, err
		}
		fmt.Printf("\n%s account has been imported successfully!\nAccount address: %s\n\n", title, account.Address)
	}
	return account, nil
}
