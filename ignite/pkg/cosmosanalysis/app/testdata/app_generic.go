package foo

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Foo[T any] struct {
	*runtime.App

	FooKeeper foo.keeper
	i         T
}

func (f Foo[T]) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return f.App.BeginBlocker(ctx, req)
}

func (f Foo[T]) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return f.App.EndBlocker(ctx, req)
}
