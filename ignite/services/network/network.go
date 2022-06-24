package network

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/events"
)

//go:generate mockery --name CosmosClient --case underscore
type CosmosClient interface {
	Account(accountName string) (cosmosaccount.Account, error)
	Address(accountName string) (sdktypes.AccAddress, error)
	Context() client.Context
	BroadcastTx(accountName string, msgs ...sdktypes.Msg) (cosmosclient.Response, error)
	BroadcastTxWithProvision(accountName string, msgs ...sdktypes.Msg) (gas uint64, broadcast func() (cosmosclient.Response, error), err error)
	Status(ctx context.Context) (*ctypes.ResultStatus, error)
	ConsensusInfo(ctx context.Context, height int64) (cosmosclient.ConsensusInfo, error)
}

// Network is network builder.
type Network struct {
	ev            events.Bus
	cosmos        CosmosClient
	account       cosmosaccount.Account
	campaignQuery campaigntypes.QueryClient
	launchQuery   launchtypes.QueryClient
	profileQuery  profiletypes.QueryClient
	rewardQuery   rewardtypes.QueryClient
	stakingQuery  stakingtypes.QueryClient
}

//go:generate mockery --name Chain --case underscore
type Chain interface {
	ID() (string, error)
	ChainID() (string, error)
	Name() string
	SourceURL() string
	SourceHash() string
	GenesisPath() (string, error)
	GentxsPath() (string, error)
	DefaultGentxPath() (string, error)
	AppTOMLPath() (string, error)
	ConfigTOMLPath() (string, error)
	NodeID(ctx context.Context) (string, error)
	CacheBinary(launchID uint64) error
	ResetGenesisTime() error
}

type Option func(*Network)

func WithCampaignQueryClient(client campaigntypes.QueryClient) Option {
	return func(n *Network) {
		n.campaignQuery = client
	}
}

func WithProfileQueryClient(client profiletypes.QueryClient) Option {
	return func(n *Network) {
		n.profileQuery = client
	}
}

func WithLaunchQueryClient(client launchtypes.QueryClient) Option {
	return func(n *Network) {
		n.launchQuery = client
	}
}

func WithRewardQueryClient(client rewardtypes.QueryClient) Option {
	return func(n *Network) {
		n.rewardQuery = client
	}
}

func WithStakingQueryClient(client stakingtypes.QueryClient) Option {
	return func(n *Network) {
		n.stakingQuery = client
	}
}

// CollectEvents collects events from the network builder.
func CollectEvents(ev events.Bus) Option {
	return func(n *Network) {
		n.ev = ev
	}
}

// New creates a Builder.
func New(cosmos CosmosClient, account cosmosaccount.Account, options ...Option) Network {
	n := Network{
		cosmos:        cosmos,
		account:       account,
		campaignQuery: campaigntypes.NewQueryClient(cosmos.Context()),
		launchQuery:   launchtypes.NewQueryClient(cosmos.Context()),
		profileQuery:  profiletypes.NewQueryClient(cosmos.Context()),
		rewardQuery:   rewardtypes.NewQueryClient(cosmos.Context()),
		stakingQuery:  stakingtypes.NewQueryClient(cosmos.Context()),
	}
	for _, opt := range options {
		opt(&n)
	}
	return n
}

func ParseID(id string) (uint64, error) {
	objID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "error parsing launchID")
	}
	if objID == 0 {
		return 0, errors.New("launch ID must be greater than 0")
	}
	return objID, nil
}
