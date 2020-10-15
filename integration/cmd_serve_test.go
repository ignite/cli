package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestServeAppWithCosmWasm(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("blog", Launchpad)
		servers = env.RandomizeServerPorts(apath)
	)

	env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(apath),
		),
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), time.Minute*5)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Serve("should serve with CosmWasm", apath, ExecCtx(ctx))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeAppStargate(t *testing.T) {
	t.Parallel()

	var (
		env     = newEnv(t)
		apath   = env.Scaffold("stargateblog", Stargate)
		servers = env.RandomizeServerPorts(apath)
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), time.Minute*5)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Serve("should serve with Stargate version", apath, ExecCtx(ctx))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
