package plugininternal

import (
	"bytes"
	"context"
	"time"

	"google.golang.org/grpc/status"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// Execute starts and executes a plugin, then shutdowns it.
func Execute(ctx context.Context, path string, args []string, options ...plugin.APIOption) (string, error) {
	var buf bytes.Buffer
	plugins, err := plugin.Load(
		ctx,
		[]pluginsconfig.Plugin{{Path: path}},
		plugin.RedirectStdout(&buf),
	)
	if err != nil {
		return "", err
	}
	defer plugins[0].KillClient()
	if plugins[0].Error != nil {
		return "", plugins[0].Error
	}
	err = plugins[0].Interface.Execute(
		ctx,
		&plugin.ExecutedCommand{Args: args},
		plugin.NewClientAPI(options...),
	)
	if err != nil {
		// Extract the rpc status message and create a simple error from it.
		// We don't want Execute to return rpc errors.
		err = errors.New(status.Convert(err).Message())
	}
	// NOTE(tb): This pause gives enough time for go-plugin to sync the
	// output from stdout/stderr of the plugin. Without that pause, this
	// output can be discarded and absent from buf.
	time.Sleep(100 * time.Millisecond)
	return buf.String(), err
}
