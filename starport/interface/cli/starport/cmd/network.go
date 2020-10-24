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

func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Create and start Blochains collaboratively",
		Args:  cobra.ExactArgs(1),
	}
	c.AddCommand(NewNetworkCreate())
	return c
}

func NewNetworkCreate() *cobra.Command {
	c := &cobra.Command{
		Use:  "create [git-url]",
		RunE: networkCreateHandler,
		Args: cobra.ExactArgs(1),
	}
	return c
}

func networkCreateHandler(cmd *cobra.Command, args []string) error {
	var (
		ctx = clictx.From(context.Background())
		ev  = events.NewBus()
		nb  = networkbuilder.New(networkbuilder.CollectEvents(ev))
		s   = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	)
	defer s.Stop()

	go func() {
		for event := range ev {
			s.Suffix = " " + event.Text()
			if event.IsOngoing() {
				s.Start()
			} else {
				s.Stop()
				fmt.Printf("✔️  %s\n", event.Description)
			}
		}
	}()

	genesis, err := nb.Init(ctx, args[0])
	if err == context.Canceled {
		s.Stop()
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Printf("\nGenesis: \n\n%s\n\n", string(genesis))
	prompt := promptui.Prompt{
		Label:     "Do you confirm the Genesis above",
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		s.Stop()
		fmt.Println("aborted")
		return nil
	}
	if err := nb.Submit(ctx, genesis); err != nil {
		return err
	}

	fmt.Println("\n🌐 Network submited")
	return nil
}
