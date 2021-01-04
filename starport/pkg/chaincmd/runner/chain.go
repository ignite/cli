package chaincmdrunner

import (
	"bytes"
	"context"
	"strings"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
)

// Start starts the blockchain.
func (r Runner) Start(ctx context.Context, args []string) error {
	return r.run(ctx, runOptions{longRunning: true}, r.cc.StartCommand(args...))
}

// Init inits the blockchain.
func (r Runner) Init(ctx context.Context, moniker string) error {
	return r.run(ctx, runOptions{}, r.cc.InitCommand(moniker))
}

// Gentx generates a genesis tx carrying a self delegation.
func (r Runner) Gentx(ctx context.Context, validatorName, selfDelegation string, options ...chaincmd.GentxOption) error {
	return r.run(ctx, runOptions{}, r.cc.GentxCommand(validatorName, selfDelegation, options...))
}

// CollectGentxs collects gentxs.
func (r Runner) CollectGentxs(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.cc.CollectGentxsCommand())
}

// ValidateGenesis validates genesis.
func (r Runner) ValidateGenesis(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.cc.ValidateGenesisCommand())
}

// ShowNodeID shows node id.
func (r Runner) ShowNodeID(ctx context.Context) (nodeID string, err error) {
	b := &bytes.Buffer{}
	err = r.run(ctx, runOptions{stdout: b}, r.cc.ShowNodeIDCommand())
	return strings.TrimSpace(b.String()), err
}
