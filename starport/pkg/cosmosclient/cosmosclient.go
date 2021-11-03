// Package cosmosclient provides a standalone client to connect to Cosmos SDK chains.
package cosmosclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	defaultNodeAddress   = "http://localhost:26657"
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 300000
)

const (
	faucetDenom     = "token"
	faucetMinAmount = 100
)

// Client is a client to access your chain by querying and broadcasting transactions.
type Client struct {
	// RPC is Tendermint RPC.
	RPC *rpchttp.HTTP

	// Factory is a Cosmos SDK tx factory.
	Factory tx.Factory

	// context is a Cosmos SDK client context.
	Context client.Context

	// AccountRegistry is the retistry to access accounts.
	AccountRegistry cosmosaccount.Registry

	addressPrefix string

	nodeAddress string
	apiAddress  string
	out         io.Writer
	chainID     string

	useFaucet       bool
	faucetAddress   string
	faucetDenom     string
	faucetMinAmount int64

	homePath           string
	keyringServiceName string
	keyringBackend     cosmosaccount.KeyringBackend
}

// Option configures your client.
type Option func(*Client)

// WithHome sets the data dir of your chain. This option is used to access your chain's
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

// WithKeyringBackend sets your keyring backend. By default, it is `test`.
func WithKeyringBackend(backend cosmosaccount.KeyringBackend) Option {
	return func(c *Client) {
		c.keyringBackend = backend
	}
}

// WithNodeAddress sets the node address of your chain. When this option is not provided
// `http://localhost:26657` is used as default.
func WithNodeAddress(addr string) Option {
	return func(c *Client) {
		c.nodeAddress = addr
	}
}

func WithAPIAddress(addr string) Option {
	return func(c *Client) {
		c.apiAddress = addr
	}
}

func WithAddressPrefix(prefix string) Option {
	return func(c *Client) {
		c.addressPrefix = prefix
	}
}

func WithUseFaucet(faucetAddress, denom string, minAmount int64) Option {
	return func(c *Client) {
		c.useFaucet = true
		c.faucetAddress = faucetAddress
		if denom != "" {
			c.faucetDenom = denom
		}
		if minAmount != 0 {
			c.faucetMinAmount = minAmount
		}
	}
}

// New creates a new client with given options.
func New(ctx context.Context, options ...Option) (Client, error) {
	c := Client{
		nodeAddress:   defaultNodeAddress,
		addressPrefix: "cosmos",
		out:           io.Discard,
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
		home, err := os.UserHomeDir()
		if err != nil {
			return Client{}, err
		}
		c.homePath = filepath.Join(home, "."+c.chainID)
	}

	c.AccountRegistry, err = cosmosaccount.New(
		cosmosaccount.WithKeyringServiceName(c.keyringServiceName),
		cosmosaccount.WithKeyringBackend(c.keyringBackend),
	)
	if err != nil {
		return Client{}, err
	}

	c.Context = newContext(c.RPC, c.out, c.chainID, c.homePath).WithKeyring(c.AccountRegistry.Keyring)
	c.Factory = newFactory(c.Context)

	return c, nil
}

func (c Client) Account(accountName string) (cosmosaccount.Account, error) {
	return c.AccountRegistry.GetByName(accountName)
}

// Address returns the account address from account name.
func (c Client) Address(accountName string) (sdktypes.AccAddress, error) {
	account, err := c.Account(accountName)
	if err != nil {
		return sdktypes.AccAddress{}, err
	}
	return account.Info.GetAddress(), nil
}

// BroadcastTx creates and broadcasts a tx with given messages for account.
func (c Client) BroadcastTx(accountName string, msgs ...sdktypes.Msg) (*sdktypes.TxResponse, error) {
	_, broadcast, err := c.BroadcastTxWithProvision(accountName, msgs...)
	if err != nil {
		return nil, err
	}
	return broadcast()
}

// protects sdktypes.Config.
var mconf sync.Mutex

