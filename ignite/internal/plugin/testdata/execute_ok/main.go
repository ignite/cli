package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type app struct{}

func (app) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "execute_ok",
	}, nil
}

func (app) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	c, err := api.GetChainInfo(ctx)
	fmt.Printf(
		"ok args=%s chainid=%s appPath=%s configPath=%s home=%s rpcAddress=%s\n",
		cmd.Args, c.ChainId, c.AppPath, c.ConfigPath, c.Home, c.RpcAddress,
	)
	if err != nil {
		return errors.Errorf("failed to get chain info: %w", err)
	}

	i, err := api.GetIgniteInfo(ctx)
	fmt.Printf(
		"ok args=%s cliVersion=%s goVersion=%s sdkVersion=%s bufVersion=%s buildDate=%s "+
			"sourceHash=%s configVersion=%s os=%s arch=%s buildFromSource=%t\n",
		cmd.Args, i.CliVersion, i.GoVersion, i.SdkVersion, i.BufVersion, i.BuildDate, i.SourceHash, i.ConfigVersion,
		i.Os, i.Arch, i.BuildFromSource,
	)
	if err != nil {
		return errors.Errorf("failed to get ignite info: %w", err)
	}
	return nil
}

func (app) ExecuteHookPre(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookPost(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (app) ExecuteHookCleanUp(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
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
