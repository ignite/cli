package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
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

	session.StopSpinner()

	return session.Println(info)
}
