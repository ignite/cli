package gno

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnolang/gno/gno.land/pkg/integration"
	core_types "github.com/gnolang/gno/tm2/pkg/bft/rpc/core/types"
	"gotest.tools/v3/assert"

	"bytes"
	"context"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNewDeployPlan(t *testing.T) {
	withTestGnoHome(t)

	dir, modulePath, err := Scaffold(KindRealm, "deployplan", ScaffoldOptions{Path: filepath.Join(t.TempDir(), "deployplan")})
	assert.NilError(t, err)

	// exercise Deploy through the stubbed broadcast seam and inspect the plan.
	var got txPlan
	restore := stubBroadcast(func(plan txPlan) (*core_types.ResultBroadcastTxCommit, error) {
		got = plan
		return &core_types.ResultBroadcastTxCommit{}, nil
	})
	defer restore()

	assert.NilError(t, Deploy(dir, DeployOptions{}))
	assert.Equal(t, modulePath, "gno.land/r/deployplan")
	assert.Equal(t, integration.DefaultAccount_Name, got.from, "default signing key is the dev account")
	assert.Equal(t, "dev", got.maketxCfg.ChainID)
	assert.Equal(t, gnoDefaultRemoteForTest, got.maketxCfg.RootCfg.Remote)
	assert.Equal(t, int64(DefaultGasWanted), got.maketxCfg.GasWanted)
	assert.Equal(t, 1, len(got.tx.Msgs))
}

const gnoDefaultRemoteForTest = "127.0.0.1:26657"

func TestDeployErrors(t *testing.T) {
	tt := []struct {
		name string
		dir  string
		opts DeployOptions
	}{
		{name: "no gnomod", dir: t.TempDir(), opts: DeployOptions{}},
		{name: "invalid gas fee", opts: DeployOptions{TxBaseOptions: TxBaseOptions{GasFee: "nope"}}},
		{name: "unknown key", opts: DeployOptions{TxBaseOptions: TxBaseOptions{From: "no-such-key"}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			withTestGnoHome(t)
			if tc.name == "invalid gas fee" {
				tc.dir = scaffoldedRealm(t, "gasfee")
			}
			assert.ErrorContains(t, Deploy(tc.dir, tc.opts), "")
		})
	}
}

// scaffoldedRealm scaffolds a realm and returns its dir.
func scaffoldedRealm(t *testing.T, name string) string {
	t.Helper()
	dir, _, err := Scaffold(KindRealm, name, ScaffoldOptions{Path: filepath.Join(t.TempDir(), name)})
	assert.NilError(t, err)
	return dir
}

// stubBroadcast swaps the signAndBroadcast seam for the test and returns a
// restore function.
func stubBroadcast(stub func(txPlan) (*core_types.ResultBroadcastTxCommit, error)) func() {
	orig := signAndBroadcast
	signAndBroadcast = stub
	return func() { signAndBroadcast = orig }
}

func TestWatchLoopReloadsOnGnoChange(t *testing.T) {
	dir := t.TempDir()
	assert.NilError(t, os.WriteFile(filepath.Join(dir, "realm.gno"), []byte("package x\n"), 0o644))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var reloads atomic.Int32
	out := &bytes.Buffer{}
	errCh := make(chan error, 1)
	go func() {
		errCh <- watchLoop(ctx, func(context.Context) error {
			reloads.Add(1)
			return nil
		}, dir, out)
	}()

	// wait for the watcher to be registered, then fire a change
	time.Sleep(200 * time.Millisecond)
	assert.NilError(t, os.WriteFile(filepath.Join(dir, "realm.gno"), []byte("package x\n// changed\n"), 0o644))

	deadline := time.Now().Add(5 * time.Second)
	for reloads.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	cancel()
	assert.NilError(t, <-errCh)
	assert.Assert(t, reloads.Load() >= 1, "expected at least one reload, out: %s", out.String())
	assert.Assert(t, bytes.Contains(out.Bytes(), []byte("🔁")), "reload should be reported: %s", out.String())
}

func TestReloadOnChangeIgnoresNonGnoFiles(t *testing.T) {
	out := &bytes.Buffer{}
	var calls atomic.Int32
	reloadOnChange(context.Background(), func(context.Context) error {
		calls.Add(1)
		return nil
	}, fsnotify.Event{Name: "/x/main.go", Op: fsnotify.Write}, out)
	reloadOnChange(context.Background(), func(context.Context) error {
		calls.Add(1)
		return nil
	}, fsnotify.Event{Name: "/x/realm.gno", Op: fsnotify.Chmod}, out)
	assert.Equal(t, int32(0), calls.Load())
}
