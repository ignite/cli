package networkbuilder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/chain/conf"
)

type Blockchain struct {
	appPath string
	url     string
	hash    string
	chain   *chain.Chain
	app     chain.App
	builder *Builder
}

func newBlockchain(ctx context.Context, builder *Builder, chainID, appPath, url, hash string,
	mustNotInitializedBefore bool) (*Blockchain, error) {
	bc := &Blockchain{
		appPath: appPath,
		url:     url,
		hash:    hash,
		builder: builder,
	}
	return bc, bc.init(ctx, chainID, mustNotInitializedBefore)
}

// init initializes blockchain by building the binaries and running the init command and
// applies some post init configuration.
func (b *Blockchain) init(ctx context.Context, chainID string, mustNotInitializedBefore bool) error {
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))

	path, err := gomodulepath.ParseFile(b.appPath)
	if err != nil {
		return err
	}
	app := chain.App{
		ChainID: chainID,
		Name:    path.Root,
		Path:    b.appPath,
	}

	c, err := chain.New(app, false, chain.LogSilent)
	if err != nil {
		return err
	}

	if mustNotInitializedBefore {
		if _, err := os.Stat(c.Home()); !os.IsNotExist(err) {
			return &DataDirExistsError{chainID, c.Home()}
		}
	}

	// cleanup home dir of app if exists.
	for _, path := range c.StoragePaths() {
		if err := xos.RemoveAllUnderHome(path); err != nil {
			return err
		}
	}

	if err := c.Build(ctx); err != nil {
		return err
	}
	if err := c.Init(ctx); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))

	// backup initial genesis so it can be used during `start`.
	genesis, err := ioutil.ReadFile(c.GenesisPath())
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(initialGenesisPath(c.Home()), genesis, 0644); err != nil {
		return err
	}

	b.chain = c
	b.app = app
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
	genesis, err := ioutil.ReadFile(b.chain.GenesisPath())
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

// Create submits Genesis to SPN to announce a new network.
func (b *Blockchain) Create(ctx context.Context) error {
	account, err := b.builder.AccountInUse()
	if err != nil {
		return err
	}
	chainID, err := b.chain.ID()
	if err != nil {
		return err
	}
	return b.builder.spnclient.ChainCreate(ctx, account.Name, chainID, b.url, b.hash)
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
	Address  string
	Mnemonic string
	Coins    string
}

// IssueGentx creates a Genesis transaction for account with proposal.
func (b *Blockchain) IssueGentx(ctx context.Context, account Account, proposal Proposal) (gentx jsondoc.Doc, err error) {
	proposal.Validator.Name = account.Name
	if err := b.chain.AddGenesisAccount(ctx, chain.Account{
		Address: account.Address,
		Coins:   account.Coins,
	}); err != nil {
		return nil, err
	}
	gentxPath, err := b.chain.Gentx(ctx, proposal.Validator)
	if err != nil {
		return nil, err
	}
	gentx, err = ioutil.ReadFile(gentxPath)
	return gentx, err
}

// Join proposes a validator to a network.
//
// address is the ip+port combination of a p2p address of a node (does not include id).
// https://docs.tendermint.com/master/spec/p2p/config.html.
func (b *Blockchain) Join(ctx context.Context, accountAddress, publicAddress string, coins types.Coins, gentx []byte, selfDelegation types.Coin) error {
	key, err := b.chain.ShowNodeID(ctx)
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

	return b.builder.Propose(
		ctx,
		chainID,
		spn.AddAccountProposal(accountAddress, coins),
		spn.AddValidatorProposal(gentx, accountAddress, selfDelegation, p2pAddress),
	)
}

// Cleanup closes the event bus and cleanups everyting related to installed blockchain.
func (b *Blockchain) Cleanup() error {
	b.builder.ev.Shutdown()
	//return os.RemoveAll(b.appPath)
	return nil
}

type DataDirExistsError struct {
	ID   string
	Home string
}

func (e DataDirExistsError) Error() string {
	return "cannot initialize. chain's data dir already exists"
}
