package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkCampaign creates a new campaign command that holds other
// subcommands related to launching a network for a campaign.
func NewNetworkCampaign() *cobra.Command {
	c := &cobra.Command{
		Use:   "campaign",
		Short: "Handle campaigns",
	}
	c.AddCommand(
		NewNetworkCampaignPublish(),
		NewNetworkCampaignList(),
		NewNetworkCampaignShow(),
		NewNetworkCampaignUpdate(),
		NewNetworkCampaignAccount(),
	)
	return c
}
