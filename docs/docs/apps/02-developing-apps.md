---
description: Using and Developing Ignite Apps
---

# Developing Ignite Apps

It's easy to create an app and use it immediately in your project. First
choose a directory outside your project and run:

```sh
$ ignite app scaffold my-app
```

This will create a new directory `my-app` that contains the app's code
and will output some instructions about how to use your app with the
`ignite` command. An app path can be a local directory which has several
benefits:

- You don't need to use a Git repository during the development of your app.
- The app is recompiled each time you run the `ignite` binary in your
  project if the source files are older than the app binary.

Thus, app development workflow is as simple as:

1. Scaffold an app with `ignite app scaffold my-app`
2. Add it to your config via `ignite app install -g /path/to/my-app`
3. Update app code
4. Run `ignite my-app` binary to compile and run the app
5. Go back to 3

Once your app is ready you can publish it to a Git repository and the
community can use it by calling `ignite app install github.com/foo/my-app`.

Now let's detail how to update your app's code.

## App interface

Under the hood Ignite Apps are implemented using a plugin system based on
`github.com/hashicorp/go-plugin`.

All apps must implement a predefined interface:

```go title=ignite/services/plugin/interface.go
type Interface interface {
	// Manifest declares app's Command(s) and Hook(s).
	Manifest() (Manifest, error)

	// Execute will be invoked by ignite when an app Command is executed.
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

The scaffolded code already implements this interface, you just need to update
the method's body.

## Defining app's manifest

Here is the `Manifest` struct:

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

	// SharedHost enables sharing a single app server across all running instances
	// of an app. Useful if an app adds or extends long running commands
	//
	// Example: if an app defines a hook on `ignite chain serve`, a server is instanciated
	// when the command is run. Now if you want to interact with that instance from commands
	// defined in that app, you need to enable `SharedHost`, or else the commands will just
	// instantiate separate app servers.
	//
	// When enabled, all apps of the same `Path` loaded from the same configuration will
	// attach it's gRPC client to a an existing gRPC server.
	//
	// If an app instance has no other running app servers, it will create one and it
	// will be the host.
	SharedHost bool `yaml:"shared_host"`
}
```

In your app's code the `Manifest` method already returns a predefined
`Manifest` struct as an example. You must adapt it according to your need.

If your app adds one or more new commands to `ignite`, add them to the
`Commands` field.

If your app adds features to existing commands, add them to the `Hooks` field.

Of course an app can declare both, `Commands` *and* `Hooks`.

An app may also share a host process by setting `SharedHost` to `true`.
`SharedHost` is desirable if an app hooks into, or declares long running commands.
Commands executed from the same app context interact with the same app server. 
Allowing all executing commands to share the same server instance, giving shared execution context.

## Adding new commands

App commands are custom commands added to Ignite CLI by an installed app.
Commands can use any path not defined already by the CLI.

For instance, let's say your app adds a new `oracle` command to `ignite
scaffold`, then the `Manifest` method will look like :

```go
func (app) Manifest() (plugin.Manifest, error) {
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

To update the app execution, you have to change the `Execute` command. For
example:

```go
func (app) Execute(cmd plugin.ExecutedCommand) error {
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

Then, run `ignite scaffold oracle` to execute the app.

## Adding hooks

App `Hooks` allow existing CLI commands to be extended with new
functionality. Hooks are useful when you want to streamline functionality
without needing to run custom scripts after or before a command has been run.
This can streamline processes that where once error prone or forgotten all
together.

The following are hooks defined which will run on a registered `ignite`
command:

| Name     | Description                                                                                                                                           |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| Pre      | Runs before a commands main functionality is invoked in the `PreRun` scope                                                                            |
| Post     | Runs after a commands main functionality is invoked in the `PostRun` scope                                                                            |
| Clean Up | Runs after a commands main functionality is invoked. If the command returns an error it will run before the error is returned to guarantee execution. |

*Note*: If a hook causes an error in the pre step the command will not run
resulting in `post` and `clean up` not executing.

The following is an example of a `hook` definition.

```go
func (app) Manifest() (plugin.Manifest, error) {
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

func (app) ExecuteHookPre(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed before ignite chain build")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}

func (app) ExecuteHookPost(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed after ignite chain build (if no error)")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}

func (app) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	switch hook.Name {
	case "my-hook":
		fmt.Println("I'm executed after ignite chain build (regardless errors)")
	default:
		return fmt.Errorf("hook not defined")
	}
	return nil
}
```

Above we can see a similar definition to `Command` where a hook has a `Name`
and a `PlaceHookOn`. You'll notice that the `Execute*` methods map directly to
each life cycle of the hook. All hooks defined within the app will invoke these
methods.
