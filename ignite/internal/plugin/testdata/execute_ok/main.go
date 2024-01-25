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
	c, _ := api.GetChainInfo(ctx)
	fmt.Printf(
		"ok args=%s chainid=%s appPath=%s configPath=%s home=%s rpcAddress=%s\n",
		cmd.Args, c.ChainId, c.AppPath, c.ConfigPath, c.Home, c.RpcAddress,
	)
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
