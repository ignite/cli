package keeper

import (
	"context"

	"github.com/tendermint/planet/x/mars/types"
)

type Keeper struct{}

func (k Keeper) MyQuery(goCtx context.Context, req *types.QueryMyQueryRequest) (*types.QueryMyQueryResponse, error) {
	return nil, nil
}
