package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"

	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

type (
	// Prop update campaign proposal
	Prop func(*updateProp)

	// updateProp represents the update campaign proposal
	updateProp struct {
		name        string
		metadata    []byte
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

// WithCampaignTotalSupply provides a total supply proposal to update the campaign.
func WithCampaignTotalSupply(totalSupply sdk.Coins) Prop {
	return func(c *updateProp) {
		c.totalSupply = totalSupply
	}
}

// Campaign fetches the campaign from Starport Network
func (n Network) Campaign(ctx context.Context, campaignID uint64) (networktypes.Campaign, error) {
	n.ev.Send("Fetching campaign information", events.ProgressStarted())
	res, err := n.campaignQuery.Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
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

	n.ev.Send("Fetching campaigns information", events.ProgressStarted())
	res, err := n.campaignQuery.
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

// CreateCampaign creates a campaign in Starport Network
func (n Network) CreateCampaign(name, metadata string, totalSupply sdk.Coins) (uint64, error) {
	n.ev.Send(fmt.Sprintf("Creating campaign %s", name), events.ProgressStarted())

	msgCreateCampaign := campaigntypes.NewMsgCreateCampaign(
		n.account.Address(networktypes.SPN),
		name,
		totalSupply,
		[]byte(metadata),
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
	n.ev.Send("Initializing the mainnet campaign", events.ProgressStarted())
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

	n.ev.Send(fmt.Sprintf("Campaign %d initialized on mainnet", campaignID), events.ProgressFinished())

	return initMainnetRes.MainnetID, nil
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

	n.ev.Send(fmt.Sprintf("Updating the campaign %d", id), events.ProgressStarted())
	account := n.account.Address(networktypes.SPN)
	msgs := make([]sdk.Msg, 0)
	if p.name != "" || len(p.metadata) > 0 {
		msgs = append(msgs, campaigntypes.NewMsgEditCampaign(
			account,
			id,
			p.name,
			p.metadata,
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
	n.ev.Send(fmt.Sprintf("Campaign %d updated", id), events.ProgressFinished())
	return nil
}
