package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

// AddAccount creates a new account or imports an account when mnemonic is provided.
// returns with an error if the operation went unsuccessful or an account with the provided name
// already exists.
func (r Runner) AddAccount(ctx context.Context, name, mnemonic string) (Account, error) {
	b := &bytes.Buffer{}

	// check if account already exists.
	var accounts []Account
	if err := r.run(ctx, runOptions{stdout: b}, r.cc.ListKeysCommand()); err != nil {
		return Account{}, err
	}
	if err := json.NewDecoder(b).Decode(&accounts); err != nil {
		return Account{}, err
	}
	for _, account := range accounts {
		if account.Name == name {
			return Account{}, errors.New("account already exists")
		}
	}
	b.Reset()

	account := Account{
		Name:     name,
		Mnemonic: mnemonic,
	}

	// import the account when mnemonic is provided, otherwise create a new one.
	if mnemonic != "" {
		if err := r.run(
			ctx,
			runOptions{},
			r.cc.ImportKeyCommand(name),
			step.Write([]byte(mnemonic+"\n")),
		); err != nil {
			return Account{}, err
		}
	} else {
		// note that, launchpad prints account output from stderr.
		if err := r.run(ctx, runOptions{stdout: b, stderr: b}, r.cc.AddKeyCommand(name)); err != nil {
			return Account{}, err
		}
		if err := json.NewDecoder(b).Decode(&account); err != nil {
			return Account{}, err
		}

		b.Reset()
	}

	// get full details of the account.
	if err := r.run(ctx, runOptions{stdout: b}, r.cc.ShowKeyAddressCommand(name)); err != nil {
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
func (r Runner) ShowAccount(ctx context.Context, name string) (Account, error) {
	b := &bytes.Buffer{}
	if err := r.run(ctx, runOptions{stdout: b}, r.cc.ShowKeyAddressCommand(name)); err != nil {
		return Account{}, err
	}
	return Account{
		Name:    name,
		Address: strings.TrimSpace(b.String()),
	}, nil
}

// AddGenesisAccount adds account to genesis by its address.
func (r Runner) AddGenesisAccount(ctx context.Context, address, coins string) error {
	return r.run(ctx, runOptions{}, r.cc.AddGenesisAccountCommand(address, coins))
}
