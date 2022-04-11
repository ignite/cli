package ignitecmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/entrywriter"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

var CampaignSummaryHeader = []string{
	"id",
	"name",
	"coordinator id",
	"mainnet id",
}

// NewNetworkCampaignList returns a new command to list all published campaigns on Ignite.
func NewNetworkCampaignList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published campaigns",
		Args:  cobra.NoArgs,
		RunE:  networkCampaignListHandler,
	}
	return c
}

func networkCampaignListHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}
	campaigns, err := n.Campaigns(cmd.Context())
	if err != nil {
		return err
	}

	nb.Cleanup()
	return renderCampaignSummaries(campaigns, os.Stdout)
}

// renderCampaignSummaries writes into the provided out, the list of summarized campaigns
func renderCampaignSummaries(campaigns []networktypes.Campaign, out io.Writer) error {
	var campaignEntries [][]string

	for _, c := range campaigns {
		mainnetID := entrywriter.None
		if c.MainnetInitialized {
			mainnetID = fmt.Sprintf("%d", c.MainnetID)
		}

		campaignEntries = append(campaignEntries, []string{
			fmt.Sprintf("%d", c.ID),
			c.Name,
			fmt.Sprintf("%d", c.CoordinatorID),
			mainnetID,
		})
	}

	return entrywriter.MustWrite(out, CampaignSummaryHeader, campaignEntries...)
}
