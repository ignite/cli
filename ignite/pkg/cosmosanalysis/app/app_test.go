package app

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/xast"
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
	//go:embed testdata/app_v2.go
	AppV2 []byte
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
			name:       "app v2",
			appFile:    AppV2,
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
	}

	cases := []struct {
		name            string
		path            string
		expectedModules []string
	}{
		{
			name:            "new basic manager with only a app.go",
			path:            "testdata/modules/single_app",
			expectedModules: basicModules,
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
				"github.com/cosmos/cosmos-sdk/x/auth/tx/config",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/consensus",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/group/module",
				"github.com/cosmos/cosmos-sdk/x/mint",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/ignite/mars/x/mars",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/username/test/x/foo",
			},
		},
		{
			name: "gaia",
			path: "testdata/modules/gaia",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/cosmos/cosmos-sdk/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/cosmos-sdk/x/group",
				"github.com/cosmos/ibc-go/v5/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
				"github.com/gravity-devs/liquidity/v2/x/liquidity",
				"github.com/strangelove-ventures/packet-forward-middleware/v2/router",
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"github.com/cosmos/gaia/v8/x/icamauth",
				"github.com/cosmos/cosmos-sdk/client/docs/statik",
			},
		},
		{
			name: "crescent",
			path: "testdata/modules/crescent",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/crescent-network/crescent/v3/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/ibc-go/v2/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v2/modules/apps/transfer",
				"github.com/tendermint/budget/x/budget",
				"github.com/crescent-network/crescent/v3/x/farming",
				"github.com/crescent-network/crescent/v3/x/liquidity",
				"github.com/crescent-network/crescent/v3/x/liquidstaking",
				"github.com/crescent-network/crescent/v3/x/liquidfarming",
				"github.com/crescent-network/crescent/v3/x/claim",
				"github.com/crescent-network/crescent/v3/x/marketmaker",
				"github.com/crescent-network/crescent/v3/x/lpfarm",
				"github.com/crescent-network/crescent/v3/client/docs/statik",
			},
		},
		{
			name: "spn",
			path: "testdata/modules/spn",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/ignite/modules/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/ibc-go/v6/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v6/modules/apps/transfer",
				"github.com/tendermint/spn/x/participation",
				"github.com/ignite/modules/x/claim",
				"github.com/tendermint/spn/x/profile",
				"github.com/tendermint/spn/x/launch",
				"github.com/tendermint/spn/x/campaign",
				"github.com/tendermint/spn/x/monitoringc",
				"github.com/tendermint/spn/x/monitoringp",
				"github.com/tendermint/spn/x/reward",
				"github.com/tendermint/fundraising/x/fundraising",
			},
		},
		{
			name: "juno",
			path: "testdata/modules/juno",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/CosmosContracts/juno/v10/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/ibc-go/v3/modules/core",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v3/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/CosmWasm/wasmd/x/wasm",
				"github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts",
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

func TestFindKeepersModules(t *testing.T) {
	basicModules := []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/cosmos/cosmos-sdk/x/gov",
		"github.com/username/test/x/foo",
	}

	cases := []struct {
		name            string
		path            string
		expectedModules []string
	}{
		{
			name:            "new basic manager with only a app.go",
			path:            "testdata/modules/single_app",
			expectedModules: basicModules,
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
		{
			name: "gaia",
			path: "testdata/modules/gaia",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/cosmos/cosmos-sdk/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/cosmos-sdk/x/group",
				"github.com/cosmos/ibc-go/v5/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
				"github.com/gravity-devs/liquidity/v2/x/liquidity",
				"github.com/strangelove-ventures/packet-forward-middleware/v2/router",
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"github.com/cosmos/gaia/v8/x/icamauth",
			},
		},
		{
			name: "crescent",
			path: "testdata/modules/crescent",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/crescent-network/crescent/v3/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/ibc-go/v2/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v2/modules/apps/transfer",
				"github.com/tendermint/budget/x/budget",
				"github.com/crescent-network/crescent/v3/x/farming",
				"github.com/crescent-network/crescent/v3/x/liquidity",
				"github.com/crescent-network/crescent/v3/x/liquidstaking",
				"github.com/crescent-network/crescent/v3/x/liquidfarming",
				"github.com/crescent-network/crescent/v3/x/claim",
				"github.com/crescent-network/crescent/v3/x/marketmaker",
				"github.com/crescent-network/crescent/v3/x/lpfarm",
			},
		},
		{
			name: "spn",
			path: "testdata/modules/spn",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/ignite/modules/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/cosmos/ibc-go/v6/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v6/modules/apps/transfer",
				"github.com/tendermint/spn/x/participation",
				"github.com/ignite/modules/x/claim",
				"github.com/tendermint/spn/x/profile",
				"github.com/tendermint/spn/x/launch",
				"github.com/tendermint/spn/x/campaign",
				"github.com/tendermint/spn/x/monitoringc",
				"github.com/tendermint/spn/x/monitoringp",
				"github.com/tendermint/spn/x/reward",
				"github.com/tendermint/fundraising/x/fundraising",
			},
		},
		{
			name: "juno",
			path: "testdata/modules/juno",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/CosmosContracts/juno/v10/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/ibc-go/v3/modules/core",
				"github.com/cosmos/cosmos-sdk/x/feegrant",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v3/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/authz",
				"github.com/CosmWasm/wasmd/x/wasm",
				"github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			appPkg, _, err := xast.ParseDir(tt.path)
			require.NoError(t, err)

			got := make([]string, 0)
			for _, f := range appPkg.Files {
				fileImports := goanalysis.FormatImports(f)
				modules, err := FindKeepersModules(f, fileImports)
				require.NoError(t, err)
				if modules != nil {
					got = append(got, modules...)
				}
			}
			require.ElementsMatch(t, tt.expectedModules, got)
		})
	}
}

