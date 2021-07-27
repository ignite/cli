---
order: 5
description: Say "Hello Starport"
---

# Say "Hello, Starport"

To prepare your Cosmos SDK blockchain to say "Hello", you take these actions:

- Modify a protocol buffer file
- Create a keeper query function that returns data
- Register a query function

Protocol buffer files contain proto `rpc` calls that define Cosmos SDK queries and message handlers. The buffer files also contain proto `message` types that define Cosmos SDK types. The `rpc` calls are also responsible for exposing an HTTP API.

To create a keeper query function that returns data, first understand what a keeper is. In Cosmos SDK modules, a keeper is an abstraction for modifying the state of the blockchain. Keeper functions let you query or write to the state. To learn more about module-specific keepers, see the Cosmos SDK [Keeper](https://docs.cosmos.network/v0.42/building-modules/keeper.html) documentation. 
After you add the first query to your chain, you must register the query function.

When you build your blockchain app, the typical workflow is to work with proto files first. The proto files define the Cosmos SDK modules:

- [Messages](https://docs.cosmos.network/v0.42/building-modules/msg-services.html)
- [Queries](https://docs.cosmos.network/v0.42/building-modules/query-services.html)
- [Message handlers](https://docs.cosmos.network/v0.39/building-modules/handler.html)

After you define the messages, queries, and message handlers, you are well prepared to implement the logic of these queries and message handlers in keeper functions.

Start by creating a `posts` query:

```sh
starport scaffold query posts --response title,body
```

The `query` parameter accepts a name of the query (in this case, `posts`), an optional list of request parameters (in this case, empty), and an optional comma-separated list of response fields with a `--response` flag (in this case, `body,title`).

The `query` command has created and modified several files:

- modified `proto/blog/query.proto`
- created `x/blog/keeper/grpc_query_posts.go`
- modified `x/blog/client/cli/query.go`
- created `x/blog/client/cli/query_posts.go`

Now, examine some of these changes. For clarity, placeholder comments are skipped in the following code blocks.

**Note:** The placeholder comments are generated during the scaffolding process and are required. Even though the comments are omitted in these examples, don't delete them. Keep these placeholders to continue using the Starport scaffolding functionality.

In `proto/blog/query.proto` a `Posts` `rpc` has been added to the `Query` `service`.

```go
service Query {
  rpc Posts(QueryPostsRequest) returns (QueryPostsResponse) {
    option (google.api.http).get = "/alice/blog/blog/posts";
  }
}
```

`Posts` `rpc` will be responsible for returning a list of all the posts on chain. It accepts request parameters (`QueryPostsRequest`) and returns response of type `QueryPostsResponse`. `option` defines the endpoint that will be used by gRPC to generate an HTTP API.

You see both the request and response types in this code:

```go
message QueryPostsRequest {
}

message QueryPostsResponse {
  string title = 1;
  string body = 2;
}
```

- `QueryPostsRequest` is empty because requesting all posts does not require any parameters.
- `QueryPostsResponse` contains `title` and `body` that is returned from the chain.

The `x/blog/keeper/grpc_query_posts.go` file contains the `Posts` keeper function that handles the query and returns data.

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

The `Posts` function makes a basic check on the request and throws an error if it's `nil`, stores context in a `ctx` variable, and returns a response of type `QueryPostsResponse`. The context contains information about the environment of the request.

Right now, the response is empty.

From the `query.proto` file, you know that the response can contain `title` and `body`, so modify the last line of the function to return a "Hello!".

```go
func (k Keeper) Posts(c context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
  //...
  return &types.QueryPostsResponse{Title: "Hello!", Body: "Starport"}, nil
}
```

If you start your chain right now and visit the posts endpoint, you see a "Not Implemented" error.

To populate the API endpoint, you must register the query handlers with gRPC.

Inside the `x/blog/module.go` file, import `"context"`, search for `RegisterGRPCGatewayRoutes`, and register query handlers:

```go
import (
  //...
  "context"
)

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
  types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}
```

Now you are ready to start your blockchain:

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
