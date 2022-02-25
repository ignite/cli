package starportcmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/clispinner"
)

// NewNetworkCampaignPublish returns a new command to publish a new campaigns on Starport Network.
func NewNetworkCampaignPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "publish [name] [total-supply]",
		Short: "Publish a campaign",
		Args:  cobra.ExactArgs(2),
		RunE:  networkCampaignPublishHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkCampaignPublishHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	totalSupply, err := sdk.ParseCoinsNormalized(args[1])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	campaignID, err := n.CreateCampaign(args[0], totalSupply)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Campaign ID: %d \n", clispinner.Bullet, campaignID)
	return nil
}
