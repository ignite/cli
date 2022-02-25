package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/tendermint/starport/starport/pkg/cosmoserror"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesisURL  string
	chainID     string
	campaignID  uint64
	noCheck     bool
	mainnet     bool
	totalSupply sdk.Coins
}

// PublishOption configures chain creation.
type PublishOption func(*publishOptions)

// WithCampaign add a campaign id.
func WithCampaign(id uint64) PublishOption {
	return func(o *publishOptions) {
		o.campaignID = id
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

// WithTotalSupply add a total supply to campaign
func WithTotalSupply(totalSupply sdk.Coins) PublishOption {
	return func(o *publishOptions) {
		o.totalSupply = totalSupply
	}
}

// Mainnet initialize a published chain into the mainnet
func Mainnet() PublishOption {
	return func(o *publishOptions) {
		o.mainnet = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (n Network) Publish(ctx context.Context, c Chain, options ...PublishOption) (launchID, campaignID, mainnetID uint64, err error) {
	o := publishOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var genesisHash string

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if o.noCheck && o.genesisURL != "" {
		if _, genesisHash, err = cosmosutil.GenesisAndHashFromURL(ctx, o.genesisURL); err != nil {
			return 0, 0, 0, err
		}
	}

	chainID := o.chainID
	if chainID == "" {
		chainID, err = c.ID()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	coordinatorAddress := n.account.Address(networktypes.SPN)
	campaignID = o.campaignID

	n.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = profiletypes.
		NewQueryClient(n.cosmos.Context).
		CoordinatorByAddress(ctx, &profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: coordinatorAddress,
		})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrInvalidRequest {
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			coordinatorAddress,
			"",
			"",
			"",
		)
		if _, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateCoordinator); err != nil {
			return 0, 0, 0, err
		}
	} else if err != nil {
		return 0, 0, 0, err
	}

	if campaignID != 0 {
		_, err = campaigntypes.
			NewQueryClient(n.cosmos.Context).
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				CampaignID: o.campaignID,
			})
		if err != nil {
			return 0, 0, 0, err
		}
	} else {
		campaignID, err = n.CreateCampaign(c.Name(), o.totalSupply)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	msgCreateChain := launchtypes.NewMsgCreateChain(
		n.account.Address(networktypes.SPN),
		chainID,
		c.SourceURL(),
		c.SourceHash(),
		o.genesisURL,
		genesisHash,
		true,
		campaignID,
		nil,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateChain)
	if err != nil {
		return 0, 0, 0, err
	}

	var createChainRes launchtypes.MsgCreateChainResponse
	if err := res.Decode(&createChainRes); err != nil {
		return 0, 0, 0, err
	}

	if o.mainnet {
		mainnetID, err = n.InitializeMainnet(campaignID, c.SourceURL(), c.SourceHash(), chainID)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	return createChainRes.LaunchID, campaignID, mainnetID, nil
}
