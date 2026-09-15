package gno

import "github.com/gnolang/gno/gno.land/pkg/integration"

// defaultDevAccountName is the well-known pre-funded dev chain account.
const defaultDevAccountName = integration.DefaultAccount_Name

// ensureKey returns the key info for name, auto-importing the well-known
// dev account when it is missing (its mnemonic is public).
func ensureKey(name string) (KeyInfo, error) {
	info, err := ShowKey(name)
	if err == nil {
		return info, nil
	}
	if name != defaultDevAccountName {
		return KeyInfo{}, err
	}
	return RecoverKey(defaultDevAccountName, integration.DefaultAccount_Seed, "", 0, 0)
}
