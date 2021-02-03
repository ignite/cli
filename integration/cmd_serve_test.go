// +build !relayer

package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestServeLaunchpadAppWithWasm(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("blog", Launchpad)
		servers = env.RandomizeServerPorts(apath)
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(apath),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Wasm", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeLaunchpadAppWithCustomHomes(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("blog2", Launchpad)
		servers = env.RandomizeServerPorts(apath)
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Wasm", apath, "./home", "./clihome", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeLaunchpadAppWithConfigHomes(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("blog3", Launchpad)
		servers = env.RandomizeServerPorts(apath)
	)

	// Set config homes
	env.SetRandomHomeConfig(apath)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Wasm", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithWasm(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog", Stargate)
		servers = env.RandomizeServerPorts(apath)
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(apath),
		)),
	))

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomHome(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog2", Stargate)
		servers = env.RandomizeServerPorts(apath)
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "./home", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithConfigHome(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("sgblog3", Stargate)
		servers = env.RandomizeServerPorts(apath)
	)

	// Set config homes
	env.SetRandomHomeConfig(apath)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), serveTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
