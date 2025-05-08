package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/v29/ignite/services/plugin"
)

type p struct{}

func (p) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "example-plugin",
		Commands: []*plugin.Command{
			{
				Use:   "example-plugin",
				Short: "Explain what the command is doing...",
				Long:  "Long description goes here...",
				Flags: plugin.Flags{
					{Name: "my-flag", Type: plugin.FlagTypeString, Usage: "my flag description"},
				},
				PlaceCommandUnder: "ignite",
			},
		},
		Hooks: []*plugin.Hook{},
	}, nil
}

func (p) Execute(ctx context.Context, cmd *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	fmt.Printf("Hello I'm the example-plugin plugin\n")
	fmt.Printf("My executed command: %q\n", cmd.Path)
	fmt.Printf("My args: %v\n", cmd.Args)

	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	myFlag, _ := flags.GetString("my-flag")
	fmt.Printf("My flags: my-flag=%q\n", myFlag)
	fmt.Printf("My config parameters: %v\n", cmd.With)

	fmt.Println(api.GetChainInfo(ctx))
	fmt.Println(api.GetIgniteInfo(ctx))

	return nil
}

func (p) ExecuteHookPre(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("Executing hook pre %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookPost(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("Executing hook post %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookCleanUp(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
	fmt.Printf("Executing hook cleanup %q\n", h.Hook.GetName())
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"example-plugin": plugin.NewGRPC(&p{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
