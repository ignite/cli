package cosmosgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "test module path",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6",
			want: "CosmosSdk",
		},
		{
			name: "test chain path",
			path: "/Users/danilopantani/Desktop/go/src/github.com/ignite/venus",
			want: "Venus",
		},
		{
			name: "test module path without version",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway",
			want: "GrpcGateway",
		},
		{
			name: "test module path with broken version",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.$",
			want: "GrpcGateway",
		},
		{
			name: "test module path",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/cosmossdk.io/x/evidence@v0.1.0",
			want: "Evidence",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractName(tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}
