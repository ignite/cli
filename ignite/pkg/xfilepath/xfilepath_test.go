package xfilepath_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
)

func TestJoin(t *testing.T) {
	retriever := xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	p, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), p)

	retriever = xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}

func TestJoinFromHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	retriever := xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	p, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		home,
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), p)

	retriever = xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}

func TestList(t *testing.T) {
	retriever := xfilepath.List()
	list, err := retriever()
	require.NoError(t, err)
	require.Equal(t, []string(nil), list)

	retriever1 := xfilepath.Join(
		xfilepath.Path("foo/bar"),
	)
	retriever2 := xfilepath.Join(
		xfilepath.Path("bar/foo"),
	)
	retriever = xfilepath.List(retriever1, retriever2)
	list, err = retriever()
	require.NoError(t, err)
	require.Equal(t, []string{
		filepath.Join("foo", "bar"),
		filepath.Join("bar", "foo"),
	}, list)

	retrieverError := xfilepath.PathWithError("foo", errors.New("foo"))
	retriever = xfilepath.List(retriever1, retrieverError, retriever2)
	_, err = retriever()
	require.Error(t, err)
}

func TestMkdir(t *testing.T) {
	newdir := path.Join(t.TempDir(), "hey")

	dir, err := xfilepath.Mkdir(xfilepath.Path(newdir))()

	require.NoError(t, err)
	require.Equal(t, newdir, dir)
	require.DirExists(t, dir)
}

func TestRelativePath(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	rootRelative, err := filepath.Rel(pwd, "/")
	require.NoError(t, err)

	tests := []struct {
		name    string
		appPath string
		want    string
		err     error
	}{
		{
			name:    "same directory",
			appPath: filepath.Join(pwd, "file.go"),
			want:    "file.go",
		},
		{
			name:    "previous directory",
			appPath: filepath.Join(filepath.Dir(pwd), "file.go"),
			want:    "../file.go",
		},
		{
			name:    "root directory",
			appPath: "/file.go",
			want:    filepath.Join(rootRelative, "file.go"),
		},
		{
			name:    "absolute path",
			appPath: pwd,
			want:    ".",
		},
		{
			name:    "NonExistentPath",
			appPath: filepath.Join(filepath.Base(pwd), "file.go"),
			want:    "",
			err:     errors.Errorf("Rel: can't make xfilepath/file.go relative to %s", pwd),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xfilepath.RelativePath(tt.appPath)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsDir(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing directory",
			path: ".",
			want: true,
		},
		{
			name: "existing sub directory",
			path: "./testdata",
			want: true,
		},
		{
			name: "existing file",
			path: "./testdata/testfile",
			want: false,
		},
		{
			name: "non-existing directory",
			path: "nonexistent",
			want: false,
		},
		{
			name: "parent directory",
			path: "..",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xfilepath.IsDir(tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMustAbs(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name string
		path string
		want string
		err  error
	}{
		{
			name: "already absolute path",
			path: "/absolute/path",
			want: "/absolute/path",
			err:  nil,
		},
		{
			name: "relative path",
			path: "relative/path",
			want: filepath.Join(pwd, "relative/path"),
			err:  nil,
		},
		{
			name: "current directory",
			path: ".",
			want: pwd,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xfilepath.MustAbs(tt.path)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
