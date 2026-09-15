package testing

import (
	"chain"
	"chain/runtime"
	"time"
)

func getContext() (
	originCaller string,
	origSendDenoms []string, origSendAmounts []int64,
	origSpendDenoms []string, origSpendAmounts []int64,
	chainID string,
	height int64,
	timeUnix int64, timeNano int64,
)

func setContext(
	originCaller string,
	currRealmAddr string, currRealmPkgPath string,
	origSendDenoms []string, origSendAmounts []int64,
	origSpendDenoms []string, origSpendAmounts []int64,
	chainID string,
	height int64,
	timeUnix int64, timeNano int64,
)

func GetContext() Context {
	originCaller, origSendDenoms, origSendAmounts, origSpendDenoms, origSpendAmounts, chainID, height, timeUnix, timeNano := getContext()
	return Context{
		OriginCaller: address(originCaller),
		OriginSend:   compactNative(origSendDenoms, origSendAmounts),
		OriginSpend:  compactNative(origSpendDenoms, origSpendAmounts),
		ChainID:      chainID,
		Height:       height,
		Time:         time.Unix(timeUnix, timeNano),
	}
}

func SetContext(ctx Context) {
	originSendDenom, originSendAmount := expandNative(ctx.OriginSend)
	originSpendDenom, originSpendAmount := expandNative(ctx.OriginSpend)

	var timestamp, timestampNano int64

	if !ctx.Time.IsZero() {
		unixNano := ctx.Time.UnixNano()
		timestamp = unixNano / int64(time.Second)
		timestampNano = unixNano % int64(time.Second)
	}

	setContext(
		ctx.OriginCaller.String(),
		ctx.CurrentRealm.Address().String(), ctx.CurrentRealm.PkgPath(),
		originSendDenom, originSendAmount,
		originSpendDenom, originSpendAmount,
		ctx.ChainID,
		ctx.Height,
		timestamp, timestampNano,
	)
}

func testIssueCoins(addr string, denom []string, amt []int64)

func SetOriginCaller(origCaller address) {
	ctx := GetContext()
	ctx.OriginCaller = origCaller
	ctx.CurrentRealm = NewUserRealm(origCaller)
	SetContext(ctx)
}

func SetOriginSend(send chain.Coins) {
	ctx := GetContext()
	ctx.OriginSend = send
	SetContext(ctx)
}

func SetOriginSpend(spend chain.Coins) {
	ctx := GetContext()
	ctx.OriginSpend = spend
	SetContext(ctx)
}

func SetHeight(height int64) {
	ctx := GetContext()
	ctx.Height = height
	SetContext(ctx)
}

// SetRealm sets the realm for the current frame.
// After calling SetRealm, calling CurrentRealm() in the test function will yield the value of
// rlm, while if a realm function is called, using PreviousRealm() will yield rlm.
// If rlm.PkgPath() is "", implies rlm.Address() is the user address.
// OriginCaller will also be set. This is identical to calling
// ctx.SetOriginCaller(userAddress)
func SetRealm(rlm runtime.Realm) {
	ctx := GetContext()
	ctx.CurrentRealm = rlm
	if rlm.PkgPath() == "" {
		ctx.OriginCaller = rlm.Address()
	}
	SetContext(ctx)
}

func SkipHeights(count int64) {
	ctx := GetContext()
	toHeight := ctx.Height + count
	toTime := ctx.Time.Add(time.Duration(count*5) * time.Second)
	SetContext(Context{Height: toHeight, Time: toTime})
}

func IssueCoins(addr address, coins chain.Coins) {
	denom, amt := expandNative(coins)
	testIssueCoins(addr.String(), denom, amt)
}

// expandNative expands for usage within natively bound functions.
func expandNative(coins chain.Coins) (denoms []string, amounts []int64) {
	denoms = make([]string, len(coins))
	amounts = make([]int64, len(coins))
	for i, coin := range coins {
		denoms[i] = coin.Denom
		amounts[i] = coin.Amount
	}

	return denoms, amounts
}

// compactNative compacts coins for usage within natively bound functions.
func compactNative(denoms []string, amounts []int64) chain.Coins {
	coins := make(chain.Coins, len(denoms))
	for i := range coins {
		coins[i] = chain.Coin{Denom: denoms[i], Amount: amounts[i]}
	}
	return coins
}

// NewUserRealm creates a new "user" realm with the given address.
func NewUserRealm(user address) runtime.Realm {
	return newRealm(string(user), "")
}

// NewCodeRealm creates a new realm for a published code realm with the given
// pkgPath.
func NewCodeRealm(pkgPath string) runtime.Realm {
	if !isRealm(pkgPath) {
		panic("should only be called for Realms")
	}
	return newRealm(string(chain.PackageAddress(pkgPath)), pkgPath)
}

func newRealm(addr, pkgPath string) runtime.Realm // injected
func isRealm(pkgPath string) bool

// MakeRealm constructs a uverse `realm` value with explicit addr/pkgPath/
// prev. Use when SetRealm composition semantics don't match the scenario:
//
//	aliceRealm := testing.MakeRealm(alice, "", testing.OriginRealm())
//	o := New(0, aliceRealm)
//	crossed := testing.MakeRealm(crypto.AddressFromPkgPath(p), p, aliceRealm)
//	o.OwnedByPrevious(0, crossed) // observes Previous() == aliceRealm
//
// Unlike SetRealm, this does not mutate any frame's captured cur.
func MakeRealm(addr address, pkgPath string, prev realm) realm {
	return makeRealm(string(addr), pkgPath, prev)
}

// OriginRealm returns the EOA-origin realm shape (addr=OriginCaller,
// pkgPath="", prev=truly-nil). Use as the seed prev when assembling
// explicit cur chains in tests.
func OriginRealm() realm {
	return originRealm()
}

func makeRealm(addr, pkgPath string, prev realm) realm // injected
func originRealm() realm                               // injected
