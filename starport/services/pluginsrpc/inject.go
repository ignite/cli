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
		if err := m.extractPlugins(ctx, rootCmd, args); err != nil {
			return false, err
		}
	}

	// We have to enable the combination of hooks and commands
	finished, err := m.InjectCommands(ctx, rootCmd, args)
	if err != nil {
		return finished, err
	}

	err = m.InjectHooks(ctx, rootCmd)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (m *Manager) InjectCommands(ctx context.Context, rootCmd *cobra.Command, args []string) (bool, error) {
	outputDir, err := formatPluginHome(m.ChainId, "output")
	if err != nil {
		return false, err
	}

	if len(m.cmdPlugins) == 0 {
		if err := m.extractCommandPlugins(ctx, outputDir, rootCmd); err != nil {
			return false, err
		}
	}

	for _, cmdPlugin := range m.cmdPlugins {
		targetCommand, _, err := rootCmd.Find(cmdPlugin.ParentCommand)
		if err != nil {
			return false, err
		}

		if targetCommand != nil && len(args) > 0 {
			c := &cobra.Command{
				Use:   cmdPlugin.Usage,
				Short: cmdPlugin.ShortDesc,
				Long:  cmdPlugin.LongDesc,
				Args:  cobra.ExactArgs(cmdPlugin.NumArgs),
				RunE:  cmdPlugin.Exec,
			}

			baseUsage := strings.Split(cmdPlugin.Usage, " ")[0]
			if args[0] != baseUsage {
				return false, ErrCommandNotFound
			}

			// Cancel the root command, execute the new command
			targetCommand.AddCommand(c)

			reloadedTargetCommand, _, err := targetCommand.Find([]string{baseUsage})
			if err != nil {
				return false, err
			}

			// This is what calls prerun again
			err = reloadedTargetCommand.Execute()
			if err != nil {
				return true, err
			}

			return true, nil
		}
	}

	return false, nil
}

func (m *Manager) InjectHooks(ctx context.Context, rootCmd *cobra.Command) error {
	outputDir, err := formatPluginHome(m.ChainId, "output")
	if err != nil {
		return err
	}

	if len(m.hookPlugins) == 0 {
		// hook specific extract
		if err := m.extractHookPlugins(ctx, outputDir, rootCmd); err != nil {
			return err
		}
	}

	for _, hookPlugin := range m.hookPlugins {
		targetCommand, _, err := rootCmd.Find(hookPlugin.ParentCommand)
		if err != nil {
			return err
		}

		if targetCommand != nil {
			switch HookType(hookPlugin.HookType) {
			case PreRunHook:
				targetCommand.PreRunE = hookPlugin.PreRun
			case PostRunHook:
				targetCommand.PostRunE = hookPlugin.PostRun
			default:
				targetCommand.PreRunE = hookPlugin.PreRun
				targetCommand.PostRunE = hookPlugin.PostRun
			}
		}
	}

	return nil
}
