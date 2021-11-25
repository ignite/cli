package plugins

import (
	"context"

	"github.com/spf13/cobra"
)

func (m *Manager) InjectPlugins(ctx context.Context, rootCmd *cobra.Command) error {
	if len(m.cmdPlugins) == 0 || len(m.hookPlugins) == 0 {
		if err := m.extractPlugins(ctx, rootCmd); err != nil {
			return err
		}
	}

	// command.Find or command.Traverse?
	for _, cmdPlugin := range m.cmdPlugins {
		for _, cmd := range cmdPlugin.Registry() {
			targetCommand, _, err := rootCmd.Find(cmd.ParentCommand())
			if err != nil {
				return err
			}

			if targetCommand != nil {
				c := &cobra.Command{
					Use:   cmd.Usage(),
					Short: cmd.ShortDesc(),
					Long:  cmd.LongDesc(),
					Args:  cobra.ExactArgs(cmd.NumArgs()),
					RunE:  cmd.Exec,
				}

				targetCommand.AddCommand(c)
			}
		}
	}

	for _, hookPlugin := range m.hookPlugins {
		for _, hook := range hookPlugin.Registry() {
			targetCommand, _, err := rootCmd.Find(hook.ParentCommand())
			if err != nil {
				return err
			}

			if targetCommand != nil {
				switch hook.Type() {
				case "pre":
					targetCommand.PreRunE = hook.PreRun
				case "post":
					targetCommand.PostRunE = hook.PostRun
				default:
					targetCommand.PreRunE = hook.PreRun
					targetCommand.PostRunE = hook.PostRun
				}
			}
		}
	}

	return nil
}