func Test_mergeImports(t *testing.T) {
	tests := []struct {
		name         string
		blankImports []string
		discovered   []string
		want         []string
	}{
		{
			name:         "test nil imports",
			blankImports: nil,
			discovered:   nil,
			want:         nil,
		},
		{
			name:         "test empty imports",
			blankImports: []string{},
			discovered:   []string{},
			want:         []string{},
		},
		{
			name:         "test only one blank import",
			blankImports: []string{"github.com/cosmos/cosmos-sdk/x/auth"},
			discovered:   []string{},
			want:         []string{"github.com/cosmos/cosmos-sdk/x/auth"},
		},
		{
			name:         "test only one discovered import",
			blankImports: []string{},
			discovered:   []string{"github.com/cosmos/cosmos-sdk/x/auth"},
			want:         []string{"github.com/cosmos/cosmos-sdk/x/auth"},
		},
		{
			name:         "test only one import",
			blankImports: []string{"github.com/cosmos/cosmos-sdk/x/auth"},
			discovered:   []string{"github.com/cosmos/cosmos-sdk/x/auth/keeper"},
			want:         []string{"github.com/cosmos/cosmos-sdk/x/auth"},
		},
		{
			name:         "test only one keeper import",
			blankImports: []string{"github.com/cosmos/cosmos-sdk/x/auth/module"},
			discovered:   []string{"github.com/cosmos/cosmos-sdk/x/auth/keeper"},
			want:         []string{"github.com/cosmos/cosmos-sdk/x/auth/module"},
		},
		{
			name: "test two keeper import",
			blankImports: []string{
				"github.com/cosmos/cosmos-sdk/x/auth/module",
				"github.com/cosmos/cosmos-sdk/x/bank/module",
			},
			discovered: []string{
				"github.com/cosmos/cosmos-sdk/x/auth/keeper",
				"github.com/cosmos/cosmos-sdk/x/bank/keeper",
			},
			want: []string{
				"github.com/cosmos/cosmos-sdk/x/auth/module",
				"github.com/cosmos/cosmos-sdk/x/bank/module",
			},
		},
		{
			name: "test two keeper import",
			blankImports: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
			},
			discovered: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper",
				"github.com/cosmos/cosmos-sdk/x/bank/keeper",
			},
			want: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"github.com/cosmos/cosmos-sdk/x/bank/keeper",
			},
		},
		{
			name: "test keeper imports",
			blankImports: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"cosmossdk.io/x/feegrant/module",
			},
			discovered: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/keeper",
				"github.com/cosmos/cosmos-sdk/x/bank/keeper",
				"cosmossdk.io/x/feegrant/types",
				"cosmossdk.io/x/feegrant",
				"cosmossdk.io/x/foo",
			},
			want: []string{
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"github.com/cosmos/cosmos-sdk/x/bank/keeper",
				"cosmossdk.io/x/feegrant/module",
				"cosmossdk.io/x/foo",
			},
		},
		{
			name: "test three keeper import",
			blankImports: []string{
				"github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts",
				"github.com/cosmos/ibc-go/modules/capability",
			},
			discovered: []string{
				"github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts",
				"github.com/cosmos/ibc-go/modules/capability",
			},
			want: []string{
				"github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts",
				"github.com/cosmos/ibc-go/modules/capability",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeImports(tt.blankImports, tt.discovered)
			require.ElementsMatch(t, tt.want, got)
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
