package plugins

import "github.com/lukerhoads/plugintypes"

type CmdPlugin interface {
	Module
	Registry() map[string]plugintypes.Command
}

type HookPlugin interface {
	Module
	Registry() map[string]plugintypes.Hook
}
