package foo

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/app/testdata/modules/registration_not_in_app_go"
)

type Foo[T any] struct {
	FooKeeper foo.keeper
	i         T
}

func (f Foo[T]) TxConfig() client.TxConfig                            { return nil }
func (f Foo[T]) RegisterAPIRoutes()                                   {}
func (f Foo[T]) RegisterTxService()                                   {}
func (f Foo[T]) RegisterTendermintService()                           {}
func (f Foo[T]) AppCodec() codec.Codec                                { return app.appCodec }
func (f Foo[T]) Name() string                                         { return app.BaseApp.Name() }
func (f Foo[T]) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (f Foo[T]) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (f Foo[T]) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (f Foo[T]) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (f Foo[T]) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (f Foo[T]) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
