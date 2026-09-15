package gno

import (
	"fmt"
	"path/filepath"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// DefaultGasWanted and DefaultGasFee are the defaults used when deploying.
const (
	DefaultGasWanted = 5_000_000
	DefaultGasFee    = 1000000 // in ugnot
)

// DeployOptions configures Deploy.
type DeployOptions struct {
	TxBaseOptions
	// PkgPath overrides the module path read from gnomod.toml.
	PkgPath string
	// MaxDeposit is the max storage deposit for the package.
	MaxDeposit string
}

// Deploy publishes the gno package at dir to a gno.land chain. The module
// path is read from dir/gnomod.toml unless overridden.
func Deploy(dir string, opts DeployOptions) error {
	opts.TxBaseOptions = opts.TxBaseOptions.withDefaults()

	pkgPath := opts.PkgPath
	if pkgPath == "" {
		mod, err := parseGnoModDir(dir)
		if err != nil {
			return errors.Errorf("reading gnomod.toml in %s: %w (scaffold one with `ignite scaffold realm`)", dir, err)
		}
		pkgPath = mod
	}

	memPkg := gno.MustReadMemPackage(dir, pkgPath, gno.MPUserAll)
	if memPkg.IsEmpty() {
		return errors.Errorf("no .gno files found in %s", dir)
	}

	var maxDeposit std.Coins
	if opts.MaxDeposit != "" {
		var err error
		if maxDeposit, err = std.ParseCoins(opts.MaxDeposit); err != nil {
			return errors.Errorf("parsing max deposit: %w", err)
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
			vm.MsgAddPackage{
				Creator:    crypto.MustAddressFromString(addr),
				Package:    memPkg,
				MaxDeposit: maxDeposit,
			},
		},
		Fee: std.NewFee(opts.GasWanted, gasFee),
	}

	return broadcast(opts.newTxPlan(tx), fmt.Sprintf("🚀 deployed %s", pkgPath))
}

// parseGnoModDir reads the module path from dir/gnomod.toml.
func parseGnoModDir(dir string) (string, error) {
	return parseGnoMod(filepath.Join(dir, "gnomod.toml"))
}

// parseGnoMod reads the module path from a gnomod.toml file.
func parseGnoMod(path string) (string, error) {
	mod, err := gnomod.ParseFilepath(path)
	if err != nil {
		return "", err
	}
	return mod.Module, nil
}
