package genesis

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/jsonfile"
)

const (
	genesisFilename = "genesis.json"
	paramStakeDenom = "app_state.staking.params.bond_denom"
	paramChainID    = "chain_id"
	paramAccounts   = "app_state.auth.accounts"
)

type (
	// Genesis represents the genesis reader
	Genesis struct {
		*jsonfile.JSONFile
	}
	accounts []struct {
		Address string `json:"address"`
	}
)

// FromPath parse genesis object from path
func FromPath(path string) (*Genesis, error) {
	file, err := jsonfile.FromPath(path)
	return &Genesis{
		JSONFile: file,
	}, err
}

// FromURL fetches the genesis from the given URL and returns its content.
func FromURL(ctx context.Context, url, path string) (*Genesis, error) {
	file, err := jsonfile.FromURL(ctx, url, path, genesisFilename)
	return &Genesis{
		JSONFile: file,
	}, err
}

// StakeDenom returns the stake denom from the genesis
func (g *Genesis) StakeDenom() (denom string, err error) {
	_, err = g.Param(paramStakeDenom, &denom)
	return
}

// ChainID returns the chain id from the genesis
func (g *Genesis) ChainID() (chainID string, err error) {
	_, err = g.Param(paramChainID, &chainID)
	return
}

// Accounts returns the auth accounts from the genesis
func (g *Genesis) Accounts() ([]string, error) {
	var accs accounts
	_, err := g.Param(paramAccounts, &accs)
	accountList := make([]string, len(accs))
	for i, acc := range accs {
		accountList[i] = acc.Address
	}
	return accountList, err
}
