package gno

import (
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/gnoland/ugnot"
	core_types "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// TxBaseOptions holds the chain connection and gas settings shared by all
// transaction-building commands (deploy, call, send).
type TxBaseOptions struct {
	// From is the key name (or bech32 address) signing the tx.
	From string
	// Remote is the chain RPC address (default: 127.0.0.1:26657).
	Remote string
	// ChainID of the target chain (default: dev).
	ChainID string
	// GasWanted / GasFee for the tx.
	GasWanted int64
	GasFee    string
	// Passphrase unlocks the signing key (empty for dev keys created without one).
	Passphrase string
}

// withDefaults fills unset tx fields with dev-chain friendly defaults.
func (b TxBaseOptions) withDefaults() TxBaseOptions {
	if b.Remote == "" {
		b.Remote = "127.0.0.1:26657"
	}
	if b.ChainID == "" {
		b.ChainID = "dev"
	}
	if b.GasWanted == 0 {
		b.GasWanted = DefaultGasWanted
	}
	if b.GasFee == "" {
		b.GasFee = fmt.Sprintf("%d%s", DefaultGasFee, ugnot.Denom)
	}
	if b.From == "" {
		b.From = defaultDevAccountName
	}
	return b
}

// txPlan is a resolved transaction ready to sign and broadcast.
type txPlan struct {
	tx        std.Tx
	maketxCfg *client.MakeTxCfg
	from      string
	pass      string
}

// newTxPlan builds the signing config for the given unsigned tx.
func (b TxBaseOptions) newTxPlan(tx std.Tx) txPlan {
	return txPlan{
		tx: tx,
		maketxCfg: &client.MakeTxCfg{
			RootCfg: &client.BaseCfg{
				BaseOptions: client.BaseOptions{
					Home:   HomeDir(),
					Remote: b.Remote,
					Quiet:  true,
				},
			},
			GasWanted: b.GasWanted,
			GasFee:    b.GasFee,
			Broadcast: true,
			ChainID:   b.ChainID,
		},
		from: b.From,
		pass: b.Passphrase,
	}
}

// signAndBroadcast is the seam used by tests to stub the real tx signing
// and broadcasting.
var signAndBroadcast = func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
	return client.SignAndBroadcastHandler(plan.maketxCfg, plan.from, plan.tx, plan.pass, commands.NewDefaultIO())
}

// broadcast signs and broadcasts plan, returning a formatted result or an
// error with the check/deliver failure details.
func broadcast(plan txPlan, describe string) error {
	bres, err := signAndBroadcast(plan)
	if err != nil {
		return fmt.Errorf("broadcasting tx: %w", err)
	}
	if bres.CheckTx.IsErr() {
		return fmt.Errorf("check tx failed: %w, log: %s", bres.CheckTx.Error, bres.CheckTx.Log)
	}
	if bres.DeliverTx.IsErr() {
		return fmt.Errorf("deliver tx failed: %w, log: %s", bres.DeliverTx.Error, bres.DeliverTx.Log)
	}
	fmt.Printf("%s (gas used: %d)\n", describe, bres.DeliverTx.GasUsed)
	return nil
}

// callerAddress resolves the address of the signing key.
func (b TxBaseOptions) callerAddress() (string, error) {
	info, err := ensureKey(b.From)
	if err != nil {
		return "", fmt.Errorf("key %q not found in the gno keybase: %w (create one with `ignite account create`)", b.From, err)
	}
	return info.Address, nil
}

// parseGasFee parses the configured gas fee coin.
func (b TxBaseOptions) parseGasFee() (std.Coin, error) {
	gasFee, err := std.ParseCoin(b.GasFee)
	if err != nil {
		return std.Coin{}, fmt.Errorf("parsing gas fee: %w", err)
	}
	return gasFee, nil
}
