package module

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

func TestDiscover(t *testing.T) {
	type args struct {
		sourcePath string
		protoDir   string
	}
	tests := []struct {
		name string
		args args
		want []Module
	}{
		{
			name: "test valid",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "proto",
			},
			want: []Module{
				{Pkg: protoanalysis.Package{}},
			},
		}, {
			name: "test no proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "",
			},
			want: []Module{
				{Pkg: protoanalysis.Package{}},
			},
		}, {
			name: "test invalid proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "invalid",
			},
		}, {
			name: "test invalid folder",
			args: args{
				sourcePath: "testdata/invalid",
				protoDir:   "",
			},
			want: []Module{},
		}, {
			name: "test invalid main and proto folder",
			args: args{
				sourcePath: "../../..",
				protoDir:   "proto",
			},
			want: []Module{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Discover(context.Background(), tt.args.sourcePath, tt.args.protoDir)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
