package account

import (
	"fmt"
	"strings"
)

type Option func(account *Account)

type Account struct {
	name     string
	address  string
	mnemonic string
}

func WithMnemonic(menmonic string) Option {
	return func(a *Account) {
		a.mnemonic = menmonic
	}
}

func NewAccount(name, address string, options ...Option) Account {
	a := Account{
		name:    name,
		address: address,
	}

	for _, apply := range options {
		apply(&a)
	}
	return a
}

func (acc Account) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Account added: %s \n", acc.name))
	b.WriteString(fmt.Sprintf("Address: %s \n", acc.address))
	if acc.mnemonic != "" {
		b.WriteString(fmt.Sprintf("Mnemonic: %s \n", breakMnemonicIntoLines(acc.mnemonic, 8)))
	}
	return b.String()
}

type Accounts []Account

func Collection(accounts ...Account) Accounts {
	return append(Accounts{}, accounts...)
}

func (accs Accounts) String() string {
	b := strings.Builder{}
	b.WriteString("\n")
	for i, acc := range accs {
		b.WriteString(acc.String())
		if len(accs)-i != 1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func breakMnemonicIntoLines(mnemonic string, breakAfterN int) string {
	splitted := strings.Split(mnemonic, " ")
	b := strings.Builder{}
	for i, s := range splitted {
		b.WriteString(s)
		b.WriteString(" ")
		if (i+1)%breakAfterN == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
