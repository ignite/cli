package plugininternal

import (
	"context"
	"testing"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/stretchr/testify/require"
)

func TestPluginExecute(t *testing.T) {
	tests := []struct {
		name           string
		scaffoldPlugin func(t *testing.T) string
		expectedError  string
	}{
		{
			name: "fail: plugin doesnt exist",
			scaffoldPlugin: func(t *testing.T) string {
				return "/not/exists"
			},
			expectedError: "local app path \"/not/exists\" not found: stat /not/exists: no such file or directory",
		},
		{
			name: "ok: plugin exists",
			scaffoldPlugin: func(t *testing.T) string {
				path, err := plugin.Scaffold(context.Background(), t.TempDir(), "foo", false)
				require.NoError(t, err)
				return path
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.scaffoldPlugin(t)

			err := Execute(context.Background(), path, "args")

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}
