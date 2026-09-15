//go:build !relayer

package gno_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v29/integration"
)

// freeRPCListener returns a free localhost listener address for the dev chain.
func freeRPCListener(t *testing.T) string {
	t.Helper()
	ports, err := availableport.Find(1)
	require.NoError(t, err)
	return fmt.Sprintf("tcp://127.0.0.1:%d", ports[0])
}

// waitRPC blocks until the dev chain RPC answers TCP connections.
func waitRPC(t *testing.T, ctx context.Context, listener string) error {
	t.Helper()
	addr := strings.TrimPrefix(listener, "tcp://")
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			return conn.Close()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
	return fmt.Errorf("dev chain RPC at %s did not come up in time", addr)
}

// startServe starts `ignite chain serve` in the background and kills it at
// test cleanup. It blocks until the RPC endpoint accepts connections.
func startServe(t *testing.T, env envtest.Env, workdir, listener string) {
	t.Helper()

	ctx, cancel := context.WithCancel(env.Ctx())
	t.Cleanup(cancel)

	logs, err := os.Create(filepath.Join(t.TempDir(), "serve.log"))
	require.NoError(t, err)

	cmd := exec.Command(envtest.IgniteApp, "chain", "serve", "--remote", listener)
	cmd.Dir = workdir
	cmd.Stdout = logs
	cmd.Stderr = logs
	require.NoError(t, cmd.Start())

	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})

	require.NoError(t, waitRPC(t, ctx, listener))
}

func TestGnoScaffoldRealm(t *testing.T) {
	env := envtest.New(t)
	tmp := env.TmpDir()

	env.Exec("scaffold a realm",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "realm", "counter"),
			step.Workdir(tmp),
		)),
	)

	for _, f := range []string{"gnomod.toml", "counter.gno", "counter_test.gno"} {
		_, err := os.Stat(filepath.Join(tmp, "counter", f))
		require.NoError(t, err, f)
	}

	gnomod, err := os.ReadFile(filepath.Join(tmp, "counter", "gnomod.toml"))
	require.NoError(t, err)
	require.Contains(t, string(gnomod), `module = "gno.land/r/counter"`)

	env.Exec("scaffold a realm with a full path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "realm", "gno.land/r/demo/pixel"),
			step.Workdir(tmp),
		)),
	)

	gnomod, err = os.ReadFile(filepath.Join(tmp, "pixel", "gnomod.toml"))
	require.NoError(t, err)
	require.Contains(t, string(gnomod), `module = "gno.land/r/demo/pixel"`)
}

func TestGnoScaffoldPackage(t *testing.T) {
	env := envtest.New(t)
	tmp := env.TmpDir()

	env.Exec("scaffold a package",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "package", "utils"),
			step.Workdir(tmp),
		)),
	)

	gnomod, err := os.ReadFile(filepath.Join(tmp, "utils", "gnomod.toml"))
	require.NoError(t, err)
	require.Contains(t, string(gnomod), `module = "gno.land/p/utils"`)
}

func TestGnoScaffoldFailsOnInvalidName(t *testing.T) {
	env := envtest.New(t)
	tmp := env.TmpDir()

	env.Exec("scaffold with an invalid name fails",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "realm", "Invalid!Name"),
			step.Workdir(tmp),
		)),
		envtest.ExecShouldError(),
	)
}

func TestGnoAccountLifecycle(t *testing.T) {
	var (
		env  = envtest.New(t)
		tmp  = env.TmpDir()
		args = []string{"--home", tmp}
	)

	env.Exec("create an account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{"account", "create", "alice"}, args...)...),
		)),
	)

	var listOut strings.Builder
	env.Exec("list accounts",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{"account", "list"}, args...)...),
			step.Stdout(&listOut),
		)),
	)
	require.Contains(t, listOut.String(), "alice")
	require.Contains(t, listOut.String(), "g1") // bech32 address prefix

	env.Exec("show the account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{"account", "show", "alice"}, args...)...),
		)),
	)

	env.Exec("export the account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{
				"account", "export", "alice", "--output", filepath.Join(tmp, "alice.asc"),
			}, args...)...),
		)),
	)

	env.Exec("delete the account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{"account", "delete", "alice", "-y"}, args...)...),
		)),
	)

	var afterDelete strings.Builder
	env.Exec("list accounts after delete",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{"account", "list"}, args...)...),
			step.Stdout(&afterDelete),
		)),
	)
	require.NotContains(t, afterDelete.String(), "alice")

	env.Exec("import the account from armor",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, append([]string{
				"account", "import", "alice2", filepath.Join(tmp, "alice.asc"),
			}, args...)...),
		)),
	)
}

// TestGnoChainE2E covers the full dev loop: serve the dev chain with a
// scaffolded realm preloaded, then call, query, send and deploy against it.
func TestGnoChainE2E(t *testing.T) {
	var (
		env  = envtest.New(t)
		tmp  = env.TmpDir()
		gnoH = filepath.Join(t.TempDir(), "gnohome")
	)

	env.Exec("scaffold a realm",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "realm", "counter"),
			step.Workdir(tmp),
		)),
	)

	listener := freeRPCListener(t)
	startServe(t, env, filepath.Join(tmp, "counter"), listener)
	remote := strings.TrimPrefix(listener, "tcp://")

	execGno := func(msg string, args ...string) {
		env.Must(env.Exec(msg,
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, args...),
				step.Workdir(tmp),
			)),
		))
	}

	execGno("call a realm function",
		"chain", "call", "gno.land/r/counter", "Increment",
		"--remote", remote, "--home", gnoH,
	)

	execGno("call it again",
		"chain", "call", "gno.land/r/counter", "Increment",
		"--remote", remote, "--home", gnoH,
	)

	var queryOut strings.Builder
	env.Must(env.Exec("query the realm state",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"chain", "query", "gno.land/r/counter.Get()",
				"--remote", remote,
			),
			step.Workdir(tmp),
			step.Stdout(&queryOut),
		)),
	))
	require.Contains(t, queryOut.String(), "2", "counter should be 2 after two increments: %s", queryOut.String())

	// fund a fresh account by key name
	execGno("create an account",
		"account", "create", "bob", "--home", gnoH,
	)
	execGno("send coins to bob",
		"chain", "send", "bob", "100ugnot",
		"--remote", remote, "--home", gnoH,
	)

	// deploy a second package at runtime
	env.Exec("scaffold a package",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "scaffold", "package", "greetings"),
			step.Workdir(tmp),
		)),
	)
	execGno("deploy the package",
		"chain", "deploy", "greetings",
		"--remote", remote, "--home", gnoH,
	)
}
