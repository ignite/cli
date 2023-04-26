package app

import (
	"cosmossdk.io/client/v2/autocli"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
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

func (Foo) Name() string                                    { return "foo" }
func (Foo) InterfaceRegistry() codectypes.InterfaceRegistry { return nil }
func (Foo) TxConfig() client.TxConfig                       { return nil }
func (Foo) AutoCliOpts() autocli.AppOptions                 { return autocli.AppOptions{} }
