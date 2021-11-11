package network

import (
	"os"
	"path/filepath"
	"strconv"
)

// ChainHomeRoot is the root dir for spn chain homes
const (
	ChainHomeRoot = "spn"
	gentxPath     = "config/gentx/gentx.json"
	genesisPath   = "config/genesis.json"
)

// IsChainHomeExist checks if a home with the provided launchID already exist
func IsChainHomeExist(launchID uint64) (string, bool, error) {
	home, err := ChainHome(launchID)
	if err != nil {
		return home, false, err
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		return home, false, nil
	}
	return home, true, err
}

// ChainHome returns the default home dir used for a chain from SPN
func ChainHome(launchID uint64) (string, error) {
	home, err := os.UserHomeDir()
	return filepath.Join(home, ChainHomeRoot, strconv.FormatUint(launchID, 10)), err
}

// Gentx returns the default gentx path into the home dir
func Gentx(home string) string {
	return filepath.Join(home, gentxPath)
}

// getChainGenesis return the chain genesis path
func getChainGenesis(home string) (ChainGenesis, bool, error) {
	genesisPath := filepath.Join(home, genesisPath)
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

// CheckGenesisAddress returns true if the address exist into the genesis file
func CheckGenesisAddress(home, addr string) (bool, error) {
	genesis, exist, err := getChainGenesis(home)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	return genesis.HasAccount(addr), nil
}
