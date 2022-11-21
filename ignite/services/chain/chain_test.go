package chain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moby/moby/pkg/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceVersion(t *testing.T) {
	t.Run("tagged latest commit", func(t *testing.T) {
		c, err := New(tempSource(t, "testdata/version/mars.v0.2.tar.gz"))
		require.NoError(t, err)

		assert.Equal(t, "0.2", c.sourceVersion.tag)
		assert.Equal(t, "503123b1ac552437c7db3d17f816fd4121ff400d", c.sourceVersion.hash)
	})

	t.Run("tagged older commit", func(t *testing.T) {
		c, err := New(tempSource(t, "testdata/version/mars.v0.2-3-gaae48b7.tar.gz"))
		require.NoError(t, err)

		assert.Equal(t, "0.2-aae48b7f", c.sourceVersion.tag)
		assert.Equal(t, "aae48b7ffa4991bbe229f0969db8fe8623bf1fd4", c.sourceVersion.hash)
	})
}

func tempSource(t *testing.T, tarPath string) (path string) {
	f, err := os.Open(tarPath)
	require.NoError(t, err)

	defer f.Close()

	dir := t.TempDir()

	require.NoError(t, archive.Untar(f, dir, &archive.TarOptions{NoLchown: true}))

	dirs, err := os.ReadDir(dir)
	require.NoError(t, err)

	return filepath.Join(dir, dirs[0].Name())
}
