package plugins

import (
	"context"

	"github.com/spf13/cobra"
)

func (m *Manager) inject(ctx context.Context, command *cobra.Command) error {
	for _, cmdPlugin := range m.cmdPlugins {
		for _, command := range cmdPlugin.Registry() {
			c := &cobra.Command{
				Use:   command.Usage(),
				Short: command.ShortDesc(),
				Long:  command.LongDesc(),
				Args:  cobra.ExactArgs(command.NumArgs()),
			}
		}
	}

	for _, hookPlugin := range m.hookPlugins {
		for _, hook := range hookPlugin.Registry() {

		}
	}

	return nil
}
