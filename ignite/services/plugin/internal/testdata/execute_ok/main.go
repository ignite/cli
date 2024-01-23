package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(ctx context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "execute_ok",
	}, nil
}

func (app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	chaininfo, _ := api.GetChainInfo(ctx)
	fmt.Printf("ok args=%s chaininfo=%s\n", cmd.Args, chaininfo.String())
	return nil
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
			"execute_ok": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
