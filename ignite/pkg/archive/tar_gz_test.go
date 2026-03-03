package archive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateArchiveAndExtractArchive(t *testing.T) {
	root := t.TempDir()
	oldWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(oldWD))
	})

	src := "src"
	require.NoError(t, os.MkdirAll(filepath.Join(src, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "a.txt"), []byte("alpha"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(src, "nested", "b.txt"), []byte("beta"), 0o600))

	var buf bytes.Buffer
	require.NoError(t, CreateArchive(src, &buf))

	dst := filepath.Join(root, "out")
	require.NoError(t, os.MkdirAll(dst, 0o755))
	require.NoError(t, ExtractArchive(dst, bytes.NewReader(buf.Bytes())))

	gotA, err := os.ReadFile(filepath.Join(dst, "src", "a.txt"))
	require.NoError(t, err)
	require.Equal(t, "alpha", string(gotA))

	gotB, err := os.ReadFile(filepath.Join(dst, "src", "nested", "b.txt"))
	require.NoError(t, err)
	require.Equal(t, "beta", string(gotB))
}

func TestAddToArchiveReturnsErrorForMissingFile(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, CreateArchive(t.TempDir(), &buf))

	err := addToArchive(nil, filepath.Join(t.TempDir(), "does-not-exist"))
	require.Error(t, err)
}
