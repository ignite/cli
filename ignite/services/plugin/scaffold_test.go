package plugin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestScaffold(t *testing.T) {
	// Arrange
	tmp := t.TempDir()
	ctx := context.Background()

	// Act
	path, err := Scaffold(ctx, tmp, "github.com/foo/bar", false)

	// Assert
	require.NoError(t, err)
	require.DirExists(t, path)
	require.FileExists(t, filepath.Join(path, "go.mod"))
	require.FileExists(t, filepath.Join(path, "main.go"))

	appYML, err := os.ReadFile(filepath.Join(path, "app.ignite.yml"))
	require.NoError(t, err)
	var config AppsConfig
	err = yaml.Unmarshal(appYML, &config)
	require.NoError(t, err)
	require.EqualValues(t, 1, config.Version)
	require.Len(t, config.Apps, 1)
}
