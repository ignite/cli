package accountview

import (
	"fmt"
	"strings"
)

type Option func(account *Account)

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

func (acc Account) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Account added: %s \n", acc.Name))
	b.WriteString(fmt.Sprintf("Address: %s \n", acc.Address))
	if acc.Mnemonic != "" {
		b.WriteString(fmt.Sprintf("Mnemonic: %s \n", breakMnemonicIntoLines(acc.Mnemonic, 8)))
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
