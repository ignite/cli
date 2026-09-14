package gno

import (
	"fmt"
	"path/filepath"

	"github.com/gnolang/gno/gno.land/pkg/gnoland/ugnot"
	"github.com/gnolang/gno/gno.land/pkg/integration"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"
	core_types "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"github.com/gnolang/gno/tm2/pkg/commands"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/client"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// DefaultGasWanted and DefaultGasFee are the defaults used when deploying.
const (
	DefaultGasWanted = 5_000_000
	DefaultGasFee    = 1000000 // in ugnot
)

// deployPlan is a resolved deployment: everything needed to sign and
// broadcast the addpkg transaction.
type deployPlan struct {
	pkgPath    string
	from       string
	passphrase string
	tx         std.Tx
	maketxCfg  *client.MakeTxCfg
}

// DeployOptions configures Deploy.
type DeployOptions struct {
	// From is the key name (or bech32 address) signing the deployment.
	From string
	// PkgPath overrides the module path read from gnomod.toml.
	PkgPath string
	// Remote is the chain RPC address (default: 127.0.0.1:26657).
	Remote string
	// ChainID of the target chain (default: dev).
	ChainID string
	// GasWanted / GasFee for the tx.
	GasWanted int64
	GasFee    string
	// Passphrase unlocks the key (empty for dev keys created without one).
	Passphrase string
	// MaxDeposit is the max storage deposit for the package.
	MaxDeposit string
}

// signAndBroadcast is the seam used by tests to stub the real tx signing
// and broadcasting.
var signAndBroadcast = func(plan deployPlan) (*core_types.ResultBroadcastTxCommit, error) {
	return client.SignAndBroadcastHandler(plan.maketxCfg, plan.from, plan.tx, plan.passphrase, commands.NewDefaultIO())
}

// Deploy publishes the gno package at dir to a gno.land chain. The module
// path is read from dir/gnomod.toml unless overridden.
func Deploy(dir string, opts DeployOptions) error {
	plan, err := newDeployPlan(dir, opts)
	if err != nil {
		return err
	}

	bres, err := signAndBroadcast(plan)
	if err != nil {
		return fmt.Errorf("broadcasting deployment tx: %w", err)
	}
	if bres.CheckTx.IsErr() {
		return fmt.Errorf("check tx failed: %w, log: %s", bres.CheckTx.Error, bres.CheckTx.Log)
	}
	if bres.DeliverTx.IsErr() {
		return fmt.Errorf("deliver tx failed: %w, log: %s", bres.DeliverTx.Error, bres.DeliverTx.Log)
	}

	fmt.Printf("🚀 deployed %s (gas used: %d)\n", plan.pkgPath, bres.DeliverTx.GasUsed)
	return nil
}

// newDeployPlan resolves options, reads gnomod.toml, ensures the signing
// key exists (auto-importing the well-known dev account when missing) and
// builds the unsigned addpkg tx.
func newDeployPlan(dir string, opts DeployOptions) (deployPlan, error) {
	if opts.Remote == "" {
		opts.Remote = "127.0.0.1:26657"
	}
	if opts.ChainID == "" {
		opts.ChainID = "dev"
	}
	if opts.GasWanted == 0 {
		opts.GasWanted = DefaultGasWanted
	}
	if opts.GasFee == "" {
		opts.GasFee = fmt.Sprintf("%d%s", DefaultGasFee, ugnot.Denom)
	}
	if opts.From == "" {
		opts.From = integration.DefaultAccount_Name
	}

	pkgPath := opts.PkgPath
	if pkgPath == "" {
		mod, err := parseGnoModDir(dir)
		if err != nil {
			return deployPlan{}, fmt.Errorf("reading gnomod.toml in %s: %w (scaffold one with `ignite scaffold realm`)", dir, err)
		}
		pkgPath = mod
	}

	memPkg := gno.MustReadMemPackage(dir, pkgPath, gno.MPUserAll)
	if memPkg.IsEmpty() {
		return deployPlan{}, fmt.Errorf("no .gno files found in %s", dir)
	}

	gasFee, err := std.ParseCoin(opts.GasFee)
	if err != nil {
		return deployPlan{}, fmt.Errorf("parsing gas fee: %w", err)
	}
	var maxDeposit std.Coins
	if opts.MaxDeposit != "" {
		if maxDeposit, err = std.ParseCoins(opts.MaxDeposit); err != nil {
			return deployPlan{}, fmt.Errorf("parsing max deposit: %w", err)
		}
	}

	info, err := ensureKey(opts.From)
	if err != nil {
		return deployPlan{}, fmt.Errorf("key %q not found in the gno keybase: %w (create one with `ignite account create`)", opts.From, err)
	}

	return deployPlan{
		pkgPath:    pkgPath,
		from:       opts.From,
		passphrase: opts.Passphrase,
		tx: std.Tx{
			Msgs: []std.Msg{
				vm.MsgAddPackage{
					Creator:    crypto.MustAddressFromString(info.Address),
					Package:    memPkg,
					MaxDeposit: maxDeposit,
				},
			},
			Fee: std.NewFee(opts.GasWanted, gasFee),
		},
		maketxCfg: &client.MakeTxCfg{
			RootCfg: &client.BaseCfg{
				BaseOptions: client.BaseOptions{
					Home:   HomeDir(),
					Remote: opts.Remote,
					Quiet:  true,
				},
			},
			GasWanted: opts.GasWanted,
			GasFee:    opts.GasFee,
			Broadcast: true,
			ChainID:   opts.ChainID,
		},
	}, nil
}

// ensureKey returns the key info for name, auto-importing the well-known
// dev account when it is missing (its mnemonic is public).
func ensureKey(name string) (KeyInfo, error) {
	info, err := ShowKey(name)
	if err == nil {
		return info, nil
	}
	if name != integration.DefaultAccount_Name {
		return KeyInfo{}, err
	}
	return RecoverKey(integration.DefaultAccount_Name, integration.DefaultAccount_Seed, "", 0, 0)
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
