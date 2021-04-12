package services

import (
	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

const (
	starportConfDir = ".starport"
)

var (
	// StarportConfPath returns the Starport Configuration directory
	StarportConfPath = xfilepath.JoinFromHome(xfilepath.Path(starportConfDir))
)
