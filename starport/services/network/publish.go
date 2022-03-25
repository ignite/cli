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
	"github.com/tendermint/starport/starport/pkg/o"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// PublishOptions holds info about how to create a chain.
type PublishOptions struct {
	genesisURL  string
	chainID     string
	campaignID  uint64
	noCheck     bool
	metadata    string
	totalShares campaigntypes.Shares
	totalSupply sdk.Coins
	mainnet     bool
}

// WithCampaign add a campaign id.
func WithCampaign(id uint64) o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.campaignID = id
	}
}

// WithChainID use a custom chain id.
func WithChainID(chainID string) o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.chainID = chainID
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.noCheck = true
	}
}

// WithCustomGenesis enables using a custom genesis during publish.
func WithCustomGenesis(url string) o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.genesisURL = url
	}
}

// WithTotalShares provides a campaign total shares
func WithTotalShares(totalShares campaigntypes.Shares) o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.totalShares = totalShares
	}
}

// WithMetadata provides a meta data proposal to update the campaign.
func WithMetadata(metadata string) o.Option[PublishOptions] {
	return func(c *PublishOptions) {
		c.metadata = metadata
	}
}

// WithTotalSupply provides a total supply proposal to update the campaign.
func WithTotalSupply(totalSupply sdk.Coins) o.Option[PublishOptions] {
	return func(c *PublishOptions) {
		c.totalSupply = totalSupply
	}
}

// Mainnet initialize a published chain into the mainnet
func Mainnet() o.Option[PublishOptions] {
	return func(o *PublishOptions) {
		o.mainnet = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (n Network) Publish(ctx context.Context, c Chain, options ...o.Option[PublishOptions]) (launchID, campaignID, mainnetID uint64, err error) {
	opt := PublishOptions{}

	o.Apply(&opt, options...)

	var (
		genesisHash string
		genesisFile []byte
		genesis     cosmosutil.ChainGenesis
	)

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if opt.genesisURL != "" {
		genesisFile, genesisHash, err = cosmosutil.GenesisAndHashFromURL(ctx, opt.genesisURL)
		if err != nil {
			return 0, 0, 0, err
		}
		genesis, err = cosmosutil.ParseChainGenesis(genesisFile)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	chainID := genesis.ChainID
	// use chain id flag always in the highest priority.
	if opt.chainID != "" {
		chainID = opt.chainID
	}
	// if the chain id is empty, use a default one.
	if chainID == "" {
		chainID, err = c.ChainID()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	coordinatorAddress := n.account.Address(networktypes.SPN)
	campaignID = opt.campaignID

	n.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = profiletypes.
		NewQueryClient(n.cosmos.Context).
		CoordinatorByAddress(ctx, &profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: coordinatorAddress,
		})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
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
				CampaignID: opt.campaignID,
			})
		if err != nil {
			return 0, 0, 0, err
		}
	} else {
		campaignID, err = n.CreateCampaign(c.Name(), opt.metadata, opt.totalSupply)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	msgCreateChain := launchtypes.NewMsgCreateChain(
		n.account.Address(networktypes.SPN),
		chainID,
		c.SourceURL(),
		c.SourceHash(),
		opt.genesisURL,
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

	if err := c.CacheBinary(createChainRes.LaunchID); err != nil {
		return 0, 0, 0, err

	}

	if !opt.totalShares.Empty() {
		if err := n.UpdateCampaign(campaignID, WithCampaignTotalShares(opt.totalShares)); err != nil {
			return 0, 0, 0, err
		}
	}

	if opt.mainnet {
		mainnetID, err = n.InitializeMainnet(campaignID, c.SourceURL(), c.SourceHash(), chainID)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	return createChainRes.LaunchID, campaignID, mainnetID, nil
}
