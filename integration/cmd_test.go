package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Launchpad)
	)

	_, statErr := os.Stat(filepath.Join(path, "config.yml"))
	require.False(t, os.IsNotExist(statErr), "config.yml cannot be found")

	env.EnsureAppIsSteady(path)
}

func TestGenerateAnAppWithCosmWasmAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Launchpad)
	)

	env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		),
	)

	env.Exec("should not add CosmWasm module second time",
		step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		),
		ExecShouldError(),
	)

	env.EnsureAppIsSteady(path)
}

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

func TestGenerateAnAppWithEmptyModuleAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Launchpad)
	)

	env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		),
	)

	env.EnsureAppIsSteady(path)
}

func TestGenerateAStargateAppWithEmptyModuleAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Stargate)
	)

	env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		),
	)

	env.EnsureAppIsSteady(path)
}

func TestGenerateAnAppWithTypeAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Stargate)
	)

	env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		),
	)

	env.EnsureAppIsSteady(path)
}
