package pluginsrpc

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/lukerhoads/plugintypes"
	"github.com/spf13/cobra"
)

type InitArgs struct {
	ctx context.Context
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var BasePluginMap = map[string]plugin.Plugin{
	"command_map": &plugintypes.CommandMapperPlugin{},
	"hook_map":    &plugintypes.HookMapperPlugin{},
}

type ExtractedCommandModule struct {
	ParentCommand []string
	Name          string
	Usage         string
	ShortDesc     string
	LongDesc      string
	NumArgs       int
	Exec          func(*cobra.Command, []string) error
}
