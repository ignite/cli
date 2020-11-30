package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/spn"
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
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	if err := nb.SubmitReviewals(cmd.Context(), args[0], spn.RejectProposal(int(id))); err != nil {
		return err
	}

	fmt.Printf("Proposal #%d rejected ⛔️\n", id)
	return nil
}
