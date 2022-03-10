package app

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	planet "github.com/tendermint/planet/x/planet"
)

type Foo struct {
	FooKeeper foo.keeper
}

var ModuleBasics = module.NewBasicManager(planet.AppModuleBasic{})

func (f Foo) RegisterAPIRoutes()         {}
func (f Foo) RegisterTxService()         {}
func (f Foo) RegisterTendermintService() {}
