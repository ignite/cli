package app

import (
	"cosmossdk.io/api/tendermint/abci"
	"cosmossdk.io/client/v2/autocli"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	foomodule "github.com/username/test/x/foo"
)

// App modules are defined as NewBasicManager arguments
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	gov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
	}),
	foomodule.AppModuleBasic{},
)

type Foo struct {
	*runtime.App
}

func (Foo) Name() string                                    { return "foo" }
func (Foo) InterfaceRegistry() codectypes.InterfaceRegistry { return nil }
func (Foo) TxConfig() client.TxConfig                       { return nil }
func (Foo) AutoCliOpts() autocli.AppOptions                 { return autocli.AppOptions{} }

func (Foo) BeginBlocker(sdk.Context, abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (Foo) EndBlocker(sdk.Context, abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (app *Foo) RegisterAPIRoutes(s *api.Server, cfg config.APIConfig) {
	// This module should be discovered
	foomodule.RegisterGRPCGatewayRoutes(s.ClientCtx, s.GRPCGatewayRouter)
	// Runtime app modules for the current Cosmos SDK should be discovered too
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
}

func (Foo) GetKey(storeKey string) *storetypes.KVStoreKey { return nil }

func (Foo) TxConfig() client.TxConfig { return nil }
