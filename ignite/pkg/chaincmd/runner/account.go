package chaincmdrunner

import (
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

const msgEmptyKeyring = "No records were found in keyring"

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
		input := newBuffer()
		_, err := fmt.Fprintln(input, mnemonic)
		if err != nil {
			return Account{}, err
		}

		if r.chainCmd.KeyringPassword() != "" {
			_, err = fmt.Fprintln(input, r.chainCmd.KeyringPassword())
			if err != nil {
				return Account{}, err
			}

			_, err = fmt.Fprintln(input, r.chainCmd.KeyringPassword())
			if err != nil {
				return Account{}, err
			}

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

// ImportAccount import an account from a key file.
func (r Runner) ImportAccount(ctx context.Context, name, keyFile, passphrase string) (Account, error) {
	if err := r.CheckAccountExist(ctx, name); err != nil {
		return Account{}, err
	}

	// write the passphrase as input
	// TODO: manage keyring backend other than test
	input := newBuffer()
	_, err := fmt.Fprintln(input, passphrase)
	if err != nil {
		return Account{}, err
	}

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

// ListAccounts returns the list of accounts in the keyring.
func (r Runner) ListAccounts(ctx context.Context) ([]Account, error) {
	// Get a JSON string with all accounts in the keyring
	b := newBuffer()
	if err := r.run(ctx, runOptions{stdout: b}, r.chainCmd.ListKeysCommand()); err != nil {
		return nil, err
	}

	// Make sure that the command output is not the empty keyring message.
	// This need to be checked because when the keyring is empty the command
	// finishes with exit code 0 and a plain text message.
	// This behavior was added to Cosmos SDK v0.46.2. See the link
	// https://github.com/cosmos/cosmos-sdk/blob/d01aa5b4a8/client/keys/list.go#L37
	if strings.TrimSpace(b.String()) == msgEmptyKeyring {
		return nil, nil
	}

	data, err := b.JSONEnsuredBytes()
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// CheckAccountExist returns an error if the account already exists in the chain keyring.
func (r Runner) CheckAccountExist(ctx context.Context, name string) error {
	accounts, err := r.ListAccounts(ctx)
	if err != nil {
		return err
	}

	// Search for the account name
	for _, account := range accounts {
		if account.Name == name {
			return ErrAccountAlreadyExists
		}
	}

	return nil
}

// ShowAccount shows details of an account.
func (r Runner) ShowAccount(ctx context.Context, name string) (Account, error) {
	b := newBuffer()
	opt := []step.Option{
		r.chainCmd.ShowKeyAddressCommand(name),
	}

	if r.chainCmd.KeyringPassword() != "" {
		input := newBuffer()
		_, err := fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		if err != nil {
			return Account{}, err
		}
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
