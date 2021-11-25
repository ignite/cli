package starportcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/entrywriter"
)

var LaunchSummaryHeader = []string{"Launch ID", "Chain ID", "Source", "Campaign ID"}

// LaunchSummary holds summarized information about a chain launch
type LaunchSummary struct {
	LaunchID   string
	ChainID    string
	Source     string
	CampaignID string
}

// NewNetworkChainList returns a new command to list all published chains on Starport Network
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	chains, err := nb.ChainLaunches(cmd.Context())
	if err != nil {
		return err
	}
	sums := launchSummaries(chains)

	s.Stop()
	return renderLaunchSummaries(sums, os.Stdout)
}

// launchSummaries returns the list of launch summaries from the list of chain launches
func launchSummaries(chains []launchtypes.Chain) (sums []LaunchSummary) {
	for _, chain := range chains {
		var campaignID string
		if chain.HasCampaign {
			campaignID = fmt.Sprintf("%d", chain.CampaignID)
		} else {
			campaignID = "no campaign"
		}

		sums = append(sums, LaunchSummary{
			LaunchID:   fmt.Sprintf("%d", chain.LaunchID),
			ChainID:    chain.GenesisChainID,
			Source:     chain.SourceURL,
			CampaignID: campaignID,
		})
	}
	return sums
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(launchSummaries []LaunchSummary, out io.Writer) error {
	var launchEntries [][]string

	for _, sum := range launchSummaries {
		launchEntries = append(launchEntries, []string{sum.LaunchID, sum.ChainID, sum.Source, sum.CampaignID})
	}

	if err := entrywriter.Write(out, LaunchSummaryHeader, launchEntries...); err != nil {
		return errors.Wrap(err, "error printing chain summaries")
	}
	return nil
}
