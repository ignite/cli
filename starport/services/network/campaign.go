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

// CampaignUpdateTotalShares updates the campaign total shares
func (n Network) CampaignUpdateTotalShares(campaignID uint64, totalShares campaigntypes.Shares) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
		"Updating the campaign %d total shares to %s",
		campaignID,
		totalShares.String(),
	)))

	msg := campaigntypes.NewMsgUpdateTotalShares(n.account.Address(networktypes.SPN), campaignID, totalShares)
	_, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Now the chain %d total shares is %s",
		campaignID,
		totalShares.String(),
	)))
	return nil
}

// CampaignEdit updates the campaign name or metadata
func (n Network) CampaignEdit(
	campaignID uint64,
	name string,
	metadata []byte,
	totalShares campaigntypes.Shares,
	totalSupply sdk.Coins,
) error {
	account := n.account.Address(networktypes.SPN)
	msgs := make([]sdk.Msg, 0)
	if name != "" || len(metadata) > 0 {
		n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
			"Updating the campaign %d name to %s",
			campaignID,
			name,
		)))
		msgs = append(msgs, campaigntypes.NewMsgEditCampaign(
			account,
			name,
			campaignID,
			metadata,
		))
	}
	if !totalShares.Empty() {
		n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
			"Updating the campaign %d total shares to %s",
			campaignID,
			totalShares.String(),
		)))
		msgs = append(msgs, campaigntypes.NewMsgUpdateTotalShares(
			account,
			campaignID,
			totalShares,
		))
	}

	if !totalSupply.Empty() {
		n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
			"Updating the campaign %d total supply to %s",
			campaignID,
			totalSupply.String(),
		)))
		msgs = append(msgs, campaigntypes.NewMsgUpdateTotalSupply(
			account,
			campaignID,
			totalSupply,
		))
	}

	_, err := n.cosmos.BroadcastTx(n.account.Name, msgs...)
	if err != nil {
		return err
	}
	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Campaign %d updated", campaignID,
	)))
	return nil
}
