package tarball

import (
	"io"
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
		name     string
		args     args
		want     []byte
		wantPath string
		err      error
	}{
		{
			name: "simple read",
			args: args{
				tarballPath: "testdata/example.tar.gz",
				file:        "example.json",
			},
			want:     exampleJSON,
			wantPath: "genesis/example.json",
		},
		{
			name: "read from root",
			args: args{
				tarballPath: "testdata/example-root.tar.gz",
				file:        "example.json",
			},
			want:     exampleJSON,
			wantPath: "example.json",
		},
		{
			name: "read from subfolder",
			args: args{
				tarballPath: "testdata/example-subfolder.tar.gz",
				file:        "example.json",
			},
			want:     exampleJSON,
			wantPath: "config/genesis/example.json",
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
			want: []byte{},
		},
		{
			name: "invalid file extension",
			args: args{
				tarballPath: "testdata/example.json",
				file:        "example.json",
			},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarball, err := os.Open(tt.args.tarballPath)
			require.NoError(t, err)

			gotReader, gotPath, err := ExtractFile(tarball, tt.args.file)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantPath, gotPath)

			gotReaderBytes, err := io.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, tt.want, gotReaderBytes)
		})
	}
}
