package network

import (
	"os"
	"path/filepath"
	"strconv"
)

// ChainHomeRoot is the root dir for spn chain homes
const ChainHomeRoot = "spn"

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