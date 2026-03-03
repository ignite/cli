package plugininternal

import (
	"bytes"
	"context"
	"sync"

	"google.golang.org/grpc/status"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type synchronizedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *synchronizedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}

// Execute starts and executes a plugin, then shutdowns it.
func Execute(ctx context.Context, path string, args []string, options ...plugin.APIOption) (string, error) {
	var buf synchronizedBuffer
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
		return "", errors.New(status.Convert(err).Message())
	}

	plugins[0].KillClient()
	return buf.String(), err
}
