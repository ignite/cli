package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	marskeeper "github.com/ignite/planet/x/mars/keeper"
)

type Foo struct {
	baseapp.BaseApp

	MarsKeeper marskeeper.Keeper
}
