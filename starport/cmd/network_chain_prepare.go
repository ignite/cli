package starportcmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkChainPrepare returns a new command to prepare the chain for launch
func NewNetworkChainPrepare() *cobra.Command {
	c := &cobra.Command{
		Use:   "prepare [launch-id]",
		Short: "Prepare the chain for launch",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainPrepareHandler,
	}

	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainPrepareHandler(cmd *cobra.Command, args []string) error {
	// TODO: Add routine from init
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

	_, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
	}

	// TODO: create and initialize the chain

	// TODO: call prepare

	return nil
}
