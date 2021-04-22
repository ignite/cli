package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

var (
	// ErrAccountAlreadyExists returned when an already exists account attempted to be imported.
	ErrAccountAlreadyExists = errors.New("account already exists")

	// ErrAccountDoesNotExist returned when account does not exit.
	ErrAccountDoesNotExist = errors.New("account does not exit")
)

// AddAccount creates a new account or imports an account when mnemonic is provided.
// returns with an error if the operation went unsuccessful or an account with the provided name
// already exists.
func (runner Runner) AddAccount(ctx context.Context, name, mnemonic string) (Account, error) {
	b := &bytes.Buffer{}

	// check if account already exists.
	var accounts []Account
	if err := runner.run(ctx, runOptions{stdout: b}, runner.chainCmd.ListKeysCommand()); err != nil {
		return Account{}, err
	}
	if err := json.NewDecoder(b).Decode(&accounts); err != nil {
		return Account{}, err
	}
	for _, account := range accounts {
		if account.Name == name {
			return Account{}, ErrAccountAlreadyExists
		}
	}
	b.Reset()

	account := Account{
		Name:     name,
		Mnemonic: mnemonic,
	}

	// import the account when mnemonic is provided, otherwise create a new one.
	if mnemonic != "" {
		input := &bytes.Buffer{}
		fmt.Fprintln(input, mnemonic)

		if runner.chainCmd.KeyringPassword() != "" {
			fmt.Fprintln(input, runner.chainCmd.KeyringPassword())
			fmt.Fprintln(input, runner.chainCmd.KeyringPassword())
		}

		if err := runner.run(
			ctx,
			runOptions{},
			runner.chainCmd.ImportKeyCommand(name),
			step.Write(input.Bytes()),
		); err != nil {
			return Account{}, err
		}
	} else {
		if err := runner.run(ctx, runOptions{
			stdout: io.MultiWriter(b, os.Stdout),
			stderr: os.Stderr,
			stdin:  os.Stdin,
		}, runner.chainCmd.AddKeyCommand(name)); err != nil {
			return Account{}, err
		}
		if err := json.NewDecoder(b).Decode(&account); err != nil {
			return Account{}, err
		}

		b.Reset()
	}

	// get full details of the account.
	opt := []step.Option{
		runner.chainCmd.ShowKeyAddressCommand(name),
	}

	if runner.chainCmd.KeyringPassword() != "" {
		input := &bytes.Buffer{}
		fmt.Fprintln(input, runner.chainCmd.KeyringPassword())
		opt = append(opt, step.Write(input.Bytes())) // TODO: stdin if not present
	}

	if err := runner.run(ctx, runOptions{
		stdout: io.MultiWriter(b, os.Stdout),
		stderr: os.Stderr,
		stdin:  os.Stdin,
	}, opt...); err != nil {
		return Account{}, err
	}
	account.Address = strings.TrimSpace(b.String())

	return account, nil
}

// Account represents a user account.
type Account struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

// ShowAccount shows details of an account.
func (runner Runner) ShowAccount(ctx context.Context, name string) (Account, error) {
	b := &bytes.Buffer{}

	opt := []step.Option{
		runner.chainCmd.ShowKeyAddressCommand(name),
	}

	if runner.chainCmd.KeyringPassword() != "" {
		input := &bytes.Buffer{}
		fmt.Fprintln(input, runner.chainCmd.KeyringPassword())
		opt = append(opt, step.Write(input.Bytes()))
	}

	if err := runner.run(ctx, runOptions{stdout: b}, opt...); err != nil {
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
func (runner Runner) AddGenesisAccount(ctx context.Context, address, coins string) error {
	return runner.run(ctx, runOptions{}, runner.chainCmd.AddGenesisAccountCommand(address, coins))
}
