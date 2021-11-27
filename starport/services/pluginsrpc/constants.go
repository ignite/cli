package pluginsrpc

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
