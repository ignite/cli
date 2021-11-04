package network

import (
	"os"
	"path/filepath"
	"strconv"
)

// ChainHomeRoot is the root dir for spn chain homes
const ChainHomeRoot = "spn"

// ChainHome returns the default home dir used for a chain from SPN
func ChainHome(launchID uint64) (string, error) {
	home, err := os.UserHomeDir()
	return filepath.Join(home, ChainHomeRoot, strconv.FormatUint(launchID, 10)), err
}
