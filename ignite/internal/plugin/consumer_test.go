package plugininternal

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/services/plugin"
	"github.com/ignite/cli/v29/ignite/services/plugin/mocks"
)

func TestConsumerPlugin(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setup          func(*testing.T, string)
		expectedOutput string
		expectedError  string
	}{
		{
			name:          "fail: missing arg",
			expectedError: "missing argument",
		},
		{
			name:          "fail: invalid arg",
			args:          []string{"xxx"},
			expectedError: "invalid argument \"xxx\"",
		},
		{
			name:          "fail: writeGenesis w/o priv_validator_key.json",
			args:          []string{"writeGenesis"},
			expectedError: "open .*/config/priv_validator_key.json: no such file or directory",
		},
		{
			name: "fail: writeFenesis w/o genesis.json",
			args: []string{"writeGenesis"},
			setup: func(t *testing.T, path string) {
				t.Helper()
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/consumer/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)
			},
			expectedError: ".*/config/genesis.json does not exist, run `init` first",
		},

		{
			name: "ok: writeGenesis",
			args: []string{"writeGenesis"},
			setup: func(t *testing.T, path string) {
				t.Helper()
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/consumer/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)

				// Add genesis.json to path
				bz, err = os.ReadFile("testdata/consumer/config/genesis.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "genesis.json"), bz, 0o777)
				require.NoError(t, err)
			},
		},
		{
			name:           "ok: isInitialized returns false",
			args:           []string{"isInitialized"},
			expectedOutput: "false",
		},
		{
			name: "ok: isInitialized returns true",
			args: []string{"isInitialized"},
			setup: func(t *testing.T, path string) {
				t.Helper()
				// isInitialized returns true if there's a consumer genesis with an
				// InitialValSet length != 0
				// Add priv_validator_key.json to path
				bz, err := os.ReadFile("testdata/consumer/config/priv_validator_key.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "priv_validator_key.json"), bz, 0o777)
				require.NoError(t, err)

				// Add genesis.json to path
				bz, err = os.ReadFile("testdata/consumer/config/genesis.json")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(path, "config", "genesis.json"), bz, 0o777)
				require.NoError(t, err)

				// Call writeGenesis to create the genesis
				chainer := mocks.NewChainerInterface(t)
				chainer.EXPECT().ID().Return("id", nil).Maybe()
				chainer.EXPECT().AppPath().Return("apppath").Maybe()
				chainer.EXPECT().ConfigPath().Return("configpath").Maybe()
				chainer.EXPECT().Home().Return(path, nil).Maybe()
				chainer.EXPECT().RPCPublicAddress().Return("rpcPublicAddress", nil).Maybe()
				_, err = Execute(context.Background(), PluginConsumerPath, []string{"writeGenesis"}, plugin.WithChain(chainer))
				require.NoError(t, err)
			},
			expectedOutput: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homePath := t.TempDir()
			err := os.MkdirAll(filepath.Join(homePath, "config"), 0o777)
			require.NoError(t, err)
			chainer := mocks.NewChainerInterface(t)
			chainer.EXPECT().ID().Return("id", nil).Maybe()
			chainer.EXPECT().AppPath().Return("apppath").Maybe()
			chainer.EXPECT().ConfigPath().Return("configpath").Maybe()
			chainer.EXPECT().Home().Return(homePath, nil).Maybe()
			chainer.EXPECT().RPCPublicAddress().Return("rpcPublicAddress", nil).Maybe()
			if tt.setup != nil {
				tt.setup(t, homePath)
			}

			out, err := Execute(
				context.Background(),
				PluginConsumerPath,
				tt.args,
				plugin.WithChain(chainer),
			)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Regexp(t, tt.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, out)
		})
	}
}
