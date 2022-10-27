package app

import (
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
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

type Foo struct{}

func (Foo) Name() string {
	return "foo"
}

func (Foo) BeginBlocker(sdk.Context, abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (Foo) EndBlocker(sdk.Context, abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

type App struct{}

func (App) RegisterAPIRoutes(s *api.Server, cfg config.APIConfig) {
	// These two modules should be discovered too
	authtx.RegisterGRPCGatewayRoutes(s.ClientCtx, s.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(s.ClientCtx, s.GRPCGatewayRouter)
}
