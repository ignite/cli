package starportcmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/tendermint/starport/starport/pkg/spn"

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
	status, err := cmd.Flags().GetString(statusFlag)
	if err != nil {
		return err
	}

	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	proposals, err := b.ProposalList(context.Background(), args[0], spn.ProposalStatus(status))
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 5, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(w, "id\tstatus\tcontent")
	for _, p := range proposals {
		var content string
		switch {
		case p.Account != nil:
			content = fmt.Sprintf("Add Account   | %s, %s", p.Account.Address, p.Account.Coins.String())
		case p.Validator != nil:
			content = fmt.Sprintf("Add Validator | (run 'chain describe' to see gentx), %s", p.Validator.PublicAddress)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", p.ID, strings.Title(string(p.Status)), content)
	}
	w.Flush()

	return nil
}
