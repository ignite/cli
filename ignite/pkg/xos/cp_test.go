package xos_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xos"
)

func TestCopyFolder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestCopyFile")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tempDir))
	})

	// Create temporary source and destination directories
	srcDir := filepath.Join(tempDir, "source")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	dstDir := filepath.Join(tempDir, "destination")
	err = os.MkdirAll(dstDir, 0o755)
	require.NoError(t, err)

	emptyDir := filepath.Join(tempDir, "empty")
	err = os.MkdirAll(emptyDir, 0o755)
	require.NoError(t, err)

	// Create a temporary source file
	srcFile1 := filepath.Join(srcDir, "file_1.txt")
	err = os.WriteFile(srcFile1, []byte("File content 1"), 0o644)
	require.NoError(t, err)

	srcFile2 := filepath.Join(srcDir, "file_2.txt")
	err = os.WriteFile(srcFile2, []byte("File content 2"), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name              string
		srcPath           string
		dstPath           string
		expectedErr       error
		expectedFileCount int
	}{
		{
			name:              "valid paths",
			srcPath:           srcDir,
			dstPath:           dstDir,
			expectedFileCount: 2,
		},
		{
			name:        "non existent destination",
			srcPath:     srcDir,
			dstPath:     filepath.Join(dstDir, "non-existent-destination"),
			expectedErr: os.ErrNotExist,
		},
		{
			name:        "non existent source",
			srcPath:     filepath.Join(dstDir, "non-existent-source"),
			dstPath:     dstDir,
			expectedErr: os.ErrNotExist,
		},
		{
			name:              "same source and destination",
			srcPath:           srcDir,
			dstPath:           srcDir,
			expectedFileCount: 2,
		},
		{
			name:              "empty source",
			srcPath:           emptyDir,
			dstPath:           filepath.Join(tempDir, "empty"),
			expectedFileCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := xos.CopyFolder(tt.srcPath, tt.dstPath)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)

			// Check the number of files in the destination directory
			files, err := os.ReadDir(tt.dstPath)
			require.NoError(t, err)
			require.Equal(t, tt.expectedFileCount, len(files))
		})
	}
}

func TestCopyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestCopyFile")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tempDir))
	})

	// Create temporary source and destination directories
	srcDir := filepath.Join(tempDir, "source")
	dstDir := filepath.Join(tempDir, "destination")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(dstDir, 0o755)
	require.NoError(t, err)

	// Create a temporary source file
	srcFile := filepath.Join(srcDir, "file.txt")
	err = os.WriteFile(srcFile, []byte("File content"), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name          string
		srcPath       string
		dstPath       string
		expectedErr   error
		expectedBytes int64 // Provide the expected number of bytes copied
	}{
		{
			name:          "valid path",
			srcPath:       srcFile,
			dstPath:       filepath.Join(dstDir, "test_1.txt"),
			expectedBytes: 12,
		},
		{
			name:        "non existent file",
			srcPath:     filepath.Join(srcDir, "non_existent_file.txt"),
			dstPath:     filepath.Join(dstDir, "test_2.txt"),
			expectedErr: os.ErrNotExist,
		},
		{
			name:        "non existent destination",
			srcPath:     srcFile,
			dstPath:     "/path/to/nonexistent/file.txt",
			expectedErr: os.ErrNotExist,
		},
		{
			name:    "same source and destination",
			srcPath: srcFile,
			dstPath: srcFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := xos.CopyFile(tt.srcPath, tt.dstPath)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)

			destFile, err := os.Open(tt.dstPath)
			require.NoError(t, err)

			destFileInfo, err := destFile.Stat()
			require.NoError(t, err)
			require.NoError(t, destFile.Close())
			require.Equal(t, tt.expectedBytes, destFileInfo.Size())
		})
	}
}
