package pluginsrpc

// HookType represents when a hook can be run, in this case "pre" or "post"
type HookType string

const (
	PreRunHook  HookType = "pre"
	PostRunHook HookType = "post"
)

const (
	PLUGINS_DIR         = "plugins"
	COMMAND_MODULE_NAME = "CommandModule"
	HOOK_MODULE_NAME    = "HookModule"
)
