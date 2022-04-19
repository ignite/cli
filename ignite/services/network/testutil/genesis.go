package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Genesis struct {
		ChainID  string   `json:"chain_id"`
		AppState AppState `json:"app_state"`
	}

	AppState struct {
		Auth    Auth    `json:"auth"`
		Staking Staking `json:"staking"`
	}

	Auth struct {
		Accounts []GenesisAccount `json:"accounts"`
	}

	GenesisAccount struct {
		Address string `json:"address"`
	}

	Staking struct {
		Params StakingParams `json:"params"`
	}

	StakingParams struct {
		BondDenom string `json:"bond_denom"`
	}
)

// NewGenesis creates easily modifiable genesis object for testing purposes
func NewGenesis(chainID string) *Genesis {
	return &Genesis{ChainID: chainID}
}

// AddAccount adds account to the genesis
func (g *Genesis) AddAccount(address string) *Genesis {
	g.AppState.Auth.Accounts = append(g.AppState.Auth.Accounts, GenesisAccount{Address: address})
	return g
}

// SaveTo saves genesis json representation to the specified directory and returns full path
func (g *Genesis) SaveTo(t *testing.T, dir string) string {
	encoded, err := json.Marshal(g)
	assert.NoError(t, err)
	savePath := filepath.Join(dir, "genesis.json")
	err = os.WriteFile(savePath, encoded, 0666)
	assert.NoError(t, err)
	return savePath
}
