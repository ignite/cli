package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/spn"
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
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	if err := nb.SubmitReviewals(cmd.Context(), args[0], spn.ApproveProposal(int(id))); err != nil {
		return err
	}

	fmt.Printf("Proposal #%d approved ✅\n", id)
	return nil
}
