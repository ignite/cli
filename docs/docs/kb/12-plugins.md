---
sidebar_position: 12
description: Using and Developing Plugins
---

# Using and Developing Plugins


## Using Plugins
Ignite plugins offer a way to extend the functionality of the Ignite CLI. There are two core concepts within plugins : `Commands` and `Hooks`. Where `Commands` extend the cli's functionality, and `Hooks` extend existing command functionality.

Plugins are registered in an Ignite scaffolded Blockchain project through the `config.yml`.

### Adding plugins to a project

Plugins are registered per project. When a command is executed within a directory on disk which contains an configuration file with a plugin registered, Ignite will automatically load it into itself before a command is run. If the plugin has not already been built there will be a build step before the command is executed which the cli will indicate. Below is an example configuration registering a plugin:

``` yaml
plugins:
    # must be the name same name as the plugin project name
  - plugin:
    # path the plugin project directory 
    # path can be a remote git repository or local directory
    path: https://github.com/project/your-plugin
    # parameters which can be passed to commands when executing
    with:
        foo: value
        bar: value
```

### Listing Installed Plugins

When in an ignite scaffolded blockchain you can use the command `ignite plugin list` to list all plugins and there statuses.

### Updating Plugins

When a plugin in a remote repository releases updates, running `ignite plugin update <path/to/plugin>` will update a specific plugin declared in your project's `config.yml`.

### Plugin Hooks

Plugin `Hooks` allow existing ignite commands to be extended with new functionality. Hooks are useful when you want to streamline functionality without needing to run custom scripts after or before a command has been run. this can streamline processes that where once error prone or forgotten all together.

### Hook Lifecyle Events

The following are hooks defined which will run on a registered `ignite` commands

| Name   | Description             |
|--------|-------------------------|
| Pre    | Runs before a commands main functionality is invoked in the `PreRun` scope |
| Post   | Runs after a commands main functionality is invoked in the `PostRun` scope |
| Clean Up| Runs after a commands main functionality is invoked. if the command returns an error it will run before the error is returned to guarantee execution.|

*Note*: If a hook causes an error in the pre step the command will not run resulting in `post` and `clean up` not executing.


## Developing Plugins

### Registering a local plugin

When developing a plugin, you may declare a plugins `path` in a `config.yml` from your local disk. Below is an example of registering a local plugin:

**Note** When using a local file path to register, the path **must** be absolute.

```yaml
plugins:
    # must be the name same name as the plugin project name
  - plugin:
    # path the plugin project directory 
    # path can be a remote git repository or local directory
    path: /path/to/plugin
    # parameters which can be passed to commands when executing
    with:
        foo: value
        bar: value
```

### Defining Commands

Plugin commands are custom commands added to the ignite cli by a `registered` plugin. Commands can be of any path not defined already by ignite. All plugin commands will extend of the command root `ignite`. An example command definition is outlined below:
```go
func (p) Commands() []plugin.Command {
	// TODO: write your command list here
	return []plugin.Command{
        {
            Use:               "foo",
            Short:             "Explain what the command is doing...",
            Long:              "Long description goes here...",
            PlaceCommandUnder: "ignite chain",
            // Examples of subcommands:
            Commands: []plugin.Command{
                {Use: "add"},
                {Use: "list"},
                {Use: "delete"},
            },
        },
    }
}

func (p) Execute(cmd plugin.Command, args []string) error {
	// TODO: write command execution here
	fmt.Printf("Hello I'm the test-plugin plugin!\nargs=%v, with=%v\n", args, cmd.With)

	// This is how the plugin can access the chain:
	c, err := ignitecmd.NewChainWithHomeFlags(cmd.CobraCmd)
	if err != nil {
		return err
	}
	_ = c

	// According to the number of declared commands, you may need a switch:
	switch cmd.Use {
	case "add":
		fmt.Println("Adding stuff...")
	case "list":
		fmt.Println("Listing stuff...")
	case "delete":
		fmt.Println("Deleting stuff...")
	}
	return nil
}
```

### Defining Hooks

The following is an example of a `hook` definition.

```go
func (p) Hooks() []plugin.Hook {
	return []plugin.Hook{
		{
			Name:        "my-hook",
			PlaceHookOn: "ignite chain build",
		},
	}
}

func (p) ExecuteHookPre(hook plugin.Hook, args []string) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("hello")
	default:
		return fmt.Errorf("hook not defined")
	}

	return nil
}

func (p) ExecuteHookPost(hook plugin.Hook, args []string) error {
	switch hook.Name {
  case "my-hook":
    fmt.Println("hey there")
	default:
		return fmt.Errorf("hook not defined")
	}

  return nil
}

func (p) ExecuteHookCleanUp(hook plugin.Hook, args []string) error {
	switch hook.Name {
  case "my-hook":
    fmt.Println("Cleaning Up")
	default:
		return fmt.Errorf("hook not defined")
	}

  return nil
}
```

Above we can see a similar definition to `Command` where a hook has a `Name` and a `PlaceHookOn`. You'll notice that the `Execute*` methods map directly to each life cycle of the hook. All hooks defined within the plugin will invoke these methods.


