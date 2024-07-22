package xembed

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var fsProtoTest embed.FS

func TestFileList(t *testing.T) {
	type args struct {
		efs  embed.FS
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "root folder",
			args: args{
				efs:  fsProtoTest,
				path: ".",
			},
			want: []string{
				"testdata/subtestdata/subfile.txt",
				"testdata/subtestdata/subtestdata/subfile2.txt",
				"testdata/test.txt",
			},
		},
		{
			name: "testdata folder",
			args: args{
				efs:  fsProtoTest,
				path: "testdata",
			},
			want: []string{
				"subtestdata/subfile.txt",
				"subtestdata/subtestdata/subfile2.txt",
				"test.txt",
			},
		},
		{
			name: "sub testdata folder",
			args: args{
				efs:  fsProtoTest,
				path: "testdata/subtestdata",
			},
			want: []string{
				"subfile.txt",
				"subtestdata/subfile2.txt",
			},
		},
		{
			name: "sub sub testdata folder", //nolint:dupword
			args: args{
				efs:  fsProtoTest,
				path: "testdata/subtestdata/subtestdata",
			},
			want: []string{"subfile2.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileList(tt.args.efs, tt.args.path)
			if tt.err != nil {
				require.Error(t, err)
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}
