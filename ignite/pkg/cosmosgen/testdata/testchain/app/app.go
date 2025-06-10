package app

import (
	"cosmossdk.io/client/v2/autocli"
	"github.com/cosmos/cosmos-sdk/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/client"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gogo/protobuf/codec"
	marskeeper "github.com/tendermint/planet/x/mars/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Foo struct {
	MarsKeeper marskeeper.Keeper
}

var ModuleBasics = module.NewBasicManager(mars.AppModuleBasic{})

func (Foo) Name() string                                         { return app.BaseApp.Name() }
func (Foo) RegisterAPIRoutes()                                   {}
func (Foo) RegisterTxService()                                   {}
func (Foo) RegisterTendermintService()                           {}
func (Foo) InterfaceRegistry() codectypes.InterfaceRegistry      { return nil }
func (Foo) TxConfig() client.TxConfig                            { return nil }
func (Foo) AppCodec() codec.Codec                                { return app.appCodec }
func (Foo) AutoCliOpts() autocli.AppOptions                      { return autocli.AppOptions{} }
func (Foo) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (Foo) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (Foo) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (Foo) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (Foo) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (Foo) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
