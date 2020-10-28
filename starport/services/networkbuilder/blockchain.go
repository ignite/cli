package networkbuilder

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/chain/conf"
)

type Blockchain struct {
	appPath string
	chain   *chain.Chain
	ev      events.Bus
}

// BlockchainInfo hold information about a Blokchain.
type BlockchainInfo struct {
	Genesis []byte
	Config  conf.Config
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
	return BlockchainInfo{
		Genesis: genesis,
		Config:  config,
	}, nil
}

// Create submits Genesis to SPN to announce a new network.
func (b *Blockchain) Create(ctx context.Context, genesis []byte) error { return nil }

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

// IssueGentx creates a Genesis transaction for account with proposal.
func (b *Blockchain) IssueGentx(ctx context.Context, account Account, proposal Proposal) (gentx interface{}, mnemonic string, err error) {
	proposal.Validator.Name = account.Name
	mnemonic, err = b.chain.CreateAccount(ctx, account.Name, account.Mnemonic, strings.Split(account.Coins, ","), false)
	if err != nil {
		return "", "", err
	}
	gentxPath, err := b.chain.Gentx(ctx, proposal.Validator)
	if err != nil {
		return "", "", err
	}
	gentxFile, err := os.Open(gentxPath)
	if err != nil {
		return "", "", err
	}
	defer gentxFile.Close()
	if err := json.NewDecoder(gentxFile).Decode(&gentx); err != nil {
		return "", "", err
	}
	return gentx, mnemonic, nil
}

// Join proposes a gentx via SPN to a network.
func (b *Blockchain) Join(ctx context.Context, gentx interface{}) error {
	return nil
}

// Cleanup closes the event bus and cleanups everyting related to installed blockchain.
func (b *Blockchain) Cleanup() error {
	b.ev.Shutdown()
	return os.RemoveAll(b.appPath)
}
