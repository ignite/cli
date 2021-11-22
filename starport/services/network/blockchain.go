package network

import (
	"context"
	"os"

	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/services/chain"
)

// Blockchain represents a blockchain.
type Blockchain struct {
	launchID      uint64
	appPath       string
	url           string
	hash          string
	chain         *chain.Chain
	isInitialized bool
	builder       *Builder
	genesisURL    string
	genesisHash   string
}

// setup setups the blockchain.
func (b *Blockchain) setup(
	chainID,
	home string,
	keyringBackend chaincmd.KeyringBackend,
) error {
	b.builder.ev.Send(events.New(events.StatusOngoing, "Setting up the blockchain"))

	chainOption := []chain.Option{
		chain.LogLevel(chain.LogSilent),
		chain.ID(chainID),
	}

	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	if gitpod.IsOnGitpod() {
		keyringBackend = chaincmd.KeyringBackendTest
	}

	chainOption = append(chainOption, chain.KeyringBackend(keyringBackend))

	chain, err := chain.New(b.appPath, chainOption...)
	if err != nil {
		return err
	}

	if !chain.Version.IsFamily(cosmosver.Stargate) {
		return sperrors.ErrOnlyStargateSupported
	}

	b.chain = chain
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain set up"))

	return nil
}

func (b *Blockchain) Home() (path string, err error) {
	return b.chain.Home()
}

func (b *Blockchain) IsHomeDirExist() (ok bool, err error) {
	home, err := b.chain.Home()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(home)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// createOptions holds info about how to create a chain.
type createOptions struct {
	campaignID uint64
	noCheck    bool
}

// CreateOption configures chain creation.
type CreateOption func(*createOptions)

func WithCampaign(id uint64) CreateOption {
	return func(o *createOptions) {
		o.campaignID = id
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() CreateOption {
	return func(o *createOptions) {
		o.noCheck = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (b *Blockchain) Publish(ctx context.Context, options ...CreateOption) (launchID, campaignID uint64, err error) {
	o := createOptions{}
	for _, apply := range options {
		apply(&o)
	}

	// if the initial genesis is a genesis URL and no check are performed, we simply fetch it and get its hash
	if o.noCheck && b.genesisURL != "" {
		_, b.genesisHash, err = genesisAndHashFromURL(ctx, b.genesisURL)
		if err != nil {
			return 0, 0, err
		}
	}

	chainID, err := b.chain.ID()
	if err != nil {
		return 0, 0, err
	}

	coordinatorAddress := b.builder.account.Address(SPNAddressPrefix)
	campaignID = o.campaignID

	b.builder.ev.Send(events.New(events.StatusOngoing, "Publishing the network"))

	_, err = profiletypes.
		NewQueryClient(b.builder.cosmos.Context).
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
		if _, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateCoordinator); err != nil {
			return 0, 0, err
		}
	}

	if campaignID != 0 {
		_, err = campaigntypes.
			NewQueryClient(b.builder.cosmos.Context).
			Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
				Id: o.campaignID,
			})
		if err != nil {
			return 0, 0, err
		}
	} else {
		msgCreateCampaign := campaigntypes.NewMsgCreateCampaign(
			coordinatorAddress,
			b.chain.Name(),
			nil,
			false,
		)
		res, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateCampaign)
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
		b.builder.account.Address(SPNAddressPrefix),
		chainID,
		b.url,
		b.hash,
		b.genesisURL,
		b.genesisHash,
		true,
		campaignID,
	)
	res, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateChain)
	if err != nil {
		return 0, 0, err
	}

	var createChainRes launchtypes.MsgCreateChainResponse
	if err := res.Decode(&createChainRes); err != nil {
		return 0, 0, err
	}

	return createChainRes.LaunchID, campaignID, nil
}
