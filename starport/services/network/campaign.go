package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	res, err := campaigntypes.NewQueryClient(n.cosmos.Context).
		CampaignAll(ctx, &campaigntypes.QueryAllCampaignRequest{})
	if err != nil {
		return campaigns, err
	}

	// Parse fetched campaigns
	for _, campaign := range res.Campaign {
		campaigns = append(campaigns, networktypes.ToCampaign(campaign))
	}

	return campaigns, nil
}
func (n Network) CreateCampaign(name string, totalSupply sdk.Coins) (uint64, error) {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Creating campaign %s", name)))

	msgCreateCampaign := campaigntypes.NewMsgCreateCampaign(
		n.account.Address(networktypes.SPN),
		name,
		totalSupply,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateCampaign)
	if err != nil {
		return 0, err
	}

	var createCampaignRes campaigntypes.MsgCreateCampaignResponse
	if err := res.Decode(&createCampaignRes); err != nil {
		return 0, err
	}

	return createCampaignRes.CampaignID, nil
}

// InitializeMainnet Initialize the mainnet of the campaign.
func (n Network) InitializeMainnet(
	campaignID uint64,
	sourceURL,
	sourceHash string,
	mainnetChainID string,
) (uint64, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Initializing the mainnet campaign"))
	msg := campaigntypes.NewMsgInitializeMainnet(
		n.account.Address(networktypes.SPN),
		campaignID,
		sourceURL,
		sourceHash,
		mainnetChainID,
	)

	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return 0, err
	}

	var initMainnetRes campaigntypes.MsgInitializeMainnetResponse
	if err := res.Decode(&initMainnetRes); err != nil {
		return 0, err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf("Campaign %d was initialized on mainnet", campaignID)))

	return initMainnetRes.MainnetID, nil
}
