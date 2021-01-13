package module

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

const (
	PathAppGo    = "app/app.go"
	PathExportGo = "app/export.go"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var templates = map[cosmosver.MajorVersion]*packr.Box{
	cosmosver.Launchpad: packr.New("module/templates/launchpad", "./launchpad"),
	cosmosver.Stargate:  packr.New("module/templates/stargate", "./stargate"),
}
