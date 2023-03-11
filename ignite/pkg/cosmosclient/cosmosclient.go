// Package cosmosclient provides a standalone client to connect to Cosmos SDK chains.
package cosmosclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	prototypes "github.com/cosmos/gogoproto/types"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
)

var (
	// FaucetTransferEnsureDuration is the duration that BroadcastTx will wait when a faucet transfer
	// is triggered prior to broadcasting but transfer's tx is not committed in the state yet.
	FaucetTransferEnsureDuration = time.Second * 40

	// ErrInvalidBlockHeight is returned when a block height value is not valid.
	ErrInvalidBlockHeight = errors.New("block height must be greater than 0")

	errCannotRetrieveFundsFromFaucet = errors.New("cannot retrieve funds from faucet")
)

const (
	// GasAuto allows to calculate gas automatically when sending transaction.
	GasAuto = "auto"

	defaultNodeAddress   = "http://localhost:26657"
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 300000

	defaultFaucetAddress   = "http://localhost:4500"
	defaultFaucetDenom     = "token"
	defaultFaucetMinAmount = 100

	defaultTXsPerPage = 30

	searchHeight = "tx.height"

	orderAsc = "asc"
)

// FaucetClient allows to mock the cosmosfaucet.Client.
//
//go:generate mockery --srcpkg . --name FaucetClient --structname FaucetClient --filename faucet_client.go --with-expecter
type FaucetClient interface {
	Transfer(context.Context, cosmosfaucet.TransferRequest) (cosmosfaucet.TransferResponse, error)
}

// Gasometer allows to mock the tx.CalculateGas func.
//
//go:generate mockery --srcpkg . --name Gasometer --filename gasometer.go --with-expecter
type Gasometer interface {
	CalculateGas(clientCtx gogogrpc.ClientConn, txf tx.Factory, msgs ...sdktypes.Msg) (*txtypes.SimulateResponse, uint64, error)
}

