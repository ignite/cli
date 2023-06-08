package plugin

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
}
