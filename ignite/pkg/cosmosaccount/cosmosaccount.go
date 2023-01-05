package cosmosaccount

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	dkeyring "github.com/99designs/keyring"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/go-bip39"
)

const (
	// KeyringServiceName used for the name of keyring in OS backend.
	KeyringServiceName = "starport"

	// DefaultAccount is the name of the default account.
	DefaultAccount = "default"
)

// KeyringHome used to store account related data.
var KeyringHome = os.ExpandEnv("$HOME/.ignite/accounts")

var ErrAccountExists = errors.New("account already exists")

const (
	AccountPrefixCosmos = "cosmos"
)

// KeyringBackend is the backend for where keys are stored.
type KeyringBackend string

const (
	// KeyringTest is the test keyring backend. With this backend, your keys will be
	// stored under your app's data dir.
	KeyringTest KeyringBackend = "test"

	// KeyringOS is the OS keyring backend. With this backend, your keys will be
	// stored in your operating system's secured keyring.
	KeyringOS KeyringBackend = "os"

	// KeyringMemory is in memory keyring backend, your keys will be stored in application memory.
	KeyringMemory KeyringBackend = "memory"
)

// Registry for accounts.
type Registry struct {
	homePath           string
	keyringServiceName string
	keyringBackend     KeyringBackend

	Keyring keyring.Keyring
}

// Option configures your registry.
type Option func(*Registry)

func WithHome(path string) Option {
	return func(c *Registry) {
		c.homePath = path
	}
}

func WithKeyringServiceName(name string) Option {
	return func(c *Registry) {
		c.keyringServiceName = name
	}
}

func WithKeyringBackend(backend KeyringBackend) Option {
	return func(c *Registry) {
		c.keyringBackend = backend
	}
}

// New creates a new registry to manage accounts.
func New(options ...Option) (Registry, error) {
	r := Registry{
		keyringServiceName: sdktypes.KeyringServiceName(),
		keyringBackend:     KeyringTest,
		homePath:           KeyringHome,
	}

	for _, apply := range options {
		apply(&r)
	}

	var err error
	inBuf := bufio.NewReader(os.Stdin)
	interfaceRegistry := types.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	r.Keyring, err = keyring.New(r.keyringServiceName, string(r.keyringBackend), r.homePath, inBuf, cdc)
	if err != nil {
		return Registry{}, err
	}

	return r, nil
}

func NewStandalone(options ...Option) (Registry, error) {
	return New(
		append([]Option{
			WithKeyringServiceName(KeyringServiceName),
			WithHome(KeyringHome),
		}, options...)...,
	)
}

func NewInMemory(options ...Option) (Registry, error) {
	return New(
		append([]Option{
			WithKeyringBackend(KeyringMemory),
		}, options...)...,
	)
}

// Account represents a Cosmos SDK account.
type Account struct {
	// Name of the account.
	Name string

	// Record holds additional info about the account.
	Record *keyring.Record
}

// Address returns the address of the account from given prefix.
func (a Account) Address(accPrefix string) (string, error) {
	if accPrefix == "" {
		accPrefix = AccountPrefixCosmos
	}

	pk, err := a.Record.GetPubKey()
	if err != nil {
		return "", err
	}

	return toBech32(accPrefix, pk.Address())
}

// PubKey returns a public key for account.
func (a Account) PubKey() (string, error) {
	pk, err := a.Record.GetPubKey()
	if err != nil {
		return "", nil
	}

	return pk.String(), nil
}

func toBech32(prefix string, addr []byte) (string, error) {
	bech32Addr, err := bech32.ConvertAndEncode(prefix, addr)
	if err != nil {
		return "", err
	}
	return bech32Addr, nil
}

// EnsureDefaultAccount ensures that default account exists.
func (r Registry) EnsureDefaultAccount() error {
	_, err := r.GetByName(DefaultAccount)

	var accErr *AccountDoesNotExistError
	if errors.As(err, &accErr) {
		_, _, err = r.Create(DefaultAccount)
		return err
	}

	return err
}

