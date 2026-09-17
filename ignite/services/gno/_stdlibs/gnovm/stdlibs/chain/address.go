package chain

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/crypto"
)

func X_packageAddress(pkgPath string) string {
	return string(gno.DerivePkgBech32Addr(pkgPath))
}

func X_deriveStorageDepositAddr(pkgPath string) string {
	return string(gno.DeriveStorageDepositBech32Addr(pkgPath))
}

// X_pubKeyAddress derives the bech32 address from a bech32-encoded
// public key. Returns ("", errStr) on error — the gno wrapper translates
// to a Go-side error or address cast.
//
// Used by r/sys/params.GetValsetEntries to parse "<bech32-pubkey>:<power>"
// entries from valset:current. The function is narrow on purpose: bech32
// decode + amino unmarshal + Address() derivation are not safely doable
// in pure gno (amino-decode is not exposed to gno), so this is the only
// chain-side primitive needed for that use case.
func X_pubKeyAddress(bech32PubKey string) (addr string, errStr string) {
	pk, err := crypto.PubKeyFromBech32(bech32PubKey)
	if err != nil {
		return "", err.Error()
	}
	if pk == nil {
		// Co-regression with C2: PubKeyFromBech32 can return (nil, nil)
		// for empty-payload bech32. Refuse here too rather than nil-deref
		// on pk.Address().
		return "", "nil pubkey from bech32 " + bech32PubKey
	}
	return pk.Address().String(), ""
}
