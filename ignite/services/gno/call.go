package gno

import (
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// CallOptions configures Call.
type CallOptions struct {
	TxBaseOptions
	// PkgPath is the realm package path (required).
	PkgPath string
	// Func is the function name to call (required).
	Func string
	// Args are the function arguments (gno literal expressions as strings).
	Args []string
	// Send attaches coins to the call.
	Send string
}

// Call invokes a realm function on a gno.land chain and waits for the
// result.
func Call(opts CallOptions) error {
	opts.TxBaseOptions = opts.TxBaseOptions.withDefaults()

	if opts.PkgPath == "" {
		return errors.Errorf("package path is required")
	}
	if opts.Func == "" {
		return errors.Errorf("function name is required")
	}

	var send std.Coins
	if opts.Send != "" {
		var err error
		if send, err = std.ParseCoins(opts.Send); err != nil {
			return errors.Errorf("parsing send coins: %w", err)
		}
	}

	addr, err := opts.callerAddress()
	if err != nil {
		return err
	}
	gasFee, err := opts.parseGasFee()
	if err != nil {
		return err
	}

	tx := std.Tx{
		Msgs: []std.Msg{
			vm.MsgCall{
				Caller:  crypto.MustAddressFromString(addr),
				PkgPath: opts.PkgPath,
				Func:    opts.Func,
				Args:    opts.Args,
				Send:    send,
			},
		},
		Fee: std.NewFee(opts.GasWanted, gasFee),
	}

	return broadcast(opts.newTxPlan(tx), fmt.Sprintf("📣 called %s.%s", opts.PkgPath, opts.Func))
}
