---
description: Using and Developing plugins
---

# Developing plugins

## Using Plugins

Ignite plugins offer a way to extend the functionality of the Ignite CLI. There
are two core concepts within plugins : `Commands` and `Hooks`. Where `Commands`
extend the cli's functionality, and `Hooks` extend existing command
functionality.

Plugins are registered in an Ignite scaffolded Blockchain project through the
`config.yml`.

### Adding plugins to a project

Plugins are registered per project, using the `config.yml` file. To use a plugin
within your project, add a `plugins` section with the following:

```yaml title=config.yml
plugins:
- path: github.com/project/cli-plugin
```

Now the next time the `ignite` command is run under your project, the declared
plugin will be fetched, compiled and ran. This will result in more avaiable
commands, and/or hooks attached to existing commands.

### Listing installed plugins

When in an ignite scaffolded blockchain you can use the command `ignite plugin
list` to list all plugins and there statuses.

### Updating plugins

When a plugin in a remote repository releases updates, running `ignite plugin
update <path/to/plugin>` will update a specific plugin declared in your
project's `config.yml`.

## Developing Plugins

It's easy to create a plugin and use it immediately in your project. First
choose a directory outside your project and run :

```sh
$ ignite plugin scaffold my-plugin
```

This will create a new directory `my-plugin` that contains the plugin's code,
and will output some instructions about how to declare your plugin in your
project. Indeed it's possible to declare a local directory in your project's
`config.yml`, which has several benefits:

- you don't need to use a git repository during the development of your plugin.
- the plugin is recompiled each time you run the `ignite` binary in your
project, if the source files are older than the plugin binary.

Thus, the plugin development workflow is as simple as :

1. scaffold a plugin
2. declare it in the `config.yml` of a chain (which can be a fresh new chain
created via `ignite scaffold chain my-chain`)
3. update plugin code
4. run `ignite my-command` binary in your chain to compile and run the plugin,
where `my-command` is a command added by your plugin, or an existing `ignite`
command with hooks added by your plugin.
5. go back to 3.

Once your plugin is ready, you can publish it to a git repository, and the
community can use it by declaring this git repository path in their chain's
`config.yml`.

Now let's detail how to update your plugin's code.

### The plugin interface

The `ignite` plugin system uses `github.com/hashicorp/go-plugin` under the hood,
which implies to implement a predefined interface:

```go title=ignite/services/plugin/interface.go
// An ignite plugin must implements the Plugin interface.
type Interface interface {
	// Manifest declares the plugin's Command(s) and Hook(s).
	Manifest() (Manifest, error)

	// Execute will be invoked by ignite when a plugin Command is executed.
	// It is global for all commands declared in Manifest, if you have declared
	// multiple commands, use cmd.Path to distinguish them.
	Execute(cmd ExecutedCommand) error

	// ExecuteHookPre is invoked by ignite when a command specified by the Hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPre(hook ExecutedHook) error

	// ExecuteHookPost is invoked by ignite when a command specified by the hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPost(hook ExecutedHook) error

	// ExecuteHookCleanUp is invoked by ignite when a command specified by the
	// hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
	// execution status of the command and hooks.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookCleanUp(hook ExecutedHook) error
}
```

The code scaffolded already implements this interface, you just need to update
the methods' body.


### Defining plugin's manifest

Here is the `Manifest` struct :

```go title=ignite/services/plugin/interface.go
type Manifest struct {
	Name string
	// Commands contains the commands that will be added to the list of ignite
	// commands. Each commands are independent, for nested commands use the
	// inner Commands field.
	Commands []Command
	// Hooks contains the hooks that will be attached to the existing ignite
	// commands.
	Hooks []Hook
}
```

In your plugin's code, the `Manifest` method already returns a predefined
`Manifest` struct as an example. Adapt it according to your need.

If your plugin adds one or more new commands to `ignite`, feeds the `Commands`
field.

If your plugin adds features to existing commands, feeds the `Hooks` field.

Of course a plugin can declare `Commands` *and* `Hooks`.

### Adding new command

Plugin commands are custom commands added to the ignite cli by a registered
plugin. Commands can be of any path not defined already by ignite. All plugin
commands will extend of the command root `ignite`. 

For instance, let's say your plugin adds a new `oracle` command to `ignite
scaffold`, the `Manifest()` method will look like :

```go
func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "oracle",
		Commands: []plugin.Command{
			{
				Use:   "oracle [name]",
				Short: "Scaffold an oracle module",
				Long:  "Long description goes here...",
				// Optionnal flags is required
				Flags: []plugin.Flag{
					{Name: "source", Type: plugin.FlagTypeString, Usage: "the oracle source"},
				},
				// Attach the command to `scaffold`
				PlaceCommandUnder: "ignite scaffold",
			},
		},
	}, nil
}
```

To update the plugin execution, you have to change the plugin `Execute` command,
for instance :

```go
func (p) Execute(cmd plugin.ExecutedCommand) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("oracle name missing")
	}
	var (
		name      = cmd.Args[0]
		source, _ = cmd.Flags().GetString("source")
	)
	// Read chain information
	c, err := getChain(cmd)
	if err != nil {
		return err
	}

	//...
}
```

Then, run `ignite scaffold oracle` in a chain where the plugin is registered to
compile and execute the plugin.

### Adding hooks

Plugin `Hooks` allow existing ignite commands to be extended with new
functionality. Hooks are useful when you want to streamline functionality
without needing to run custom scripts after or before a command has been run.
this can streamline processes that where once error prone or forgotten all
together.

The following are hooks defined which will run on a registered `ignite` commands

| Name     | Description                                                                                                                                           |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| Pre      | Runs before a commands main functionality is invoked in the `PreRun` scope                                                                            |
| Post     | Runs after a commands main functionality is invoked in the `PostRun` scope                                                                            |
| Clean Up | Runs after a commands main functionality is invoked. if the command returns an error it will run before the error is returned to guarantee execution. |

*Note*: If a hook causes an error in the pre step the command will not run
resulting in `post` and `clean up` not executing.

The following is an example of a `hook` definition.

```go
func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "oracle",
		Hooks: []plugin.Hook{
			{
				Name:        "my-hook",
				PlaceHookOn: "ignite chain build",
			},
		},
	}, nil
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed before ignite chain build")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed after ignite chain build (if no error)")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed after ignite chain build (regardless errors)")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}
```

Above we can see a similar definition to `Command` where a hook has a `Name` and
a `PlaceHookOn`. You'll notice that the `Execute*` methods map directly to each
life cycle of the hook. All hooks defined within the plugin will invoke these
methods.
