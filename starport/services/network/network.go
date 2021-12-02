package network

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

// Network is network builder.
type Network struct {
	ev      events.Bus
	cosmos  cosmosclient.Client
	account cosmosaccount.Account
}

type Chain interface {
	ID() (string, error)
	Name() string
	SourceURL() string
	SourceHash() string
	GentxPath() (string, error)
	GenesisPath() (string, error)
	Peer(ctx context.Context, addr string) (string, error)
}

type Option func(*Network)

// CollectEvents collects events from the network builder.
func CollectEvents(ev events.Bus) Option {
	return func(b *Network) {
		b.ev = ev
	}
}

// New creates a Builder.
func New(cosmos cosmosclient.Client, account cosmosaccount.Account, options ...Option) (Network, error) {
	n := Network{
		cosmos:  cosmos,
		account: account,
	}
	for _, opt := range options {
		opt(&n)
	}
	return n, nil
}

type LaunchInfo = networkchain.Launch

// LaunchInfo fetches the chain launch from Starport Network by launch id.
func (n Network) LaunchInfo(ctx context.Context, id uint64) (LaunchInfo, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))

	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: id,
	})
	if err != nil {
		return LaunchInfo{}, err
	}

	info := LaunchInfo{
		ID:         id,
		ChainID:    res.Chain.GenesisChainID,
		SourceURL:  res.Chain.SourceURL,
		SourceHash: res.Chain.SourceHash,
	}

	// check if custom genesis URL is provided.
	if customGenesisURL := res.Chain.InitialGenesis.GetGenesisURL(); customGenesisURL != nil {
		info.GenesisURL = customGenesisURL.Url
		info.GenesisHash = customGenesisURL.Hash
	}

	n.ev.Send(events.New(events.StatusOngoing, "Chain information fetched"))
	return info, nil
}

func ParseLaunchID(id string) (uint64, error) {
	launchID, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return 0, errors.New("launch ID must be greater than 0")
	}
	return launchID, nil
}