// Create creates a new account with name.
func (r Registry) Create(name string) (acc Account, mnemonic string, err error) {
	acc, err = r.GetByName(name)
	if err == nil {
		return Account{}, "", ErrAccountExists
	}
	var accErr *AccountDoesNotExistError
	if !errors.As(err, &accErr) {
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
	record, err := r.Keyring.NewAccount(name, mnemonic, "", r.hdPath(), algo)
	if err != nil {
		return Account{}, "", err
	}

	acc = Account{
		Name:   name,
		Record: record,
	}

	return acc, mnemonic, nil
}

// Import imports an existing account with name and passphrase and secret where secret can be a
// mnemonic or a private key.
func (r Registry) Import(name, secret, passphrase string) (Account, error) {
	_, err := r.GetByName(name)
	if err == nil {
		return Account{}, ErrAccountExists
	}
	var accErr *AccountDoesNotExistError
	if !errors.As(err, &accErr) {
		return Account{}, err
	}

	if bip39.IsMnemonicValid(secret) {
		algo, err := r.algo()
		if err != nil {
			return Account{}, err
		}
		_, err = r.Keyring.NewAccount(name, secret, passphrase, r.hdPath(), algo)
		if err != nil {
			return Account{}, err
		}
	} else if err := r.Keyring.ImportPrivKey(name, secret, passphrase); err != nil {
		return Account{}, err
	}

	return r.GetByName(name)
}

// Export exports an account as a private key.
func (r Registry) Export(name, passphrase string) (key string, err error) {
	if _, err = r.GetByName(name); err != nil {
		return "", err
	}

	return r.Keyring.ExportPrivKeyArmor(name, passphrase)
}

// ExportHex exports an account as a private key in hex.
func (r Registry) ExportHex(name, passphrase string) (hex string, err error) {
	if _, err = r.GetByName(name); err != nil {
		return "", err
	}

	return unsafeExportPrivKeyHex(r.Keyring, name, passphrase)
}

func unsafeExportPrivKeyHex(kr keyring.Keyring, uid, passphrase string) (privKey string, err error) {
	priv, err := kr.ExportPrivKeyArmor(uid, passphrase)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString([]byte(priv)), nil
}

// GetByName returns an account by its name.
func (r Registry) GetByName(name string) (Account, error) {
	record, err := r.Keyring.Key(name)
	if errors.Is(err, dkeyring.ErrKeyNotFound) || errors.Is(err, sdkerrors.ErrKeyNotFound) {
		return Account{}, &AccountDoesNotExistError{name}
	}
	if err != nil {
		return Account{}, err
	}

	acc := Account{
		Name:   name,
		Record: record,
	}

	return acc, nil
}

// GetByAddress returns an account by its address.
func (r Registry) GetByAddress(address string) (Account, error) {
	sdkAddr, err := sdktypes.AccAddressFromBech32(address)
	if err != nil {
		return Account{}, err
	}
	record, err := r.Keyring.KeyByAddress(sdkAddr)
	if errors.Is(err, dkeyring.ErrKeyNotFound) || errors.Is(err, sdkerrors.ErrKeyNotFound) {
		return Account{}, &AccountDoesNotExistError{address}
	}
	if err != nil {
		return Account{}, err
	}
	return Account{
		Name:   record.Name,
		Record: record,
	}, nil
}

// List lists all accounts.
func (r Registry) List() ([]Account, error) {
	records, err := r.Keyring.List()
	if err != nil {
		return nil, err
	}

	var accounts []Account

	for _, record := range records {
		accounts = append(accounts, Account{
			Name:   record.Name,
			Record: record,
		})
	}

	return accounts, nil
}

// DeleteByName deletes an account by name.
func (r Registry) DeleteByName(name string) error {
	err := r.Keyring.Delete(name)
	if errors.Is(err, dkeyring.ErrKeyNotFound) {
		return &AccountDoesNotExistError{name}
	}
	return err
}

func (r Registry) hdPath() string {
	return hd.CreateHDPath(sdktypes.GetConfig().GetCoinType(), 0, 0).String()
}

func (r Registry) algo() (keyring.SignatureAlgo, error) {
	algos, _ := r.Keyring.SupportedAlgorithms()
	return keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
}

type AccountDoesNotExistError struct {
	Name string
}

func (e *AccountDoesNotExistError) Error() string {
	return fmt.Sprintf("account %q does not exist", e.Name)
}
