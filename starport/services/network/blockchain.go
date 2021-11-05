package network

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

// setup setups blockchain.
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

// Init initializes blockchain by building the binaries and running the init command and
// applies some post init configuration.
func (b *Blockchain) Init(ctx context.Context) error {
	chainHome, err := b.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err := os.RemoveAll(chainHome); err != nil {
		return err
	}

	// build the chain and initialize it with a new validator key
	b.builder.ev.Send(events.New(events.StatusOngoing, "Compile the blockchain"))
	if _, err := b.chain.Build(ctx, ""); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain compiled"))
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))
	if err := b.chain.Init(ctx, false); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))

	// write the custom genesis if a genesis URL is provided
	if b.genesisURL != "" {
		genesis, hash, err := genesisAndHashFromURL(ctx, b.genesisURL)
		if err != nil {
			return err
		}
		if hash != b.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. Expected hash %s, actual hash %s", b.genesisURL, b.genesisHash, hash)
		}
		genesisPath, err := b.chain.GenesisPath()
		if err != nil {
			return err
		}
		if err := os.WriteFile(genesisPath, genesis, 0644); err != nil {
			return err
		}
	}

	b.isInitialized = true

	return nil
}

// InitAccount initializes an account for the blockchain and issue a gentx in config/gentx/gentx.json
func (b *Blockchain) InitAccount(ctx context.Context, v chain.Validator, keyName, mnemonic string) (chaincmdrunner.Account, string, error) {
	if !b.isInitialized {
		return chaincmdrunner.Account{}, "", errors.New("the blockchain must be initialized to initialize an account")
	}

	// If no name is specified for the key, moniker is used
	if keyName == "" {
		keyName = v.Moniker
	}
	v.Name = keyName

	// create the chain account
	chainCmd, err := b.chain.Commands(ctx)
	if err != nil {
		return chaincmdrunner.Account{}, "", err
	}
	acc, err := chainCmd.AddAccount(ctx, keyName, mnemonic, "")
	if err != nil {
		return acc, "", err
	}

	// add account into the genesis
	err = chainCmd.AddGenesisAccount(ctx, acc.Address, v.StakingAmount)
	if err != nil {
		return acc, "", err
	}

	// create the gentx
	issuedGentxPath, err := b.chain.IssueGentx(ctx, v)
	if err != nil {
		return acc, "", err
	}

	// rename the issued gentx into gentx.json
	gentxPath := filepath.Join(filepath.Dir(issuedGentxPath), gentxFilename)
	return acc, gentxPath, os.Rename(issuedGentxPath, gentxPath)
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
