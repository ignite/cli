package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	_ "embed"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/app"
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
	//go:embed testdata/app_full.go
	AppFullFile []byte
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "app.go")
			err := os.WriteFile(tmpFile, tt.appFile, 0644)
			require.NoError(t, err)

			err = app.CheckKeeper(tmpDir, tt.keeperName)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetRegisteredModules(t *testing.T) {
	tmpDir := t.TempDir()

	tmpFile := filepath.Join(tmpDir, "app.go")
	err := os.WriteFile(tmpFile, AppFullFile, 0644)
	require.NoError(t, err)

	tmpNoAppFile := filepath.Join(tmpDir, "someOtherFile.go")
	err = os.WriteFile(tmpNoAppFile, NoAppFile, 0644)
	require.NoError(t, err)

	registeredModules, err := app.FindRegisteredModules(tmpDir)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/genutil",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/capability",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/cosmos/cosmos-sdk/x/mint",
		"github.com/cosmos/cosmos-sdk/x/distribution",
		"github.com/cosmos/cosmos-sdk/x/gov",
		"github.com/cosmos/cosmos-sdk/x/params",
		"github.com/cosmos/cosmos-sdk/x/crisis",
		"github.com/cosmos/cosmos-sdk/x/slashing",
		"github.com/cosmos/cosmos-sdk/x/feegrant/module",
		"github.com/cosmos/ibc-go/v5/modules/core",
		"github.com/cosmos/cosmos-sdk/x/upgrade",
		"github.com/cosmos/cosmos-sdk/x/evidence",
		"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
		"github.com/cosmos/cosmos-sdk/x/auth/vesting",
		"github.com/tendermint/testchain/x/testchain",
		"github.com/tendermint/testchain/x/queryonlymod",
		"github.com/cosmos/cosmos-sdk/x/auth/tx",
		"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
	}, registeredModules)
}
