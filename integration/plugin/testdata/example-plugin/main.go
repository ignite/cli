package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/plugin"
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
				Flags: []*plugin.Flag{
					{Name: "my-flag", Type: plugin.FlagTypeString, Usage: "my flag description"},
				},
				PlaceCommandUnder: "ignite",
			},
		},
		Hooks: []*plugin.Hook{},
	}, nil
}

func (p) Execute(_ context.Context, cmd *plugin.ExecutedCommand) error {
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
	return nil
}

func (p) ExecuteHookPre(_ context.Context, h *plugin.ExecutedHook) error {
	fmt.Printf("Executing hook pre %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookPost(_ context.Context, h *plugin.ExecutedHook) error {
	fmt.Printf("Executing hook post %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookCleanUp(_ context.Context, h *plugin.ExecutedHook) error {
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
