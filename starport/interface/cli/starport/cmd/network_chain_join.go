package starportcmd

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kr/pretty"
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
		nb  = networkbuilder.New(networkbuilder.CollectEvents(ev))
		s   = spinner.New(spinner.CharSets[42], 100*time.Millisecond)
	)
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

	//info, err := blockchain.Info()
	//if err != nil {
	//return err
	//}

	var proposal networkbuilder.Proposal

	questions := []cliquiz.Question{
		cliquiz.NewQuestion("Moniker", "mynode", &proposal.Moniker),
		cliquiz.NewQuestion("Staking amount", 10, &proposal.StakingAmount),
	}

	if err := cliquiz.Ask(questions...); err != nil {
		return err
	}
	pretty.Println(proposal)

	fmt.Println("\n📜 Proposed validator to join to network")
	return nil
}
