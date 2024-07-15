package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(ctx context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "run_command",
	}, nil
}

func (app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	return api.RunCommand(ctx, "version")
}

func (app) ExecuteHookPre(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(ctx context.Context, h *plugin.ExecutedHook, api plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"run_command": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
