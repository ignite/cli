// Package cosmosclient provides a standalone client to connect to Cosmos SDK chains.
package cosmosclient

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/spn/app/params"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	defaultNodeAddress   = "http://localhost:26657"
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 300000
)

// KeyringBackend is the backend for where keys are stored.
type KeyringBackend string

const (
	// KeyringTest is the test keyring backend. with this backend, your keys will be
	// stored under your app's data dir,
	KeyringTest KeyringBackend = "test"

	// KeyringOS is the OS keyring backend. with this backend, your keys will be
	// stored in your operating system's secured keyring.
	KeyringOS KeyringBackend = "os"
)

// Client is a client to access your chain by querying and broadcasting transactions.
type Client struct {
	// RPC is Tendermint RPC.
	RPC *rpchttp.HTTP

	// Factory is a Cosmos SDK tx factory.
	Factory tx.Factory

	// context is a Cosmos SDK client context.
	Context client.Context

	// Keyring is a Cosmos SDK keyring.
	Keyring keyring.Keyring

	nodeAddress        string
	out                io.Writer
	chainID            string
	homePath           string
	keyringServiceName string
	keyringBackend     KeyringBackend
}

// Option configures your client.
type Option func(*Client)

// WithNodeAddress sets the node address of your chain. when this option is not provided
// `http://localhost:26657` used as default.
func WithNodeAddress(addr string) Option {
	return func(c *Client) {
		c.nodeAddress = addr
	}
}

// WithHome sets the data dir of your chain. this option is used to access to your chain's
// file based keyring which is only needed when you deal with creating and signing transactions.
// when it is not provided, your data dir will be assumed as `$HOME/.your-chain-id`.
func WithHome(path string) Option {
	return func(c *Client) {
		c.homePath = path
	}
}

// WithKeyringServiceName used as the keyring's name when you are using OS keyring backend.
// by default it is `cosmos`.
func WithKeyringServiceName(name string) Option {
	return func(c *Client) {
		c.keyringServiceName = name
	}
}

// WithKeyringBackend sets your keyring backend. by default it is `test`.
func WithKeyringBackend(backend KeyringBackend) Option {
	return func(c *Client) {
		c.keyringBackend = backend
	}
}

// New creates a new client with given options.
func New(ctx context.Context, options ...Option) (Client, error) {
	c := Client{
		nodeAddress:        defaultNodeAddress,
		out:                io.Discard,
		keyringServiceName: sdktypes.KeyringServiceName(),
		keyringBackend:     KeyringTest,
	}

	var err error

	for _, apply := range options {
		apply(&c)
	}

	if c.RPC, err = rpchttp.New(c.nodeAddress, "/websocket"); err != nil {
		return Client{}, err
	}

	statusResp, err := c.RPC.Status(ctx)
	if err != nil {
		return Client{}, err
	}

	c.chainID = statusResp.NodeInfo.Network

	if c.homePath == "" {
		c.homePath = os.ExpandEnv(fmt.Sprintf("$HOME/.%s", c.chainID))
	}

	if c.Keyring, err = keyring.New(c.keyringServiceName, string(c.keyringBackend), c.homePath, os.Stdin); err != nil {
		return Client{}, err
	}

	c.Context = newContext(c.RPC, c.out, c.chainID, c.homePath).WithKeyring(c.Keyring)
	c.Factory = newFactory(c.Context)

	return c, nil
}

// Address returns the account address from account name.
func (c Client) Address(accountName string) (sdktypes.AccAddress, error) {
	accountInfo, err := c.Keyring.Key(accountName)
	if err != nil {
		return nil, err
	}
	return accountInfo.GetAddress(), nil
}

// BroadcastTx creates and broadcasts a tx with given messages for account.
func (c Client) BroadcastTx(accountName string, messages ...sdktypes.Msg) (*sdktypes.TxResponse, error) {
	for _, message := range messages {
		if err := message.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	accountAddress, err := c.Address(accountName)
	if err != nil {
		return nil, err
	}

	context := c.Context.
		WithFromName(accountName).
		WithFromAddress(accountAddress)

	txf, err := tx.PrepareFactory(context, c.Factory)
	if err != nil {
		return nil, err
	}

	txUnsigned, err := tx.BuildUnsignedTx(txf, messages...)
	if err != nil {
		return nil, err
	}
	if err = tx.Sign(txf, accountName, txUnsigned, true); err != nil {
		return nil, err
	}

	txBytes, err := context.TxConfig.TxEncoder()(txUnsigned.GetTx())
	if err != nil {
		return nil, err
	}

	return context.BroadcastTx(txBytes)
}

func newContext(
	c *rpchttp.HTTP,
	out io.Writer,
	chainID,
	home string,
) client.Context {
	encodingConfig := params.MakeEncodingConfig()
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return client.Context{}.
		WithChainID(chainID).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(out).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(home).
		WithClient(c).
		WithSkipConfirmation(true)
}

func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(defaultGasLimit).
		WithGasAdjustment(defaultGasAdjustment).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
}
