package cosmosutil

import (
	"encoding/json"
	"errors"
	"os"
)

// ChainGenesis represents the stargate genesis file
type ChainGenesis struct {
	AppState struct {
		Auth struct {
			Accounts []struct {
				Address string `json:"address"`
			} `json:"accounts"`
		} `json:"auth"`
	} `json:"app_state"`
}

// HasAccount check if account exist into the genesis account
func (g ChainGenesis) HasAccount(address string) bool {
	for _, account := range g.AppState.Auth.Accounts {
		if account.Address == address {
			return true
		}
	}
	return false
}

// ParseGenesis parse ChainGenesis object from a genesis file
func ParseGenesis(genesisPath string) (genesis ChainGenesis, err error) {
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return genesis, errors.New("cannot open genesis file: " + err.Error())
	}
	return genesis, json.Unmarshal(genesisFile, &genesis)
}

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := ParseGenesis(genesisPath)
	if err != nil {
		return false, err
	}
	return genesis.HasAccount(addr), nil
}
