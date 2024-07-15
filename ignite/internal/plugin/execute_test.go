package plugininternal_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	ignitecmd "github.com/ignite/cli/v29/ignite/cmd"
	plugininternal "github.com/ignite/cli/v29/ignite/internal/plugin"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	"github.com/ignite/cli/v29/ignite/services/plugin/mocks"
)

func TestPluginExecute(t *testing.T) {
	cmd, cleanup, err := ignitecmd.New(context.Background())
	require.NoError(t, err)
	t.Cleanup(cleanup)

	tests := []struct {
		name           string
		pluginPath     string
		expectedOutput string
		expectedError  string
	}{
		{
			name:          "fail: plugin doesnt exist",
			pluginPath:    "/not/exists",
			expectedError: "local app path \"/not/exists\" not found: stat /not/exists: no such file or directory",
		},
		{
			name:           "ok: plugin execute ok",
			pluginPath:     "testdata/execute_ok",
			expectedOutput: "ok args=[arg1 arg2] chainid=id appPath=apppath configPath=configpath home=home rpcAddress=rpcPublicAddress\n",
		},
		{
			name:          "ok: plugin execute fail",
			pluginPath:    "testdata/execute_fail",
			expectedError: "fail",
		},
		{
			name:          "ok: plugin run command execute",
			pluginPath:    "testdata/run_command",
			expectedError: "....",
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

			out, err := plugininternal.Execute(
				context.Background(),
				pluginPath,
				[]string{"arg1", "arg2"},
				plugin.WithChain(chainer),
				plugin.WithCmd(cmd),
			)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, out)
		})
	}
}
