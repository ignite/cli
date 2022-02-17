package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkCampaign creates a new campaign command that holds some other
// sub commands related to launching a network for a campaign.
func NewNetworkCampaign() *cobra.Command {
	c := &cobra.Command{
		Use:   "campaign",
		Short: "Handle campaigns",
	}

	c.AddCommand(
		NewNetworkCampaignList(),
		NewNetworkCampaignShow(),
	)

	return c
}
