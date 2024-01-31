package keeper

import (
	"context"

	"github.com/tendermint/planet/x/mars/types"
)

type Keeper struct{}

func (k Keeper) MyQuery(goCtx context.Context, req *types.QueryMyQueryRequest) (*types.QueryMyQueryResponse, error) {
	return nil, nil
}

func (k Keeper) Foo(goCtx context.Context, req *types.QueryFooRequest) (*types.QueryFooResponse, error) {
	return nil, nil
}