// Signer allows to mock the tx.Sign func.
//
//go:generate mockery --srcpkg . --name Signer --filename signer.go --with-expecter
type Signer interface {
	Sign(txf tx.Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error
}

// Client is a client to access your chain by querying and broadcasting transactions.
type Client struct {
	// RPC is Tendermint RPC.
	RPC rpcclient.Client

	// TxFactory is a Cosmos SDK tx factory.
	TxFactory tx.Factory

	// context is a Cosmos SDK client context.
	context client.Context

	// AccountRegistry is the registry to access accounts.
	AccountRegistry cosmosaccount.Registry

	accountRetriever client.AccountRetriever
	bankQueryClient  banktypes.QueryClient
	faucetClient     FaucetClient
	gasometer        Gasometer
	signer           Signer

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

	gas           string
	gasPrices     string
	gasAdjustment float64
	fees          string
	generateOnly  bool
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

// WithKeyringServiceName used as the keyring name when you are using OS keyring backend.
// by default, it is `cosmos`.
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

// WithKeyringDir sets the directory of the keyring. By default, it uses cosmosaccount.KeyringHome.
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
// Set to "auto" to calculate automatically.
func WithGas(gas string) Option {
	return func(c *Client) {
		c.gas = gas
	}
}

// WithGasPrices sets the price per gas (e.g. 0.1uatom).
func WithGasPrices(gasPrices string) Option {
	return func(c *Client) {
		c.gasPrices = gasPrices
	}
}

// WithGasAdjustment sets the gas adjustment.
func WithGasAdjustment(gasAdjustment float64) Option {
	return func(c *Client) {
		c.gasAdjustment = gasAdjustment
	}
}

// WithFees sets the fees (e.g. 10uatom).
func WithFees(fees string) Option {
	return func(c *Client) {
		c.fees = fees
	}
}

// WithGenerateOnly tells if txs will be generated only.
func WithGenerateOnly(generateOnly bool) Option {
	return func(c *Client) {
		c.generateOnly = generateOnly
	}
}

// WithRPCClient sets a tendermint RPC client.
// Already set by default.
func WithRPCClient(rpc rpcclient.Client) Option {
	return func(c *Client) {
		c.RPC = rpc
	}
}

// WithAccountRetriever sets the account retriever
// Already set by default.
func WithAccountRetriever(accountRetriever client.AccountRetriever) Option {
	return func(c *Client) {
		c.accountRetriever = accountRetriever
	}
}

// WithBankQueryClient sets the bank query client.
// Already set by default.
func WithBankQueryClient(bankQueryClient banktypes.QueryClient) Option {
	return func(c *Client) {
		c.bankQueryClient = bankQueryClient
	}
}

// WithFaucetClient sets the faucet client.
// Already set by default.
func WithFaucetClient(faucetClient FaucetClient) Option {
	return func(c *Client) {
		c.faucetClient = faucetClient
	}
}

// WithGasometer sets the gasometer.
// Already set by default.
func WithGasometer(gasometer Gasometer) Option {
	return func(c *Client) {
		c.gasometer = gasometer
	}
}

// WithSigner sets the signer.
// Already set by default.
func WithSigner(signer Signer) Option {
	return func(c *Client) {
		c.signer = signer
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

	if c.RPC == nil {
		if c.RPC, err = rpchttp.New(c.nodeAddress, "/websocket"); err != nil {
			return Client{}, err
		}
	}
	// Wrap RPC client to have more contextualized errors
	c.RPC = rpcWrapper{
		Client:      c.RPC,
		nodeAddress: c.nodeAddress,
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

	c.context = c.newContext()
	c.TxFactory = newFactory(c.context)

	if c.accountRetriever == nil {
		c.accountRetriever = authtypes.AccountRetriever{}
	}
	if c.bankQueryClient == nil {
		c.bankQueryClient = banktypes.NewQueryClient(c.context)
	}
	if c.faucetClient == nil {
		c.faucetClient = cosmosfaucet.NewClient(c.faucetAddress)
	}
	if c.gasometer == nil {
		c.gasometer = gasometer{}
	}
	if c.signer == nil {
		c.signer = signer{}
	}
	// set address prefix in SDK global config
	c.SetConfigAddressPrefix()

	return c, nil
}

// LatestBlockHeight returns the latest block height of the app.
func (c Client) LatestBlockHeight(ctx context.Context) (int64, error) {
	resp, err := c.Status(ctx)
	if err != nil {
		return 0, err
	}
	return resp.SyncInfo.LatestBlockHeight, nil
}

// WaitForNextBlock waits until next block is committed.
// It reads the current block height and then waits for another block to be
// committed, or returns an error if ctx is canceled.
func (c Client) WaitForNextBlock(ctx context.Context) error {
	return c.WaitForNBlocks(ctx, 1)
}

// WaitForNBlocks reads the current block height and then waits for another n
// blocks to be committed, or returns an error if ctx is canceled.
func (c Client) WaitForNBlocks(ctx context.Context, n int64) error {
	start, err := c.LatestBlockHeight(ctx)
	if err != nil {
		return err
	}
	return c.WaitForBlockHeight(ctx, start+n)
}

// WaitForBlockHeight waits until block height h is committed, or returns an
// error if ctx is canceled.
func (c Client) WaitForBlockHeight(ctx context.Context, h int64) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		latestHeight, err := c.LatestBlockHeight(ctx)
		if err != nil {
			return err
		}
		if latestHeight >= h {
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "timeout exceeded waiting for block")
		case <-ticker.C:
		}
	}
}

// WaitForTx requests the tx from hash, if not found, waits for next block and
// tries again. Returns an error if ctx is canceled.
func (c Client) WaitForTx(ctx context.Context, hash string) (*ctypes.ResultTx, error) {
	bz, err := hex.DecodeString(hash)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode tx hash '%s'", hash)
	}
	for {
		resp, err := c.RPC.Tx(ctx, bz, false)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// Tx not found, wait for next block and try again
				err := c.WaitForNextBlock(ctx)
				if err != nil {
					return nil, errors.Wrap(err, "waiting for next block")
				}
				continue
			}
			return nil, errors.Wrapf(err, "fetching tx '%s'", hash)
		}
		// Tx found
		return resp, nil
	}
}

// Account returns the account with name or address equal to nameOrAddress.
func (c Client) Account(nameOrAddress string) (cosmosaccount.Account, error) {
	defer c.lockBech32Prefix()()

	acc, err := c.AccountRegistry.GetByName(nameOrAddress)
	if err == nil {
		return acc, nil
	}
	return c.AccountRegistry.GetByAddress(nameOrAddress)
}

