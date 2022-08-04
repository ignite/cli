package localfs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupGlobTests(t *testing.T, files []string) string {
	t.Helper()
	tmpdir := t.TempDir()

	for _, file := range files {
		fileDir := filepath.Dir(file)
		fileDir = filepath.Join(tmpdir, fileDir)
		err := os.MkdirAll(fileDir, 0o755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(tmpdir, file), []byte{}, 0o644)
		require.NoError(t, err)
	}
	return tmpdir
}

func TestSearch(t *testing.T) {
	files := []string{
		"foo/file.proto",
		"foo/bar/file1.proto",
		"foo/bar/file2.proto",
		"foo/baz/file.proto",
		"foo/file",
		"foo/baz/file",
	}
	tmpdir := setupGlobTests(t, files)
	type args struct {
		path    string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "get all proto files by pattern",
			args: args{
				path:    tmpdir,
				pattern: "*.proto",
			},
			want: []string{
				filepath.Join(tmpdir, "foo/bar/file1.proto"),
				filepath.Join(tmpdir, "foo/bar/file2.proto"),
				filepath.Join(tmpdir, "foo/baz/file.proto"),
				filepath.Join(tmpdir, "foo/file.proto"),
			},
		}, {
			name: "get only one proto file by name",
			args: args{
				path:    tmpdir,
				pattern: "file1.proto",
			},
			want: []string{filepath.Join(tmpdir, "foo/bar/file1.proto")},
		}, {
			name: "get two proto files by name",
			args: args{
				path:    tmpdir,
				pattern: "file.proto",
			},
			want: []string{
				filepath.Join(tmpdir, "foo/baz/file.proto"),
				filepath.Join(tmpdir, "foo/file.proto"),
			},
		}, {
			name: "get a specific file by name",
			args: args{
				path:    tmpdir,
				pattern: "file",
			},
			want: []string{
				filepath.Join(tmpdir, "foo/baz/file"),
				filepath.Join(tmpdir, "foo/file"),
			},
		}, {
			name: "not found directory",
			args: args{
				path:    "no-directory",
				pattern: "file",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.args.path, tt.args.pattern)
			if tt.err != nil {
				require.Error(t, err)
				require.EqualValues(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}
