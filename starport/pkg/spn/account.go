package spn

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

// Account represents an account on SPN.
type Account struct {
	Name     string
	Address  string
	Mnemonic string
}

// Keyring uses given keyring type as storage.
func Keyring(keyring string) Option {
	return func(c *options) {
		c.keyringBackend = keyring
	}
}

// AccountGet retrieves an account by name from the keyring.
func (c *Client) AccountGet(accountName string) (Account, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return Account{}, err
	}
	return toAccount(info), nil
}

// AccountList returns a list of accounts.
func (c *Client) AccountList() ([]Account, error) {
	var accounts []Account
	infos, err := c.kr.List()
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		accounts = append(accounts, toAccount(info))
	}
	return accounts, nil
}

// AccountCreate creates an account by name and mnemonic (optional) in the keyring.
func (c *Client) AccountCreate(accountName, mnemonic string) (Account, error) {
	if mnemonic == "" {
		entropySeed, err := bip39.NewEntropy(256)
		if err != nil {
			return Account{}, err
		}
		mnemonic, err = bip39.NewMnemonic(entropySeed)
		if err != nil {
			return Account{}, err
		}
	}
	algos, _ := c.kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return Account{}, err
	}
	hdPath := hd.CreateHDPath(types.GetConfig().GetCoinType(), 0, 0).String()
	info, err := c.kr.NewAccount(accountName, mnemonic, "", hdPath, algo)
	if err != nil {
		return Account{}, err
	}
	account := toAccount(info)
	account.Mnemonic = mnemonic
	return account, nil
}

func toAccount(info keyring.Info) Account {
	ko, _ := keyring.Bech32KeyOutput(info)
	return Account{
		Name:    ko.Name,
		Address: ko.Address,
	}
}

// AccountExport exports an account in the keyring by name and an encryption password into privateKey.
// password later can be used to decrypt the privateKey.
func (c *Client) AccountExport(accountName, password string) (privateKey string, err error) {
	return c.kr.ExportPrivKeyArmor(accountName, password)
}

// AccountImport imports an account to the keyring by account name, privateKey and decryption password.
func (c *Client) AccountImport(accountName, privateKey, password string) error {
	return c.kr.ImportPrivKey(accountName, privateKey, password)
}

// buildClientCtx builds the context for the client
func (c *Client) buildClientCtx(accountName string) (client.Context, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return client.Context{}, err
	}
	return c.clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress()), nil
}
