package starportcmd

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tendermint/starport/starport/pkg/spn"
	"os"

	"github.com/spf13/cobra"
)

const statusFlag = "status"

func NewNetworkProposalList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [chain-id]",
		Short: "List all pending proposals",
		RunE:  networkProposalListHandler,
		Args:  cobra.ExactArgs(1),
	}
	c.Flags().String(statusFlag, "pending", "Filter proposals by status")
	return c
}

func networkProposalListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	status, err := cmd.Flags().GetString(statusFlag)
	if err != nil {
		return err
	}

	proposals, err := nb.ProposalList(context.Background(), args[0], spn.ProposalStatus(status))
	if err != nil {
		return err
	}

	// Rendering
	proposalTable := tablewriter.NewWriter(os.Stdout)
	proposalTable.SetHeader([]string{"ID", "Status", "Type", "Content"})

	for _, proposal := range proposals {
		id := fmt.Sprintf("%d", proposal.ID)
		proposalType := "Unknown"
		content := ""

		switch {
		case proposal.Account != nil:
			proposalType = "Add Account"
			content = fmt.Sprintf("%s, %s", proposal.Account.Address, proposal.Account.Coins.String())
		case proposal.Validator != nil:
			proposalType = "Add Validator"
			content = fmt.Sprintf("(run 'chain show' to see gentx), %s", proposal.Validator.P2PAddress)
		}

		proposalTable.Append([]string{
			id,
			status,
			proposalType,
			content,
		})
	}
	proposalTable.Render()

	return nil
}
