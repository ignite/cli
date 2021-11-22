package network

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	// ChainHomeRoot is the root dir for all spn chain homes
	ChainHomeRoot = "spn"

	// ChainHomeInitDir is the chain homes directory for validator initialization and local testing
	ChainHomeInitDir = "init"

	// ChainHomeLaunchDir is the chain homes directory for chain launch preparation
	ChainHomeLaunchDir = "launch"
)

// IsChainHomeExist checks if a home with the provided launchID already exist
func IsChainHomeExist(launchID uint64, isLaunchPreparation bool) (string, bool, error) {
	home, err := ChainHome(launchID, isLaunchPreparation)
	if err != nil {
		return home, false, err
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		return home, false, nil
	}
	return home, true, err
}

// ChainHome returns the default home dir used for a chain from SPN
func ChainHome(launchID uint64, isLaunchPreparation bool) (string, error) {
	var chainHomeDir string
	if isLaunchPreparation {
		chainHomeDir = ChainHomeLaunchDir
	} else {
		chainHomeDir = ChainHomeInitDir
	}

	home, err := os.UserHomeDir()
	return filepath.Join(home, ChainHomeRoot, chainHomeDir, strconv.FormatUint(launchID, 10)), err
}
