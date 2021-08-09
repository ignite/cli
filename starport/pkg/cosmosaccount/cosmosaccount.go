package cosmosaccount

import (
	"errors"
	"os"

	dkeyring "github.com/99designs/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

// KeyringServiceName used for the name of keyring in OS backend.
const KeyringServiceName = "starport"

// KeyringHome used to store account related data.
var KeyringHome = os.ExpandEnv("$HOME/.starport/accounts")

var (
	ErrAccountExists       = errors.New("account already exists")
	ErrAccountDoesNotExist = errors.New("account does not exist")
)

const (
	accountPrefixCosmos = "cosmos"
	pubKeyPrefix        = "pub"
)

// Registry for accounts.
type Registry struct {
	kr keyring.Keyring
}

// New creates a new registry to manage accounts.
func New(backend string) (Registry, error) {
	kr, err := keyring.New(KeyringServiceName, backend, KeyringHome, os.Stdin)
	if err != nil {
		return Registry{}, err
	}

	r := Registry{
		kr: kr,
	}

	return r, nil
}

// Account represents an CosmosSDK account.
type Account struct {
	// Name of the account.
	Name string

	// Info holds additional info about the account.
	Info keyring.Info
}

// Address returns an address for acocunt for given prefix.
func (a Account) Address(accPrefix string) string {
	if accPrefix == "" {
		accPrefix = accountPrefixCosmos
	}

	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(accPrefix, pubKeyPrefix)

	ko, err := keyring.Bech32KeyOutput(a.Info)
	if err != nil {
		panic(err)
	}
	return ko.Address
}

// PubKey returns a public key for given account prefix.
func (a Account) PubKey(accPrefix string) string {
	if accPrefix == "" {
		accPrefix = accountPrefixCosmos
	}

	conf := types.GetConfig()
	conf.SetBech32PrefixForAccount(accPrefix, pubKeyPrefix)

	o, err := keyring.Bech32KeyOutput(a.Info)
	if err != nil {
		panic(err)
	}
	return o.PubKey
}

// Create creates a new account with name.
func (r Registry) Create(name string) (acc Account, mnemonic string, err error) {
	acc, err = r.GetByName(name)
	if err == nil {
		return Account{}, "", ErrAccountExists
	}
	if !errors.Is(err, ErrAccountDoesNotExist) {
		return Account{}, "", err
	}

	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return Account{}, "", err
	}
	mnemonic, err = bip39.NewMnemonic(entropySeed)
	if err != nil {
		return Account{}, "", err
	}

	algo, err := r.algo()
	if err != nil {
		return Account{}, "", err
	}
	info, err := r.kr.NewAccount(name, mnemonic, "", r.hdPath(), algo)
	if err != nil {
		return Account{}, "", err
	}

	acc = Account{
		Name: name,
		Info: info,
	}

	return acc, mnemonic, nil
}

// Import imports an existing account with name andd passphrase and secret where secret can be a
// mnemonic or private key.
func (r Registry) Import(name, secret, passphrase string) (Account, error) {
	_, err := r.GetByName(name)
	if err == nil {
		return Account{}, ErrAccountExists
	}
	if !errors.Is(err, ErrAccountDoesNotExist) {
		return Account{}, err
	}

	if bip39.IsMnemonicValid(secret) {
		algo, err := r.algo()
		if err != nil {
			return Account{}, err
		}
		_, err = r.kr.NewAccount(name, secret, passphrase, r.hdPath(), algo)
		if err != nil {
			return Account{}, err
		}
	} else if err := r.kr.ImportPrivKey(name, secret, passphrase); err != nil {
		return Account{}, err
	}

	return r.GetByName(name)
}

// Export exports an account as a private key.
func (r Registry) Export(name, passphrase string) (key string, err error) {
	if _, err = r.GetByName(name); err != nil {
		return "", err
	}

	return r.kr.ExportPrivKeyArmor(name, passphrase)

}

// GetByName returns an account by its name.
func (r Registry) GetByName(name string) (Account, error) {
	info, err := r.kr.Key(name)
	if err == dkeyring.ErrKeyNotFound {
		return Account{}, ErrAccountDoesNotExist
	}
	if err != nil {
		return Account{}, nil
	}

	acc := Account{
		Name: name,
		Info: info,
	}

	return acc, nil
}

// List lists all accounts.
func (r Registry) List() ([]Account, error) {
	info, err := r.kr.List()
	if err != nil {
		return nil, err
	}

	var accounts []Account

	for _, accinfo := range info {
		accounts = append(accounts, Account{
			Name: accinfo.GetName(),
			Info: accinfo,
		})
	}

	return accounts, nil
}

// DeleteByName deletes an account by name.
func (r Registry) DeleteByName(name string) error {
	err := r.kr.Delete(name)
	if err == dkeyring.ErrKeyNotFound {
		return ErrAccountDoesNotExist
	}
	return err
}

func (r Registry) hdPath() string {
	return hd.CreateHDPath(types.GetConfig().GetCoinType(), 0, 0).String()
}

func (r Registry) algo() (keyring.SignatureAlgo, error) {
	algos, _ := r.kr.SupportedAlgorithms()
	return keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
}
