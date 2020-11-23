package starportcmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewNetworkProposalReject() *cobra.Command {
	c := &cobra.Command{
		Use:   "reject [chain-id] [number]",
		Short: "Reject a proposal",
		RunE:  networkProposalRejectHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalRejectHandler(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	err = nb.ProposalReject(context.Background(), args[0], int(id))
	if err != nil {
		return err
	}

	fmt.Printf("Proposal #%d rejected ⛔️\n", id)
	return nil
}
