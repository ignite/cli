package keeper

import (
	"github.com/ignite/mars/x/mars/types"
)

var _ types.QueryServer = Keeper{}
