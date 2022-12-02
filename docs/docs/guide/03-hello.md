---
sidebar_position: 2
description: Step-by-step guidance to build your first blockchain and your first Cosmos SDK module.
---

# Hello, World!

In this tutorial you will create a blockchain with only one custom feature: when
queried it will respond with `"Hello, %s"`, where `%s` is a name provided in a
query.

By the end of the tutorial, you will have a better understanding of how to
create a blockchain with very simple custom functionality.

This tutorial is divided into two sections. In the first section, we will use
Ignite's query scaffold command to generate most of the code and modify the
query to return `"Hello, %s"`. In the second section, we will implement the same
functionality manually to gain a better understanding of how queries work.

First, create a new `hello` blockchain with Ignite CLI:

```
ignite scaffold chain hello
```

Let's add a query to the blockchain we just created.

In the Cosmos SDK, a query is a request for information from the blockchain.
Queries are used to retrieve data from the blockchain, such as the current state
of the ledger or the details of a specific transaction. The Cosmos SDK provides
a number of built-in query methods that can be used to retrieve data from the
blockchain, and developers can also create custom queries to access specific
data or perform complex operations. Queries are processed by the blockchain's
nodes and the results are returned to the querying client.

## Create a query with Ignite

To add a query, run the following command:

```
ignite scaffold query say-hello name --response name
```

The `ignite scaffold query` command is a tool used to quickly create new
queries. When you run this command, it makes changes to your source code to add
the new query and make it available in your API. This command accepts a query
name (`"say-hello"`) and a list of request fields (in our case only `name`). The
optional `--reponse` flag specifies the return values of the query.

This command made the following changes to the source code.

The `proto/hello/hello/query.proto` file was modified to define the request and
response for a query, as well as to add the `SayHello` query in the `Query`
service. 

The `x/hello/keeper/grpc_query_say_hello.go` file was created with a keeper
method called `SayHello`. This method is responsible for handling the "say
hello" query, which can be called by a client using the command-line interface
(CLI) or an API. When the "say hello" query is executed, the `SayHello` method
is called to perform the necessary actions and return a response to the client.
The `SayHello` method may retrieve data from the application's database, process
the data, and return a result to the client in a specific format, such as a
string of text or a data structure.

The `x/hello/client/cli/query.go file` was created and added to the project.
This file contains a CLI command `CmdSayHello` that allows users to submit a
"say hello" query to the blockchain. This command allows users to interact with
the blockchain in a more user-friendly way, allowing them to easily submit
queries and receive responses from the blockchain.

The `x/hello/client/cli/query.go` was modified to add the `CmdSayHello` command
to the CLI of the blockchain.

<!-- 

This tutorial is a great place to start your journey into the Cosmos ecosystem.
Instead of wondering how to build a blockchain, follow these steps to build your
first blockchain and your first Cosmos SDK module.

## Get started

