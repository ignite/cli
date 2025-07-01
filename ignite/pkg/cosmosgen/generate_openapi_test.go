package cosmosgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func Test_extractRootModulePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "test cosmos-sdk path",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6/proto/cosmos/distribution/v1beta1",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6",
		},
		{
			name: "test cosmos-sdk module proto path",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6/x/bank",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.50.6",
		},
		{
			name: "test ibc path",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/ibc-go/v8@v8.2.0/proto/ibc/applications/interchain_accounts/controller/v1",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/cosmos/ibc-go/v8@v8.2.0",
		},
		{
			name: "test chain path",
			path: "/Users/ignite/Desktop/go/src/github.com/ignite/venus",
			want: "/Users/ignite/Desktop/go/src/github.com/ignite/venus",
		},
		{
			name: "test module path without version",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/proto/applications",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/proto/applications",
		},
		{
			name: "test module path with broken version",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.$/controller",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.$/controller",
		},
		{
			name: "test module path with v2 version",
			path: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.19.1/proto/files",
			want: "/Users/ignite/Desktop/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.19.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRootModulePath(tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateOpenAPI(t *testing.T) {
	require := require.New(t)
	testdataDir := "testdata"
	appDir := filepath.Join(testdataDir, "testchain")
	openAPIFile := filepath.Join(appDir, "docs", "static", "openapi.json")

	cacheStorage, err := cache.NewStorage(filepath.Join(t.TempDir(), "cache.db"))
	require.NoError(err)

	buf, err := cosmosbuf.New(cacheStorage, t.Name())
	require.NoError(err)

	// Use module discovery to collect test module proto.
	m, err := module.Discover(t.Context(), appDir, appDir, module.WithProtoDir("proto"))
	require.NoError(err, "failed to discover module")
	require.Len(m, 1, "expected exactly one module to be discovered")

	g := &generator{
		appPath:      appDir,
		protoDir:     "proto",
		goModPath:    "go.mod",
		cacheStorage: cacheStorage,
		buf:          buf,
		appModules:   m,
		opts: &generateOptions{
			specOut: openAPIFile,
		},
	}

	err = g.generateOpenAPISpec(t.Context())
	if err != nil && !errors.Is(err, dirchange.ErrNoFile) {
		require.NoError(err, "failed to generate OpenAPI spec")
	}

	// compare generated OpenAPI spec with golden files
	goldenFile := filepath.Join(testdataDir, "expected_files", "openapi", "openapi.json")
	gold, err := os.ReadFile(goldenFile)
	require.NoError(err, "failed to read golden file: %s", goldenFile)

	gotBytes, err := os.ReadFile(openAPIFile)
	require.NoError(err, "failed to read generated file: %s", openAPIFile)
	require.Equal(string(gold), string(gotBytes), "generated OpenAPI spec does not match golden file")
}
