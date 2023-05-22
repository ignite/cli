package dirchange_test

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/dirchange"
)

const (
	TmpPattern  = "starport-dirchange"
	ChecksumKey = "checksum"
)

func randomBytes(t *testing.T, n int) []byte {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	require.NoError(t, err)

	return bytes
}

func TestHasDirChecksumChanged(t *testing.T) {
	tempDir := os.TempDir()
	cacheDir := os.TempDir()

	cacheStorage, err := cache.NewStorage(filepath.Join(cacheDir, "testcache.db"))
	require.NoError(t, err)
	c := cache.New[[]byte](cacheStorage, "testnamespace")

	// Create directory tree
	dir1 := filepath.Join(tempDir, "foo1")
	err = os.MkdirAll(dir1, 0o700)
	require.NoError(t, err)
	defer os.RemoveAll(dir1)
	dir2 := filepath.Join(tempDir, "foo2")
	err = os.MkdirAll(dir2, 0o700)
	require.NoError(t, err)
	defer os.RemoveAll(dir2)
	dir3 := filepath.Join(tempDir, "foo3")
	err = os.MkdirAll(dir3, 0o700)
	require.NoError(t, err)
	defer os.RemoveAll(dir3)

	dir11, err := os.MkdirTemp(dir1, TmpPattern)
	require.NoError(t, err)
	dir12, err := os.MkdirTemp(dir1, TmpPattern)
	require.NoError(t, err)
	dir21, err := os.MkdirTemp(dir2, TmpPattern)
	require.NoError(t, err)

	// Create files
	err = os.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir11, "foo"), randomBytes(t, 15), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir12, "foo"), randomBytes(t, 20), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir21, "foo"), randomBytes(t, 20), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(t, 10), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir3, "foo2"), randomBytes(t, 10), 0o644)
	require.NoError(t, err)

	// Check checksum
	paths := []string{dir1, dir2, dir3}
	checksum, err := dirchange.ChecksumFromPaths("", paths...)
	require.NoError(t, err)
	// md5 checksum is 16 bytes
	require.Len(t, checksum, 16)

	// Checksum remains the same if a file is deleted and recreated with the same content
	err = os.Remove(filepath.Join(dir1, "foo"))
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"), 0o644)
	require.NoError(t, err)
	tmpChecksum, err := dirchange.ChecksumFromPaths("", paths...)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Can compute the checksum from a specific workdir
	pathNames := []string{"foo1", "foo2", "foo3"}
	tmpChecksum, err = dirchange.ChecksumFromPaths(tempDir, pathNames...)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Ignore non existent dir
	pathNames = append(pathNames, "nonexistent")
	tmpChecksum, err = dirchange.ChecksumFromPaths(tempDir, pathNames...)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Checksum from a subdir is different
	tmpChecksum, err = dirchange.ChecksumFromPaths("", dir1, dir2)
	require.NoError(t, err)
	require.NotEqual(t, checksum, tmpChecksum)

	// Checksum changes if a file is modified
	err = os.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(t, 10), 0o644)
	require.NoError(t, err)
	newChecksum, err := dirchange.ChecksumFromPaths("", paths...)
	require.NoError(t, err)
	require.NotEqual(t, checksum, newChecksum)

	// Error if no files in the specified dirs
	empty1 := filepath.Join(tempDir, "empty1")
	err = os.MkdirAll(empty1, 0o700)
	require.NoError(t, err)
	defer os.RemoveAll(empty1)
	empty2 := filepath.Join(tempDir, "empty2")
	err = os.MkdirAll(empty2, 0o700)
	require.NoError(t, err)
	defer os.RemoveAll(empty2)
	_, err = dirchange.ChecksumFromPaths("", empty1, empty2)
	require.Error(t, err)

	// SaveDirChecksum saves the checksum in the cache
	saveDir, err := os.MkdirTemp(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(saveDir)
	err = dirchange.SaveDirChecksum(c, ChecksumKey, "", paths...)
	require.NoError(t, err)
	savedChecksum, err := c.Get(ChecksumKey)
	require.NoError(t, err)
	require.Equal(t, newChecksum, savedChecksum)

	// Error if the paths contains no file
	err = dirchange.SaveDirChecksum(c, ChecksumKey, "", empty1, empty2)
	require.Error(t, err)

	// HasDirChecksumChanged returns false if the directory has not changed
	changed, err := dirchange.HasDirChecksumChanged(c, ChecksumKey, "", paths...)
	require.NoError(t, err)
	require.False(t, changed)

	// Return true if cache entry doesn't exist
	err = c.Delete(ChecksumKey)
	require.NoError(t, err)
	changed, err = dirchange.HasDirChecksumChanged(c, ChecksumKey, "", paths...)
	require.NoError(t, err)
	require.True(t, changed)

	// Return true if the paths contains no file
	changed, err = dirchange.HasDirChecksumChanged(c, ChecksumKey, "", empty1, empty2)
	require.NoError(t, err)
	require.True(t, changed)

	// Return true if it has been changed
	err = os.WriteFile(filepath.Join(dir21, "bar"), randomBytes(t, 20), 0o644)
	require.NoError(t, err)
	changed, err = dirchange.HasDirChecksumChanged(c, ChecksumKey, "", paths...)
	require.NoError(t, err)
	require.True(t, changed)
}
