package ignitecmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
)

func TestEnsureDefaultPlugins(t *testing.T) {
	tests := []struct {
		name                 string
		cfg                  *pluginsconfig.Config
		expectAddedInCommand bool
	}{
		{
			name:                 "empty config",
			cfg:                  &pluginsconfig.Config{},
			expectAddedInCommand: true,
		},
		{
			name: "config with default plugin",
			cfg: &pluginsconfig.Config{
				Plugins: []pluginsconfig.Plugin{{
					Path: "github.com/ignite/cli-plugin-network@v42",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "ignite"}

			ensureDefaultPlugins(cmd, tt.cfg)

			expectedCmd := findCommandByPath(cmd, "network")
			if tt.expectAddedInCommand {
				assert.NotNil(t, expectedCmd)
			} else {
				assert.Nil(t, expectedCmd)
			}
		})
	}
}
