package keeper

import (
	"context"

	"github.com/ignite/planet/x/mars/types"
)

type Keeper struct{}

func (k Keeper) QuerySimple(goCtx context.Context, req *types.QuerySimpleRequest) (*types.QuerySimpleResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySimpleParams(goCtx context.Context, req *types.QuerySimpleParamsRequest) (*types.QuerySimpleParamsResponse, error) {
	return nil, nil
}

func (k Keeper) QueryParamsWithPagination(goCtx context.Context, req *types.QueryWithPaginationRequest) (*types.QueryWithPaginationResponse, error) {
	return nil, nil
}

func (k Keeper) QueryWithQueryParams(goCtx context.Context, req *types.QueryWithQueryParamsRequest) (*types.QueryWithQueryParamsResponse, error) {
	return nil, nil
}

func (k Keeper) QueryWithQueryParamsWithPagination(goCtx context.Context, req *types.QueryWithQueryParamsWithPaginationRequest) (*types.QueryWithQueryParamsWithPaginationResponse, error) {
	return nil, nil
}
