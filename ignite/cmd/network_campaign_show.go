package ignitecmd

import (
	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/yaml"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/spf13/cobra"
)

// NewNetworkCampaignShow returns a new command to show published campaign on Ignite
func NewNetworkCampaignShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [campaign-id]",
		Short: "Show published campaign",
		Args:  cobra.ExactArgs(1),
		RunE:  networkCampaignShowHandler,
	}
	return c
}

func networkCampaignShowHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	// parse campaign ID
	campaignID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	campaign, err := n.Campaign(cmd.Context(), campaignID)
	if err != nil {
		return err
	}

	info, err := yaml.Marshal(cmd.Context(), campaign)
	if err != nil {
		return err
	}

	return session.Println(info)
}
