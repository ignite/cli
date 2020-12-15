package starportcmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

const searchFlag = "search"

// NewNetworkChainList creates a new chain list command to list
// chains on SPN.
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all chains with proposals summary",
		RunE:  networkChainListHandler,
		Args:  cobra.NoArgs,
	}
	c.Flags().String(searchFlag, "", "List chains with the specified prefix in chain id")
	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	// Parse search flag
	prefix, err := cmd.Flags().GetString(searchFlag)
	if err != nil {
		return err
	}

	// Get the chain summaries
	chainSummaries, err := nb.ChainList(cmd.Context(), prefix)
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
