package module_create

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var templates = map[cosmosver.MajorVersion]*packr.Box{
	cosmosver.Launchpad: packr.New("module/create/templates/launchpad", "./launchpad"),
	cosmosver.Stargate:  packr.New("module/create/templates/stargate", "./stargate"),
}
