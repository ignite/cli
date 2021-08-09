---
order: 2
---

# Hello, World!

Throughout the tutorials unless noted otherwise you'll be using a specific version of Starport. To install Starport v0.17.1 use the following command:

```
curl https://get.starport.network/starport@v0.17.1! | bash
```

Starport comes with a number of scaffolding commands that are designed to make development easier by creating everything that's necessary to start working on a particular task. One of these tasks is a `scaffold scaffold chain` command which provides you with a foundation of a fresh Cosmos SDK blockchain so that you don't have to write it yourself.

To use this command, open a terminal, navigate to a directory where you have permissions to create files, and run:

```
starport scaffold chain github.com/cosmonaut/hello
```

This will create a new Cosmos SDK blockchain called Hello in a `hello` directory. The source code inside the `hello` directory contains a fully functional ready-to-use blockchain. This new blockchain imports standard Cosmos SDK modules, such as [`staking`](https://docs.cosmos.network/v0.42/modules/staking/) (for delegated proof of stake), [`bank`](https://docs.cosmos.network/v0.42/modules/bank/) (for fungible token transfers between accounts), [`gov`](https://docs.cosmos.network/v0.42/modules/gov/) (for on-chain governance) and [other modules](https://docs.cosmos.network/v0.42/modules/).

Note: You can see all the command line options that Starport provides by running `starport scaffold chain --help`.

After you create the blockchain, switch to its directory:

```
cd hello
```

The `hello` directory will have a number of generated files and directories that make up the structure of a Cosmos SDK blockchain. Most of the work in this tutorial will happen in the `x` directory, but here is a quick overview of files and directories that Starport creates by default:

| File/directory | Purpose                                                                                                                                                                        |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| app/           | Contains files that wire together the blockchain. The most important file is app.go that contains type definition of the blockchain and functions to create and initialize it. |
| cmd/           | Contains the main package responsible for the CLI of compiled binary.                                                                                                          |
| docs/          | Directory for project's documentation. By default an OpenAPI spec is generated.                                                                                                |
| proto/         | Protocol buffer files describing                                                                                                                                               |
| testutil/      | Contains helper functions for testing                                                                                                                                          |
| vue/           | Contains a Vue 3 web app template                                                                                                                                              |
| x/             | Contains custom modules                                                                                                                                                        |
| config.yml     | A configuration file for customising a chain in development                                                                                                                    |

To get started, let's get our blockchain up and running locally on a single node.

## Starting a Blockchain

You actually have a fully-functional blockchain already. To start it on your development machine, run the following command in the `hello` directory

```
starport chain serve
```

This will download dependencies, compile the source code into a binary called `hellod` (by default Starport uses the name of the repo + `d`), use this binary to initialize a single validator node, and start the node.

![./images/api.png](./images/api.png)

A validator node exposes two endpoints: [http://localhost:26657](http://localhost:26657) for the low-level Tendermint API and [http://localhost:1317](http://localhost:1317) for the high-level blockchain API.

When you want to stop your blockchain, press Ctrl+C in the terminal window where it's running. In the development environment, Starport doesn't require you to restart the blockchain; changes you make in files will be automatically picked up by Starport.

## Say "Hello, Starport"

To get your Cosmos SDK blockchain to say "Hello", you will need to modify a protocol buffer file, create a keeper query function that returns data, and register a query function. Protocol buffer files contain proto `rpc`s that define Cosmos SDK queries and message handlers and proto `message`s that define Cosmos SDK types. `rpc`s are also responsible for exposing an HTTP API. [Keeper](https://docs.cosmos.network/v0.42/building-modules/keeper.html) is an abstraction for modifying the state of the blockchain and keeper functions let you query or write to the state. Registering a query needs to happen only once after you add the first query to your chain.

In terms of workflow, developers usually work with proto files first to define Cosmos SDK [messages](https://docs.cosmos.network/v0.42/building-modules/msg-services.html), [queries](https://docs.cosmos.network/v0.42/building-modules/query-services.html), message handlers and then implement the logic of these queries and message handlers in keeper functions.

Create a `posts` query:

```
starport scaffold query posts --response title,body
```

`query` accepts a name of the query (in our case, `posts`), an optional list of request parameters (in our case, empty), and an optional comma-separated list of response fields with a `--response` flag (in our case, `body,title`).

The `query` command has created and modified several files:

- modified `proto/hello/query.proto`
- created `x/hello/keeper/grpc_query_posts.go`
- modified `x/hello/client/cli/query.go`
- created `x/hello/client/cli/query_posts.go`

Let's examine some of these changes. For clarity, in the following code blocks we'll skip placeholder comments Starport uses to scaffold code. Don't delete these placeholders, however, to be able to continue using Starport's scaffolding functionality.

In `proto/hello/query.proto` a `Posts` `rpc` has been added to the `Query` `service`.

```
service Query {
  rpc Posts(QueryPostsRequest) returns (QueryPostsResponse) {
    option (google.api.http).get = "/cosmonaut/hello/hello/posts";
  }
}
```

`Posts` `rpc` will be responsible for returning a list of all the posts on chain. It accepts request parameters (`QueryPostsRequest`) and returns response of type `QueryPostsResponse`. `option` defines the endpoint that will be used by gRPC to generate an HTTP API.

Below you can see both request and response types.

```
message QueryPostsRequest {
}

message QueryPostsResponse {
  string title = 1;
  string body = 2;
}
```

`QueryPostsRequest` is empty because requesting all posts doesn't require are parameters. `QueryPostsResponse` contains `title` and `body` that will be returned from the chain.

`x/hello/keeper/grpc_query_posts.go` contains `Posts` keeper function that handles the query and returns data.

```go
func (k Keeper) Posts(c context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
  if req == nil {
    return nil, status.Error(codes.InvalidArgument, "invalid request")
  }
  ctx := sdk.UnwrapSDKContext(c)
  _ = ctx
  return &types.QueryPostsResponse{}, nil
}
```

`Posts` function makes a basic check on the request and throws an error if it's `nil`, stores context in a `ctx` variable (context contains information about the environment of the request) and returns a response of type `QueryPostsResponse`. Right now the response is empty.

From the `query.proto` we know that response may contain `title` and `body`, so let's modify the last line of the function to return a "Hello!".

```go
func (k Keeper) Posts(c context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
  //...
  return &types.QueryPostsResponse{Title: "Hello!", Body: "Starport"}, nil
}
```

If we start our chain right now and visit the posts endpoint, we would get a "Not Implemented" error. To fix that we need to wire up our API by registering query handlers with gRPC.

Inside `x/hello/module.go` import `"context"`, search for `RegisterGRPCGatewayRoutes` and register query handlers:

```go
import (
  //...
  "context"
)

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
  types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}
```

Now we're ready to start our blockchain:

```go
starport chain serve
```

Once the chain has been started, visit [http://localhost:1317/cosmonaut/hello/hello/posts](http://localhost:1317/cosmonaut/hello/hello/posts) and see our text displayed!

```go
{
  "title": "Hello!",
  "body": "Starport"
}
```

The `query` command has also scaffolded `x/hello/client/cli/query_posts.go` that implements a CLI equivalent of the posts query and mounted this command `x/hello/client/cli/query_posts.go` . Run the following command and get the same JSON response:

```go
hellod q hello posts
```