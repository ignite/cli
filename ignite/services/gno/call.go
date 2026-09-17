package gno

import (
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/cli/v30/ignite/pkg/errors"
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
func Call(opts CallOptions) (*BroadcastResult, error) {
	opts.TxBaseOptions = opts.TxBaseOptions.withDefaults()

	if opts.PkgPath == "" {
		return nil, errors.Errorf("package path is required")
	}
	if opts.Func == "" {
		return nil, errors.Errorf("function name is required")
	}

	var send std.Coins
	if opts.Send != "" {
		var err error
		if send, err = std.ParseCoins(opts.Send); err != nil {
			return nil, errors.Errorf("parsing send coins: %w", err)
		}
	}

	caller, err := opts.callerAddress()
	if err != nil {
		return nil, err
	}
	gasFee, err := opts.parseGasFee()
	if err != nil {
		return nil, err
	}

	tx := std.Tx{
		Msgs: []std.Msg{
			vm.MsgCall{
				Caller:  caller,
				PkgPath: opts.PkgPath,
				Func:    opts.Func,
				Args:    opts.Args,
				Send:    send,
			},
		},
		Fee: std.NewFee(opts.GasWanted, gasFee),
	}

	return broadcast(opts.newTxPlan(tx))
}
