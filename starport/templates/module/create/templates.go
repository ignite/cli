package modulecreate

import (
	"github.com/gobuffalo/packr/v2"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var (
	launchpadTemplate = packr.New("module/create/templates/launchpad", "./launchpad")
	stargateTemplate  = packr.New("module/create/templates/stargate", "./stargate")
	ibcTemplate       = packr.New("module/create/templates/ibc", "./ibc")
)
