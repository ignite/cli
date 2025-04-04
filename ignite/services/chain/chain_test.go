package chain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/archive"
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

func TestBech32Prefix(t *testing.T) {
	t.Run("default prefix when not specified", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)
		c, err := New(dir)
		require.NoError(t, err)

		prefix, err := c.Bech32Prefix()
		require.NoError(t, err)

		// Should return the default Cosmos prefix
		assert.Equal(t, "cosmos", prefix)
	})

	t.Run("returns custom prefix when specified", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)

		// Create mock app.go with custom prefix
		mockAppGo := `package app

		const (
			AccountAddressPrefix = "mars"
		)
		`
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app", "app.go"), []byte(mockAppGo), 0o644))

		c, err := New(dir)
		require.NoError(t, err)

		prefix, err := c.Bech32Prefix()
		require.NoError(t, err)

		assert.Equal(t, "mars", prefix)
	})

	t.Run("handles alternate prefix declaration format", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)

		// Create mock app.go with custom prefix in alternate format
		mockAppGo := `package app

		const AccountAddressPrefix string = "jupiter" // Some comment
		`
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app", "app.go"), []byte(mockAppGo), 0o644))

		c, err := New(dir)
		require.NoError(t, err)

		prefix, err := c.Bech32Prefix()
		require.NoError(t, err)

		assert.Equal(t, "jupiter", prefix)
	})
}

func TestCoinType(t *testing.T) {
	t.Run("default coin type when not specified", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)

		c, err := New(dir)
		require.NoError(t, err)

		coinType, err := c.CoinType()
		require.NoError(t, err)

		assert.Equal(t, uint32(118), coinType)
	})

	t.Run("returns custom coin type when specified", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)

		// Create mock app.go with custom coin type
		mockAppGo := `package app

		const (
			ChainCoinType = 529
		)
		`
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app", "app.go"), []byte(mockAppGo), 0o644))

		c, err := New(dir)
		require.NoError(t, err)

		coinType, err := c.CoinType()
		require.NoError(t, err)

		assert.Equal(t, uint32(529), coinType)
	})

	t.Run("handles coin type with comments", func(t *testing.T) {
		dir, err := tempSourceWithApp(t)
		require.NoError(t, err)

		mockAppGo := `package app

		// ChainCoinType is the coin type for this chain
		const ChainCoinType = 330 // Custom coin type for test
		`
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app", "app.go"), []byte(mockAppGo), 0o644))

		c, err := New(dir)
		require.NoError(t, err)

		coinType, err := c.CoinType()
		require.NoError(t, err)

		assert.Equal(t, uint32(330), coinType)
	})
}

func tempSource(t *testing.T, tarPath string) (path string) {
	t.Helper()

	f, err := os.Open(tarPath)
	require.NoError(t, err)

	defer f.Close()

	dir := t.TempDir()

	require.NoError(t, archive.ExtractArchive(dir, f))

	dirs, err := os.ReadDir(dir)
	require.NoError(t, err)

	return filepath.Join(dir, dirs[0].Name())
}

func tempSourceWithApp(t *testing.T) (string, error) {
	t.Helper()

	tmpDir := t.TempDir()

	emptyFilesPaths := []string{
		filepath.Join(tmpDir, "go.mod"),
		filepath.Join(tmpDir, "app", "app.go"),
	}

	if err := os.WriteFile(emptyFilesPaths[0], []byte("module my-new-chain"), 0o755); err != nil {
		return "", err
	}

	for _, f := range emptyFilesPaths[1:] {
		if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
			return "", err
		}

		if err := os.WriteFile(f, []byte("package my-new-chain"), 0o755); err != nil {
			return "", err
		}

	}

	return tmpDir, nil
}
