package testing

import (
	"chain"
	"chain/runtime"
	"time"
)

type Context struct {
	CurrentRealm runtime.Realm
	OriginCaller address
	OriginSend   chain.Coins
	OriginSpend  chain.Coins
	ChainID      string
	Height       int64
	Time         time.Time
}
