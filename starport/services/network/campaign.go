package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

type (
	// Prop update campaign proposal
	Prop func(*updateProp)

	// updateProp represents the update campaign proposal
	updateProp struct {
		name        string
		metadata    []byte
		totalShares campaigntypes.Shares
		totalSupply sdk.Coins
	}
)

// WithCampaignName provides a name proposal to update the campaign.
func WithCampaignName(name string) Prop {
	return func(c *updateProp) {
		c.name = name
	}
}

// WithCampaignMetadata provides a meta data proposal to update the campaign.
func WithCampaignMetadata(metadata string) Prop {
	return func(c *updateProp) {
		c.metadata = []byte(metadata)
	}
}

// WithCampaignTotalShares provides a total shares proposal to update the campaign.
func WithCampaignTotalShares(totalShares campaigntypes.Shares) Prop {
	return func(c *updateProp) {
		c.totalShares = totalShares
	}
}

// WithCampaignTotalSupply provides a total supply proposal to update the campaign.
func WithCampaignTotalSupply(totalSupply sdk.Coins) Prop {
	return func(c *updateProp) {
		c.totalSupply = totalSupply
	}
}

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
	if _, err := n.cosmos.BroadcastTx(n.account.Name, msg); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Now the chain %d total shares is %s",
		campaignID,
		totalShares.String(),
	)))
	return nil
}

// UpdateCampaign updates the campaign name or metadata
func (n Network) UpdateCampaign(
	id uint64,
	props ...Prop,
) error {
	// Apply the options provided by the user
	p := updateProp{}
	for _, apply := range props {
		apply(&p)
	}

	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Updating the campaign %d", id)))
	account := n.account.Address(networktypes.SPN)
	msgs := make([]sdk.Msg, 0)
	if p.name != "" || len(p.metadata) > 0 {
		msgs = append(msgs, campaigntypes.NewMsgEditCampaign(
			account,
			p.name,
			id,
			p.metadata,
		))
	}
	if !p.totalShares.Empty() {
		msgs = append(msgs, campaigntypes.NewMsgUpdateTotalShares(
			account,
			id,
			p.totalShares,
		))
	}
	if !p.totalSupply.Empty() {
		msgs = append(msgs, campaigntypes.NewMsgUpdateTotalSupply(
			account,
			id,
			p.totalSupply,
		))
	}

	if _, err := n.cosmos.BroadcastTx(n.account.Name, msgs...); err != nil {
		return err
	}
	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Campaign %d updated", id,
	)))
	return nil
}
