package spn

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	chattypes "github.com/tendermint/spn/x/chat/types"
	"github.com/tendermint/starport/starport/pkg/xurl"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var spn = "spn"
var homedir = os.ExpandEnv("$HOME/spnd")

// Account represents an account on SPN.
type Account struct {
	Name     string
	Address  string
	Mnemonic string
}

// Client is client to interact with SPN.
type Client struct {
	kr        keyring.Keyring
	factory   tx.Factory
	clientCtx client.Context
}

type options struct {
	keyringBackend string
}

// Option configures Client options.
type Option func(*options)

// Keyring uses given keyring type as storage.
func Keyring(keyring string) Option {
	return func(c *options) {
		c.keyringBackend = keyring
	}
}

// New creates a new SPN Client with nodeAddress of a full SPN node.
// by default, OS is used as keyring backend.
func New(nodeAddress string, option ...Option) (Client, error) {
	opts := &options{
		keyringBackend: keyring.BackendOS,
	}
	for _, o := range option {
		o(opts)
	}
	kr, err := keyring.New(types.KeyringServiceName(), opts.keyringBackend, homedir, os.Stdin)
	if err != nil {
		return Client{}, err
	}

	client, err := rpchttp.New(xurl.TCP(nodeAddress), "/websocket")
	if err != nil {
		return Client{}, err
	}
	clientCtx := NewClientCtx(kr, client)
	factory := NewFactory(clientCtx)
	return Client{
		kr:        kr,
		factory:   factory,
		clientCtx: clientCtx,
	}, nil
}

// AccountGet retrieves an account by name from the keyring.
func (c Client) AccountGet(accountName string) (Account, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return Account{}, err
	}
	return toAccount(info), nil
}

// AccountList returns a list of accounts.
func (c Client) AccountList() ([]Account, error) {
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

// AccountCreate creates an account by name in the keyring.
func (c Client) AccountCreate(accountName string) (Account, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return Account{}, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return Account{}, err
	}
	algos, _ := c.kr.SupportedAlgorithms()
	if err != nil {
		return Account{}, err
	}
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
func (c Client) AccountExport(accountName, password string) (privateKey string, err error) {
	return c.kr.ExportPrivKeyArmor(accountName, password)
}

// AccountImport imports an account to the keyring by account name, privateKey and decryption password.
func (c Client) AccountImport(accountName, privateKey, password string) error {
	return c.kr.ImportPrivKey(accountName, privateKey, password)
}

// ChainCreate creates a new chain.
// TODO right now this uses chat module, use genesis.
func (c Client) ChainCreate(accountName, chainID, genesis, sourceURL, sourceHash string) error {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return err
	}
	clientCtx := c.clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress())
	msg, err := chattypes.NewMsgCreateChannel(
		clientCtx.GetFromAddress(),
		chainID,
		sourceURL,
		[]byte(genesis),
	)
	if err != nil {
		return err
	}
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.BroadcastTx(clientCtx, c.factory, msg)
}

// Chain represents a chain in Genesis module of SPN.
type Chain struct {
	URL  string
	Hash string
}

// TODO ShowChain shows chain info.
func (c Client) ShowChain(accountName, chainID string) (Chain, error) {
	return Chain{
		URL:  "https://github.com/tendermint/spn",
		Hash: "df49c9256dfcbd0096fd0a8acdd4907ba3332cd5",
	}, nil
}
