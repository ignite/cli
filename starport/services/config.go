package services

import (
	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

var (
	// StarportConfPath returns the Starport Configuration directory
	StarportConfPath = xfilepath.JoinFromHome(xfilepath.Path(".starport"))
)
