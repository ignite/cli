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
		source string
		file   string
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
				source: "testdata/example.tar.gz",
				file:   "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "read from root",
			args: args{
				source: "testdata/example-root.tar.gz",
				file:   "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "read from subfolder",
			args: args{
				source: "testdata/example-subfolder.tar.gz",
				file:   "example.json",
			},
			want: exampleJSON,
		},
		{
			name: "not found file",
			args: args{
				source: "testdata/not-found.tar.gz",
				file:   "example.json",
			},
			err: ErrFileNotFound,
		},
		{
			name: "empty folders",
			args: args{
				source: "testdata/example-empty.tar.gz",
				file:   "example.json",
			},
			err: ErrGzipFileNotFound,
		},
		{
			name: "invalid file",
			args: args{
				source: "testdata/invalid_file",
				file:   "example.json",
			},
			err: ErrInvalidGzipFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.args.source, tt.args.file)
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
