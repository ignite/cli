---
description: Debugging your Cosmos SDK blockchain
---

# Debugging a chain

Ignite chain debug command can help you find issues during development. It uses
[Delve](https://github.com/go-delve/delve) debugger which enables you to
interact with your blockchain app by controlling the execution of the process,
evaluating variables, and providing information of thread / goroutine state, CPU
register state and more.

## Debug Command

The debug command requires that the blockchain app binary is build with
debugging support by removing optimizations and inlining. A debug binary is
built by default by the `ignite chain serve` command or can optionally be
created using the `--debug` flag when running `ignite chain init` or `ignite
chain build` sub-commands.

To start a debugging session in the terminal run:

```
ignite chain debug
```

The command runs your blockchan app in the background, attaches to it and
launches a terminal debugger shell:

```
Type 'help' for list of commands.
(dlv)
```

At this point the blockchain app blocks execution, so you can set one or more
breakpoints before continuing execution.

Use the
[break](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md#break)
(alias `b`) command to set any number of breakpoints using, for example the
`<filename>:<line>` notation:

```
(dlv) break x/hello/client/cli/query_say_hello.go:14
```

This command adds a breakpoint to the `x/hello/client/cli/query_say_hello.go`
file at line 14.

Once all breakpoints are set resume blockchain execution using the
[continue](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md#continue)
(alias `c`) command:

```
(dlv) continue
```

The debugger will launch the shell and stop blockchain execution again when a
breakpoint is triggered.

Within the debugger shell use the `quit` (alias `q`) or `exit` commands to stop
the blockchain app and exit the debugger.

## Debug Server

A debug server can optionally be started in cases where the default terminal
client is not desirable. When the server starts it first runs the blockchain
app, attaches to it and finally waits for a client connection. The default
server address is *tcp://127.0.0.1:30500* and it accepts both JSON-RPC or DAP
client connections.

To start a debug server use the following flag:

```
ignite chain debug --server
```

To start a debug server with a custom address use the following flags:

```
ignite chain debug --server --server-address 127.0.0.1:30500
```

The debug server stops automatically when the client connection is closed.

## Debugging Clients

### Gdlv: Multiplatform Delve UI

[Gdlv](https://github.com/aarzilli/gdlv) is a graphical frontend to Delve for
Linux, Windows and macOS.

Using it as debugging client is straightforward as it doesn't require any
configuration. Once the debug server is running and listening for client
requests connect to it by running:

```
gdlv connect 127.0.0.1:30500
```

Setting breakpoints and continuing execution is done in the same way as Delve,
by using the `break` and `continue` commands.

### Visual Studio Code

Using [Visual Studio Code](https://code.visualstudio.com/) as debugging client
requires an initial configuration to allow it to connect to the debug server.

Make sure that the [Go](https://code.visualstudio.com/docs/languages/go)
extension is installed.

VS Code debugging is configured using the `launch.json` file which is usually
located inside the `.vscode` folder in your workspace.

You can use the following launch configuration to set up VS Code as debugging
client:

```json title=launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to Debug Server",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 30500,
            "host": "127.0.0.1"
        }
    ]
}
```

Alternatively it's possible to create a custom `launch.json` file from the "Run
and Debug" panel. When prompted choose the Go debugger option labeled "Go:
Connect to Server" and enter the debug host address and then the port number.

## Example: Debugging a Blockchain App

In this short example we will be using Ignite CLI to create a new blockchain and
a query to be able to trigger a debugging breakpoint when the query is called.

Create a new blockchain:

```
ignite scaffold chain hello
```

Scaffold a new query in the `hello` directory:

```
ignite scaffold query say-hello name --response name
```

The next step initializes the blockchain's data directory and compiles a debug
binary:

```
ignite chain init --debug
```

Once the initialization finishes launch the debugger shell:

```
ignite chain debug
```

Within the debugger shell create a breakpoint that will be triggered when the
`SayHello` function is called and then continue execution:

```
(dlv) break x/hello/keeper/query_say_hello.go:12
(dlv) continue
```

From a different terminal use the `hellod` binary to call the query:

```
hellod query hello say-hello bob
```

A debugger shell will be launched when the breakpoint is triggered:

```
     7:		"google.golang.org/grpc/codes"
     8:		"google.golang.org/grpc/status"
     9:		"hello/x/hello/types"
    10:	)
    11:
=>  12:	func (k Keeper) SayHello(goCtx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
    13:		if req == nil {
    14:			return nil, status.Error(codes.InvalidArgument, "invalid request")
    15:		}
    16:
    17:		ctx := sdk.UnwrapSDKContext(goCtx)
```

From then on you can use Delve commands like `next` (alias `n`) or `print`
(alias `p`) to control execution and print values. For example, to print the
*name* argument value use the `print` command followed by "req.Name":

```
(dlv) print req.Name
"bob"
```

Finally, use `quit` (alias `q`) to stop the blockchain app and finish the
debugging session.
