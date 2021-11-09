package network

import (
	"context"
	"os"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	gentxFilename = "gentx.json"
)

// Blockchain represents a blockchain.
type Blockchain struct {
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
	genesisURL string
	noCheck    bool
}

// CreateOption configures chain creation.
type CreateOption func(*createOptions)

// WithCustomGenesisFromURL creates the chain with a custom one living at u.
func WithCustomGenesisFromURL(u string) CreateOption {
	return func(o *createOptions) {
		o.genesisURL = u
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() CreateOption {
	return func(o *createOptions) {
		o.noCheck = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (b *Blockchain) Publish(ctx context.Context, options ...CreateOption) error {
	o := createOptions{}
	for _, apply := range options {
		apply(&o)
	}

	var err error
	var genesisHash string

	// if check must be performed, we initialize the chain to check the initial genesis
	if !o.noCheck {
		if err := b.Init(ctx); err != nil {
			return err
		}
		genesisHash = b.genesisHash
	} else if b.genesisURL != "" {
		// if the initial genesis is a genesis URL and no check are performed, we simply fetched and get its hash
		_, genesisHash, err = genesisAndHashFromURL(ctx, b.genesisURL)
		if err != nil {
			return err
		}
	}

	chainID, err := b.chain.ID()
	if err != nil {
		return err
	}

	_, err = profiletypes.
		NewQueryClient(b.builder.cosmos.Context).
		CoordinatorByAddress(ctx, &profiletypes.QueryGetCoordinatorByAddressRequest{
			Address: b.builder.account.Address(SPNAddressPrefix),
		})

	// TODO check for not found and only then create a new coordinator, otherwise return the err.
	if err != nil {
		msgCreateCoordinator := profiletypes.NewMsgCreateCoordinator(
			b.builder.account.Address(SPNAddressPrefix),
			"",
			"",
			"",
		)
		if _, err := b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateCoordinator); err != nil {
			return err
		}
	}

	msgCreateChain := launchtypes.NewMsgCreateChain(
		b.builder.account.Address(SPNAddressPrefix),
		chainID,
		b.url,
		b.hash,
		o.genesisURL,
		genesisHash,
		false,
		0,
	)
	_, err = b.builder.cosmos.BroadcastTx(b.builder.account.Name, msgCreateChain)
	return err
}
