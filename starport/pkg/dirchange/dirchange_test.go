package dirchange

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const TmpPattern = "starport-dirchange"

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func TestHasDirChecksumChanged(t *testing.T) {
	tempDir := os.TempDir()

	// Create directory tree
	dir1, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(dir1)
	dir2, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(dir2)
	dir3, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(dir3)

	dir11, err := ioutil.TempDir(dir1, TmpPattern)
	require.NoError(t, err)
	dir12, err := ioutil.TempDir(dir1, TmpPattern)
	require.NoError(t, err)
	dir21, err := ioutil.TempDir(dir2, TmpPattern)
	require.NoError(t, err)

	// Create files
	err = ioutil.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"),0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir11, "foo"), randomBytes(15),0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir12, "foo"), randomBytes(20),0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir21, "foo"), randomBytes(20),0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(10),0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir3, "foo2"), randomBytes(10),0644)
	require.NoError(t, err)

	// Check checksum
	paths := []string{dir1, dir2, dir3}
	checksum, err := checksumFromPaths(paths)
	require.NoError(t, err)
	// md5 checksum is 16 bytes
	require.Len(t, checksum, 16)

	// Checksum remains the same if a file is deleted and recreated with the same content
	err = os.Remove(filepath.Join(dir1, "foo"))
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"),0644)
	require.NoError(t, err)
	tmpChecksum, err := checksumFromPaths(paths)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Checksum from a subdir is different
	tmpChecksum, err = checksumFromPaths([]string{dir1, dir2})
	require.NoError(t, err)
	require.NotEqual(t, checksum, tmpChecksum)

	// Checksum changes if a file is modified
	err = ioutil.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(10),0644)
	require.NoError(t, err)
	newChecksum, err := checksumFromPaths(paths)
	require.NoError(t, err)
	require.NotEqual(t, checksum, newChecksum)

	// SaveDirChecksum saves the checksum in the specified dir
	saveDir, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(saveDir)
	err = SaveDirChecksum(paths, saveDir)
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(saveDir, checksumFile))
	fileContent, err := ioutil.ReadFile(filepath.Join(saveDir, checksumFile))
	require.NoError(t, err)
	require.Equal(t, newChecksum, fileContent)

	// HasDirChecksumChanged returns false if the directory has not changed
	changed, err := HasDirChecksumChanged(paths, saveDir)
	require.NoError(t, err)
	require.False(t, changed)

	// Return true and create the checksum file if it doesn't exist
	newSaveDir, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(newSaveDir)
	changed, err = HasDirChecksumChanged(paths, newSaveDir)
	require.NoError(t, err)
	require.True(t, changed)
	require.FileExists(t, filepath.Join(newSaveDir, checksumFile))
	fileContent, err = ioutil.ReadFile(filepath.Join(newSaveDir, checksumFile))
	require.NoError(t, err)
	require.Equal(t, newChecksum, fileContent)

	// Return true and rewrite the checksum if it has been changed
	err = ioutil.WriteFile(filepath.Join(dir21, "bar"), randomBytes(20),0644)
	require.NoError(t, err)
	changed, err = HasDirChecksumChanged(paths, saveDir)
	require.NoError(t, err)
	require.True(t, changed)
	fileContent, err = ioutil.ReadFile(filepath.Join(saveDir, checksumFile))
	require.NoError(t, err)
	newChecksum, err = checksumFromPaths(paths)
	require.NoError(t, err)
	require.Equal(t, newChecksum, fileContent)
}