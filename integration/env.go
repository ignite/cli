package envtest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/httpstatuschecker"
	"github.com/ignite/cli/ignite/pkg/xexec"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

const (
	ConfigYML = "config.yml"
)

var (
	isCI, _   = strconv.ParseBool(os.Getenv("CI"))
	IgniteApp string
)

func init() {
	wd, _ := os.Getwd()
	_, appPath, err := gomodulepath.Find(wd)
	if err != nil {
		panic(err)
	}
	// Build the ignite binary
	tmp, err := os.MkdirTemp("", "integration-bin")
	if err != nil {
		panic(err)
	}
	IgniteApp = path.Join(tmp, "ignite")
	command := exec.Command("go", "build",
		"-o", IgniteApp, appPath+"/ignite/cmd/ignite")
	if err := command.Run(); err != nil {
		panic(err)
	}
}

// Env provides an isolated testing environment and what's needed to
// make it possible.
type Env struct {
	t   *testing.T
	ctx context.Context
}

// New creates a new testing environment.
func New(t *testing.T) Env {
	ctx, cancel := context.WithCancel(context.Background())
	e := Env{
		t:   t,
		ctx: ctx,
	}
	t.Cleanup(cancel)

	if !xexec.IsCommandAvailable(IgniteApp) {
		t.Fatal("ignite needs to be installed")
	}

	return e
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

// IsAppServed checks that app is served properly and servers are started to listening
// before ctx canceled.
func (e Env) IsAppServed(ctx context.Context, host chainconfig.Host) error {
	checkAlive := func() error {
		addr, err := xurl.HTTP(host.API)
		if err != nil {
			return err
		}

		ok, err := httpstatuschecker.Check(ctx, fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/node_info", addr))
		if err == nil && !ok {
			err = errors.New("app is not online")
		}
		if HasTestVerboseFlag() {
			fmt.Printf("IsAppServed at %s: %v\n", addr, err)
		}
		return err
	}

	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}

// IsFaucetServed checks that faucet of the app is served properly
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
