package keeper

import (
	"github.com/cosmonaut/blog/x/blog/types"
)

var _ types.QueryServer = Keeper{}
