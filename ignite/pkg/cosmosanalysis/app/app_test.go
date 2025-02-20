package app

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

var (
	//go:embed testdata/app_minimal.go
	AppMinimalFile []byte
	//go:embed testdata/app_generic.go
	AppGenericFile []byte
	//go:embed testdata/no_app.go
	NoAppFile []byte
	//go:embed testdata/two_app.go
	TwoAppFile []byte
	//go:embed testdata/app_di.go
	AppDepinject []byte
)

func TestCheckKeeper(t *testing.T) {
	tests := []struct {
		name          string
		appFile       []byte
		keeperName    string
		expectedError string
	}{
		{
			name:       "minimal app",
			appFile:    AppMinimalFile,
			keeperName: "FooKeeper",
		},
		{
			name:       "generic app",
			appFile:    AppGenericFile,
			keeperName: "FooKeeper",
		},
		{
			name:          "no app",
			appFile:       NoAppFile,
			keeperName:    "FooKeeper",
			expectedError: "app.go should contain a single app (got 0)",
		},
		{
			name:          "two apps",
			appFile:       TwoAppFile,
			keeperName:    "FooKeeper",
			expectedError: "app.go should contain a single app (got 2)",
		},
		{
			name:       "app depinject",
			appFile:    AppDepinject,
			keeperName: "FooKeeper",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "app.go")
			err := os.WriteFile(tmpFile, tt.appFile, 0o644)
			require.NoError(t, err)

			err = CheckKeeper(tmpDir, tt.keeperName)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFindRegisteredModules(t *testing.T) {
	basicModules := []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/cosmos/cosmos-sdk/x/gov",
		"github.com/username/test/x/foo",
		"github.com/cosmos/cosmos-sdk/x/auth/tx",
		"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
		"github.com/cosmos/cosmos-sdk/client/grpc/node",
	}

	cases := []struct {
		name            string
		path            string
		expectedModules []string
	}{
		{
			name:            "new basic manager with only a app.go",
			path:            "testdata/modules/single_app",
			expectedModules: append(basicModules, "github.com/cosmos/ibc-go/v7/modules/core"),
		},
		{
			name:            "with runtime api routes",
			path:            "testdata/modules/runtime",
			expectedModules: basicModules,
		},
		{
			name: "with app_config.go file",
			path: "testdata/modules/app_config",
			expectedModules: []string{
				"cosmossdk.io/x/circuit",
				"cosmossdk.io/x/evidence",
				"cosmossdk.io/x/feegrant/module",
				"cosmossdk.io/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/x/auth/tx/config",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/consensus",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/group/module",
				"github.com/cosmos/cosmos-sdk/x/mint",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/ignite/mars/x/mars",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/username/test/x/foo",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
				"github.com/cosmos/cosmos-sdk/client/grpc/node",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindRegisteredModules(tt.path)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.expectedModules, got)
		})
	}
}

func TestDiscoverModules(t *testing.T) {
	basicModules := []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/cosmos/cosmos-sdk/x/gov",
		"github.com/username/test/x/foo",
		"github.com/cosmos/cosmos-sdk/x/auth/tx",
		"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
		"github.com/cosmos/cosmos-sdk/client/grpc/node",
	}

	cases := []struct {
		name            string
		path            string
		expectedModules []string
	}{
		{
			name:            "new basic manager with only a app.go",
			path:            "testdata/modules/single_app",
			expectedModules: append(basicModules, "github.com/cosmos/ibc-go/v7/modules/core"),
		},
		{
			name:            "with app_config.go file",
			path:            "testdata/modules/app_config",
			expectedModules: basicModules,
		},
		{
			name:            "with runtime api routes",
			path:            "testdata/modules/runtime",
			expectedModules: basicModules,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			appPkg, _, err := xast.ParseDir(tt.path)
			require.NoError(t, err)

			got := make([]string, 0)
			for _, f := range appPkg.Files {
				fileImports := goanalysis.FormatImports(f)
				modules, err := DiscoverModules(f, tt.path, fileImports)
				require.NoError(t, err)
				if modules != nil {
					got = append(got, modules...)
				}
			}
			require.ElementsMatch(t, tt.expectedModules, got)
		})
	}
}

func Test_removeKeeperPkgPath(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "test controller keeper",
			arg:  "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper",
			want: "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
		},
		{
			name: "test controller",
			arg:  "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller",
			want: "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
		},
		{
			name: "test keeper",
			arg:  "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/keeper",
			want: "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
		},
		{
			name: "test controller keeper",
			arg:  "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/controller/keeper",
			want: "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
		},
		{
			name: "test host controller keeper",
			arg:  "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/host/keeper",
			want: "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeKeeperPkgPath(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}
