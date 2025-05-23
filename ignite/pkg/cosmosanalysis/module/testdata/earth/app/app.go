package app

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	marskeeper "github.com/tendermint/planet/x/mars/keeper"
)

type Foo struct {
	*runtime.App

	AuthKeeper    authkeeper.Keeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper stakingkeeper.Keeper
	GovKeeper     govkeeper.Keeper
	MarsKeeper    marskeeper.Keeper
}

var ModuleBasics = module.NewBasicManager(foo.AppModuleBasic{})
