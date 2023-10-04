package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/gogo/protobuf/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	fookeeper "github.com/username/test/x/foo/keeper"
)

type Foo struct {
	AuthKeeper    authkeeper.Keeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper stakingkeeper.Keeper
	GovKeeper     govkeeper.Keeper
	FooKeeper     fookeeper.Keeper
}

func (Foo) Name() string {
	return "foo"
}

func (Foo) GetKey(storeKey string) *storetypes.KVStoreKey { return nil }

func (Foo) TxConfig() client.TxConfig { return nil }

func (Foo) BeginBlocker(sdk.Context, abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (Foo) EndBlocker(sdk.Context, abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (Foo) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	_ = apiSvr.ClientCtx
}

func (Foo) AppCodec() codec.Codec {
	return app.appCodec
}
