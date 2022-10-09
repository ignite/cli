package accountview

import (
	"fmt"
	"strings"
)

type Option func(*Account)

type Account struct {
	Name     string
	Address  string
	Mnemonic string
}

func WithMnemonic(menmonic string) Option {
	return func(a *Account) {
		a.Mnemonic = menmonic
	}
}

func NewAccount(name, address string, options ...Option) Account {
	a := Account{
		Name:    name,
		Address: address,
	}

	for _, apply := range options {
		apply(&a)
	}

	return a
}

func (a Account) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("🙂 Added account %s with address %s", a.Name, a.Address))

	if a.Mnemonic != "" {
		b.WriteString(fmt.Sprintf(" and mnemonic: %s", a.Mnemonic))
	}

	return b.String()
}

type Accounts []Account

func (a Accounts) String() string {
	b := strings.Builder{}

	for i, account := range a {
		if i > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(account.String())
	}

	return b.String()
}
