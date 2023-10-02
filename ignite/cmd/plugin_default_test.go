package ignitecmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	appsconfig "github.com/ignite/cli/ignite/config/apps"
)

func TestEnsureDefaultPlugins(t *testing.T) {
	tests := []struct {
		name                 string
		cfg                  *appsconfig.Config
		expectAddedInCommand bool
	}{
		{
			name:                 "should add because absent from config",
			cfg:                  &appsconfig.Config{},
			expectAddedInCommand: true,
		},
		{
			name: "should not add because already present in config",
			cfg: &appsconfig.Config{
				Apps: []appsconfig.App{{
					Path: "github.com/ignite/cli-plugin-network@v42",
				}},
			},
			expectAddedInCommand: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "ignite"}

			ensureDefaultPlugins(cmd, tt.cfg)

			expectedCmd := findCommandByPath(cmd, "ignite network")
			if tt.expectAddedInCommand {
				assert.NotNil(t, expectedCmd)
			} else {
				assert.Nil(t, expectedCmd)
			}
		})
	}
}
