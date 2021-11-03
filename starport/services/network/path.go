package network

import (
	"github.com/tendermint/starport/starport/pkg/xfilepath"
	"strconv"
)

var (
	// SpnPath returns the path used to store chain home from SPN
	SpnPath = xfilepath.JoinFromHome(xfilepath.Path("spn"))
)

// ChainHome returns the default home dir used for a chain from SPN
func ChainHome(launchID uint64) (string, error) {
	return xfilepath.Join(SpnPath, xfilepath.Path(strconv.FormatUint(launchID, 10)))()
}
