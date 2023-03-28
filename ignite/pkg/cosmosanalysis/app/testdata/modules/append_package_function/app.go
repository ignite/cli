package app

import (
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
)

var ModuleBasics = module.NewBasicManager(
	append(
		[]module.AppModuleBasic{
			auth.AppModuleBasic{},
			bank.AppModuleBasic{},
			staking.AppModuleBasic{},
			gov.NewAppModuleBasic([]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
			}),
		},
		basicModules()...,
	),
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

func (Foo) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	_ = apiSvr.ClientCtx
}
