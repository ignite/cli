package plugins

// use map keys for command names soon

type CmdPlugin interface {
	Module
	Registry() map[string]Command
}

type HookPlugin interface {
	Module
	Registry() map[string]Hook
}
