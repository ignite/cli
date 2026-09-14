package runtime

import (
	"chain"
	"strings"
)

type Realm struct {
	addr    address
	pkgPath string
}

// NewRealm constructs a Realm from an address and pkgPath. Exposed so
// the chain/runtime/unsafe package (which lives outside chain/runtime)
// can return Realm values from stack-walking helpers; Realm's fields
// are unexported and must not be set outside this package's discipline.
//
// Callers should NOT use this to fabricate realm-context for
// authorization — a constructed Realm is just a (addr, pkgPath) tuple
// and carries no runtime-validated authority. Use cur.Previous() for
// caller-identity checks.
func NewRealm(addr address, pkgPath string) Realm {
	return Realm{addr, pkgPath}
}

func (r Realm) String() string {
	if r.pkgPath == "" {
		return "UserRealm{ " + r.addr.String() + " }"
	} else {
		return "CodeRealm{ " + r.addr.String() + ", " + r.pkgPath + " }"
	}
}

func (r Realm) Address() address {
	return r.addr
}

func (r Realm) PkgPath() string {
	return r.pkgPath
}

// IsEphemeral returns true if r.pkgPath is an ephemeral path (/e/).
func (r Realm) IsEphemeral() bool {
	if r.pkgPath == "" {
		return false
	}

	// Check if pkgPath matches pattern: domain/e/*
	// We need to find "/e/" and ensure it comes exactly after the domain (no "/" before it)
	// We assume pkgPath is a full path (e.g., "gno.land/e/user/run") to avoid unnecessary computation

	idx := strings.Index(r.pkgPath, "/e/")

	// no /e/ found
	if idx == -1 {
		return false
	}

	// has no content after /e/
	if len(r.pkgPath) <= idx+3 {
		return false
	}

	// has a domain right before /e/
	if strings.Index(r.pkgPath[:idx], "/") != -1 {
		return false
	}
	return true
}

// IsCode returns true if r.pkgPath is not empty.
// - true if Coderealm and !IsEphemeral, an uploaded code realm.
// - true if CodeRealm and IsEphemeral, 'gnokey maketx run'.
// - false if UserRealm, 'gnokey maketx call'.
func (r Realm) IsCode() bool {
	return r.pkgPath != ""
}

// IsUser returns true if r.Address is a user address.
// - true if UserRealm, 'gnokey maketx call',
// - true if CodeRealm and IsEphemeral, 'gnokey maketx run',
// - false if Coderealm and !IsEphemeral, an uploaded code realm.
func (r Realm) IsUser() bool {
	return r.IsUserCall() || r.IsUserRun()
}

// IsUserRun returns true if the realm represents a user run context.
// This checks if the pkgPath matches the pattern: domain/e/{addr}/run
// where {addr} is the realm's address.
func (r Realm) IsUserRun() bool {
	// Compute expected path dynamically from address
	idx := strings.Index(r.pkgPath, "/")
	if idx == -1 {
		return false
	}
	domain := r.pkgPath[:idx]
	expectedPath := domain + "/e/" + r.addr.String() + "/run"
	return r.pkgPath == expectedPath
}

// IsUserCall returns true if this is an EOA (Externally Owned Account) call.
// This is true when pkgPath is empty, which is the opposite of IsCode().
// - true for 'gnokey maketx call' from user
// - false for any code realm execution
func (r Realm) IsUserCall() bool {
	return r.pkgPath == ""
}

func (r Realm) CoinDenom(coinName string) string {
	return chain.CoinDenom(r.pkgPath, coinName)
}
