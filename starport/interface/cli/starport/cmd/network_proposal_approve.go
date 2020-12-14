package starportcmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/pkg/spn"
)

// NewNetworkProposalApprove creates a new proposal approve command to approve proposals
// for a chain.
func NewNetworkProposalApprove() *cobra.Command {
	c := &cobra.Command{
		Use:     "approve [chain-id] [number<,...>]",
		Aliases: []string{"accept"},
		Short:   "Approve proposals",
		RunE:    networkProposalApproveHandler,
		Args:    cobra.ExactArgs(2),
	}
	return c
}

func networkProposalApproveHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	s.SetText("Calculating gas...")

	var (
		chainID      = args[0]
		proposalList = args[1]
	)

	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	var reviewals []spn.Reviewal

	ids, err := numbers.ParseList(proposalList)
	if err != nil {
		return err
	}
	for _, id := range ids {
		reviewals = append(reviewals, spn.ApproveProposal(id))
	}

	gas, broadcast, err := nb.SubmitReviewals(cmd.Context(), chainID, reviewals...)
	if err != nil {
		return err
	}

	s.Stop()

	// Prompt for confirmation
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("This operation will cost about %v gas. Confirm the transaction",
			gas,
		),
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		return errors.New("transaction aborted")
	}

	s.SetText("Approving...")
	s.Start()

	// Broadcast the transaction
	if err := broadcast(); err != nil {
		return err
	}
	s.Stop()

	fmt.Printf("Proposal(s) %s approved ✅\n", numbers.List(ids, "#"))
	return nil
}
