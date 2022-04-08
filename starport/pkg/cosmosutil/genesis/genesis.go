package genesis

import (
	"context"
	"os"

	"github.com/tendermint/starport/starport/pkg/jsonfile"
)

const (
	genesisFilename     = "genesis.json"
	fieldPathStakeDenom = "app_state.staking.params.bond_denom"
	fieldPathChainID    = "chain_id"
	fieldPathAccounts   = "app_state.auth.accounts"
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

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := FromPath(genesisPath)
	if err != nil {
		return false, err
	}
	defer genesis.Close()
	return genesis.HasAccount(addr), nil
}

// HasAccount check if account exist into the genesis account
func (g Genesis) HasAccount(address string) bool {
	accounts, err := g.Accounts()
	if err != nil {
		return false
	}
	for _, account := range accounts {
		if account == address {
			return true
		}
	}
	return false
}

// StakeDenom returns the stake denom from the genesis
func (g *Genesis) StakeDenom() (denom string, err error) {
	_, err = g.Field(fieldPathStakeDenom, &denom)
	return
}

// ChainID returns the chain id from the genesis
func (g *Genesis) ChainID() (chainID string, err error) {
	_, err = g.Field(fieldPathChainID, &chainID)
	return
}

// Accounts returns the auth accounts from the genesis
func (g *Genesis) Accounts() ([]string, error) {
	var accs accounts
	_, err := g.Field(fieldPathAccounts, &accs)
	accountList := make([]string, len(accs))
	for i, acc := range accs {
		accountList[i] = acc.Address
	}
	return accountList, err
}
