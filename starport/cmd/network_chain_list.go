package starportcmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/ctxreader"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/networkbuilder"
	"golang.org/x/sync/errgroup"
)

const (
	chainsPerPageCount = 40
)

const flagSearch = "search"

// NewNetworkChainList creates a new chain list command to list
// chains on SPN.
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List all chains with proposals summary",
		RunE:  networkChainListHandler,
		Args:  cobra.NoArgs,
	}
	c.Flags().String(flagSearch, "", "List chains with the specified prefix in chain id")
	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	nb, err := newNetworkBuilder(cmd.Context())
	if err != nil {
		return err
	}

	// Parse search flag
	prefix, err := cmd.Flags().GetString(flagSearch)
	if err != nil {
		return err
	}

	var pageKey []byte

	for {
		s.SetText("Querying chains...")
		s.Start()

		chainSummaries, nextPageKey, err := listChainSummaries(cmd.Context(), nb, prefix, pageKey)
		if err != nil {
			return err
		}

		s.Stop()
		renderChainSummaries(chainSummaries)

		// check if there is a next page, if so ask to load more result.
		if nextPageKey != nil {
			pageKey = nextPageKey
		} else {
			return nil
		}

		fmt.Printf("\nPress <Enter> to show more blockchains.\n")
		buf := bufio.NewReader(ctxreader.New(cmd.Context(), os.Stdin))
		if _, err := buf.ReadBytes('\n'); err != nil {
			return err
		}
	}
}

// ChainSummary keys summarized chain info.
type ChainSummary struct {
	ChainID            string
	Source             string
	TotalValidators    int
	ApprovedValidators int
	TotalProposals     int
	ApprovedProposals  int
}

// renderChainSummaries renders chain summaries to std output.
func renderChainSummaries(chainSummaries []ChainSummary) {
	// Rendering
	chainTable := tablewriter.NewWriter(os.Stdout)
	chainTable.SetHeader([]string{"Chain ID", "Source", "Validators (approved)", "Proposals (approved)"})

	for _, chainSummary := range chainSummaries {
		validators := fmt.Sprintf("%d (%d)", chainSummary.TotalValidators, chainSummary.ApprovedValidators)
		proposals := fmt.Sprintf("%d (%d)", chainSummary.TotalProposals, chainSummary.ApprovedProposals)
		chainTable.Append([]string{chainSummary.ChainID, chainSummary.Source, validators, proposals})
	}
	chainTable.Render()
}

// listChainSummaries lists chains with their summary info by using nextPageKey as the
// pagination key to fetch the next page.
func listChainSummaries(ctx context.Context, nb *networkbuilder.Builder, prefix string, pageKey []byte) (summaries []ChainSummary,
	nextPageKey []byte, err error) {
	var chains []spn.Chain
	chains, nextPageKey, err = nb.ChainList(ctx, spn.PrefixChainListing(prefix), spn.PaginateChainListing(pageKey, chainsPerPageCount))
	if err != nil {
		return nil, nil, err
	}

	summaries = make([]ChainSummary, len(chains))

	// fetchAndSetSummary creates a summary for chain and sets it to i index in summaries.
	fetchAndSetSummary := func(i int, chain spn.Chain) error {
		proposals, err := nb.ProposalList(ctx, chain.ChainID)
		if err != nil {
			return err
		}

		summary := ChainSummary{
			ChainID:        chain.ChainID,
			Source:         chain.URL,
			TotalProposals: len(proposals),
		}

		for _, proposal := range proposals {
			if proposal.Status == spn.ProposalStatusApproved {
				summary.ApprovedProposals++
			}
			if proposal.Validator != nil {
				summary.TotalValidators++
				if proposal.Status == spn.ProposalStatusApproved {
					summary.ApprovedValidators++
				}
			}
		}

		summaries[i] = summary
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	for i, chain := range chains {
		i, chain := i, chain

		g.Go(func() error { return fetchAndSetSummary(i, chain) })
	}

	return summaries, nextPageKey, g.Wait()
}
