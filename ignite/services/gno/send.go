package gno

import (
	"fmt"
	"strings"

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
func Send(opts SendOptions) error {
	opts.TxBaseOptions = opts.TxBaseOptions.withDefaults()

	if opts.To == "" {
		return errors.Errorf("beneficiary is required")
	}
	coins, err := std.ParseCoins(opts.Amount)
	if err != nil {
		return errors.Errorf("parsing amount: %w", err)
	}
	if len(coins) == 0 {
		return errors.Errorf("amount is required")
	}

	// resolve beneficiary: allow key names for convenience
	to := opts.To
	if !strings.HasPrefix(opts.To, "g1") {
		info, err := ShowKey(opts.To)
		if err != nil {
			return errors.Errorf("%q is neither a bech32 address nor a key in the keybase: %w", opts.To, err)
		}
		to = info.Address
	}

	fromAddr, err := opts.callerAddress()
	if err != nil {
		return err
	}
	gasFee, err := opts.parseGasFee()
	if err != nil {
		return err
	}

	tx := std.Tx{
		Msgs: []std.Msg{
			bank.MsgSend{
				FromAddress: crypto.MustAddressFromString(fromAddr),
				ToAddress:   crypto.MustAddressFromString(to),
				Amount:      coins,
			},
		},
		Fee: std.NewFee(opts.GasWanted, gasFee),
	}

	return broadcast(opts.newTxPlan(tx), fmt.Sprintf("💸 sent %s to %s", opts.Amount, to))
}
