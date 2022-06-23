package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

const (
	flagMetadata = "metadata"
)

// NewNetworkCampaignPublish returns a new command to publish a new campaigns on Ignite.
func NewNetworkCampaignPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name] [total-supply]",
		Short: "Create a campaign",
		Args:  cobra.ExactArgs(2),
		RunE:  networkCampaignPublishHandler,
	}
	c.Flags().String(flagMetadata, "", "Add a metada to the chain")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkCampaignPublishHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	totalSupply, err := sdk.ParseCoinsNormalized(args[1])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	metadata, _ := cmd.Flags().GetString(flagMetadata)
	campaignID, err := n.CreateCampaign(args[0], metadata, totalSupply)
	if err != nil {
		return err
	}

	session.StopSpinner()

	return session.Printf("%s Campaign ID: %d \n", icons.Bullet, campaignID)
}
