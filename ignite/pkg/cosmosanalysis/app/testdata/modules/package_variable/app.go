package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var ModuleBasics = module.NewBasicManager(basicModules...)

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
