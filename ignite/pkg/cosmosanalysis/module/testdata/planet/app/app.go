package app

import (
	"github.com/cosmos/cosmos-sdk/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/types/module"
	planet "github.com/tendermint/planet/x/planet"
)

type Foo struct {
	FooKeeper foo.keeper
}

var ModuleBasics = module.NewBasicManager(planet.AppModuleBasic{})

func (f Foo) Name() string { return app.BaseApp.Name() }
func (f Foo) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (f Foo) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
func (f Foo) RegisterAPIRoutes()         {}
func (f Foo) RegisterTxService()         {}
func (f Foo) RegisterTendermintService() {}
