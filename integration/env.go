package envtest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
	"github.com/ignite/cli/ignite/pkg/env"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/httpstatuschecker"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

const (
	ConfigYML = "config.yml"
)

var (
	// IgniteApp hold the location of the ignite binary used in the integration
	// tests. The binary is compiled the first time the env.New() function is
	// invoked.
	IgniteApp = path.Join(os.TempDir(), "ignite-tests", "ignite")

	isCI, _           = strconv.ParseBool(os.Getenv("CI"))
	compileBinaryOnce sync.Once
)

// Env provides an isolated testing environment and what's needed to
// make it possible.
type Env struct {
	t   *testing.T
	ctx context.Context
}

// New creates a new testing environment.
func New(t *testing.T) Env {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	e := Env{
		t:   t,
		ctx: ctx,
	}
	// To avoid conflicts with the default config folder located in $HOME, we
	// set an other one thanks to env var.
	cfgDir := path.Join(t.TempDir(), ".ignite")
	env.SetConfigDir(cfgDir)

	t.Cleanup(cancel)
	compileBinaryOnce.Do(func() {
		compileBinary(ctx)
	})
	return e
}

func compileBinary(ctx context.Context) {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to get working dir: %v", err))
	}
	_, appPath, err := gomodulepath.Find(wd)
	if err != nil {
		panic(fmt.Sprintf("unable to read go module path: %v", err))
	}
	var (
		output, binary = filepath.Split(IgniteApp)
		path           = path.Join(appPath, "ignite", "cmd", "ignite")
	)
	err = gocmd.BuildPath(ctx, output, binary, path, nil)
	if err != nil {
		panic(fmt.Sprintf("error while building binary: %v", err))
	}
}

func (e Env) T() *testing.T {
	return e.t
}

// SetCleanup registers a function to be called when the test (or subtest) and all its
// subtests complete.
func (e Env) SetCleanup(f func()) {
	e.t.Cleanup(f)
}

// Ctx returns parent context for the test suite to use for cancelations.
func (e Env) Ctx() context.Context {
	return e.ctx
}

// IsAppServed checks that app is served properly and servers are started to listening before ctx canceled.
func (e Env) IsAppServed(ctx context.Context, apiAddr string) error {
	checkAlive := func() error {
		addr, err := xurl.HTTP(apiAddr)
		if err != nil {
			return err
		}

		ok, err := httpstatuschecker.Check(ctx, fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/node_info", addr))
		if err == nil && !ok {
			err = errors.New("waiting for app")
		}
		if HasTestVerboseFlag() {
			fmt.Printf("IsAppServed at %s: %v\n", addr, err)
		}
		return err
	}

	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}

// IsFaucetServed checks that faucet of the app is served properly.
func (e Env) IsFaucetServed(ctx context.Context, faucetClient cosmosfaucet.HTTPClient) error {
	checkAlive := func() error {
		_, err := faucetClient.FaucetInfo(ctx)
		return err
	}

	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}

// TmpDir creates a new temporary directory.
func (e Env) TmpDir() (path string) {
	return e.t.TempDir()
}

// Home returns user's home dir.
func (e Env) Home() string {
	home, err := os.UserHomeDir()
	require.NoError(e.t, err)
	return home
}

// AppHome returns app's root home/data dir path.
func (e Env) AppHome(name string) string {
	return filepath.Join(e.Home(), fmt.Sprintf(".%s", name))
}

// Must fails the immediately if not ok.
// t.Fail() needs to be called for the failing tests before running Must().
func (e Env) Must(ok bool) {
	if !ok {
		e.t.FailNow()
	}
}

func (e Env) HasFailed() bool {
	return e.t.Failed()
}

func (e Env) RequireExpectations() {
	e.Must(e.HasFailed())
}

func Contains(s, partial string) bool {
	return strings.Contains(s, strings.TrimSpace(partial))
}

func HasTestVerboseFlag() bool {
	return flag.Lookup("test.v").Value.String() == "true"
}
