package gno

import (
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk/bank"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// SendOptions configures Send.
type SendOptions struct {
	TxBaseOptions
	// To is the beneficiary bech32 address or key name.
	To string
	// Amount is the coins to send, e.g. "10000000ugnot".
	Amount string
}

// Send transfers coins between accounts on a gno.land chain. Handy on dev
// chains to fund accounts created with `ignite account create`.
func Send(opts SendOptions) (*BroadcastResult, error) {
	opts.TxBaseOptions = opts.TxBaseOptions.withDefaults()

	if opts.To == "" {
		return nil, errors.Errorf("beneficiary is required")
	}
	coins, err := std.ParseCoins(opts.Amount)
	if err != nil {
		return nil, errors.Errorf("parsing amount: %w", err)
	}
	if len(coins) == 0 {
		return nil, errors.Errorf("amount is required")
	}

	// resolve beneficiary: allow key names for convenience
	to, err := resolveTo(opts.To)
	if err != nil {
		return nil, err
	}

	fromAddr, err := opts.callerAddress()
	if err != nil {
		return nil, err
	}
	gasFee, err := opts.parseGasFee()
	if err != nil {
		return nil, err
	}

	tx := std.Tx{
		Msgs: []std.Msg{
			bank.MsgSend{
				FromAddress: fromAddr,
				ToAddress:   to,
				Amount:      coins,
			},
		},
		Fee: std.NewFee(opts.GasWanted, gasFee),
	}

	return broadcast(opts.newTxPlan(tx))
}

// resolveTo resolves the beneficiary to an address: key names are looked up
// in the keybase, anything else must be a valid bech32 address.
func resolveTo(nameOrAddr string) (crypto.Address, error) {
	info, keyErr := ShowKey(nameOrAddr)
	if keyErr == nil {
		addr, err := crypto.AddressFromString(info.Address)
		if err != nil {
			return crypto.Address{}, errors.Errorf("invalid address stored for key %q: %w", nameOrAddr, err)
		}
		return addr, nil
	}
	addr, addrErr := crypto.AddressFromString(nameOrAddr)
	if addrErr != nil {
		return crypto.Address{}, errors.Errorf("%q is neither a bech32 address nor a key in the keybase", nameOrAddr)
	}
	return addr, nil
}
