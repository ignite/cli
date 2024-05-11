package cosmosgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractRootModulePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "test cosmos-sdk path",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6/proto/cosmos/distribution/v1beta1",
			want: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6",
		},
		{
			name: "test ibc path",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/cosmos/ibc-go/v8@v8.2.0/proto/ibc/applications/interchain_accounts/controller/v1",
			want: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/cosmos/ibc-go/v8@v8.2.0",
		},
		{
			name: "test chain path",
			path: "/Users/danilopantani/Desktop/go/src/github.com/ignite/venus",
			want: "/Users/danilopantani/Desktop/go/src/github.com/ignite/venus",
		},
		{
			name: "test module path without version",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/proto/applications",
			want: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/proto/applications",
		},
		{
			name: "test module path with broken version",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.$/controller",
			want: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.$/controller",
		},
		{
			name: "test module path with v2 version",
			path: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.19.1/proto/files",
			want: "/Users/danilopantani/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.19.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRootModulePath(tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}
