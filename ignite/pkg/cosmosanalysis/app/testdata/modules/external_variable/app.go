package app

import (
	"app/modules"

	"cosmossdk.io/client/v2/autocli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var ModuleBasics = module.NewBasicManager(modules.Basic...)

type Foo struct{}

func (Foo) Name() string                                    { return "foo" }
func (Foo) InterfaceRegistry() codectypes.InterfaceRegistry { return nil }
func (Foo) TxConfig() client.TxConfig                       { return nil }
func (Foo) AutoCliOpts() autocli.AppOptions                 { return autocli.AppOptions{} }
