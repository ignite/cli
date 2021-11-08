package network

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain"
)

// Blockchain represents a blockchain.
type Blockchain struct {
	appPath       string
	url           string
	hash          string
	chain         *chain.Chain
	isInitialized bool
	builder       *Builder
}

// setup setups blockchain.
func (b *Blockchain) setup(
	chainID,
	home string,
	keyringBackend chaincmd.KeyringBackend,
) error {
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))

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

// Init initializes blockchain by building the binaries and running the init command and
// applies some post init configuration.
func (b *Blockchain) Init(ctx context.Context) error {
	chainHome, err := b.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err := xos.RemoveAllUnderHome(chainHome); err != nil {
		return err
	}

	if _, err := b.chain.Build(ctx, ""); err != nil {
		return err
	}

	if err := b.chain.Init(ctx, false); err != nil {
		return err
	}

	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))
	b.isInitialized = true

	return nil
}

// createOptions holds info about how to create a chain.
type createOptions struct {
	genesisURL string
	campaignID uint64
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

	var genesisHash string

	if o.genesisURL != "" {
		// download the custom given genesis, validate it and calculate its hash.
		var genesis []byte
		var err error

		genesis, genesisHash, err = genesisAndHashFromURL(ctx, o.genesisURL)
		if err != nil {
			return 0, 0, err
		}

		if !o.noCheck {
			if !b.isInitialized {
				if err := b.Init(ctx); err != nil {
					return 0, 0, err
				}
			}

			genesisPath, err := b.chain.GenesisPath()
			if err != nil {
				return 0, 0, err
			}

			if err := os.WriteFile(genesisPath, genesis, 0666); err != nil {
				return 0, 0, err
			}

			commands, err := b.chain.Commands(ctx)
			if err != nil {
				return 0, 0, err
			}

			if err := commands.ValidateGenesis(ctx); err != nil {
				return 0, 0, err
			}
		}
	}

	chainID, err := b.chain.ID()
	if err != nil {
		return 0, 0, err
	}

	coordinatorAddress := b.builder.account.Address(SPNAddressPrefix)
	campaignID = o.campaignID

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
			"default",
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
		o.genesisURL,
		genesisHash,
		false,
		0,
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

func genesisAndHashFromURL(ctx context.Context, u string) (genesis []byte, hash string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	genesis, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(genesis)); err != nil {
		return nil, "", err
	}

	hexhash := hex.EncodeToString(h.Sum(nil))

	return genesis, hexhash, nil
}
