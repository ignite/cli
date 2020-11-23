package starportcmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
)

func NewNetworkChainCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [repo]",
		Short: "Create a new network",
		RunE:  networkChainCreateHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkChainCreateHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	ctx := clictx.From(context.Background())

	blockchain, err := nb.InitBlockchainFromPath(ctx, args[0])
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

	prettyGenesis, err := info.Genesis.Pretty()
	if err != nil {
		return err
	}
	fmt.Printf("\nGenesis: \n\n%s\n\n", prettyGenesis)

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
		if strings.Contains(err.Error(), "not found") {
			return errors.New("make sure that your SPN account has enough balance")
		}
		return err
	}

	fmt.Println("\n🌐 Network submited")
	return nil
}
