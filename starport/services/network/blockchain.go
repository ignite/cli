package network

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
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
	genesisURL    string
	genesisHash   string
}

// setup setups blockchain.
func (b *Blockchain) setup(
	launchID,
	home string,
	keyringBackend chaincmd.KeyringBackend,
) error {
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))

	chainOption := []chain.Option{
		chain.LogLevel(chain.LogSilent),
		chain.ID(launchID),
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

// publishOptions holds info about how to create a chain.
type publishOptions struct {
	genesisURL string
	noCheck    bool
}

// PublishOption configures chain creation.
type PublishOption func(*publishOptions)

// WithCustomGenesisFromURL creates the chain with a custom one living at u.
func WithCustomGenesisFromURL(u string) PublishOption {
	return func(o *publishOptions) {
		o.genesisURL = u
	}
}

// WithNoCheck disables checking integrity of the chain.
func WithNoCheck() PublishOption {
	return func(o *publishOptions) {
		o.noCheck = true
	}
}

// Publish submits Genesis to SPN to announce a new network.
func (b *Blockchain) Publish(ctx context.Context, options ...PublishOption) error {
	o := publishOptions{}
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
			return err
		}

		if !o.noCheck {
			if !b.isInitialized {
				if err := b.Init(ctx); err != nil {
					return err
				}
			}

			genesisPath, err := b.chain.GenesisPath()
			if err != nil {
				return err
			}

			if err := os.WriteFile(genesisPath, genesis, 0666); err != nil {
				return err
			}

			commands, err := b.chain.Commands(ctx)
			if err != nil {
				return err
			}

			if err := commands.ValidateGenesis(ctx); err != nil {
				return err
			}
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

func genesisAndHashFromURL(ctx context.Context, url string) (genesis []byte, hash string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

// joinOptions holds info about how to create a chain.
type joinOptions struct {
	amount uint64
	gentx  []byte
}

// JoinOption configures chain creation.
type JoinOption func(*joinOptions)

// WithAmount join the chain with a custom amount for create account message.
func WithAmount(amount uint64) JoinOption {
	return func(o *joinOptions) {
		o.amount = amount
	}
}

// WithGentx join the chain with a custom gentx json file.
func WithGentx(gentx []byte) JoinOption {
	return func(o *joinOptions) {
		o.gentx = gentx
	}
}

// Join to the network.
func (b *Blockchain) Join(ctx context.Context, options ...JoinOption) error {
	o := joinOptions{}
	for _, apply := range options {
		apply(&o)
	}

	commands, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}

	key, err := commands.ShowNodeID(ctx)
	if err != nil {
		return err
	}

	//ca, err := cosmosaccount.New(cosmosaccount.WithKeyringBackend(o.)
	//if err != nil {
	//	return err
	//}
	//
	//p2pAddress := fmt.Sprintf("%s@%s", key, publicAddress)
	//
	//chainID, err := b.chain.ID()
	//if err != nil {
	//	return err
	//}
	//
	//var proposalOptions []spn.ProposalOption
	//if account != nil {
	//	coins, err := types.ParseCoinsNormalized(account.Coins)
	//	if err != nil {
	//		return err
	//	}
	//
	//	proposalOptions = append(proposalOptions, spn.AddAccountProposal(account.Address, coins))
	//}
	//
	//proposalOptions = append(proposalOptions, spn.AddValidatorProposal(gentx, validatorAddress, selfDelegation, p2pAddress))
	//
	//return b.builder.Propose(ctx, chainID, proposalOptions...)
}
