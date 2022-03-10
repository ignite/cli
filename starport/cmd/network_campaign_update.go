package starportcmd

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"

	"github.com/tendermint/starport/starport/pkg/yaml"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagCampaignName        = "name"
	flagCampaignMetadata    = "metadata"
	flagCampaignTotalShares = "total-shares"
	flagCampaignTotalSupply = "total-supply"
)

func NewNetworkCampaignUpdate() *cobra.Command {
	c := &cobra.Command{
		Use:   "update [campaign-id]",
		Short: "Update details fo the campaign of the campaign",
		Args:  cobra.ExactArgs(1),
		RunE:  networkCampaignUpdateHandler,
	}
	c.Flags().String(flagCampaignName, "", "Update the campaign name")
	c.Flags().String(flagCampaignMetadata, "", "Update the campaign metadata")
	c.Flags().String(flagCampaignTotalShares, "", "Update the shares supply for the campaign")
	c.Flags().String(flagCampaignTotalSupply, "", "Update the total of the mainnet of a campaign")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkCampaignUpdateHandler(cmd *cobra.Command, args []string) error {
	var (
		campaignName, _        = cmd.Flags().GetString(flagCampaignName)
		metadata, _            = cmd.Flags().GetString(flagCampaignMetadata)
		campaignTotalShares, _ = cmd.Flags().GetString(flagCampaignTotalShares)
		campaignTotalSupply, _ = cmd.Flags().GetString(flagCampaignTotalSupply)
	)
	totalShares, err := campaigntypes.NewShares(campaignTotalShares)
	if err != nil {
		return err
	}
	totalSupply, err := sdk.ParseCoinsNormalized(campaignTotalSupply)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse campaign ID
	campaignID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	if campaignName == "" && metadata == "" &&
		totalShares.Empty() && totalSupply.Empty() {
		return fmt.Errorf("at least one of the flags %s must be provided",
			strings.Join([]string{
				flagCampaignName,
				flagCampaignMetadata,
				flagCampaignTotalShares,
				flagCampaignTotalSupply,
			}, ", "),
		)
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	var proposals []network.Prop

	if campaignName != "" {
		proposals = append(proposals, network.WithCampaignName(campaignName))
	}
	if metadata != "" {
		proposals = append(proposals, network.WithCampaignMetadata(metadata))
	}
	if !totalShares.Empty() {
		proposals = append(proposals, network.WithCampaignTotalShares(totalShares))
	}
	if !totalSupply.Empty() {
		proposals = append(proposals, network.WithCampaignTotalSupply(totalSupply))
	}

	if err = n.UpdateCampaign(campaignID, proposals...); err != nil {
		return err
	}

	campaign, err := n.Campaign(cmd.Context(), campaignID)
	if err != nil {
		return err
	}
	fmt.Println("")

	info, err := yaml.Marshal(cmd.Context(), campaign)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Print(info)
	return nil
}