In the previous chapter you've learned how to install [Ignite
CLI](https://github.com/ignite/cli), the tool that offers everything you need to
build, test, and launch your blockchain with a decentralized worldwide
community.

This tutorial is based on a specific version of Ignite CLI, so be sure to
install the correct version. 

Use the following command:

```bash
curl https://get.ignite.com/cli@v0.25.1! | bash
```

Ignite CLI comes with a number of scaffolding commands that are designed to make
development easier by creating everything that's required to start working on a
particular task.

First, use Ignite CLI to build the foundation of a fresh Cosmos SDK blockchain.
With Ignite CLI, you don't have to write the blockchain code yourself.

Are you ready? Open a terminal window and navigate to a directory where you have
permissions to create files.

To create your blockchain with the default directory structure, run this
command:

```bash
ignite scaffold chain hello
```

This command creates a Cosmos SDK blockchain called hello in a `hello`
directory. The source code inside the `hello` directory contains a fully
functional ready-to-use blockchain.

This new blockchain imports standard Cosmos SDK modules, including:


- [`staking`](https://docs.cosmos.network/main/modules/staking/) for delegated
  Proof-of-Stake (PoS) consensus mechanism
- [`bank`](https://docs.cosmos.network/main/modules/bank/) for fungible token
  transfers between accounts
- [`gov`](https://docs.cosmos.network/main/modules/gov/) for on-chain governance
- And other Cosmos SDK [modules](https://docs.cosmos.network/main/modules/) that
provide the benefits of the extensive Cosmos SDK framework

You can get help on any command. Now that you have run your first command, take
a minute to see all of the command line options for the `scaffold` command.  

To learn about the command you just used, run:

```bash
ignite scaffold --help
```

## Blockchain directory structure

After you create the blockchain, switch to its directory:

```bash
cd hello
```

The `hello` directory contains a number of generated files and directories that
make up the structure of a Cosmos SDK blockchain. Most of the work in this
tutorial happens in the `x` directory. Here is a quick overview of files and
directories that are created by default:

* `app/`: files that wire together the blockchain. The most important file is
  `app.go` that contains type definition of the blockchain and functions to
  create and initialize it.
* `cmd/`: the main package responsible for the CLI of compiled binary.
* `docs/`: directory for project documentation. By default, an OpenAPI spec is
  generated.
* `proto/`: protocol buffer files describing the data structure.
* `testutil/`: helper functions for testing.
* `x/`: Cosmos SDK modules and custom modules.
* `config.yml`: a configuration file for customizing a chain in development.
* `readme.md`: a readme file for your sovereign application-specific blockchain
  project.

Now you can get your blockchain up and running locally on a single node.

## Start a blockchain

You already have a fully-functional blockchain. To start your chain on your
development machine, run the following command in the `hello` directory

```bash
ignite chain serve
```

This command downloads dependencies and compiles the source code into a binary
called `hellod`. By default, the binary name is the name of the repo + `d`. From
now on, use this `hellod` binary to run all of your chain commands. For example,
to initialize a single validator node and start a node.

Leave this terminal window open while your chain is running.

## HTTP API Console

By default, a validator node exposes two API endpoints:

- [http://localhost:26657](http://localhost:26657) for the low-level Tendermint
  API
- [http://localhost:1317](http://localhost:1317) for the high-level blockchain
  API

Now that you started your `hello` chain, use a web browser to see the high-level
`hello` blockchain API:

![./images/api.png](./images/api.png)

## Stop a blockchain

When you want to stop your blockchain, press Ctrl+C in the terminal window where
it's running.

In the development environment, you can experiment and instantly see updates.
You don't have to restart the blockchain after you make changes. Hot reloading
automatically detects all of the changes you make in the `hello` directory
files.

## Say "Hello, Ignite CLI"

To get your blockchain to say `Hello! Ignite CLI`, you need to make these
changes:

- Modify a protocol buffer file
- Create a keeper query function that returns data
- Register a query function

Protocol buffer files contain proto RPC calls that define Cosmos SDK queries and
message handlers, and proto messages that define Cosmos SDK types. The RPC calls
are also responsible for exposing an HTTP API.

For each Cosmos SDK module, the
[Keeper](https://docs.cosmos.network/main/building-modules/keeper.html) is an
abstraction for modifying the state of the blockchain. Keeper functions let you
query or write to the state. After you add the first query to your chain, the
next step is to register the query. You only need to register a query once.

A typical blockchain developer workflow looks something like this:

- Start with proto files to define Cosmos SDK
  [messages](https://docs.cosmos.network/main/building-modules/msg-services.html)
- Define and register
  [queries](https://docs.cosmos.network/main/building-modules/query-services.html)
- Define message handler logic
- Finally, implement the logic of these queries and message handlers in keeper
  functions

## Create a query

For all subsequent commands, use a terminal window that is different from the
window you started the chain in.

In a different terminal window, run the commands in your `hello` directory.

Create a `hello` query:

```bash
ignite scaffold query hello --response text
```

`query` accepts a name of the query (in this case, `hello`), an optional list of
request parameters (in this case, empty), and an optional comma-separated list
of response fields with a `--response` flag (in this case, `text`).

The `query` command has created and modified several files:

```
modify proto/hello/hello/query.proto
modify x/hello/client/cli/query.go
create x/hello/client/cli/query_hello.go
create x/hello/keeper/grpc_query_hello.go
```

Let's examine some of these changes. For clarity, the following code blocks do
not show the placeholder comments that Ignite CLI uses to scaffold code. Don't
delete these placeholders since they are required to continue using Ignite CLI's
scaffolding functionality.

Note: it's recommended to commit changes to a version control system (for
example, Git) after scaffolding. This allows others to easily distinguish
between code generated by Ignite and the code writen by hand.

```
git add .
git commit -am "Scaffolded a hello query with Ignite CLI"
```

### Updates to the query service

In the `proto/hello/hello/query.proto` file, the `Hello` RPC has been added to
the `Query` service.

```protobuf title="proto/hello/hello/query.proto"
service Query {
 rpc Hello(QueryHelloRequest) returns (QueryHelloResponse) {
  option (google.api.http).get = "/hello/hello/hello";
 }
}
```

Here's how the `Hello` RPC for the `Query` service works:

- Is responsible for returning a `text` string
- Accepts request parameters (`QueryHelloRequest`)
- Returns response of type `QueryHelloResponse`
- The `option` defines the endpoint that is used by gRPC to generate an HTTP API

### Request and reponse types

Now, take a look at the following request and response types:

```protobuf
message QueryHelloRequest {
}

message QueryHelloResponse {
  string text = 1;
}
```

- The `QueryHelloRequest` message is empty because this request does not require
  parameters.
- The `QueryHelloResponse` message contains `text` that is returned from the
  chain.

## Hello keeper function

The `x/hello/keeper/grpc_query_hello.go` file contains the `Hello` keeper
function that handles the query and returns data.

```go title="x/hello/keeper/grpc_query_hello.go"
func (k Keeper) Hello(goCtx context.Context, req *types.QueryHelloRequest) (*types.QueryHelloResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.QueryHelloResponse{}, nil
}
```

The `Hello` function performs these actions:

- Makes a basic check on the request and throws an error if it's `nil`
- Stores context in a `ctx` variable that contains information about the
  environment of the request
- Returns a response of type `QueryHelloResponse`

Right now the response is empty.

### Update keeper function

In the `query.proto` file, the response accepts `text`.

- Use a text editor to modify the `x/hello/keeper/grpc_query_hello.go` file that
  contains the keeper function.
- On the last line of the keeper function, change the line to return "Hello,
  Ignite CLI!":

```go title="x/hello/keeper/grpc_query_hello.go"
func (k Keeper) Hello(goCtx context.Context, req *types.QueryHelloRequest) (*types.QueryHelloResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// remove-next-line
	return &types.QueryHelloResponse{}, nil
	// highlight-next-line
	return &types.QueryHelloResponse{Text: "Hello, Ignite CLI!"}, nil
}
```

- Save the file to restart your chain.
- After the chain has been started, visit
[http://localhost:1317/hello/hello/hello](http://localhost:1317/hello/hello/hello)
and see your text displayed:

   ```json
   {
     "text": "Hello, Ignite CLI!"
   }
   ```

The `query` command has also scaffolded `x/hello/client/cli/query_hello.go` that
implements a CLI equivalent of the hello query and mounted this command in
`x/hello/client/cli/query.go`. Run the following command and get the same JSON
response:

```bash
hellod q hello hello
```

Congratulations, you have built your first blockchain and your first Cosmos SDK
module. Continue the journey to learn more about scaffolding Cosmos SDK
messages, types in protocol buffer files, the keeper, and more. -->
