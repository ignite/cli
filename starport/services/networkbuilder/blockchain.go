package networkbuilder

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/types"
	conf "github.com/tendermint/starport/starport/chainconf"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain"
)

type Blockchain struct {
	appPath string
	url     string
	hash    string
	chain   *chain.Chain
	builder *Builder
}

func newBlockchain(
	ctx context.Context,
	builder *Builder,
	chainID,
	appPath,
	url,
	hash,
	home string,
	keyringBackend chaincmd.KeyringBackend,
	mustNotInitializedBefore bool,
) (*Blockchain, error) {
	bc := &Blockchain{
		appPath: appPath,
		url:     url,
		hash:    hash,
		builder: builder,
	}
	return bc, bc.init(ctx, chainID, home, keyringBackend, mustNotInitializedBefore)
}

// init initializes blockchain by building the binaries and running the init command and
// applies some post init configuration.
func (b *Blockchain) init(
	ctx context.Context,
	chainID,
	home string,
	keyringBackend chaincmd.KeyringBackend,
	mustNotInitializedBefore bool,
) error {
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))

	chainOption := []chain.Option{
		chain.LogLevel(chain.LogSilent),
		chain.ID(chainID),
	}

	// Custom home directories
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	if gitpod.IsOnGitpod() {
		chainOption = append(chainOption, chain.KeyringBackend(chaincmd.KeyringBackendTest))
	} else {
		// Otherwise use the keyring backend specified by the user
		chainOption = append(chainOption, chain.KeyringBackend(keyringBackend))
	}

	chain, err := chain.New(ctx, b.appPath, chainOption...)
	if err != nil {
		return err
	}

	if !chain.Version.Major().Is(cosmosver.Stargate) {
		return sperrors.ErrOnlyStargateSupported
	}
	chainHome, err := chain.Home()
	if err != nil {
		return err
	}

	if mustNotInitializedBefore {
		if _, err := os.Stat(chainHome); !os.IsNotExist(err) {
			return &DataDirExistsError{chainID, chainHome}
		}
	}

	// cleanup home dir of app if exists.
	if err := xos.RemoveAllUnderHome(chainHome); err != nil {
		return err
	}

	if _, err := chain.Build(ctx); err != nil {
		return err
	}
	if err := chain.Init(ctx, false); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))

	// backup initial genesis so it can be used during `start`.
	genesisPath, err := chain.GenesisPath()
	if err != nil {
		return err
	}
	genesis, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(initialGenesisPath(chainHome), genesis, 0644); err != nil {
		return err
	}

	b.chain = chain
	return nil
}

func genesisPath(appHome string) string {
	return fmt.Sprintf("%s/config/genesis.json", appHome)
}

func initialGenesisPath(appHome string) string {
	return fmt.Sprintf("%s/config/initial_genesis.json", appHome)
}

// BlockchainInfo hold information about a Blokchain.
type BlockchainInfo struct {
	Genesis          jsondoc.Doc
	Config           conf.Config
	RPCPublicAddress string
}

// Info returns information about the blockchain.
func (b *Blockchain) Info() (BlockchainInfo, error) {
	genesisPath, err := b.chain.GenesisPath()
	if err != nil {
		return BlockchainInfo{}, err
	}
	genesis, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return BlockchainInfo{}, err
	}
	config, err := b.chain.Config()
	if err != nil {
		return BlockchainInfo{}, err
	}
	paddr, err := b.chain.RPCPublicAddress()
	if err != nil {
		return BlockchainInfo{}, err
	}
	return BlockchainInfo{
		Genesis:          genesis,
		Config:           config,
		RPCPublicAddress: paddr,
	}, nil
}

// createOptions holds info about how to create a chain.
type createOptions struct {
	genesisURL string
}

// CreateOption configures chain creation.
type CreateOption func(*createOptions)

// WithCustomGenesisFromURL creates the chain with a custom one living at u.
func WithCustomGenesisFromURL(u string) CreateOption {
	return func(o *createOptions) {
		o.genesisURL = u
	}
}

