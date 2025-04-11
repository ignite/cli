//go:build !relayer

package chain_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xos"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestServeWithCustomHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/sgbloga")
		servers = app.RandomizeServerPorts()
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()
	app.MustServe(ctx)

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeWithConfigHome(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/sgblogb")
		servers = app.RandomizeServerPorts()
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()
	app.MustServe(ctx)

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func TestServeWithCustomConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/sgblogc")
	)
	// Move config
	newConfig := "new_config.yml"
	newConfigPath := filepath.Join(tmpDir, newConfig)
	err := xos.Rename(filepath.Join(app.SourcePath(), "config.yml"), newConfigPath)
	require.NoError(t, err)
	app.SetConfigPath(newConfigPath)

	servers := app.RandomizeServerPorts()

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()
	app.MustServe(ctx)

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

// TestServeWithName tests serving a new chain scaffolded using a local name instead of a GitHub URI.
func TestServeWithName(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("sgblogd")
		servers = app.RandomizeServerPorts()
	)

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)

	var isBackendAliveErr error

	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()
	app.MustServe(ctx)

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
