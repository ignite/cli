package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func NewGenesis(chainID string) *Genesis {
	return &Genesis{ChainID: chainID}
}

func (g *Genesis) AddAccount(address string) {
	g.AppState.Auth.Accounts = append(g.AppState.Auth.Accounts, GenesisAccount{Address: address})
}

func (g *Genesis) SaveTo(dir string) (string, error) {
	encoded, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	savePath := filepath.Join(dir, "genesis.json")
	return savePath, os.WriteFile(savePath, encoded, 0666)
}