// Address returns the account address from account name.
func (c Client) Address(accountName string) (string, error) {
	a, err := c.AccountRegistry.GetByName(accountName)
	if err != nil {
		return "", err
	}
	return a.Address(c.addressPrefix)
}

// Context returns client context.
func (c Client) Context() client.Context {
	return c.context
}

// SetConfigAddressPrefix sets the account prefix in the SDK global config.
func (c Client) SetConfigAddressPrefix() {
	// TODO find a better way if possible.
	// https://github.com/ignite/cli/issues/2744
	mconf.Lock()
	defer mconf.Unlock()
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(c.addressPrefix, c.addressPrefix+"pub")
}

// Response of your broadcasted transaction.
type Response struct {
	Codec codec.Codec

	// TxResponse is the underlying tx response.
	*sdktypes.TxResponse
}

// Decode decodes the proto func response defined in your Msg service into your message type.
// message needs to be a pointer. and you need to provide the correct proto message(struct) type to the Decode func.
//
// e.g., for the following CreateChain func the type would be: `types.MsgCreateChainResponse`.
//
// ```proto
//
//	service Msg {
//	  rpc CreateChain(MsgCreateChain) returns (MsgCreateChainResponse);
//	}
//
// ```
//
//nolint:godot,nolintlint
func (r Response) Decode(message proto.Message) error {
	data, err := hex.DecodeString(r.Data)
	if err != nil {
		return err
	}

	var txMsgData sdktypes.TxMsgData
	if err := r.Codec.Unmarshal(data, &txMsgData); err != nil {
		return err
	}

	// check deprecated Data
	if len(txMsgData.Data) != 0 {
		resData := txMsgData.Data[0]
		return prototypes.UnmarshalAny(&prototypes.Any{
			// TODO get type url dynamically(basically remove `+ "Response"`) after the following issue has solved.
			// https://github.com/ignite/cli/issues/2098
			// https://github.com/cosmos/cosmos-sdk/issues/10496
			TypeUrl: resData.MsgType + "Response",
			Value:   resData.Data,
		}, message)
	}

	resData := txMsgData.MsgResponses[0]
	return prototypes.UnmarshalAny(&prototypes.Any{
		TypeUrl: resData.TypeUrl,
		Value:   resData.Value,
	}, message)
}

// Status returns the node status.
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

func (c Client) BroadcastTx(ctx context.Context, account cosmosaccount.Account, msgs ...sdktypes.Msg) (Response, error) {
	txService, err := c.CreateTx(ctx, account, msgs...)
	if err != nil {
		return Response{}, err
	}

	return txService.Broadcast(ctx)
}

func (c Client) CreateTx(goCtx context.Context, account cosmosaccount.Account, msgs ...sdktypes.Msg) (TxService, error) {
	defer c.lockBech32Prefix()()

	if c.useFaucet && !c.generateOnly {
		addr, err := account.Address(c.addressPrefix)
		if err != nil {
			return TxService{}, errors.WithStack(err)
		}
		if err := c.makeSureAccountHasTokens(goCtx, addr); err != nil {
			return TxService{}, err
		}
	}

	sdkaddr, err := account.Record.GetAddress()
	if err != nil {
		return TxService{}, errors.WithStack(err)
	}

	ctx := c.context.
		WithFromName(account.Name).
		WithFromAddress(sdkaddr)

	txf, err := c.prepareFactory(ctx)
	if err != nil {
		return TxService{}, err
	}

	var gas uint64
	if c.gas != "" && c.gas != GasAuto {
		gas, err = strconv.ParseUint(c.gas, 10, 64)
		if err != nil {
			return TxService{}, errors.WithStack(err)
		}
	} else {
		_, gas, err = c.gasometer.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return TxService{}, errors.WithStack(err)
		}
		// the simulated gas can vary from the actual gas needed for a real transaction
		// we add an amount to ensure sufficient gas is provided
		gas += 20000
	}
	txf = txf.WithGas(gas)
	txf = txf.WithFees(c.fees)

	if c.gasPrices != "" {
		txf = txf.WithGasPrices(c.gasPrices)
	}

	if c.gasAdjustment != 0 && c.gasAdjustment != defaultGasAdjustment {
		txf = txf.WithGasAdjustment(c.gasAdjustment)
	}

	txUnsigned, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return TxService{}, errors.WithStack(err)
	}

	txUnsigned.SetFeeGranter(ctx.GetFeeGranterAddress())

	return TxService{
		client:        c,
		clientContext: ctx,
		txBuilder:     txUnsigned,
		txFactory:     txf,
	}, nil
}

