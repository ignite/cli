---
order: 5
description: Say "Hello World"
---

### Say "Hello, Starport"

To get your Cosmos SDK blockchain to say "Hello", you will need to modify a protocol buffer file, create a keeper query function that returns data, and register a query function. Protocol buffer files contain proto `rpc`s that define Cosmos SDK queries and message handlers and proto `message`s that define Cosmos SDK types. `rpc`s are also responsible for exposing an HTTP API. [Keeper](https://docs.cosmos.network/v0.42/building-modules/keeper.html) is an abstraction for modifying the state of the blockchain and keeper functions let you query or write to the state. Registering a query needs to happen only once after you add the first query to your chain.

In terms of workflow, developers usually work with proto files first to define Cosmos SDK [messages](https://docs.cosmos.network/v0.42/building-modules/msg-services.html), [queries](https://docs.cosmos.network/v0.42/building-modules/query-services.html), message handlers and then implement the logic of these queries and message handlers in keeper functions.

Let's start by creating a `posts` query:

```
starport scaffold query posts --response title,body
```

`query` accepts a name of the query (in our case, `posts`), an optional list of request parameters (in our case, empty), and an optional comma-separated list of response fields with a `--response` flag (in our case, `body,title`).

The `query` command has created and modified several files:

- modified `proto/blog/query.proto`
- created `x/blog/keeper/grpc_query_posts.go`
- modified `x/blog/client/cli/query.go`
- created `x/blog/client/cli/query_posts.go`

Let's examine some of these changes. For clarity, in the following code blocks we'll skip placeholder comments Starport uses to scaffold code. Don't delete these placeholders, however, to be able to continue using Starport's scaffolding functionality.

In `proto/blog/query.proto` a `Posts` `rpc` has been added to the `Query` `service`.

```
service Query {
  rpc Posts(QueryPostsRequest) returns (QueryPostsResponse) {
    option (google.api.http).get = "/alice/blog/blog/posts";
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

`x/blog/keeper/grpc_query_posts.go` contains `Posts` keeper function that handles the query and returns data.

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

Inside `x/blog/module.go` import `"context"`, search for `RegisterGRPCGatewayRoutes` and register query handlers:

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

Once the chain has been started, visit [http://localhost:1317/alice/blog/blog/posts](http://localhost:1317/alice/blog/blog/posts) and see our text displayed!

```go
{
  "title": "Hello!",
  "body": "Starport"
}
```

The `query` command has also scaffolded `x/blog/client/cli/query_posts.go` that implements a CLI equivalent of the posts query and mounted this command `x/blog/client/cli/query_posts.go` . Run the following command and get the same JSON response:

```go
blogd q blog posts
```
