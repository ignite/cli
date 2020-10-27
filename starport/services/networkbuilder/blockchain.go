package networkbuilder

import (
	"context"
	"io/ioutil"
	"os"

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
	Moniker       string
	StakingAmount int32
	Account       conf.Account
}

// Join proposes validator to be added to network via SPN.
func (b *Blockchain) Join(ctx context.Context, proposal Proposal) error { return nil }

// Cleanup closes the event bus and cleanups everyting related to installed blockchain.
func (b *Blockchain) Cleanup() error {
	b.ev.Shutdown()
	return os.RemoveAll(b.appPath)
}
