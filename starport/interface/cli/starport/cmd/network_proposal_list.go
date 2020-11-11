package starportcmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/tendermint/starport/starport/pkg/spn"

	"github.com/spf13/cobra"
)

func NewNetworkProposalList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [chain-id]",
		Short: "List all pending proposals",
		RunE:  networkProposalListHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkProposalListHandler(cmd *cobra.Command, args []string) error {
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	proposals, err := b.ProposalList(context.Background(), args[0], spn.ProposalPending)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 5, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(w, "id\tstatus\tcontent")
	for _, p := range proposals {
		var content string
		switch {
		case p.Account != nil:
			content = fmt.Sprintf("account   | %s, %s", p.Account.Address, p.Account.Coins.String())
		case p.Validator != nil:
			content = fmt.Sprintf("validator | (run 'chain describe' to see gentx), %s", p.Validator.PublicAddress)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", p.ID, p.Status, content)
	}
	w.Flush()

	return nil
}
