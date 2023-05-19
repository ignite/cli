package xos_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xos"
)

func TestFindFiles(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		extension string
		want      []string
		err       error
	}{
		{
			name:      "test 3 json files",
			files:     []string{"file1.json", "file2.txt", "file3.json", "file4.json"},
			extension: "json",
			want:      []string{"file1.json", "file3.json", "file4.json"},
			err:       nil,
		},
		{
			name:      "test 1 txt files",
			files:     []string{"file1.json", "file2.txt", "file3.json", "file4.json"},
			extension: "txt",
			want:      []string{"file2.txt"},
			err:       nil,
		},
		{
			name:      "test 1 json files",
			files:     []string{"file1.json"},
			extension: "json",
			want:      []string{"file1.json"},
			err:       nil,
		},
		{
			name:      "test no files",
			files:     []string{"file1.json", "file2.json", "file3.json", "file4.json"},
			extension: "txt",
			want:      []string{},
			err:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirName := strings.ReplaceAll(t.Name(), "/", "_")
			tempDir, err := os.MkdirTemp("", dirName)
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, os.RemoveAll(tempDir))
			})

			for _, filename := range tt.files {
				filePath := filepath.Join(tempDir, filename)
				file, err := os.Create(filePath)
				require.NoError(t, err)
				require.NoError(t, file.Close())
			}

			gotFiles, err := xos.FindFiles(tempDir, tt.extension)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			want := make([]string, len(tt.want))
			for i, filename := range tt.want {
				want[i] = filepath.Join(tempDir, filename)
			}
			require.EqualValues(t, want, gotFiles)
		})
	}
}

func TestFileExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "TestCopyFile")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tempDir))
	})

	srcDir := filepath.Join(tempDir, "source")
	err = os.MkdirAll(srcDir, 0o755)
	require.NoError(t, err)

	srcFile := filepath.Join(srcDir, "file.txt")
	err = os.WriteFile(srcFile, []byte("File content"), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "existing file",
			filename: srcFile,
			want:     true,
		},
		{
			name:     "non existing file",
			filename: "non_existing_file.txt",
			want:     false,
		},
		{
			name:     "directory",
			filename: srcDir,
			want:     false,
		},
		{
			name:     "empty filename",
			filename: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xos.FileExists(tt.filename)
			require.EqualValues(t, tt.want, got)
		})
	}
}
