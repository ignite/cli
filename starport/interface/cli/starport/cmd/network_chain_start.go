package starportcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clictx"
)

func NewNetworkChainStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start [chain-id] [-- <flags>...]",
		Short: "Start network",
		RunE:  networkChainStartHandler,
		Args:  cobra.MinimumNArgs(1),
	}
	return c
}

func networkChainStartHandler(cmd *cobra.Command, args []string) error {
	var startFlags []string
	chainID := args[0]
	if len(args) > 1 { // first arg is always `chain-id`.
		startFlags = args[1:]
	}

	ctx := clictx.From(context.Background())

	err := nb.StartChain(ctx, chainID, startFlags)
	if err == context.Canceled {
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
