package chain

import (
	"errors"
	"strings"
)

func packageAddress(string) string

// subRealmSep joins a host pkgpath and a subpath in a sub-realm token's
// synthesized pkgpath ("host#subpath"). Must match gnolang.subRealmSep.
// '#' is reserved in package paths and kept distinct from gnoweb's ':'.
const subRealmSep = "#"

// Returns a crypto hash derived pkgPath, unless pkgPath is a MsgRun run path,
// in which case the address is extracted from the path.
func PackageAddress(pkgPath string) address {
	addr := packageAddress(pkgPath)
	return address(addr)
}

// DerivePkgSubAddr derives the address of a sub-realm identity minted
// by cur.Sub(subpath) on the host realm at pkgpath. Pure computation;
// callable from anywhere — lets other realms and off-chain tooling
// compute a sub-address without consulting the host. Enforces the same
// subpath rules as Sub(); the host-"#"-free assert also subsumes the
// run-path derivation assert (the run-path shape is "#"-free).
func DerivePkgSubAddr(pkgpath, subpath string) address {
	assertValidSubpath(pkgpath, subpath)
	return PackageAddress(pkgpath + subRealmSep + subpath)
}

// SplitPkgSubPath splits a possibly-synthesized pkgpath into host and
// subpath. ok reports whether p is a synthesized sub-realm pkgpath
// ("#" is reserved in real package paths and appears only via
// realm.Sub synthesis).
func SplitPkgSubPath(p string) (host, subpath string, ok bool) {
	return strings.Cut(p, subRealmSep)
}

// assertValidSubpath mirrors gnolang.subRealmPathError so
// DerivePkgSubAddr can never derive an address that cur.Sub would
// refuse to mint (which would strand funds). Keep in sync with
// subRealmPathError / isValidSubpath in gnovm/pkg/gnolang/uverse.go.
func assertValidSubpath(pkgpath, subpath string) {
	switch {
	case subpath == "":
		panic("Sub: subpath cannot be empty")
	case len(pkgpath)+1+len(subpath) > 256:
		panic("Sub: synthesized pkgpath too long (max 256 bytes)")
	case strings.Contains(pkgpath, subRealmSep):
		panic("Sub: receiver pkgpath is already synthesized or invalid")
	case !isValidSubpath(subpath):
		panic("Sub: subpath must be '/'-separated segments of [a-z0-9], with '_.-' allowed inside a segment")
	}
}

// isValidSubpath is segment ("/" segment)*, segment =
// [a-z0-9] ([a-z0-9_.-]* [a-z0-9])?. Byte-for-byte transliteration of
// gnolang.isValidSubpath.
func isValidSubpath(s string) bool {
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			if !isValidSubpathSegment(s[start:i]) {
				return false
			}
			start = i + 1
		}
	}
	return true
}

func isValidSubpathSegment(seg string) bool {
	if seg == "" {
		return false
	}
	for i := 0; i < len(seg); i++ {
		c := seg[i]
		alnum := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if i == 0 || i == len(seg)-1 {
			if !alnum {
				return false
			}
			continue
		}
		if !alnum && c != '_' && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

func deriveStorageDepositAddr(string) string

// Returns a crypto hash derived pkgPath, unless pkgPath is a MsgRun run path,
// in which case the address is extracted from the path.
func DeriveStorageDepositAddr(pkgPath string) address {
	addr := deriveStorageDepositAddr(pkgPath)
	return address(addr)
}

func pubKeyAddress(string) (addr string, errStr string)

// PubKeyAddress derives the bech32 address from a bech32-encoded public key.
// Returns an error if the input is malformed or amino-decodes to a nil
// concrete pubkey (e.g., empty-payload "gpub1mdgqmw").
//
// This is the inverse half of the on-chain valset entry format
// "<bech32-pubkey>:<power>": split the entry, then call this helper to
// recover the address from the pubkey component.
func PubKeyAddress(bech32PubKey string) (address, error) {
	addr, errStr := pubKeyAddress(bech32PubKey)
	if errStr != "" {
		return "", errors.New(errStr)
	}
	return address(addr), nil
}
