package network

import (
	"context"

	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesisURL   string
	chainID      string
	campaignName string
	campaignID   uint64
	noCheck      bool
}

// PublishOption configures chain creation.
type PublishOption func(*publishOptions)

// WithCampaign add a campaign id.
func WithCampaign(id uint64) PublishOption {
	return func(o *publishOptions) {
		o.campaignID = id
	}
}

// WithCampaignName use a custom campaign name.
func WithCampaignName(campaignName string) PublishOption {
	return func(o *publishOptions) {
		o.campaignName = campaignName
	}
}

// WithChainID use a custom chain id.
func WithChainID(chainID string) PublishOption {
	return func(o *publishOptions) {
		o.chainID = chainID
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() PublishOption {
	return func(o *publishOptions) {
		o.noCheck = true
	}
}

// WithCustomGenesis enables using a custom genesis during publish.
func WithCustomGenesis(url string) PublishOption {
	return func(o *publishOptions) {
		o.genesisURL = url
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (n Network) Publish(ctx context.Context, c Chain, options ...PublishOption) (launchID, campaignID uint64, err error) {
	o := publishOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var genesisHash string

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if o.noCheck && o.genesisURL != "" {
		if _, genesisHash, err = cosmosutil.GenesisAndHashFromURL(ctx, o.genesisURL); err != nil {
			return 0, 0, err
		}
	}

	chainID := o.chainID
	if chainID == "" {
		chainID, err = c.ID()
		if err != nil {
			return 0, 0, err
		}
	}

	coordinatorAddress := n.account.Address(networkchain.SPN)
	campaignID = o.campaignID

	n.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = profiletypes.
		NewQueryClient(n.cosmos.Context).
		CoordinatorByAddress(ctx, &profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: coordinatorAddress,
		})

	// TODO check for not found and only then create a new coordinator, otherwise return the err.
	if err != nil {
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			coordinatorAddress,
			"",
			"",
			"",
		)
		if _, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateCoordinator); err != nil {
			return 0, 0, err
		}
	}

	if campaignID != 0 {
		_, err = campaigntypes.
			NewQueryClient(n.cosmos.Context).
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				Id: o.campaignID,
			})
		if err != nil {
			return 0, 0, err
		}
	} else {
		campaignName := o.campaignName
		if campaignName == "" {
			campaignName = c.Name()
		}
		msgCreateCampaign := campaigntypes.NewMsgCreateCampaign(
			coordinatorAddress,
			campaignName,
			nil,
			false,
		)
		res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateCampaign)
		if err != nil {
			return 0, 0, err
		}

		var createCampaignRes campaigntypes.MsgCreateCampaignResponse
		if err := res.Decode(&createCampaignRes); err != nil {
			return 0, 0, err
		}
		campaignID = createCampaignRes.CampaignID
	}

	msgCreateChain := launchtypes.NewMsgCreateChain(
		n.account.Address(networkchain.SPN),
		chainID,
		c.SourceURL(),
		c.SourceHash(),
		o.genesisURL,
		genesisHash,
		true,
		campaignID,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateChain)
	if err != nil {
		return 0, 0, err
	}

	var createChainRes launchtypes.MsgCreateChainResponse
	if err := res.Decode(&createChainRes); err != nil {
		return 0, 0, err
	}

	return createChainRes.LaunchID, campaignID, nil
}
