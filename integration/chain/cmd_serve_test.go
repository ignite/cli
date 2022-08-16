//go:build !relayer

package chain_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/xos"
	envtest "github.com/ignite/cli/integration"
)

func TestServeStargateWithWasm(t *testing.T) {
	t.Skip()

	var (
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/sgblog")
		servers = app.RandomizeServerPorts()
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "wasm", "--yes"),
			step.Workdir(app.SourcePath()),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/sgblog2")
		servers = app.RandomizeServerPorts()
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithConfigHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/sgblog3")
		servers = app.RandomizeServerPorts()
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/sgblog4")
	)
	// Move config
	newConfig := "new_config.yml"
	newConfigPath := filepath.Join(tmpDir, newConfig)
	err := xos.Rename(filepath.Join(apath, "config.yml"), newConfigPath)
	require.NoError(t, err)
	app.SetConfigPath(newConfigPath)

	servers := app.RandomizeServerPorts()

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

// TestServeStargateWithName tests serving a new chain scaffolded using a local name instead of a GitHub URI.
func TestServeStargateWithName(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.Scaffold("sgblog5")
		servers = app.RandomizeServerPorts()
	)

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)

	var isBackendAliveErr error

	go func() {
		defer cancel()

		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()

	env.Must(app.Serve("should serve with Stargate version", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
