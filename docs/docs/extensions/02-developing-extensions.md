---
description: Using and Developing Ignite Extensions
---

# Developing Ignite Extensions

It's easy to create an extension and use it immediately in your project. First
choose a directory outside your project and run:

```sh
ignite extension scaffold my-extension
```

This will create a new directory `my-extension` that contains the extension's code
and will output some instructions about how to use your extension with the
`ignite` command. An extension path can be a local directory which has several
benefits:

- You don't need to use a Git repository during the development of your extension.
- The extension is recompiled each time you run the `ignite` binary in your
  project if the source files are older than the extension binary.

Thus, extension development workflow is as simple as:

1. Scaffold an extension with `ignite extension scaffold my-extension`
2. Add it to your config via `ignite extension install -g /path/to/my-extension`
3. Update extension code
4. Run `ignite my-extension` binary to compile and run the extension
5. Go back to 3

Once your extension is ready you can publish it to a Git repository and the
community can use it by calling `ignite extension install github.com/foo/my-extension`.

Now let's detail how to update your extension's code.

## extension interface

Under the hood Ignite Extensions are implemented using a plugin system based on
`github.com/hashicorp/go-plugin`.

All extensions must implement a predefined interface:

```go title=ignite/services/plugin/interface.go
type Interface interface {
 // Manifest declares extension's Command(s) and Hook(s).
 Manifest(context.Context) (*Manifest, error)

 // Execute will be invoked by ignite when an extension Command is executed.
 // It is global for all commands declared in Manifest, if you have declared
 // multiple commands, use cmd.Path to distinguish them.
 // The ClientAPI argument can be used by plugins to get chain extension analysis info.
 Execute(context.Context, *ExecutedCommand, ClientAPI) error

 // ExecuteHookPre is invoked by ignite when a command specified by the Hook
 // path is invoked.
 // It is global for all hooks declared in Manifest, if you have declared
 // multiple hooks, use hook.Name to distinguish them.
 // The ClientAPI argument can be used by plugins to get chain extension analysis info.
 ExecuteHookPre(context.Context, *ExecutedHook, ClientAPI) error

 // ExecuteHookPost is invoked by ignite when a command specified by the hook
 // path is invoked.
 // It is global for all hooks declared in Manifest, if you have declared
 // multiple hooks, use hook.Name to distinguish them.
 // The ClientAPI argument can be used by plugins to get chain extension analysis info.
 ExecuteHookPost(context.Context, *ExecutedHook, ClientAPI) error

 // ExecuteHookCleanUp is invoked by ignite when a command specified by the
 // hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
 // execution status of the command and hooks.
 // It is global for all hooks declared in Manifest, if you have declared
 // multiple hooks, use hook.Name to distinguish them.
 // The ClientAPI argument can be used by plugins to get chain extension analysis info.
 ExecuteHookCleanUp(context.Context, *ExecutedHook, ClientAPI) error
}
```

The scaffolded code already implements this interface, you just need to update
the method's body.

## Defining extension's manifest

Here is the `Manifest` proto message definition:

```protobuf title=proto/ignite/services/plugin/grpc/v1/types.proto
message Manifest {
  // extension name.
  string name = 1;

  // Commands contains the commands that will be added to the list of ignite commands.
  // Each commands are independent, for nested commands use the inner Commands field.
  bool shared_host = 2;

  // Hooks contains the hooks that will be attached to the existing ignite commands.
  repeated Command commands = 3;

  // Enables sharing a single extension server across all running instances of an Ignite Extension.
  // Useful if an extension adds or extends long running commands.
  //
  // Example: if an extension defines a hook on `ignite chain serve`, a server is instantiated
  // when the command is run. Now if you want to interact with that instance
  // from commands defined in that extension, you need to enable shared host, or else the
  // commands will just instantiate separate extension servers.
  //
  // When enabled, all extensions of the same path loaded from the same configuration will
  // attach it's RPC client to a an existing RPC server.
  //
  // If an extension instance has no other running extension servers, it will create one and it
  // will be the host.
  repeated Hook hooks = 4;
}
```

In your extension's code the `Manifest` method already returns a predefined
`Manifest` struct as an example. You must adapt it according to your need.

If your extension adds one or more new commands to `ignite`, add them to the
`Commands` field.

If your extension adds features to existing commands, add them to the `Hooks` field.

Of course an extension can declare both, `Commands` *and* `Hooks`.

An extension may also share a host process by setting `SharedHost` to `true`.
`SharedHost` is desirable if an extension hooks into, or declares long running commands.
Commands executed from the same extension context interact with the same extension server.
Allowing all executing commands to share the same server instance, giving shared execution context.

## Adding new commands

extension commands are custom commands added to Ignite CLI by an installed extension.
Commands can use any path not defined already by the CLI.

For instance, let's say your extension adds a new `oracle` command to `ignite
scaffold`, then the `Manifest` method will look like :

```go
func (extension) Manifest(context.Context) (*plugin.Manifest, error) {
 return &plugin.Manifest{
  Name: "oracle",
  Commands: []*plugin.Command{
   {
    Use:   "oracle [name]",
    Short: "Scaffold an oracle module",
    Long:  "Long description goes here...",
    // Optional flags is required
    Flags: []*plugin.Flag{
     {Name: "source", Type: plugin.FlagTypeString, Usage: "the oracle source"},
    },
    // Attach the command to `scaffold`
    PlaceCommandUnder: "ignite scaffold",
   },
  },
 }, nil
}
```

To update the extension execution, you have to change the `Execute` command. For
example:

```go
func (extension) Execute(_ context.Context, cmd *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
 if len(cmd.Args) == 0 {
  return fmt.Errorf("oracle name missing")
 }

 flags, err := cmd.NewFlags()
 if err != nil {
  return err
 }

 var (
  name      = cmd.Args[0]
  source, _ = flags.GetString("source")
 )

 // Read chain information
 c, err := getChain(cmd)
 if err != nil {
  return err
 }

 //...
}
```

Then, run `ignite scaffold oracle` to execute the extension.

## Adding hooks

extension `Hooks` allow existing CLI commands to be extended with new
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
func (extension) Manifest(context.Context) (*plugin.Manifest, error) {
 return &plugin.Manifest{
  Name: "oracle",
  Hooks: []*plugin.Hook{
   {
    Name:        "my-hook",
    PlaceHookOn: "ignite chain build",
   },
  },
 }, nil
}

func (extension) ExecuteHookPre(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
 switch h.Hook.GetName() {
 case "my-hook":
  fmt.Println("I'm executed before ignite chain build")
 default:
  return fmt.Errorf("hook not defined")
 }
 return nil
}

func (extension) ExecuteHookPost(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
 switch h.Hook.GetName() {
 case "my-hook":
  fmt.Println("I'm executed after ignite chain build (if no error)")
 default:
  return fmt.Errorf("hook not defined")
 }
 return nil
}

func (extension) ExecuteHookCleanUp(_ context.Context, h *plugin.ExecutedHook, _ plugin.ClientAPI) error {
 switch h.Hook.GetName() {
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
each life cycle of the hook. All hooks defined within the extension will invoke these
methods.
