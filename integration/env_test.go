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
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/availableport"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
)

const (
	serveTimeout = time.Minute * 15
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
	ctx                    context.Context
	shouldErr, shouldRetry bool
	stdout, stderr         io.Writer
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

// ExecSterr captures stderr of an execution.
func ExecStderr(w io.Writer) execOption {
	return func(o *execOptions) {
		o.stderr = w
	}
}

// ExecRetry retries command until it is successful before context is canceled.
func ExecRetry() execOption {
	return func(o *execOptions) {
		o.shouldRetry = true
	}
}

// Exec executes a command step with options where msg describes the expectation from the test.
// unless calling with Must(), Exec() will not exit test runtime on failure.
func (e env) Exec(msg string, steps step.Steps, options ...execOption) (ok bool) {
	opts := &execOptions{
		ctx:    e.ctx,
		stdout: ioutil.Discard,
		stderr: ioutil.Discard,
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
		cmdrunner.DefaultStderr(io.MultiWriter(stderr, opts.stderr)),
	}
	if isCI {
		copts = append(copts, cmdrunner.EndSignal(os.Kill))
	}
	err := cmdrunner.
		New(copts...).
		Run(opts.ctx, steps...)
	if err == context.Canceled {
		err = nil
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if opts.shouldRetry && opts.ctx.Err() == nil {
			time.Sleep(time.Second)
			return e.Exec(msg, steps, options...)
		}
	}
	if err != nil {
		msg = fmt.Sprintf("%s\n\nLogs:\n\n%s\n\nError Logs:\n\n%s\n",
			msg,
			stdout.String(),
			stderr.String())
	}
	if opts.shouldErr {
		return assert.Error(e.t, err, msg)
	}
	return assert.NoError(e.t, err, msg)
}

const (
	Launchpad = "launchpad"
	Stargate  = "stargate"
)

// Scaffold scaffolds an app to a unique appPath and returns it.
func (e env) Scaffold(appName, sdkVersion string) (appPath string) {
	root := e.TmpDir()
	e.Exec("scaffold a launchpad app",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"app",
				fmt.Sprintf("github.com/test/%s", appName),
				"--sdk-version",
				sdkVersion,
			),
			step.Workdir(root),
		)),
	)
	return filepath.Join(root, appName)
}

// Serve serves an application lives under path with options where msg describes the
// expection from the serving action.
// unless calling with Must(), Serve() will not exit test runtime on failure.
func (e env) Serve(msg string, path string, options ...execOption) (ok bool) {
	return e.Exec(msg,
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"serve",
				"-v",
			),
			step.Workdir(path),
		)),
		options...,
	)
}

// EnsureAppIsSteady ensures that app living at the path can compile and its tests
// are passing.
func (e env) EnsureAppIsSteady(appPath string) {
	e.Exec("make sure app is steady",
		step.NewSteps(step.New(
			step.Exec("go", "test", "./..."),
			step.Workdir(appPath),
		)),
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
// its config.yml and returns new values.
func (e env) RandomizeServerPorts(path string) starportconf.Servers {
	// generate random server ports and servers list.
	ports, err := availableport.Find(7)
	require.NoError(e.t, err)

	genAddr := func(port int) string {
		return fmt.Sprintf("localhost:%d", port)
	}

	servers := starportconf.Servers{
		RPCAddr:      genAddr(ports[0]),
		P2PAddr:      genAddr(ports[1]),
		ProfAddr:     genAddr(ports[2]),
		GRPCAddr:     genAddr(ports[3]),
		APIAddr:      genAddr(ports[4]),
		FrontendAddr: genAddr(ports[5]),
		DevUIAddr:    genAddr(ports[6]),
	}

	// update config.yml with the generated servers list.
	configyml, err := os.OpenFile(filepath.Join(path, "config.yml"), os.O_RDWR, 0755)
	require.NoError(e.t, err)
	defer configyml.Close()

	var conf starportconf.Config
	require.NoError(e.t, yaml.NewDecoder(configyml).Decode(&conf))

	conf.Servers = servers
	require.NoError(e.t, configyml.Truncate(0))
	_, err = configyml.Seek(0, 0)
	require.NoError(e.t, err)
	require.NoError(e.t, yaml.NewEncoder(configyml).Encode(conf))

	return servers
}

// Must fails the immediately if not ok.
// t.Fail() needs to be called for the failing tests before running Must().
func (e env) Must(ok bool) {
	if !ok {
		e.t.FailNow()
	}
}

// Home returns user's home dir.
func (e env) Home() string {
	home, err := os.UserHomeDir()
	require.NoError(e.t, err)
	return home
}

// AppHome returns appd's home dir.
func (e env) AppdHome(name, sdkVersion string) string {
	switch sdkVersion {
	case Stargate:
		return filepath.Join(e.Home(), fmt.Sprintf("%sd", name))
	case Launchpad:
		return filepath.Join(e.Home(), fmt.Sprintf(".%sd", name))
	}
	return ""
}
