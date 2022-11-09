package accountview

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
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
	msg := fmt.Sprintf("%s Added account %s with address %s", icons.OK, colors.Name(a.Name), a.Address)

	b.WriteString(msg)

	if a.Mnemonic != "" {
		s := wordwrap.String(a.Mnemonic, 80)
		s = indent.String(s, 2)

		b.WriteString(fmt.Sprintf(" and mnemonic:\n%s\n", colors.Mnemonic(s)))
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

	b.WriteRune('\n')

	return b.String()
}
