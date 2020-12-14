package starportcmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

// NewNetworkChainList creates a new chain list command to list
// chains on SPN.
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all chains with proposals summary",
		RunE:  networkChainListHandler,
		Args:  cobra.NoArgs,
	}
	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	// Get the chain summaries
	chainSummaries, err := nb.ChainList(cmd.Context())
	if err != nil {
		return err
	}

	// Rendering
	chainTable := tablewriter.NewWriter(os.Stdout)
	chainTable.SetHeader([]string{"Chain ID", "Source", "Validators (approved)", "Proposals (approved)"})

	for _, chainSummary := range chainSummaries {
		validators := fmt.Sprintf("%d (%d)", chainSummary.TotalValidators, chainSummary.ApprovedValidators)
		proposals := fmt.Sprintf("%d (%d)", chainSummary.TotalProposals, chainSummary.ApprovedProposals)
		chainTable.Append([]string{chainSummary.ChainID, chainSummary.Source, validators, proposals})
	}
	chainTable.Render()

	return nil
}
