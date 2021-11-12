package starportcmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/cmd/mocks"
)

func Test_ExecPluginInstall(t *testing.T) {
	tests := []struct {
		Desc        string
		Args        []string
		CallHandler bool
	}{
		{
			Desc:        "Success",
			Args:        []string{"plugin", "install", "test-repo"},
			CallHandler: true,
		},

		{
			Desc:        "No extra arguments",
			Args:        []string{"plugin", "install"},
			CallHandler: false,
		},
	}

	for _, test := range tests {
		// Prepare test
		ctx := context.Background()
		c := New(ctx)

		c.SetArgs(test.Args)

		// Mocks
		mockPluginCmdHandler := mocks.PluginCmdHandler{}

		if test.CallHandler {
			installCmd := findCommand(c, test.Args[:len(test.Args)-1])
			mockPluginCmdHandler.
				On("HandleInstall", installCmd, []string{test.Args[len(test.Args)-1]}).
				Return(nil)
		}

		// Apply mock instance
		pluginHandler = &mockPluginCmdHandler

		// Test
		c.Execute()

		// Asserts
		mockPluginCmdHandler.AssertExpectations(t)
	}
}

func findCommand(cmd *cobra.Command, args []string) *cobra.Command {
	for _, c := range cmd.Commands() {
		if c.Use == args[0] {
			if len(args) == 1 {
				return c
			}
			return findCommand(c, args[1:])
		}
	}

	return nil
}
