package ignitecmd

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
)

const (
	flagCampaignName        = "name"
	flagCampaignMetadata    = "metadata"
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
	c.Flags().String(flagCampaignTotalSupply, "", "Update the total of the mainnet of a campaign")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkCampaignUpdateHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	var (
		campaignName, _        = cmd.Flags().GetString(flagCampaignName)
		metadata, _            = cmd.Flags().GetString(flagCampaignMetadata)
		campaignTotalSupply, _ = cmd.Flags().GetString(flagCampaignTotalSupply)
	)
	totalSupply, err := sdk.ParseCoinsNormalized(campaignTotalSupply)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse campaign ID
	campaignID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	if campaignName == "" && metadata == "" && totalSupply.Empty() {
		return fmt.Errorf("at least one of the flags %s must be provided",
			strings.Join([]string{
				flagCampaignName,
				flagCampaignMetadata,
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
	session.Println()

	info, err := yaml.Marshal(cmd.Context(), campaign)
	if err != nil {
		return err
	}

	session.StopSpinner()

	return session.Print(info)
}
