package network

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
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
	GenesisPath() (string, error)
	GentxsPath() (string, error)
	DefaultGentxPath() (string, error)
	AppTOMLPath() (string, error)
	ConfigTOMLPath() (string, error)
	NodeID(ctx context.Context) (string, error)
	CacheBinary(launchID uint64) error
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
