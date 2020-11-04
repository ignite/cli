package starportcmd

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkChainCreate() *cobra.Command {
	c := &cobra.Command{
		Use:  "create [repo]",
		RunE: networkChainCreateHandler,
		Args: cobra.ExactArgs(1),
	}
	return c
}

func networkChainCreateHandler(cmd *cobra.Command, args []string) error {
	var (
		ctx = clictx.From(context.Background())
		ev  = events.NewBus()
		s   = spinner.New(spinner.CharSets[42], 100*time.Millisecond)
	)
	b, err := newNetworkBuilder(networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}
	s.Color("blue")
	defer s.Stop()

	go printEvents(ev, s)

	blockchain, err := b.InitBlockchain(ctx, args[0])
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

	fmt.Printf("\nGenesis: \n\n%s\n\n", string(info.Genesis))
	prompt := promptui.Prompt{
		Label:     "Do you confirm the Genesis above",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		s.Stop()
		fmt.Println("said no")
		return nil
	}
	if err := blockchain.Create(ctx, info.Genesis); err != nil {
		return err
	}

	fmt.Println("\n🌐 Network submited")
	return nil
}