// GetBlockTXs returns the transactions in a block.
// The list of transactions can be empty if there are no transactions in the block
// at the moment this method is called.
// Tendermint might index a limited number of block so trying to fetch transactions
// from a block that is not indexed would return an error.
func (c Client) GetBlockTXs(ctx context.Context, height int64) (txs []TX, err error) {
	if height == 0 {
		return nil, ErrInvalidBlockHeight
	}

	r, err := c.RPC.Block(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", height, err)
	}

	query := createTxSearchByHeightQuery(height)

	// TODO: improve to fetch pages in parallel (requires fetching page 1 to calculate n. of pages)
	page := 1
	perPage := defaultTXsPerPage
	blockTime := r.Block.Time
	for {
		res, err := c.RPC.TxSearch(ctx, query, false, &page, &perPage, orderAsc)
		if err != nil {
			return nil, err
		}

		for _, tx := range res.Txs {
			txs = append(txs, TX{
				BlockTime: blockTime,
				Raw:       tx,
			})
		}

		// Stop when the last page is fetched
		if res.TotalCount <= (page * perPage) {
			break
		}

		page++
	}

	return txs, nil
}

// CollectTXs collects transactions from multiple consecutive blocks.
// Transactions from a single block are send to the channel only if all transactions
// from that block are collected successfully.
// Blocks are traversed sequentially starting from a height until the latest block height
// available at the moment this method is called.
// The channel might contain the transactions collected successfully up until that point
// when an error is returned.
func (c Client) CollectTXs(ctx context.Context, fromHeight int64, tc chan<- []TX) error {
	defer close(tc)

	latestHeight, err := c.LatestBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block height: %w", err)
	}

	if fromHeight == 0 {
		fromHeight = 1
	}

	for height := fromHeight; height <= latestHeight; height++ {
		txs, err := c.GetBlockTXs(ctx, height)
		if err != nil {
			return err
		}

		// Ignore blocks without transactions
		if txs == nil {
			continue
		}

		// Make sure that collection finishes if the context
		// is done when the transactions channel is full
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tc <- txs:
		}
	}

	return nil
}

// makeSureAccountHasTokens makes sure the address has a positive balance.
// It requests funds from the faucet if the address has an empty balance.
func (c *Client) makeSureAccountHasTokens(ctx context.Context, address string) error {
	if err := c.checkAccountBalance(ctx, address); err == nil {
		return nil
	}

	// request coins from the faucet.
	faucetResp, err := c.faucetClient.Transfer(ctx, cosmosfaucet.TransferRequest{AccountAddress: address})
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
	resp, err := c.bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
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

// handleBroadcastResult handles the result of broadcast messages result and checks if an error occurred.
func handleBroadcastResult(resp *sdktypes.TxResponse, err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("make sure that your account has enough balance")
		}
		return err
	}

	if resp.Code > 0 {
		return errors.Errorf("error code: '%d' msg: '%s'", resp.Code, resp.RawLog)
	}
	return nil
}

func (c *Client) prepareFactory(clientCtx client.Context) (tx.Factory, error) {
	var (
		from = clientCtx.GetFromAddress()
		txf  = c.TxFactory
	)

	if err := c.accountRetriever.EnsureExists(clientCtx, from); err != nil {
		return txf, errors.WithStack(err)
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := c.accountRetriever.GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, errors.WithStack(err)
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

func (c Client) newContext() client.Context {
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
		WithChainID(c.chainID).
		WithInterfaceRegistry(interfaceRegistry).
		WithCodec(marshaler).
		WithTxConfig(txConfig).
		WithLegacyAmino(amino).
		WithInput(os.Stdin).
		WithOutput(c.out).
		WithAccountRetriever(c.accountRetriever).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(c.homePath).
		WithClient(c.RPC).
		WithSkipConfirmation(true).
		WithKeyring(c.AccountRegistry.Keyring).
		WithGenerateOnly(c.generateOnly)
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

func createTxSearchByHeightQuery(height int64) string {
	params := url.Values{}
	params.Set(searchHeight, strconv.FormatInt(height, 10))
	return params.Encode()
}
