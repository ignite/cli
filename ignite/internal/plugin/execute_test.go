package plugininternal

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/ignite/cli/v28/ignite/services/plugin/mocks"
)

func TestPluginExecute(t *testing.T) {
	tests := []struct {
		name          string
		pluginPath    string
		expectedOut   string
		expectedError string
	}{
		{
			name:          "fail: plugin doesnt exist",
			pluginPath:    "/not/exists",
			expectedError: "local app path \"/not/exists\" not found: stat /not/exists: no such file or directory",
		},
		{
			name:        "ok: plugin execute ok ",
			pluginPath:  "testdata/execute_ok",
			expectedOut: "ok args=[arg1 arg2] chainid=id appPath=apppath configPath=configpath home=home rpcAddress=rpcPublicAddress\n",
		},
		{
			name:          "ok: plugin execute fail ",
			pluginPath:    "testdata/execute_fail",
			expectedError: "fail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pluginPath := tt.pluginPath
			if !strings.HasPrefix(pluginPath, "/") {
				// add working dir to relative paths
				wd, err := os.Getwd()
				require.NoError(t, err)
				pluginPath = filepath.Join(wd, pluginPath)
			}
			chainer := mocks.NewChainerInterface(t)
			chainer.EXPECT().ID().Return("id", nil).Maybe()
			chainer.EXPECT().AppPath().Return("apppath").Maybe()
			chainer.EXPECT().ConfigPath().Return("configpath").Maybe()
			chainer.EXPECT().Home().Return("home", nil).Maybe()
			chainer.EXPECT().RPCPublicAddress().Return("rpcPublicAddress", nil).Maybe()

			out, err := Execute(
				context.Background(),
				pluginPath,
				[]string{"arg1", "arg2"},
				plugin.WithChain(chainer),
			)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOut, out)
		})
	}
}
