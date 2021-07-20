package module

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

func TestDiscover(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	testDataPath := filepath.Join(goPath,
		"src/github.com/tendermint/starport/starport/pkg/cosmosanalysis/module/testdata")

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
				sourcePath: filepath.Join(testDataPath, "planet"),
				protoDir:   "proto",
			},
			want: []Module{
				{Pkg: protoanalysis.Package{}},
			},
		}, {
			name: "test no proto folder",
			args: args{
				sourcePath: filepath.Join(testDataPath, "planet"),
				protoDir:   "",
			},
			want: []Module{
				{Pkg: protoanalysis.Package{}},
			},
		}, {
			name: "test invalid proto folder",
			args: args{
				sourcePath: filepath.Join(testDataPath, "planet"),
				protoDir:   "invalid",
			},
			want: nil,
		}, {
			name: "test invalid folder",
			args: args{
				sourcePath: filepath.Join(testDataPath, "invalid"),
				protoDir:   "",
			},
			want: []Module{},
		}, {
			name: "test invalid main and proto folder",
			args: args{
				sourcePath: filepath.Join(testDataPath, "../../.."),
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
