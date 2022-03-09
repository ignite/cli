//go:build !relayer
// +build !relayer

package app_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	envtest "github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestServeStargateWithWasm(t *testing.T) {
	t.Skip()

	var (
		env     = envtest.New(t)
		apath   = env.Scaffold("sgblog")
		servers = env.RandomizeServerPorts(apath, "")
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "wasm"),
			step.Workdir(apath),
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
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		apath   = env.Scaffold("sgblog2")
		servers = env.RandomizeServerPorts(apath, "")
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "./home", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithConfigHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		apath   = env.Scaffold("sgblog3")
		servers = env.RandomizeServerPorts(apath, "")
	)

	// Set config homes
	env.SetRandomHomeConfig(apath, "")

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeStargateWithCustomConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	var (
		env   = envtest.New(t)
		apath = env.Scaffold("sgblog4")
	)
	// Move config
	newConfig := "new_config.yml"
	newConfigPath := filepath.Join(tmpDir, newConfig)
	err := os.Rename(filepath.Join(apath, "config.yml"), newConfigPath)
	require.NoError(t, err)

	servers := env.RandomizeServerPorts(tmpDir, newConfig)

	// Set config homes
	env.SetRandomHomeConfig(tmpDir, newConfig)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", apath, "", newConfigPath, envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
