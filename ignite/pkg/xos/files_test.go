package xos

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

			gotFiles, err := FindFiles(tempDir, tt.extension)
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
