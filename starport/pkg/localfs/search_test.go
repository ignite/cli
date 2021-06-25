package localfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGlobTests(t *testing.T, files []string) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpdir, err := os.MkdirTemp(dir, "glob-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpdir)
	})
	for _, file := range files {
		fileDir := filepath.Dir(file)
		fileDir = filepath.Join(tmpdir, fileDir)
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, 0755)
		}
		err = os.WriteFile(filepath.Join(tmpdir, file), []byte{}, 0644)
		if err != nil {
			t.Fatal(err)
		}
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
				tmpdir + "/foo/bar/file1.proto",
				tmpdir + "/foo/bar/file2.proto",
				tmpdir + "/foo/baz/file.proto",
				tmpdir + "/foo/file.proto",
			},
		}, {
			name: "get only one proto file by name",
			args: args{
				path:    tmpdir,
				pattern: "file1.proto",
			},
			want: []string{tmpdir + "/foo/bar/file1.proto"},
		}, {
			name: "get two proto files by name",
			args: args{
				path:    tmpdir,
				pattern: "file.proto",
			},
			want: []string{tmpdir + "/foo/baz/file.proto", tmpdir + "/foo/file.proto"},
		}, {
			name: "get a specific file by name",
			args: args{
				path:    tmpdir,
				pattern: "file",
			},
			want: []string{tmpdir + "/foo/baz/file", tmpdir + "/foo/file"},
		}, {
			name: "error invalid directory",
			args: args{
				path:    "no-directory",
				pattern: "file",
			},
			err: &fs.PathError{Op: "stat", Path: "no-directory", Err: syscall.ENOENT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.args.path, tt.args.pattern)
			if tt.err != nil {
				require.Error(t, err)
				assert.EqualValues(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
