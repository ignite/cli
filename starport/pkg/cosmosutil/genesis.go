package cosmosutil

import (
	"os"
	"path/filepath"
)

const (
	gentxPath   = "config/gentx/gentx.json"
	genesisPath = "config/genesis.json"
)

// GenesisPath returns the default genesis path into the home dir
func GenesisPath(home string) string {
	return filepath.Join(home, genesisPath)
}

// GentxPath returns the default gentx path into the home dir
func GentxPath(home string) string {
	return filepath.Join(home, gentxPath)
}

// getChainGenesis return the chain genesis path
func getChainGenesis(home string) (ChainGenesis, bool, error) {
	genesisPath := GenesisPath(home)
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
