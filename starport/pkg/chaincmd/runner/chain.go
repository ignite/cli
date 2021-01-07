package chaincmdrunner

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
)

// Start starts the blockchain.
func (r Runner) Start(ctx context.Context, args ...string) error {
	return r.run(ctx, runOptions{longRunning: true}, r.cc.StartCommand(args...))
}

// LaunchpadStartRestServer start launchpad rest server.
func (r Runner) LaunchpadStartRestServer(ctx context.Context, apiAddress, rpcAddress string) error {
	return r.run(ctx, runOptions{longRunning: true}, r.cc.LaunchpadRestServerCommand(apiAddress, rpcAddress))
}

// Init inits the blockchain.
func (r Runner) Init(ctx context.Context, moniker string) error {
	return r.run(ctx, runOptions{}, r.cc.InitCommand(moniker))
}

// KV holds a key, value pair.
type KV struct {
	key   string
	value string
}

// NewKV returns a new key, value pair.
func NewKV(key, value string) KV {
	return KV{key, value}
}

// LaunchpadSetConfigs updates configurations for a launchpad app.
func (r Runner) LaunchpadSetConfigs(ctx context.Context, kvs ...KV) error {
	for _, kv := range kvs {
		if err := r.run(ctx, runOptions{}, r.cc.LaunchpadSetConfigCommand(kv.key, kv.value)); err != nil {
			return err
		}
	}
	return nil
}

var gentxRe = regexp.MustCompile(`(?m)"(.+?)"`)

// Gentx generates a genesis tx carrying a self delegation.
func (r Runner) Gentx(ctx context.Context, validatorName, selfDelegation string, options ...chaincmd.GentxOption) (gentxPath string, err error) {
	b := &bytes.Buffer{}

	// note: launchpad outputs from stderr.
	if err := r.run(ctx, runOptions{stdout: b, stderr: b}, r.cc.GentxCommand(validatorName, selfDelegation, options...)); err != nil {
		return "", err
	}

	return gentxRe.FindStringSubmatch(b.String())[1], nil
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
	nodeID = strings.TrimSpace(b.String())
	return
}
