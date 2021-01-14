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

	// DefaultAmount specifies the default amount to transfer to an account
	// on each request.
	DefaultAmount = 10000000

	// DefaultMaxAmount specifies the maximum amount that can be tranffered to an
	// account in all times.
	DefaultMaxAmount = 100000000
)

type Faucet struct {
	// runner used to intereact with blockchain's binary to transfer tokens.
	runner chaincmdrunner.Runner

	// accountName to transfer tokens from.
	accountName string

	// accountMnemonic is the mnemonic of the account.
	accountMnemonic string

	// coins keeps a list of coins that can be distributed by the faucet.
	coins []coin
}

type coin struct {
	//amount is the amount of the coin can be distributed per request.
	amount uint64

	// maxAmount is the maximum amount of the coin that can be sent to a single account.
	maxAmount uint64

	// denom is denomination of the coin to be distributed by the faucet.
	denom string
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

// Coin adds a new coin to coins list to distribute by the faucet.
// the first coin added to the list considered as the default coin during transfer requests.
//
// amount is the amount of the coin can be distributed per request.
// maxAmount is the maximum amount of the coin that can be sent to a single account.
// denom is denomination of the coin to be distributed by the faucet.
func Coin(amount, maxAmount uint64, denom string) Option {
	return func(f *Faucet) {
		f.coins = append(f.coins, coin{amount, maxAmount, denom})
	}
}

// New creates a new faucet with ccr (to access and use blockchain's CLI) and given options.
func New(ctx context.Context, ccr chaincmdrunner.Runner, options ...Option) (Faucet, error) {
	f := Faucet{
		runner:      ccr,
		accountName: DefaultAccountName,
	}

	for _, apply := range options {
		apply(&f)
	}

	if len(f.coins) == 0 {
		f.coins = append(f.coins, coin{DefaultAmount, DefaultMaxAmount, DefaultDenom})
	}

	// import the account if mnemonic is provided.
	if f.accountMnemonic != "" {
		_, err := f.runner.AddAccount(ctx, f.accountName, f.accountMnemonic)
		if err != nil && err != chaincmdrunner.ErrAccountAlreadyExists {
			return Faucet{}, err
		}
	}

	return f, nil
}
