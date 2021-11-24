package cosmosutil

import (
	"os"
	"path/filepath"
)

// genesisPath returns the default genesis path into the home dir
func genesisPath(home string) string {
	return filepath.Join(home, "config/genesis.json")
}

// getChainGenesis return the chain genesis path
func getChainGenesis(home string) (ChainGenesis, bool, error) {
	genesisPath := genesisPath(home)
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return ChainGenesis{}, false, nil
	} else if err != nil {
		return ChainGenesis{}, false, err
	}
	net, err := ParseGenesis(genesisPath)
	if err != nil {
		return ChainGenesis{}, false, err
	}
	return net, true, nil
}

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(home, addr string) (bool, error) {
	genesis, exist, err := getChainGenesis(home)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	return genesis.HasAccount(addr), nil
}