func (c Client) BroadcastTxWithProvision(accountName string, msgs ...sdktypes.Msg) (
	gas uint64, broadcast func() (*sdktypes.TxResponse, error), err error) {
	if err := c.prepareBroadcast(context.Background(), accountName, msgs); err != nil {
		return 0, nil, err
	}

	// TODO find a better way if possible.
	mconf.Lock()
	defer mconf.Unlock()
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(c.addressPrefix, c.addressPrefix+"pub")

	accountAddress, err := c.Address(accountName)
	if err != nil {
		return 0, nil, err
	}

	context := c.Context.
		WithFromName(accountName).
		WithFromAddress(accountAddress)

	txf, err := c.Factory.Prepare(context)
	if err != nil {
		return 0, nil, err
	}

	_, gas, err = tx.CalculateGas(context, txf, msgs...)
	if err != nil {
		return 0, nil, err
	}
	// the simulated gas can vary from the actual gas needed for a real transaction
	// we add an additional amount to endure sufficient gas is provided
	gas += 10000
	txf = txf.WithGas(gas)

	// Return the provision function
	return gas, func() (*sdktypes.TxResponse, error) {
		txUnsigned, err := txf.BuildUnsignedTx(msgs...)
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

		resp, err := context.BroadcastTx(txBytes)
		return resp, handleBroadcastResult(resp, err)
	}, nil
}

// prepareBroadcast performs checks and operations before broadcasting messages
func (c *Client) prepareBroadcast(ctx context.Context, accountName string, _ []sdktypes.Msg) error {
	// TODO uncomment after https://github.com/tendermint/spn/issues/363
	// validate msgs.
	//  for _, msg := range msgs {
	//  if err := msg.ValidateBasic(); err != nil {
	//  return err
	//  }
	//  }

	account, err := c.Account(accountName)
	if err != nil {
		return err
	}

	// make sure that account has enough balances before broadcasting.
	if c.useFaucet {
		if err := c.makeSureAccountHasTokens(ctx, account.Address(c.addressPrefix)); err != nil {
			return err
		}
	}

	return nil
}

// makeSureAccountHasTokens makes sure the address has a positive balance
// it requests funds from the faucet if the address has an empty balance
func (c *Client) makeSureAccountHasTokens(ctx context.Context, address string) error {
	// check the balance.
	balancesEndpoint := fmt.Sprintf("%s/bank/balances/%s", c.apiAddress, address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, balancesEndpoint, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var balances struct {
		Result []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return err
	}

	// if the balance is enough do nothing.
	if len(balances.Result) > 0 {
		for _, c := range balances.Result {
			amount, err := strconv.ParseInt(c.Amount, 10, 32)
			if err != nil {
				return err
			}
			if c.Denom == faucetDenom && amount >= faucetMinAmount {
				return nil
			}
		}
	}

	// request coins from the faucet.
	fc := cosmosfaucet.NewClient(c.faucetAddress)
	faucetResp, err := fc.Transfer(ctx, cosmosfaucet.TransferRequest{AccountAddress: address})
	if err != nil {
		return errors.Wrap(err, "faucet server request failed")
	}
	if faucetResp.Error != "" {
		return fmt.Errorf("cannot retrieve tokens from faucet: %s", faucetResp.Error)
	}
	for _, transfer := range faucetResp.Transfers {
		if transfer.Error != "" {
			return fmt.Errorf("cannot retrieve tokens from faucet: %s", transfer.Error)
		}
	}

	return nil
}

// handleBroadcastResult handles the result of broadcast messages result and checks if an error occurred
func handleBroadcastResult(resp *sdktypes.TxResponse, err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("make sure that your SPN account has enough balance")
		}

		return err
	}

	if resp.Code > 0 {
		return fmt.Errorf("SPN error with '%d' code: %s", resp.Code, resp.RawLog)
	}
	return nil
}
func newContext(
	c *rpchttp.HTTP,
	out io.Writer,
	chainID,
	home string,
) client.Context {
	var (
		amino             = codec.NewLegacyAmino()
		interfaceRegistry = codectypes.NewInterfaceRegistry()
		marshaler         = codec.NewProtoCodec(interfaceRegistry)
		txConfig          = authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)
	)

	authtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	sdktypes.RegisterInterfaces(interfaceRegistry)
	staking.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)

	return client.Context{}.
		WithChainID(chainID).
		WithInterfaceRegistry(interfaceRegistry).
		WithCodec(marshaler).
		WithTxConfig(txConfig).
		WithLegacyAmino(amino).
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
