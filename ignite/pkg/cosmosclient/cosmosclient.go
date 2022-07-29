// Package cosmosclient provides a standalone client to connect to Cosmos SDK chains.
package cosmosclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
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
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
)

// FaucetTransferEnsureDuration is the duration that BroadcastTx will wait when a faucet transfer
// is triggered prior to broadcasting but transfer's tx is not committed in the state yet.
var FaucetTransferEnsureDuration = time.Second * 40

var errCannotRetrieveFundsFromFaucet = errors.New("cannot retrieve funds from faucet")

const (
	defaultNodeAddress   = "http://localhost:26657"
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 300000
)

const (
	defaultFaucetAddress   = "http://localhost:4500"
	defaultFaucetDenom     = "token"
	defaultFaucetMinAmount = 100
)

// Client is a client to access your chain by querying and broadcasting transactions.
type Client struct {
	// RPC is Tendermint RPC.
	RPC *rpchttp.HTTP

	// Factory is a Cosmos SDK tx factory.
	Factory tx.Factory

	// context is a Cosmos SDK client context.
	context client.Context

	// AccountRegistry is the retistry to access accounts.
	AccountRegistry cosmosaccount.Registry

	addressPrefix string

	nodeAddress string
	out         io.Writer
	chainID     string

	useFaucet       bool
	faucetAddress   string
	faucetDenom     string
	faucetMinAmount uint64

	homePath           string
	keyringServiceName string
	keyringBackend     cosmosaccount.KeyringBackend
	keyringDir         string

	gas       string
	gasPrices string
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

// WithKeyringDir sets the directory of the keyring. By default, it uses cosmosaccount.KeyringHome
func WithKeyringDir(keyringDir string) Option {
	return func(c *Client) {
		c.keyringDir = keyringDir
	}
}

// WithNodeAddress sets the node address of your chain. When this option is not provided
// `http://localhost:26657` is used as default.
func WithNodeAddress(addr string) Option {
	return func(c *Client) {
		c.nodeAddress = addr
	}
}

func WithAddressPrefix(prefix string) Option {
	return func(c *Client) {
		c.addressPrefix = prefix
	}
}

func WithUseFaucet(faucetAddress, denom string, minAmount uint64) Option {
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

// WithGas sets an explicit gas-limit on transactions.
// Set to "auto" to calculate automatically
func WithGas(gas string) Option {
	return func(c *Client) {
		c.gas = gas
	}
}

// WithGasPrices sets the price per gas (e.g. 0.1uatom)
func WithGasPrices(gasPrices string) Option {
	return func(c *Client) {
		c.gasPrices = gasPrices
	}
}

// New creates a new client with given options.
func New(ctx context.Context, options ...Option) (Client, error) {
	c := Client{
		nodeAddress:     defaultNodeAddress,
		keyringBackend:  cosmosaccount.KeyringTest,
		addressPrefix:   "cosmos",
		faucetAddress:   defaultFaucetAddress,
		faucetDenom:     defaultFaucetDenom,
		faucetMinAmount: defaultFaucetMinAmount,
		out:             io.Discard,
		gas:             strconv.Itoa(defaultGasLimit),
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

	if c.keyringDir == "" {
		c.keyringDir = c.homePath
	}

	c.AccountRegistry, err = cosmosaccount.New(
		cosmosaccount.WithKeyringServiceName(c.keyringServiceName),
		cosmosaccount.WithKeyringBackend(c.keyringBackend),
		cosmosaccount.WithHome(c.keyringDir),
	)
	if err != nil {
		return Client{}, err
	}

	c.context = newContext(c.RPC, c.out, c.chainID, c.homePath).WithKeyring(c.AccountRegistry.Keyring)
	c.Factory = newFactory(c.context)

	return c, nil
}

// Account returns the account with name or address equal to nameOrAddress.
func (c Client) Account(nameOrAddress string) (cosmosaccount.Account, error) {
	a, err := c.AccountRegistry.GetByName(nameOrAddress)
	if err == nil {
		return a, nil
	}
	return c.AccountRegistry.GetByAddress(nameOrAddress)
}

// Address returns the account address from account name.
func (c Client) Address(accountName string) (sdktypes.AccAddress, error) {
	account, err := c.Account(accountName)
	if err != nil {
		return sdktypes.AccAddress{}, err
	}
	return account.Info.GetAddress(), nil
}

func (c Client) Context() client.Context {
	return c.context
}

// Response of your broadcasted transaction.
type Response struct {
	Codec codec.Codec

	// TxResponse is the underlying tx response.
	*sdktypes.TxResponse
}

// Decode decodes the proto func response defined in your Msg service into your message type.
// message needs be a pointer. and you need to provide the correct proto message(struct) type to the Decode func.
//
// e.g., for the following CreateChain func the type would be: `types.MsgCreateChainResponse`.
//
// ```proto
// service Msg {
//   rpc CreateChain(MsgCreateChain) returns (MsgCreateChainResponse);
// }
// ```
func (r Response) Decode(message proto.Message) error {
	data, err := hex.DecodeString(r.Data)
	if err != nil {
		return err
	}

	var txMsgData sdktypes.TxMsgData
	if err := r.Codec.Unmarshal(data, &txMsgData); err != nil {
		return err
	}

	resData := txMsgData.Data[0]

	return prototypes.UnmarshalAny(&prototypes.Any{
		// TODO get type url dynamically(basically remove `+ "Response"`) after the following issue has solved.
		// https://github.com/cosmos/cosmos-sdk/issues/10496
		TypeUrl: resData.MsgType + "Response",
		Value:   resData.Data,
	}, message)
}

// Status returns the node status
func (c Client) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	return c.RPC.Status(ctx)
}

// protects sdktypes.Config.
var mconf sync.Mutex

func (c Client) lockBech32Prefix() (unlockFn func()) {
	mconf.Lock()
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(c.addressPrefix, c.addressPrefix+"pub")

	return mconf.Unlock
}

func performQuery[T any](c Client, q func() (T, error)) (T, error) {
	defer c.lockBech32Prefix()()

	return q()
}

func (c Client) BroadcastTx(account cosmosaccount.Account, msgs ...sdktypes.Msg) (Response, error) {
	txService, err := c.CreateTx(account, msgs...)
	if err != nil {
		return Response{}, err
	}

	return txService.Broadcast()
}

func (c Client) CreateTx(account cosmosaccount.Account, msgs ...sdktypes.Msg) (TxService, error) {
	defer c.lockBech32Prefix()()

	ctx := c.context.
		WithFromName(account.Name).
		WithFromAddress(account.Info.GetAddress())

	txf, err := prepareFactory(ctx, c.Factory)
	if err != nil {
		return TxService{}, err
	}

	var gas uint64
	if c.gas != "" && c.gas != "auto" {
		gas, err = strconv.ParseUint(c.gas, 10, 64)
		if err != nil {
			return TxService{}, err
		}
	} else {
		_, gas, err = tx.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return TxService{}, err
		}
		// the simulated gas can vary from the actual gas needed for a real transaction
		// we add an additional amount to endure sufficient gas is provided
		gas += 10000
	}
	txf = txf.WithGas(gas)

	if c.gasPrices != "" {
		txf = txf.WithGasPrices(c.gasPrices)
	}

	txUnsigned, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return TxService{}, err
	}

	txUnsigned.SetFeeGranter(ctx.GetFeeGranterAddress())

	return TxService{
		client:        c,
		clientContext: ctx,
		txBuilder:     txUnsigned,
		txFactory:     txf,
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
	if err := c.checkAccountBalance(ctx, address); err == nil {
		return nil
	}

	// request coins from the faucet.
	fc := cosmosfaucet.NewClient(c.faucetAddress)
	faucetResp, err := fc.Transfer(ctx, cosmosfaucet.TransferRequest{AccountAddress: address})
	if err != nil {
		return errors.Wrap(errCannotRetrieveFundsFromFaucet, err.Error())
	}
	if faucetResp.Error != "" {
		return errors.Wrap(errCannotRetrieveFundsFromFaucet, faucetResp.Error)
	}

	// make sure funds are retrieved.
	ctx, cancel := context.WithTimeout(ctx, FaucetTransferEnsureDuration)
	defer cancel()

	return backoff.Retry(func() error {
		return c.checkAccountBalance(ctx, address)
	}, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}

func (c *Client) checkAccountBalance(ctx context.Context, address string) error {
	resp, err := banktypes.NewQueryClient(c.context).Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   c.faucetDenom,
	})
	if err != nil {
		return err
	}

	if resp.Balance.Amount.Uint64() >= c.faucetMinAmount {
		return nil
	}

	return fmt.Errorf("account has not enough %q balance, min. required amount: %d", c.faucetDenom, c.faucetMinAmount)
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

func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
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
	banktypes.RegisterInterfaces(interfaceRegistry)

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