// Create submits Genesis to SPN to announce a new network.
func (b *Blockchain) Create(ctx context.Context, options ...CreateOption) error {
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

	account, err := b.builder.AccountInUse()
	if err != nil {
		return err
	}
	chainID, err := b.chain.ID()
	if err != nil {
		return err
	}
	return b.builder.spnclient.ChainCreate(
		ctx,
		account.Name,
		chainID,
		b.url,
		b.hash,
		o.genesisURL,
		genesisHash,
	)
}

// Proposal holds proposal info of validator candidate to join to a network.
type Proposal struct {
	Validator chain.Validator
	Meta      ProposalMeta
}

type ProposalMeta struct {
	Website  string
	Identity string
	Details  string
}

type Account struct {
	Name     string
	Mnemonic string
	Coins    string
}

func (b *Blockchain) CreateAccount(ctx context.Context, account chain.Account) (chain.Account, error) {
	commands, err := b.chain.Commands(ctx)
	if err != nil {
		return chain.Account{}, err
	}

	acc, err := commands.AddAccount(ctx, account.Name, account.Mnemonic)
	if err != nil {
		return chain.Account{}, err
	}

	return chain.Account{
		Name:     acc.Name,
		Address:  acc.Address,
		Mnemonic: acc.Mnemonic,
	}, nil
}

// IssueGentx creates a Genesis transaction for account with proposal.
func (b *Blockchain) IssueGentx(ctx context.Context, account chain.Account, proposal Proposal) (gentx jsondoc.Doc, err error) {
	commands, err := b.chain.Commands(ctx)
	if err != nil {
		return nil, err
	}

	if err := commands.AddGenesisAccount(ctx, account.Address, account.Coins); err != nil {
		return nil, err
	}

	gentxPath, err := commands.Gentx(
		ctx,
		account.Name,
		proposal.Validator.StakingAmount,
		chaincmd.GentxWithMoniker(proposal.Validator.Moniker),
		chaincmd.GentxWithCommissionRate(proposal.Validator.CommissionRate),
		chaincmd.GentxWithCommissionMaxRate(proposal.Validator.CommissionMaxRate),
		chaincmd.GentxWithCommissionMaxChangeRate(proposal.Validator.CommissionMaxChangeRate),
		chaincmd.GentxWithMinSelfDelegation(proposal.Validator.MinSelfDelegation),
		chaincmd.GentxWithGasPrices(proposal.Validator.GasPrices),
	)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(gentxPath)
}

// Join proposes a validator to a network.
//
// address is the ip+port combination of a p2p address of a node (does not include id).
// https://docs.tendermint.com/master/spec/p2p/config.html.
func (b *Blockchain) Join(
	ctx context.Context,
	account *chain.Account,
	validatorAddress,
	publicAddress string,
	gentx []byte,
	selfDelegation types.Coin,
) error {
	commands, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}

	key, err := commands.ShowNodeID(ctx)
	if err != nil {
		return err
	}

	if xchisel.IsEnabled() {
		publicAddress = xchisel.ServerAddr()
	}

	p2pAddress := fmt.Sprintf("%s@%s", key, publicAddress)

	chainID, err := b.chain.ID()
	if err != nil {
		return err
	}

	var proposalOptions []spn.ProposalOption
	if account != nil {
		coins, err := types.ParseCoinsNormalized(account.Coins)
		if err != nil {
			return err
		}

		proposalOptions = append(proposalOptions, spn.AddAccountProposal(account.Address, coins))
	}

	proposalOptions = append(proposalOptions, spn.AddValidatorProposal(gentx, validatorAddress, selfDelegation, p2pAddress))

	return b.builder.Propose(ctx, chainID, proposalOptions...)
}

// Cleanup closes the event bus and cleanups everything related to installed blockchain.
func (b *Blockchain) Cleanup() error {
	b.builder.ev.Shutdown()
	return nil
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

type DataDirExistsError struct {
	ID   string
	Home string
}

func (e DataDirExistsError) Error() string {
	return "cannot initialize. chain's data dir already exists"
}
