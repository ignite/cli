package plugininternal

import (
	"context"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// Execute starts and executes a plugin, then shutdowns it.
func Execute(ctx context.Context, path string, args ...string) error {
	plugins, err := plugin.Load(ctx, []pluginsconfig.Plugin{{Path: path}})
	if err != nil {
		return err
	}
	defer plugins[0].KillClient()
	if plugins[0].Error != nil {
		return plugins[0].Error
	}
	return plugins[0].Interface.Execute(ctx, &plugin.ExecutedCommand{Args: args}, plugin.NewClientAPI())
}
