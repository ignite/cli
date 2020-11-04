package networkbuilder

import (
	"io/ioutil"
	"os"

	"github.com/tendermint/starport/starport/pkg/spn"
)

// AccountUse sets the account to be used while working with SPN.
func (b *Builder) AccountUse(name string) error {
	if err := os.MkdirAll(starportConfDir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(confPath, []byte(name), 0755)
}

// AccountInUse gets the account in use while working with SPN.
// if there is no account in use, it creates "spn" account if not exists
// and puts into use.
func (b *Builder) AccountInUse() (spn.Account, error) {
	var name string
	nameb, err := ioutil.ReadFile(confPath)
	if err != nil {
		name = "spn"
		b.AccountCreate(name)
		if err := b.AccountUse(name); err != nil {
			return spn.Account{}, err
		}
	} else {
		name = string(nameb)
	}
	return b.AccountGet(name)
}

// AccountList lists all accounts in OS keyring.
func (b *Builder) AccountList() ([]spn.Account, error) {
	return b.spnclient.AccountList()
}

// AccountGet gets an account by name in OS keyring.
func (b *Builder) AccountGet(name string) (spn.Account, error) {
	return b.spnclient.AccountGet(name)
}

// AccountCreate creates a new account in OS keyring.
func (b *Builder) AccountCreate(name string) (spn.Account, error) {
	return b.spnclient.AccountCreate(name)
}

// AccountExport exports an account in OS keyring with name and password.
func (b *Builder) AccountExport(name, password string) (privateKey string, err error) {
	return b.spnclient.AccountExport(name, password)
}

// AccountImport imports account to OS keyring with name, password and privateKey.
func (b *Builder) AccountImport(name, privateKey, password string) error {
	return b.spnclient.AccountImport(name, privateKey, password)
}
