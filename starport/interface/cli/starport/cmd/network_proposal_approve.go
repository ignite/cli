package starportcmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewNetworkProposalApprove() *cobra.Command {
	c := &cobra.Command{
		Use:   "approve [chain-id] [number]",
		Short: "Approve a proposal",
		RunE:  networkProposalApproveHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalApproveHandler(cmd *cobra.Command, args []string) error {
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	err = b.ProposalApprove(context.Background(), args[0], int(id))
	if err != nil {
		return err
	}

	fmt.Println("👌 Proposal approved.")
	return nil
}
