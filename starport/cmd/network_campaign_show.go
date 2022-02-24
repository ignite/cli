package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/yaml"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkCampaignShow returns a new command to show published campaign on Starport Network
func NewNetworkCampaignShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [campaign-id]",
		Short: "Show published campaign",
		Args:  cobra.ExactArgs(1),
		RunE:  networkCampaignShowHandler,
	}
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkCampaignShowHandler(cmd *cobra.Command, args []string) error {
	// parse campaign ID
	campaignID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

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

	nb.Spinner.Stop()
	fmt.Print(info)
	return nil
}
