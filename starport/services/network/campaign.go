package network

import (
	"context"

	campaigntypes "github.com/tendermint/spn/x/campaign/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// Campaign fetches the campaign from Starport Network
func (n Network) Campaign(ctx context.Context, campaignID uint64) (networktypes.Campaign, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching campaign information"))
	res, err := campaigntypes.NewQueryClient(n.cosmos.Context).Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: campaignID,
	})
	if err != nil {
		return networktypes.Campaign{}, err
	}
	return networktypes.ToCampaign(res.Campaign), nil
}

// Campaigns fetches the campaigns from Starport Network
func (n Network) Campaigns(ctx context.Context) ([]networktypes.Campaign, error) {
	var campaigns []networktypes.Campaign

	n.ev.Send(events.New(events.StatusOngoing, "Fetching campaigns information"))
	res, err := campaigntypes.NewQueryClient(n.cosmos.Context).CampaignAll(ctx, &campaigntypes.QueryAllCampaignRequest{})
	if err != nil {
		return campaigns, err
	}

	// Parse fetched campaigns
	for _, campaign := range res.Campaign {
		campaigns = append(campaigns, networktypes.ToCampaign(campaign))
	}

	return campaigns, nil
}
