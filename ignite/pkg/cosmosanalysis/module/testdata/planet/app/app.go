package app

import (
	"cosmossdk.io/client/v2/autocli"
	"github.com/cosmos/cosmos-sdk/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	app "github.com/ignite/cli/ignite/pkg/cosmosanalysis/app/testdata/modules/registration_not_in_app_go"

	"github.com/tendermint/planet/x/mars"
)

type Foo struct {
	FooKeeper foo.keeper
}

var ModuleBasics = module.NewBasicManager(mars.AppModuleBasic{})

func (Foo) Name() string                                    { return app.BaseApp.Name() }
func (Foo) GetKey(storeKey string) *storetypes.KVStoreKey   { return nil }
func (Foo) RegisterAPIRoutes()                              {}
func (Foo) RegisterTxService()                              {}
func (Foo) RegisterTendermintService()                      {}
func (Foo) InterfaceRegistry() codectypes.InterfaceRegistry { return nil }
func (Foo) TxConfig() client.TxConfig                       { return nil }
func (Foo) AutoCliOpts() autocli.AppOptions                 { return autocli.AppOptions{} }
func (Foo) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (Foo) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
