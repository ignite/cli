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

type ExtractedHookModule struct {
	ParentCommand []string
	Name          string
	HookType      string
	PreRun        func(*cobra.Command, []string) error
	PostRun       func(*cobra.Command, []string) error
}

type PluginState uint32

const (
	Undefined PluginState = iota
	Configured
	Downloaded
	Built
)

func PluginStateFromString(state string) PluginState {
	switch state {
	case "configured":
		return Configured
	case "downloaded":
		return Downloaded
	case "built":
		return Built
	}

	return Undefined
}

func (p PluginState) String() string {
	switch p {
	case Configured:
		return "configured"
	case Downloaded:
		return "downloaded"
	case Built:
		return "built"
	}

	return "undefined"
}
