package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
)

var (
	// ErrAccountAlreadyExists returned when an already exists account attempted to be imported.
	ErrAccountAlreadyExists = errors.New("account already exists")

	// ErrAccountDoesNotExist returned when account does not exit.
	ErrAccountDoesNotExist = errors.New("account does not exit")
)

// Account represents a user account.
type Account struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

// AddAccount creates a new account or imports an account when mnemonic is provided.
// returns with an error if the operation went unsuccessful or an account with the provided name
// already exists.
func (r Runner) AddAccount(ctx context.Context, name, mnemonic, coinType string) (Account, error) {
	if err := r.CheckAccountExist(ctx, name); err != nil {
		return Account{}, err
	}
	b := newBuffer()

	account := Account{
		Name:     name,
		Mnemonic: mnemonic,
	}

	// import the account when mnemonic is provided, otherwise create a new one.
	if mnemonic != "" {
		input := &bytes.Buffer{}
		fmt.Fprintln(input, mnemonic)

		if r.chainCmd.KeyringPassword() != "" {
			fmt.Fprintln(input, r.chainCmd.KeyringPassword())
			fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		}

		if err := r.run(
			ctx,
			runOptions{},
			r.chainCmd.RecoverKeyCommand(name, coinType),
			step.Write(input.Bytes()),
		); err != nil {
			return Account{}, err
		}
	} else {
		if err := r.run(ctx, runOptions{
			stdout: b,
			stderr: b,
			stdin:  os.Stdin,
		}, r.chainCmd.AddKeyCommand(name, coinType)); err != nil {
			return Account{}, err
		}

		data, err := b.JSONEnsuredBytes()
		if err != nil {
			return Account{}, err
		}
		if err := json.Unmarshal(data, &account); err != nil {
			return Account{}, err
		}
	}

	// get the address of the account.
	retrieved, err := r.ShowAccount(ctx, name)
	if err != nil {
		return Account{}, err
	}
	account.Address = retrieved.Address

	return account, nil
}

// ImportAccount import an account from a key file
func (r Runner) ImportAccount(ctx context.Context, name, keyFile, passphrase string) (Account, error) {
	if err := r.CheckAccountExist(ctx, name); err != nil {
		return Account{}, err
	}

	// write the passphrase as input
	// TODO: manage keyring backend other than test
	input := &bytes.Buffer{}
	fmt.Fprintln(input, passphrase)

	if err := r.run(
		ctx,
		runOptions{},
		r.chainCmd.ImportKeyCommand(name, keyFile),
		step.Write(input.Bytes()),
	); err != nil {
		return Account{}, err
	}

	return r.ShowAccount(ctx, name)
}

// CheckAccountExist returns an error if the account already exists in the chain keyring
func (r Runner) CheckAccountExist(ctx context.Context, name string) error {
	b := newBuffer()

	// get and decodes all accounts of the chains
	var accounts []Account
	if err := r.run(ctx, runOptions{stdout: b}, r.chainCmd.ListKeysCommand()); err != nil {
		return err
	}

	data, err := b.JSONEnsuredBytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &accounts); err != nil {
		return err
	}

	// search for the account name
	for _, account := range accounts {
		if account.Name == name {
			return ErrAccountAlreadyExists
		}
	}
	return nil
}

// ShowAccount shows details of an account.
func (r Runner) ShowAccount(ctx context.Context, name string) (Account, error) {
	b := &bytes.Buffer{}

	opt := []step.Option{
		r.chainCmd.ShowKeyAddressCommand(name),
	}

	if r.chainCmd.KeyringPassword() != "" {
		input := &bytes.Buffer{}
		fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		opt = append(opt, step.Write(input.Bytes()))
	}

	if err := r.run(ctx, runOptions{stdout: b}, opt...); err != nil {
		if strings.Contains(err.Error(), "item could not be found") ||
			strings.Contains(err.Error(), "not a valid name or address") {
			return Account{}, ErrAccountDoesNotExist
		}
		return Account{}, err
	}

	return Account{
		Name:    name,
		Address: strings.TrimSpace(b.String()),
	}, nil
}

// AddGenesisAccount adds account to genesis by its address.
func (r Runner) AddGenesisAccount(ctx context.Context, address, coins string) error {
	return r.run(ctx, runOptions{}, r.chainCmd.AddGenesisAccountCommand(address, coins))
}

// AddVestingAccount adds vesting account to genesis by its address.
func (r Runner) AddVestingAccount(
	ctx context.Context,
	address,
	originalCoins,
	vestingCoins string,
	vestingEndTime int64,
) error {
	return r.run(ctx, runOptions{}, r.chainCmd.AddVestingAccountCommand(address, originalCoins, vestingCoins, vestingEndTime))
}
