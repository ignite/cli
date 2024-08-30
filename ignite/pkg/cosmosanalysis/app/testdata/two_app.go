package foo

import (
	paramstypes "cosmossdk.io/x/params/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/app/testdata/modules/registration_not_in_app_go"
)

type Foo struct {
	FooKeeper foo.keeper
}

func (f Foo) TxConfig() client.TxConfig                            { return nil }
func (f Foo) RegisterAPIRoutes()                                   {}
func (f Foo) RegisterTxService()                                   {}
func (f Foo) RegisterTendermintService()                           {}
func (f Foo) Name() string                                         { return app.BaseApp.Name() }
func (f Foo) AppCodec() codec.Codec                                { return app.appCodec }
func (F Foo) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (F Foo) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (F Foo) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (F Foo) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (f Foo) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (f Foo) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

type Bar struct {
	FooKeeper foo.keeper
}

func (f Bar) TxConfig() client.TxConfig                            { return nil }
func (f Bar) RegisterAPIRoutes()                                   {}
func (f Bar) RegisterTxService()                                   {}
func (f Bar) RegisterTendermintService()                           {}
func (f Bar) Name() string                                         { return app.BaseApp.Name() }
func (f Bar) AppCodec() codec.Codec                                { return app.appCodec }
func (f Bar) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (f Bar) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (f Bar) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (f Bar) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (f Bar) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (f Bar) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
