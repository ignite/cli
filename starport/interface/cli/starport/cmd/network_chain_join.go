package starportcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:  "join [chain-id]",
		RunE: networkChainJoinHandler,
		Args: cobra.ExactArgs(1),
	}
	c.Flags().StringP("repo", "r", ".", "repo of the blockchain")
	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	var (
		ctx = clictx.From(context.Background())
		ev  = events.NewBus()
		s   = spinner.New(spinner.CharSets[42], 100*time.Millisecond)
	)
	nb, err := networkbuilder.New(spnAddress, networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	s.Color("blue")
	defer s.Stop()

	go printEvents(ev, s)

	repo, _ := cmd.Flags().GetString("repo")
	blockchain, err := nb.InitBlockchain(ctx, repo)
	if err == context.Canceled {
		s.Stop()
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	info, err := blockchain.Info()
	if err != nil {
		return err
	}

	var (
		proposal networkbuilder.Proposal
		account  networkbuilder.Account
		denom    string
		address  string
	)
	if info.Config.Validator.Staked != "" {
		if c, err := types.ParseCoin(info.Config.Validator.Staked); err == nil {
			denom = c.Denom
		}
	}

	acc, _ := info.Config.AccountByName(info.Config.Validator.Name)

	questions := []cliquiz.Question{
		cliquiz.NewQuestion("Public address", info.RPCPublicAddress, &address),
		cliquiz.NewQuestion("Account name", acc.Name, &account.Name),
		cliquiz.NewQuestion("Account mnemonic", acc.Mnemonic, &account.Mnemonic),
		cliquiz.NewQuestion("Account coins", strings.Join(acc.Coins, ","), &account.Coins),
		cliquiz.NewQuestion("Staking amount", info.Config.Validator.Staked, &proposal.Validator.StakingAmount),
		cliquiz.NewQuestion("Moniker", "mynode", &proposal.Validator.Moniker),
		cliquiz.NewQuestion("Commission rate", "0.10", &proposal.Validator.CommissionRate),
		cliquiz.NewQuestion("Commission max rate", "0.20", &proposal.Validator.CommissionMaxRate),
		cliquiz.NewQuestion("Commission max change rate", "0.01", &proposal.Validator.CommissionMaxChangeRate),
		cliquiz.NewQuestion("Min self delegation", "1", &proposal.Validator.MinSelfDelegation),
		cliquiz.NewQuestion("Gas prices", "0.025"+denom, &proposal.Validator.GasPrices),
		cliquiz.NewQuestion("Website", "", &proposal.Meta.Website),
		cliquiz.NewQuestion("Identity", "", &proposal.Meta.Identity),
		cliquiz.NewQuestion("Details", "", &proposal.Meta.Details),
	}

	if err := cliquiz.Ask(questions...); err != nil {
		return err
	}
	gentx, mnemonic, err := blockchain.IssueGentx(ctx, account, proposal)
	if err != nil {
		return err
	}

	gentxFormatted, err := json.MarshalIndent(gentx, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("\nGentx: \n\n%s\n\n", gentxFormatted)
	prompt := promptui.Prompt{
		Label:     "Do you confirm the Gentx above",
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		s.Stop()
		fmt.Println("said no")
		return nil
	}

	if err := blockchain.Join(ctx, address, gentx); err != nil {
		return err
	}

	if mnemonic != "" {
		fmt.Printf("\n*** IMPORTANT - Save your mnemonic in a secret place:\n%s\n", mnemonic)
	}

	fmt.Println("\n📜 Proposed validator to join to network")
	return nil
}
