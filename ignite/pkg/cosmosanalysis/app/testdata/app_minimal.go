package foo

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Foo struct {
	*baseapp.BaseApp

	FooKeeper foo.keeper
}

func (f Foo) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return f.BaseApp.BeginBlocker(ctx, req)
}

func (f Foo) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return f.BaseApp.EndBlocker(ctx, req)
}
