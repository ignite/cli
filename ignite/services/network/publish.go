package network

import (
	"context"
	"os"
	"path/filepath"

	cosmosgenesis "github.com/ignite/cli/ignite/pkg/cosmosutil/genesis"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesisURL       string
	genesisConfig    string
	chainID          string
	campaignID       uint64
	metadata         string
	totalSupply      sdk.Coins
	sharePercentages SharePercents
	mainnet          bool
	accountBalance   sdk.Coins
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

// WithCustomGenesisURL enables using a custom genesis during publish.
func WithCustomGenesisURL(url string) PublishOption {
	return func(o *publishOptions) {
		o.genesisURL = url
	}
}

// WithCustomGenesisConfig enables using a custom genesis during publish.
func WithCustomGenesisConfig(configFile string) PublishOption {
	return func(o *publishOptions) {
		o.genesisConfig = configFile
	}
}

// WithMetadata provides a meta data proposal to update the campaign.
func WithMetadata(metadata string) PublishOption {
	return func(c *publishOptions) {
		c.metadata = metadata
	}
}

// WithTotalSupply provides a total supply proposal to update the campaign.
func WithTotalSupply(totalSupply sdk.Coins) PublishOption {
	return func(c *publishOptions) {
		c.totalSupply = totalSupply
	}
}

// WithPercentageShares enables minting vouchers for shares.
func WithPercentageShares(sharePercentages []SharePercent) PublishOption {
	return func(c *publishOptions) {
		c.sharePercentages = sharePercentages
	}
}

// WithAccountBalance set a balance used for all genesis account of the chain
func WithAccountBalance(accountBalance sdk.Coins) PublishOption {
	return func(c *publishOptions) {
		c.accountBalance = accountBalance
	}
}

// Mainnet initialize a published chain into the mainnet
func Mainnet() PublishOption {
	return func(o *publishOptions) {
		o.mainnet = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (n Network) Publish(ctx context.Context, c Chain, options ...PublishOption) (launchID, campaignID uint64, err error) {
	o := publishOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var (
		genesisHash string
		genesis     *cosmosgenesis.Genesis
		chainID     string
	)

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if o.genesisURL != "" {
		genesis, err = cosmosgenesis.FromURL(ctx, o.genesisURL, filepath.Join(os.TempDir(), "genesis.json"))
		if err != nil {
			return 0, 0, err
		}
		genesisHash, err = genesis.Hash()
		if err != nil {
			return 0, 0, err
		}
		chainID, err = genesis.ChainID()
		if err != nil {
			return 0, 0, err
		}
	}

	// use chain id flag always in the highest priority.
	if o.chainID != "" {
		chainID = o.chainID
	}
	// if the chain id is empty, use a default one.
	if chainID == "" {
		chainID, err = c.ChainID()
		if err != nil {
			return 0, 0, err
		}
	}

	coordinatorAddress, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return 0, 0, err
	}
	campaignID = o.campaignID

	n.ev.Send("Publishing the network", events.ProgressStart())

	// a coordinator profile is necessary to publish a chain
	// if the user doesn't have an associated coordinator profile, we create one
	if _, err := n.CoordinatorIDByAddress(ctx, coordinatorAddress); err == ErrObjectNotFound {
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			coordinatorAddress,
			"",
			"",
			"",
		)
		if _, err := n.cosmos.BroadcastTx(ctx, n.account, msgCreateCoordinator); err != nil {
			return 0, 0, err
		}
	} else if err != nil {
		return 0, 0, err
	}

	// check if a campaign associated to the chain is provided
	if campaignID != 0 {
		_, err = n.campaignQuery.
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				CampaignID: o.campaignID,
			})
		if err != nil {
			return 0, 0, err
		}
	} else if o.mainnet {
		// a mainnet is always associated to a campaign
		// if no campaign is provided, we create one, and we directly initialize the mainnet
		campaignID, err = n.CreateCampaign(ctx, c.Name(), o.metadata, o.totalSupply)
		if err != nil {
			return 0, 0, err
		}
	}

	// mint vouchers
	if campaignID != 0 && !o.sharePercentages.Empty() {
		totalSharesResp, err := n.campaignQuery.TotalShares(ctx, &campaigntypes.QueryTotalSharesRequest{})
		if err != nil {
			return 0, 0, err
		}

		var coins []sdk.Coin
		for _, percentage := range o.sharePercentages {
			coin, err := percentage.Share(totalSharesResp.TotalShares)
			if err != nil {
				return 0, 0, err
			}
			coins = append(coins, coin)
		}
		// TODO consider moving to UpdateCampaign, but not sure, may not be relevant.
		// It is better to send multiple message in a single tx too.
		// consider ways to refactor to accomplish a better API and efficiency.

		addr, err := n.account.Address(networktypes.SPN)
		if err != nil {
			return 0, 0, err
		}

		msgMintVouchers := campaigntypes.NewMsgMintVouchers(
			addr,
			campaignID,
			campaigntypes.NewSharesFromCoins(sdk.NewCoins(coins...)),
		)
		_, err = n.cosmos.BroadcastTx(ctx, n.account, msgMintVouchers)
		if err != nil {
			return 0, 0, err
		}
	}

	// depending on mainnet flag initialize mainnet or testnet
	if o.mainnet {
		launchID, err = n.InitializeMainnet(ctx, campaignID, c.SourceURL(), c.SourceHash(), chainID)
		if err != nil {
			return 0, 0, err
		}
	} else {
		addr, err := n.account.Address(networktypes.SPN)
		if err != nil {
			return 0, 0, err
		}

		// get initial genesis
		initialGenesis := launchtypes.NewDefaultInitialGenesis()
		switch {
		case o.genesisURL != "":
			initialGenesis = launchtypes.NewGenesisURL(
				o.genesisURL,
				genesisHash,
			)
		case o.genesisConfig != "":
			initialGenesis = launchtypes.NewConfigGenesis(
				o.genesisConfig,
			)
		}

		msgCreateChain := launchtypes.NewMsgCreateChain(
			addr,
			chainID,
			c.SourceURL(),
			c.SourceHash(),
			initialGenesis,
			campaignID != 0,
			campaignID,
			o.accountBalance,
			nil,
		)
		res, err := n.cosmos.BroadcastTx(ctx, n.account, msgCreateChain)
		if err != nil {
			return 0, 0, err
		}
		var createChainRes launchtypes.MsgCreateChainResponse
		if err := res.Decode(&createChainRes); err != nil {
			return 0, 0, err
		}
		launchID = createChainRes.LaunchID
	}
	if err := c.CacheBinary(launchID); err != nil {
		return 0, 0, err
	}

	return launchID, campaignID, nil
}
