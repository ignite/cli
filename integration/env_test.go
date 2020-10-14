// Package integration_test integration test Starport and scaffolded apps.
package integration_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
)

var isCI, _ = strconv.ParseBool(os.Getenv("CI"))

// env provides an isolated testing environment and what's needed to
// make it possible.
type env struct {
	t   *testing.T
	ctx context.Context
}

// env creates a new testing environment.
func newEnv(t *testing.T) env {
	ctx, cancel := context.WithCancel(context.Background())
	e := env{
		t:   t,
		ctx: ctx,
	}
	t.Cleanup(cancel)
	return e
}

// Ctx returns parent context for the test suite to use for cancelations.
func (e env) Ctx() context.Context {
	return e.ctx
}

type execOptions struct {
	ctx       context.Context
	shouldErr bool
	stdout    io.Writer
}

type execOption func(*execOptions)

// ExecShouldError sets the expectations of a command's execution to end with a failure.
func ExecShouldError() execOption {
	return func(o *execOptions) {
		o.shouldErr = true
	}
}

// ExecCtx sets cancelation context for the execution.
func ExecCtx(ctx context.Context) execOption {
	return func(o *execOptions) {
		o.ctx = ctx
	}
}

// ExecStdout captures stdout of an execution.
func ExecStdout(w io.Writer) execOption {
	return func(o *execOptions) {
		o.stdout = w
	}
}

// Exec executes a command step with options where msg describes the expectation from the test.
func (e env) Exec(msg string, step *step.Step, options ...execOption) {
	opts := &execOptions{
		ctx:    e.ctx,
		stdout: ioutil.Discard,
	}
	for _, o := range options {
		o(opts)
	}
	var (
		stdout = &bytes.Buffer{}
		stderr = &bytes.Buffer{}
	)
	copts := []cmdrunner.Option{
		cmdrunner.DefaultStdout(io.MultiWriter(stdout, opts.stdout)),
		cmdrunner.DefaultStderr(stderr),
	}
	if isCI {
		copts = append(copts, cmdrunner.EndSignal(os.Kill))
	}
	err := cmdrunner.
		New(copts...).
		Run(opts.ctx, step)
	if err == context.Canceled {
		err = nil
	}
	if err != nil {
		msg = fmt.Sprintf("%s\n\nLogs:\n\n%s\n\nError Logs:\n\n%s\n",
			msg,
			stdout.String(),
			stderr.String())
	}
	if opts.shouldErr {
		require.Error(e.t, err, msg)
	} else {
		require.NoError(e.t, err, msg)
	}
}

const (
	Launchpad = "launchpad"
	Stargate  = "stargate"
)

// Scaffold scaffolds an app to a unique appPath and returns it.
func (e env) Scaffold(appName, sdkVersion string) (appPath string) {
	root := e.TmpDir()
	e.Exec("scaffold a launchpad app",
		step.New(
			step.Exec(
				"starport",
				"app",
				fmt.Sprintf("github.com/test/%s", appName),
				"--sdk-version",
				sdkVersion,
			),
			step.Workdir(root),
		),
	)
	return filepath.Join(root, appName)
}

// Serve serves an application lives under path with options where msg describes the
// expection from the serving action.
func (e env) Serve(msg string, path string, options ...execOption) {
	e.Exec(msg,
		step.New(
			step.Exec(
				"starport",
				"serve",
				"-v",
			),
			step.Workdir(path),
		),
		options...,
	)
}

// EnsureAppIsSteady ensures that app living at the path can compile and its tests
// are passing.
func (e env) EnsureAppIsSteady(appPath string) {
	e.Exec("make sure app is steady",
		step.New(
			step.Exec("go", "test", "./..."),
			step.Workdir(appPath),
		),
	)
}

// IsAppServed checks that app is served properly and servers are started to listening
// before ctx canceled.
func (e env) IsAppServed(ctx context.Context, servers starportconf.Servers) error {
	checkAlive := func() error {
		ok, err := httpstatuschecker.Check(ctx, xurl.HTTP(servers.APIAddr)+"/node_info")
		if err == nil && !ok {
			err = errors.New("app is not online")
		}
		return err
	}
	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}

// TmpDir creates a new temporary directory.
func (e env) TmpDir() (path string) {
	path, err := ioutil.TempDir("", "integration")
	require.NoError(e.t, err, "create a tmp dir")
	e.t.Cleanup(func() { os.RemoveAll(path) })
	return path
}

// RandomizeServerPorts randomizes server ports for the app at path, updates
// its config.yml and returns new values. TODO
func (e env) RandomizeServerPorts(path string) starportconf.Servers {
	return starportconf.Servers{
		APIAddr: "localhost:1317",
	}
}
