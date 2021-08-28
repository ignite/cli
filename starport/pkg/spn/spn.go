package spn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

var spnHomePath = xfilepath.JoinFromHome(xfilepath.Path("spnd"))

const (
	spn             = "spn"
	faucetDenom     = "token"
	faucetMinAmount = 100
)

// Client is client to interact with SPN.
type Client struct {
	apiAddress    string
	faucetAddress string
	cosmos        cosmosclient.Client
}

type options struct {
	keyringBackend string
}

// Option configures Client options.
type Option func(*options)

// New creates a new SPN Client with nodeAddress of a full SPN node.
// by default, OS is used as keyring backend.
func New(ctx context.Context, nodeAddress, apiAddress, faucetAddress string, option ...Option) (*Client, error) {
	opts := &options{
		keyringBackend: keyring.BackendOS,
	}
	for _, o := range option {
		o(opts)
	}

	homePath, err := spnHomePath()
	if err != nil {
		return nil, err
	}

	cosmos, err := cosmosclient.New(
		ctx,
		cosmosclient.WithNodeAddress(nodeAddress),
		cosmosclient.WithHome(homePath),
		cosmosclient.WithKeyringBackend(cosmosclient.KeyringBackend(opts.keyringBackend)),
	)
	return &Client{
		cosmos:        cosmos,
		apiAddress:    apiAddress,
		faucetAddress: faucetAddress,
	}, nil
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
func handleBroadcastResult(resp *types.TxResponse, err error) error {
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

// prepareBroadcast performs checks and operations before broadcasting messages
func (c *Client) prepareBroadcast(ctx context.Context, accountName string, msgs ...types.Msg) error {
	// validate msgs.
	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
	}

	addr, err := c.cosmos.Address(accountName)
	if err != nil {
		return err
	}

	// make sure that account has enough balances before broadcasting.
	if err := c.makeSureAccountHasTokens(ctx, addr.String()); err != nil {
		return err
	}

	return nil
}

// broadcast directly broadcasts the messages into spn handlers
func (c *Client) broadcast(ctx context.Context, accountName string, msgs ...types.Msg) error {
	if err := c.prepareBroadcast(ctx, accountName, msgs...); err != nil {
		return err
	}

	// broadcast tx.
	return handleBroadcastResult(c.cosmos.BroadcastTx(accountName, msgs...))
}

// broadcastProvision provides a provision function to broadcast the messages with returned amount of gas
func (c *Client) broadcastProvision(ctx context.Context, accountName string, msgs ...types.Msg) (gas uint64, broadcast func() error, err error) {
	if err := c.prepareBroadcast(ctx, accountName, msgs...); err != nil {
		return 0, nil, err
	}

	// calculate the necessary gas for the transaction
	txf, err := tx.PrepareFactory(c.cosmos.Context, c.cosmos.Factory)
	if err != nil {
		return 0, nil, err
	}
	_, gas, err = tx.CalculateGas(c.cosmos.Context.QueryWithData, txf, msgs...)
	if err != nil {
		return 0, nil, err
	}
	// the simulated gas can vary from the actual gas needed for a real transaction
	// we add an additional amount to endure sufficient gas is provided
	gas += 10000
	txf = txf.WithGas(gas)

	// Return the provision function
	return gas, func() error {
		// broadcast tx.
		return handleBroadcastResult(c.cosmos.BroadcastTx(accountName, msgs...))
	}, nil
}
