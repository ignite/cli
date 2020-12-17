package networkbuilder

import (
	"github.com/tendermint/starport/starport/pkg/spn"
)

// AccountUse sets the account to be used while working with SPN.
func (b *Builder) AccountUse(name string) error {
	c, err := ConfigGet()
	if err != nil {
		return err
	}
	c.SPNAccount = name
	return ConfigSave(c)
}

// AccountInUse gets the account in use while working with SPN.
func (b *Builder) AccountInUse() (spn.Account, error) {
	c, err := ConfigGet()
	if err != nil {
		return spn.Account{}, nil
	}
	return b.AccountGet(c.SPNAccount)
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
func (b *Builder) AccountCreate(name, mnemonic string) (spn.Account, error) {
	return b.spnclient.AccountCreate(name, mnemonic)
}

// AccountExport exports an account in OS keyring with name and password.
func (b *Builder) AccountExport(name, password string) (privateKey string, err error) {
	return b.spnclient.AccountExport(name, password)
}

// AccountImport imports account to OS keyring with name, password and privateKey.
func (b *Builder) AccountImport(name, privateKey, password string) error {
	return b.spnclient.AccountImport(name, privateKey, password)
}
