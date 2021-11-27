package pluginsrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (m *Manager) InjectPlugins(ctx context.Context, rootCmd *cobra.Command, args []string) (bool, error) {
	fmt.Println("💉 Injecting plugins...")

	if len(m.cmdPlugins) == 0 || len(m.hookPlugins) == 0 {
		if err := m.extractPlugins(ctx, rootCmd); err != nil {
			return false, err
		}
	}

	for _, cmd := range m.cmdPlugins {
		targetCommand, _, err := rootCmd.Find(cmd.ParentCommand)
		if err != nil {
			return false, err
		}

		if targetCommand != nil && len(args) > 0 {
			c := &cobra.Command{
				Use:   cmd.Usage,
				Short: cmd.ShortDesc,
				Long:  cmd.LongDesc,
				Args:  cobra.ExactArgs(cmd.NumArgs),
				RunE:  cmd.Exec,
			}

			// Cancel the root command, execute the new command
			targetCommand.AddCommand(c)

			baseUsage := strings.Split(cmd.Usage, " ")[0]
			if args[0] != baseUsage {
				return false, ErrCommandNotFound
			}

			reloadedTargetCommand, _, err := targetCommand.Find([]string{baseUsage})
			if err != nil {
				return false, err
			}

			err = reloadedTargetCommand.Execute()
			if err != nil {
				return true, err
			}

			return true, nil
		}
	}

	for _, hook := range m.hookPlugins {
		targetCommand, _, err := rootCmd.Find(hook.ParentCommand)
		if err != nil {
			return false, err
		}

		if targetCommand != nil {
			switch HookType(hook.HookType) {
			case PreRunHook:
				targetCommand.PreRunE = hook.PreRun
			case PostRunHook:
				targetCommand.PostRunE = hook.PostRun
			default:
				targetCommand.PreRunE = hook.PreRun
				targetCommand.PostRunE = hook.PostRun
			}
		}
	}

	return false, nil
}
