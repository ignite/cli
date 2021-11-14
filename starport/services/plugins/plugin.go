package plugins

type CmdPlugin interface {
	Module
	Registry() map[string]Command
}

type HookPlugin interface {
	Module
	Registry() map[string]Hook
}
