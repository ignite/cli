// Package cosmosfaucet is a faucet to request tokens for sdk accounts.
package cosmosfaucet

import (
	"context"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
)

const (
	// DefaultAccountName is the default account to transfer tokens from.
	DefaultAccountName = "faucet"

	// DefaultDenom is the default denomination to distribute.
	DefaultDenom = "uatom"

	// DefaultCreditAmount specifies the default amount to transfer to an account
	// on each request.
	DefaultCreditAmount = 10000000

	// DefaultMaxCredit specifies the maximum amount that can be tranffered to an
	// account in all times.
	DefaultMaxCredit = 100000000
)

type Faucet struct {
	// runner used to intereact with blockchain's binary to transfer tokens.
	runner chaincmdrunner.Runner

	// accountName to transfer tokens from.
	accountName string

	// accountMnemonic is the mnemonic of the account.
	accountMnemonic string

	// denom is denomination of the coin to be distributed by the faucet.
	denom string

	// creditAmount is the amount to credit in each request.
	creditAmount uint64

	// maxCredit is the maximum credit per account.
	maxCredit uint64
}

// Option configures the faucetOptions.
type Option func(*Faucet)

// Account provides the account information to transfer tokens from.
// when mnemonic isn't provided, account assumed to be exists in the keyring.
func Account(name, mnemonic string) Option {
	return func(f *Faucet) {
		f.accountName = name
		f.accountMnemonic = mnemonic
	}
}

// Denom set the denomination of the coin to be distributed by the faucet.
func Denom(denom string) Option {
	return func(f *Faucet) {
		f.denom = denom
	}
}

// CreditAmount sets the amount to credit in each request.
func CreditAmount(credit uint64) Option {
	return func(f *Faucet) {
		f.creditAmount = credit
	}
}

// MaxCredit sets the maximum credit per account.
func MaxCredit(credit uint64) Option {
	return func(f *Faucet) {
		f.maxCredit = credit
	}
}

// New creates a new faucet with ccr (to access and use blockchain's CLI) and given options.
func New(ctx context.Context, ccr chaincmdrunner.Runner, options ...Option) (Faucet, error) {
	f := Faucet{
		runner:       ccr,
		accountName:  DefaultAccountName,
		denom:        DefaultDenom,
		creditAmount: DefaultCreditAmount,
		maxCredit:    DefaultMaxCredit,
	}

	for _, apply := range options {
		apply(&f)
	}

	if f.accountMnemonic != "" {
		_, err := f.runner.AddAccount(ctx, f.accountName, f.accountMnemonic)
		if err != nil && err != chaincmdrunner.ErrAccountAlreadyExists {
			return Faucet{}, err
		}
	}

	return f, nil
}
