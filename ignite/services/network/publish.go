package network

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// SharePercentage represent percent of total share
type SharePercentage struct {
	denom   string
	percent float64
}

// NewSharePercentage creates new share percentage
func NewSharePercentage(denom string, percent float64) SharePercentage {
	return SharePercentage{denom: denom, percent: percent}
}

// ParseSharePercentages parsers SharePercentage list from string
// format: 12.4%foo,10%bar,0.133%baz
func ParseSharePercentages(percentagesString string) ([]SharePercentage, error) {
	var rePercentageRequired = regexp.MustCompile(`^[0-9]+.[0-9]*%`)
	rawPercentages := strings.Split(percentagesString, ",")
	percentages := make([]SharePercentage, len(rawPercentages))
	for i, percentage := range rawPercentages {
		// validate raw percentage format
		if len(rePercentageRequired.FindStringIndex(percentage)) == 0 {
			return nil, fmt.Errorf("invalid percentage format %s", percentage)
		}

		foo := strings.Split(percentage, "%")
		denom := foo[1]
		percent, err := strconv.ParseFloat(foo[0], 64)
		if err != nil {
			return nil, err
		}
		if percent > 100 {
			return nil, fmt.Errorf("%q can not be bigger than 100", denom)
		}
		percentages[i] = NewSharePercentage(denom, percent)
	}

	return percentages, nil
}

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesisURL       string
	chainID          string
	campaignID       uint64
	noCheck          bool
	metadata         string
	totalSupply      sdk.Coins
	sharePercentages []SharePercentage
	mainnet          bool
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
func WithPercentageShares(sharePercentages []SharePercentage) PublishOption {
	return func(c *publishOptions) {
		c.sharePercentages = sharePercentages
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

	var (
		genesisHash string
		genesisFile []byte
		genesis     cosmosutil.ChainGenesis
	)

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash.
	if o.genesisURL != "" {
		genesisFile, genesisHash, err = cosmosutil.GenesisAndHashFromURL(ctx, o.genesisURL)
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
	if o.chainID != "" {
		chainID = o.chainID
	}
	// if the chain id is empty, use a default one.
	if chainID == "" {
		chainID, err = c.ChainID()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	coordinatorAddress := n.account.Address(networktypes.SPN)
	campaignID = o.campaignID

	n.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = n.profileQuery.
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
		_, err = n.campaignQuery.
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				CampaignID: o.campaignID,
			})
		if err != nil {
			return 0, 0, 0, err
		}
	} else {
		campaignID, err = n.CreateCampaign(c.Name(), o.metadata, o.totalSupply)
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

	msgs := []sdk.Msg{msgCreateChain}

	if len(o.sharePercentages) != 0 {
		totalSharesResp, err := n.campaignQuery.TotalShares(ctx, &campaigntypes.QueryTotalSharesRequest{})
		if err != nil {
			return 0, 0, 0, err
		}

		var coins []sdk.Coin

		for _, share := range o.sharePercentages {
			amount := int64(share.percent * float64(totalSharesResp.TotalShares/100))
			coins = append(coins, sdk.NewInt64Coin(share.denom, amount))
		}
		// TODO consider moving to UpdateCampaign, but not sure, may not be relevant.
		// It is better to send multiple message in a single tx too.
		// consider ways to refactor to accomplish a better API and efficiency.
		msgMintVouchers := campaigntypes.NewMsgMintVouchers(
			n.account.Address(networktypes.SPN),
			campaignID,
			campaigntypes.NewSharesFromCoins(sdk.NewCoins(coins...)),
		)
		msgs = append(msgs, msgMintVouchers)
	}

	res, err := n.cosmos.BroadcastTx(n.account.Name, msgs...)
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

	if o.mainnet {
		mainnetID, err = n.InitializeMainnet(campaignID, c.SourceURL(), c.SourceHash(), chainID)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	return createChainRes.LaunchID, campaignID, mainnetID, nil
}
