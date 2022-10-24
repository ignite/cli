package xexec_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xexec"
)

func TestIsExec(t *testing.T) {
	cases := []struct {
		name, path string
		want       bool
	}{
		{
			name: "executable",
			path: "testdata/bin.sh",
			want: true,
		},
		{
			name: "not_executable",
			path: "testdata/nobin",
			want: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ok, err := xexec.IsExec(tt.path)

			// Assert
			require.NoError(t, err)
			require.Equal(t, tt.want, ok)
		})
	}
}

func TestResolveAbsPath(t *testing.T) {
	// Get the absolute path to the testdata directory
	testdata, err := filepath.Abs("testdata")
	require.NoError(t, err)

	cases := []struct {
		name, path, want string
		env              []string
	}{
		{
			name: "relative",
			path: "testdata/bin.sh",
			want: filepath.Join(testdata, "bin.sh"),
		},
		{
			name: "path",
			path: "bin.sh",
			env:  []string{"PATH", testdata},
			want: filepath.Join(testdata, "bin.sh"),
		},
		{
			name: "go bin path",
			path: "bin.sh",
			env:  []string{"GOBIN", testdata},
			want: filepath.Join(testdata, "bin.sh"),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			if tt.env != nil {
				t.Setenv(tt.env[0], tt.env[1])
			}

			// Act
			path, err := xexec.ResolveAbsPath(tt.path)

			// Assert
			require.NoError(t, err)
			require.Equal(t, tt.want, path)
		})
	}
}

func TestResolveAbsPathError(t *testing.T) {
	// Arrange
	fileName := "invalid-file.ko"

	// Act
	_, err := xexec.ResolveAbsPath(fileName)

	// Assert
	require.Errorf(t, err, `exec: "%s": executable file not found in $PATH`, fileName)
}

func TestTryResolveAbsPath(t *testing.T) {
	// Get the absolute path to the testdata directory
	testdata, err := filepath.Abs("testdata")
	require.NoError(t, err)

	cases := []struct {
		name, path, want string
		env              []string
	}{
		{
			name: "valid file",
			path: "testdata/bin.sh",
			want: filepath.Join(testdata, "bin.sh"),
		},
		{
			name: "invalid file",
			path: "invalid-file.ko",
			want: "invalid-file.ko",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			path := xexec.TryResolveAbsPath(tt.path)

			// Assert
			require.NoError(t, err)
			require.Equal(t, tt.want, path)
		})
	}
}

func TestIsCommandAvailable(t *testing.T) {
	cases := []struct {
		name, path string
		want       bool
	}{
		{
			name: "available",
			path: "testdata/bin.sh",
			want: true,
		},
		{
			name: "not_available",
			path: "invalid-file.ko",
			want: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ok := xexec.IsCommandAvailable(tt.path)

			// Assert
			require.Equal(t, tt.want, ok)
		})
	}
}
