package tarball

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	exampleJSON, err := os.ReadFile("testdata/example.json")
	require.NoError(t, err)

	type args struct {
		tarballPath string
		file        string
	}
	tests := []struct {
		name string
		args args
		want []byte
		err  error
	}{
		{
			name: "simple read",
			args: args{
				tarballPath: "testdata/example.tar.gz",
				file:        "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "read from root",
			args: args{
				tarballPath: "testdata/example-root.tar.gz",
				file:        "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "read from subfolder",
			args: args{
				tarballPath: "testdata/example-subfolder.tar.gz",
				file:        "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "empty folders",
			args: args{
				tarballPath: "testdata/example-empty.tar.gz",
				file:        "example.json",
			},
			err: ErrGzipFileNotFound,
		},
		{
			name: "invalid file",
			args: args{
				tarballPath: "testdata/invalid_file",
				file:        "example.json",
			},
			err: ErrInvalidGzipFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarball, err := os.ReadFile(tt.args.tarballPath)
			require.NoError(t, err)

			got, err := ReadFile(tarball, tt.args.file)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsTarball(t *testing.T) {
	tests := []struct {
		name        string
		tarballPath string
		err         error
	}{
		{
			name:        "simple read",
			tarballPath: "testdata/example.tar.gz",
		},
		{
			name:        "read from root",
			tarballPath: "testdata/example-root.tar.gz",
		},
		{
			name:        "read from subfolder",
			tarballPath: "testdata/example-subfolder.tar.gz",
		},
		{
			name:        "invalid file",
			tarballPath: "testdata/invalid_file",
			err:         ErrInvalidGzipFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarball, err := os.ReadFile(tt.tarballPath)
			require.NoError(t, err)

			err = IsTarball(tarball)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
