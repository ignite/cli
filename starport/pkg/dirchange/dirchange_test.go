package dirchange

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	TmpPattern   = "starport-dirchange"
	ChecksumFile = "checksum.txt"
)

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func TestHasDirChecksumChanged(t *testing.T) {
	tempDir := os.TempDir()

	// Create directory tree
	dir1 := filepath.Join(tempDir, "foo1")
	err := os.MkdirAll(dir1, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(dir1)
	dir2 := filepath.Join(tempDir, "foo2")
	err = os.MkdirAll(dir2, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(dir2)
	dir3 := filepath.Join(tempDir, "foo3")
	err = os.MkdirAll(dir3, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(dir3)

	dir11, err := ioutil.TempDir(dir1, TmpPattern)
	require.NoError(t, err)
	dir12, err := ioutil.TempDir(dir1, TmpPattern)
	require.NoError(t, err)
	dir21, err := ioutil.TempDir(dir2, TmpPattern)
	require.NoError(t, err)

	// Create files
	err = ioutil.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir11, "foo"), randomBytes(15), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir12, "foo"), randomBytes(20), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir21, "foo"), randomBytes(20), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(10), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir3, "foo2"), randomBytes(10), 0644)
	require.NoError(t, err)

	// Check checksum
	paths := []string{dir1, dir2, dir3}
	checksum, err := checksumFromPaths("", paths)
	require.NoError(t, err)
	// md5 checksum is 16 bytes
	require.Len(t, checksum, 16)

	// Checksum remains the same if a file is deleted and recreated with the same content
	err = os.Remove(filepath.Join(dir1, "foo"))
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir1, "foo"), []byte("some bytes"), 0644)
	require.NoError(t, err)
	tmpChecksum, err := checksumFromPaths("", paths)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Can compute the checksum from a specific workdir
	pathNames := []string{"foo1", "foo2", "foo3"}
	tmpChecksum, err = checksumFromPaths(tempDir, pathNames)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Ignore non existent dir
	pathNames = append(pathNames, "nonexistent")
	tmpChecksum, err = checksumFromPaths(tempDir, pathNames)
	require.NoError(t, err)
	require.Equal(t, checksum, tmpChecksum)

	// Checksum from a subdir is different
	tmpChecksum, err = checksumFromPaths("", []string{dir1, dir2})
	require.NoError(t, err)
	require.NotEqual(t, checksum, tmpChecksum)

	// Checksum changes if a file is modified
	err = ioutil.WriteFile(filepath.Join(dir3, "foo1"), randomBytes(10), 0644)
	require.NoError(t, err)
	newChecksum, err := checksumFromPaths("", paths)
	require.NoError(t, err)
	require.NotEqual(t, checksum, newChecksum)

	// Error if no files in the specified dirs
	empty1 := filepath.Join(tempDir, "empty1")
	err = os.MkdirAll(empty1, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(empty1)
	empty2 := filepath.Join(tempDir, "empty2")
	err = os.MkdirAll(empty2, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(empty2)
	_, err = checksumFromPaths("", []string{empty1, empty2})
	require.Error(t, err)

	// SaveDirChecksum saves the checksum in the specified dir
	saveDir, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(saveDir)
	err = SaveDirChecksum("", paths, saveDir, ChecksumFile)
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(saveDir, ChecksumFile))
	fileContent, err := ioutil.ReadFile(filepath.Join(saveDir, ChecksumFile))
	require.NoError(t, err)
	require.Equal(t, newChecksum, fileContent)

	// Error if the paths contains no file
	err = SaveDirChecksum("", []string{empty1, empty2}, saveDir, ChecksumFile)
	require.Error(t, err)

	// HasDirChecksumChanged returns false if the directory has not changed
	changed, err := HasDirChecksumChanged("", paths, saveDir, ChecksumFile)
	require.NoError(t, err)
	require.False(t, changed)

	// Return true if checksum file doesn't exist
	newSaveDir, err := ioutil.TempDir(tempDir, TmpPattern)
	require.NoError(t, err)
	defer os.RemoveAll(newSaveDir)
	changed, err = HasDirChecksumChanged("", paths, newSaveDir, ChecksumFile)
	require.NoError(t, err)
	require.True(t, changed)

	// Return true if the paths contains no file
	changed, err = HasDirChecksumChanged("", []string{empty1, empty2}, saveDir, ChecksumFile)
	require.NoError(t, err)
	require.True(t, changed)

	// Return true if it has been changed
	err = ioutil.WriteFile(filepath.Join(dir21, "bar"), randomBytes(20), 0644)
	require.NoError(t, err)
	changed, err = HasDirChecksumChanged("", paths, saveDir, ChecksumFile)
	require.NoError(t, err)
	require.True(t, changed)
}
