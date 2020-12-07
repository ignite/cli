package spn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/manifoldco/promptui"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var spn = "spn"
var homedir = os.ExpandEnv("$HOME/spnd")

const (
	faucetDenom     = "token"
	faucetMinAmount = 100
)

// Client is client to interact with SPN.
type Client struct {
	kr            keyring.Keyring
	factory       tx.Factory
	clientCtx     client.Context
	apiAddress    string
	faucetAddress string
	out           *bytes.Buffer
}

type options struct {
	keyringBackend string
}

// Option configures Client options.
type Option func(*options)

// New creates a new SPN Client with nodeAddress of a full SPN node.
// by default, OS is used as keyring backend.
func New(nodeAddress, apiAddress, faucetAddress string, option ...Option) (*Client, error) {
	opts := &options{
		keyringBackend: keyring.BackendOS,
	}
	for _, o := range option {
		o(opts)
	}
	kr, err := keyring.New(types.KeyringServiceName(), opts.keyringBackend, homedir, os.Stdin)
	if err != nil {
		return nil, err
	}

	client, err := rpchttp.New(nodeAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	clientCtx := NewClientCtx(kr, client, out)
	factory := NewFactory(clientCtx)
	return &Client{
		kr:            kr,
		factory:       factory,
		clientCtx:     clientCtx,
		apiAddress:    apiAddress,
		faucetAddress: faucetAddress,
		out:           out,
	}, nil
}

func (c *Client) buildClientCtx(accountName string) (client.Context, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return client.Context{}, err
	}
	return c.clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress()), nil
}

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

	// request amounts from faucet.
	body, err := json.Marshal(struct {
		Address string `json:"address"`
	}{address})
	if err != nil {
		return err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.faucetAddress, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("faucet server request failed: %v", resp.Status)
	}

	var result struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Status != "ok" {
		return fmt.Errorf("cannot retrieve tokens from faucet: %s", result.Error)
	}

	return nil
}

func (c *Client) broadcast(ctx context.Context, clientCtx client.Context, confirmPrompt bool, msgs ...types.Msg) error {
	// validate msgs.
	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
	}

	// make sure that account has enough balances before broadcasting.
	if err := c.makeSureAccountHasTokens(ctx, clientCtx.GetFromAddress().String()); err != nil {
		return err
	}

	c.out.Reset()

	// calculate the necessary gas for the transaction
	txf, err := tx.PrepareFactory(clientCtx, c.factory)
	if err != nil {
		return err
	}
	_, gas, err := tx.CalculateGas(clientCtx.QueryWithData, txf, msgs...)
	if err != nil {
		return nil
	}
	// The simulated gas can vary from the actual gas needed for a real transaction
	// We add an additional amount to endure sufficient gas is provided
	gas += 10000
	txf = txf.WithGas(gas)

	// Prompt for confirmation
	if confirmPrompt {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("This operation will cost about %v gas. Confirm the transaction?",
				gas,
			),
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err != nil {
			return errors.New("transaction aborted")
		}
	}

	// broadcast tx.
	if err := tx.BroadcastTx(clientCtx, txf, msgs...); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("make sure that your SPN account has enough balance")
		}
		return err
	}

	// handle results.
	out := struct {
		Code int    `json:"code"`
		Log  string `json:"raw_log"`
	}{}
	if err := json.NewDecoder(c.out).Decode(&out); err != nil {
		return err
	}
	if out.Code > 0 {
		return fmt.Errorf("SPN error with '%d' code: %s", out.Code, out.Log)
	}
	return nil
}