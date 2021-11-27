# Plugins

Plugins introduce functionality to Starport that promotes modular development. Use custom commands and run hooks, without customizing Starport from source.

Because of the limiting nature of the built in go plugins, this service uses go-plugin (https://github.com/hashicorp/go-plugin) in order to achieve acceptable functionality.

## Usage

Base command:
```
starport plugins
```

Reload plugins:
```
starport plugins reload
```

When starting a chain with `starport chain serve`, plugins will automatically be applied.

## Writing a plugin

To write a plugin, you must conform to the interfaces defined in `github.com/lukerhoads/plugintypes` repo. You must also set up a main function, which serves RPC calls that the plugin service uses to get information on your custom plugin. 

To display output, you must use the log library.

An example is included here: https://github.com/lukerhoads/testplugin

For a command plugin:
```golang
type Command interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	NumArgs() int
	Exec(*cobra.Command, []string) error
}
```

For a hook plugin:
```golang
type Hook interface {
	ParentCommand() []string
	Name() string
	Type() string

	PreRun(*cobra.Command, []string) error
	PostRun(*cobra.Command, []string) error
}
```

### Field Guide
- ParentCommand: Command to nest your plugin under.
- Name: Name your plugin is referred to as. Not required.
- Usage: Printed as a brief show of how the command is formatted.
- ShortDesc: Short description of command.
- LongDesc: Long description of command.
- NumArgs: Number of command arguments.
- HookType: either "pre" or "post", designate which command is to be ran.

To register your plugin, it must be in the main package, and have a main function.

Example main function for hook plugin:
```golang
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugintypes.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"hook": &plugintypes.HookPlugin{Impl: &testHooks{}},
		},
	})
}
```

**NOTE**: the key of the map MUST either be command or hook, based on the type of plugin.

Keep note that the plugin service uses primitives defined by the `github.com/hashicorp/go-plugin` package.

## Todo for Contributors
- Add more validation to plugins, have a clearer specification of what is permitted