package xos_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/errors"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xos"
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

func TestValidateFolderCopy(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestValidateFolderCopy")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tempDir))
	})

	// Create source and destination directories and files
	var (
		srcPath         = filepath.Join(tempDir, "source")
		srcFile         = filepath.Join(srcPath, "test.txt")
		dstPath         = filepath.Join(tempDir, "destination")
		dstFile         = filepath.Join(dstPath, "test.txt")
		emptyPath       = filepath.Join(tempDir, "empty")
		nonExistentPath = filepath.Join(tempDir, "nonexistent")
	)

	err = os.MkdirAll(srcPath, 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(dstPath, 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(emptyPath, 0o755)
	require.NoError(t, err)
	err = os.WriteFile(srcFile, []byte("source test"), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(dstFile, []byte("destination test"), 0o644)
	require.NoError(t, err)

	type args struct {
		srcPath string
		dstPath string
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "valid paths",
			args: args{
				srcPath: srcPath,
				dstPath: dstPath,
			},
			want: []string{"test.txt"},
		},
		{
			name: "same source and destination",
			args: args{
				srcPath: srcPath,
				dstPath: srcPath,
			},
			want: []string{},
			err:  errors.Errorf("source and destination paths are the same %s", srcPath),
		},
		{
			name: "empty directory",
			args: args{
				srcPath: emptyPath,
				dstPath: dstPath,
			},
			want: []string{},
		},
		{
			name: "non existent source",
			args: args{
				srcPath: nonExistentPath,
				dstPath: dstPath,
			},
			err: errors.Errorf("source path does not exist: %s", nonExistentPath),
		},
		{
			name: "non existent destination",
			args: args{
				srcPath: srcPath,
				dstPath: nonExistentPath,
			},
			err: errors.Errorf("destination path does not exist: %s", nonExistentPath),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xos.ValidateFolderCopy(tt.args.srcPath, tt.args.dstPath)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
