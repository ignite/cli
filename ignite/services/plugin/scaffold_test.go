package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const fooBarAppURI = "github.com/foo/bar"

func TestScaffold(t *testing.T) {
	// Arrange
	tmp := t.TempDir()
	ctx := context.Background()

	// Act
	session := cliui.New(cliui.WithoutUserInteraction(true))
	path, err := plugin.Scaffold(ctx, session, tmp, fooBarAppURI, false)

	// Assert
	require.NoError(t, err)
	require.DirExists(t, path)
	require.FileExists(t, filepath.Join(path, "go.mod"))
	require.FileExists(t, filepath.Join(path, "main.go"))
}

func TestScaffoldedConfig(t *testing.T) {
	// Arrange
	ctx := context.Background()
	path := scaffoldApp(ctx, t, fooBarAppURI)

	// Act
	cfg := readConfig(t, path)

	// Assert
	require.EqualValues(t, 1, cfg.Version)
	require.Len(t, cfg.Apps, 1)
}

func TestScaffoldedTests(t *testing.T) {
	// Arrange
	ctx := context.Background()
	path := scaffoldApp(ctx, t, fooBarAppURI)
	path = filepath.Join(path, "integration")

	// Act
	err := gocmd.Test(ctx, path, []string{
		"-timeout",
		"10m",
		"-run",
		"^TestBar$",
	})

	// Assert
	require.NoError(t, err)
}

func scaffoldApp(ctx context.Context, t *testing.T, path string) string {
	t.Helper()

	session := cliui.New(cliui.WithoutUserInteraction(true))
	path, err := plugin.Scaffold(ctx, session, t.TempDir(), path, false)
	require.NoError(t, err)
	return path
}

func readConfig(t *testing.T, path string) (cfg plugin.AppsConfig) {
	t.Helper()

	bz, err := os.ReadFile(filepath.Join(path, "app.ignite.yml"))
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(bz, &cfg))
	return
}
