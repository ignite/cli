package xos_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

func TestFindFiles(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		extension []string
		prefix    string
		want      []string
		err       error
	}{
		{
			name:  "test zero files",
			files: []string{},
			want:  []string{},
			err:   nil,
		},
		{
			name:  "test one file",
			files: []string{"file.json"},
			want:  []string{"file.json"},
			err:   nil,
		},
		{
			name:  "test 3 files",
			files: []string{"file1.json", "file2.txt", "file3.json"},
			want:  []string{"file1.json", "file2.txt", "file3.json"},
			err:   nil,
		},
		{
			name:   "test file prefix",
			files:  []string{"file.prefix.test.json"},
			prefix: "file.prefix",
			want:   []string{"file.prefix.test.json"},
			err:    nil,
		},
		{
			name:   "test bigger file prefix",
			files:  []string{"file.prefix.test.json"},
			prefix: "file.prefix.test",
			want:   []string{"file.prefix.test.json"},
			err:    nil,
		},
		{
			name:   "test 3 files prefix",
			files:  []string{"test.file1.json", "test.file2.txt", "test.file3.json"},
			prefix: "test.file",
			want:   []string{"test.file1.json", "test.file2.txt", "test.file3.json"},
			err:    nil,
		},
		{
			name:      "test 3 extension json files",
			files:     []string{"file1.json", "file2.txt", "file3.json", "file4.json"},
			extension: []string{"json"},
			want:      []string{"file1.json", "file3.json", "file4.json"},
			err:       nil,
		},
		{
			name:      "test 3 extension json files with subfolder",
			files:     []string{"testdata/file1.json", "file2.txt", "foo/file3.json", "file4.json"},
			extension: []string{"json"},
			want:      []string{"testdata/file1.json", "foo/file3.json", "file4.json"},
			err:       nil,
		},
		{
			name:      "test 1 extension txt files",
			files:     []string{"file1.json", "file2.txt", "file3.json", "file4.json"},
			extension: []string{"txt"},
			want:      []string{"file2.txt"},
			err:       nil,
		},
		{
			name:      "test 1 extension json files",
			files:     []string{"file1.json"},
			extension: []string{"json"},
			want:      []string{"file1.json"},
			err:       nil,
		},
		{
			name:      "test invalid files extension",
			files:     []string{"file1.json", "file2.json", "file3.json", "file4.json"},
			extension: []string{"txt"},
			want:      []string{},
			err:       nil,
		},
		{
			name:      "test file prefix and extension",
			files:     []string{"test.file1.json", "test.file2.txt", "test.file3.json"},
			prefix:    "test.file",
			extension: []string{"json"},
			want:      []string{"test.file1.json", "test.file3.json"},
			err:       nil,
		},
		{
			name:      "test 2 different extensions",
			files:     []string{"file1.json", "file2.txt", "file3.json", "file4.json", "file.yaml"},
			extension: []string{"txt", "yaml"},
			want:      []string{"file2.txt", "file.yaml"},
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
				require.NoError(t, os.MkdirAll(filepath.Dir(filePath), 0o755))
				file, err := os.Create(filePath)
				require.NoError(t, err)
				require.NoError(t, file.Close())
			}

			opts := make([]xos.FindFileOptions, 0)
			if tt.prefix != "" {
				opts = append(opts, xos.WithPrefix(tt.prefix))
			}

			for _, ext := range tt.extension {
				opts = append(opts, xos.WithExtension(ext))
			}

			gotFiles, err := xos.FindFiles(tempDir, opts...)
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
			require.ElementsMatch(t, want, gotFiles)
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
