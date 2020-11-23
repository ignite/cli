package starportcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/manifoldco/promptui"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
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
	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	nb, err := newNetworkBuilder(networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	ctx := clictx.From(context.Background())

	blockchain, err := nb.InitBlockchainFromChainID(ctx, args[0], false)

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

	ip, err := ipify.GetIp()
	if err != nil {
		return err
	}
	publicAddr := fmt.Sprintf("%s:26657", ip)

	questions := []cliquiz.Question{
		cliquiz.NewQuestion("Account name", &account.Name, cliquiz.DefaultAnswer(acc.Name)),
		cliquiz.NewQuestion("Account mnemonic", &account.Mnemonic, cliquiz.DefaultAnswer(acc.Mnemonic)),
		cliquiz.NewQuestion("Public address", &publicAddress, cliquiz.DefaultAnswer(publicAddr)),
		cliquiz.NewQuestion("Account coins", &account.Coins, cliquiz.DefaultAnswer(strings.Join(acc.Coins, ","))),
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

	s.Stop()

	if err := cliquiz.Ask(questions...); err != nil {
		return err
	}
	gentx, address, mnemonic, err := blockchain.IssueGentx(ctx, account, proposal)
	if err != nil {
		return err
	}

	prettyGentx, err := gentx.Pretty()
	if err != nil {
		return err
	}
	fmt.Printf("\nGentx: \n\n%s\n\n", prettyGentx)

	prompt := promptui.Prompt{
		Label:     "Do you confirm the Gentx above",
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		s.Stop()
		fmt.Println("said no")
		return nil
	}

	coins, err := types.ParseCoins(account.Coins)
	if err != nil {
		return err
	}
	if err := blockchain.Join(ctx, address, publicAddress, coins, gentx); err != nil {
		return err
	}

	if mnemonic != "" {
		fmt.Printf("\n*** IMPORTANT - Save your mnemonic in a secret place:\n%s\n", mnemonic)
	}

	fmt.Println("\n📜 Proposed validator to join to network")
	return nil
}
