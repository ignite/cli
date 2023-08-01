package accountview

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

var (
	fmtExistingAccount = "%s %s's account address: %s\n"
	fmtNewAccount      = "%s Added account %s with address %s and mnemonic:\n%s\n"
)

type Option func(*Account)

type Account struct {
	Name     string
	Address  string
	Mnemonic string
}

func WithMnemonic(mnemonic string) Option {
	return func(a *Account) {
		a.Mnemonic = mnemonic
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
	name := colors.Name(a.Name)

	// The account is new when the mnemonic is available
	if a.Mnemonic != "" {
		return fmt.Sprintf(fmtNewAccount, icons.OK, name, a.Address, colors.Mnemonic(a.Mnemonic))
	}

	return fmt.Sprintf(fmtExistingAccount, icons.User, name, a.Address)
}

type Accounts []Account

func (a Accounts) String() string {
	b := strings.Builder{}

	for i, acc := range a {
		// Make sure accounts are separated by an
		// empty line when the mnemonic is available.
		if i > 0 && acc.Mnemonic != "" {
			b.WriteRune('\n')
		}

		b.WriteString(acc.String())
	}

	b.WriteRune('\n')

	return b.String()
}

func (a Accounts) Append(acc Account) Accounts {
	return append(a, acc)
}
