---
order: 1
parent:
  order: 1
  title: Getting Started with Starport
---

# Getting Started

## What is Starport?

Starport is an easy to use CLI tool for creating sovereign blockchains with Cosmos SDK. Cosmos SDK is the world's most popular framework for building blockchains. Both Starport CLI and Cosmos SDK are written in the Go programming language.

## Starting a New Starport Project

### Installing Go

Before you start using Starport, you should check that your system has Go installed. To do so run:

```
go version
```

If you see `go1.16` (or higher), then you have the right version of Go installed. If the output is `command not found` or installed version of Go is older than 1.16, [install or upgrade Go](https://golang.org/doc/install).

### Installing Starport CLI

To install Starport run the following command:

```
curl https://get.starport.network/starport! | bash
```

This command will fetch the `starport` binary and install it into `/usr/local/bin`. If this command throws a permission error, lose `!` and it will download the binary in the current directory, you can then move it manually into your `$PATH`.

### Creating a Blog-Chain

Starport comes with a number of scaffolding commands that are designed to make development easier by creating everything that's necessary to start working on a particular task. One of these tasks is a `scaffold app` command which provides you with a foundation of a fresh Cosmos SDK blockchain so that you don't have to write it yourself.

To use this command, open a terminal, navigate to a directory where you have permissions to create files, and run:

```
starport app github.com/alice/blog
```

This will create a new Cosmos SDK blockchain called Blog in a `blog` directory. The source code inside the `blog` directory contains a fully functional ready-to-use blockchain. This new blockchain imports standard Cosmos SDK modules, such as [`staking`](https://docs.cosmos.network/v0.42/modules/staking/) (for delegated proof of stake), [`bank`](https://docs.cosmos.network/v0.42/modules/bank/) (for fungible token transfers between accounts), [`gov`](https://docs.cosmos.network/v0.42/modules/gov/) (for on-chain governance) and [other modules](https://docs.cosmos.network/v0.42/modules/).

Note: You can see all the command line options that Starport provides by running `starport scaffold chain --help`.

After you create the blockchain, switch to its directory:

```
cd blog
```

The `blog` directory will have a number of generated files and directories that make up the structure of a Cosmos SDK blockchain. Most of the work in this tutorial will happen in the `x` directory, but here is a quick overview of files and directories that Starport creates by default:

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


## Hello, Starport!

To get started, let's get our blockchain up and running locally on a single node.

### Starting a Blockchain

You actually have a fully-functional blockchain already. To start it on your development machine, run the following command in the `blog` directory

```
starport serve
```

This will download dependencies, compile the source code into a binary called `blogd` (by default Starport uses the name of the repo + `d`), use this binary to initialize a single validator node, and start the node.

![./images/api.png](./images/api.png)

A validator node exposes two endpoints: [http://localhost:26657](http://localhost:26657) for the low-level Tendermint API and [http://localhost:1317](http://localhost:1317) for the high-level blockchain API.

When you want to stop your blockchain, press Ctrl+C in the terminal window where it's running. In the development environment, Starport doesn't require you to restart the blockchain; changes you make in files will be automatically picked up by Starport.

### Say "Hello, Starport"

To get your Cosmos SDK blockchain to say "Hello", you will need to modify a protocol buffer file, create a keeper query function that returns data, and register a query function. Protocol buffer files contain proto `rpc`s that define Cosmos SDK queries and message handlers and proto `message`s that define Cosmos SDK types. `rpc`s are also responsible for exposing an HTTP API. [Keeper](https://docs.cosmos.network/v0.42/building-modules/keeper.html) is an abstraction for modifying the state of the blockchain and keeper functions let you query or write to the state. Registering a query needs to happen only once after you add the first query to your chain.

In terms of workflow, developers usually work with proto files first to define Cosmos SDK messages, queries, message handlers and then implement the logic of these queries and message handlers in keeper functions.

Let's start by creating a `posts` query:

```
starport query posts --response title,body
```

`query` accepts a name of the query (in our case, `posts`), an optional list of request parameters (in our case, empty), and an optional comma-separated list of response fields with a `--response` flag (in our case, `body,title`).

The `query` command has created and modified several files:

- modified `proto/blog/query.proto`
- created `x/blog/keeper/grpc_query_posts.go`
- modified `x/blog/client/cli/query.go`
- created `x/blog/client/cli/query_posts.go`

Let's examine some of these changes. For clarity, we'll skip placeholder comments Starport uses to scaffold code.

In `proto/chain/query.proto` a `Posts` `rpc` has been added to the `Query` `service`.

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
func (k Keeper) Posts(goCtx context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
  if req == nil {
    return nil, status.Error(codes.InvalidArgument, "invalid request")
  }
  ctx := sdk.UnwrapSDKContext(goCtx)
  _ = ctx
  return &types.QueryPostsResponse{}, nil
}
```

`Posts` function makes a basic check on the request and throws an error if it's `nil`, stores context in a `ctx` variable (context contains information about the environment of the request) and returns a response of type `QueryPostsResponse`. Right now the response is empty.

From the `query.proto` we know that response may contain `title` and `body`, so let's modify the last line of the function to return a "Hello!".

```go
func (k Keeper) Posts(goCtx context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
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
starport serve
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

## Creating Posts

So far, we've discussed how to modify proto files to define a new API endpoint and modify a keeper query function to return static data back to the user. Of course, a keeper can do more than return a string of data. Its purpose is to manage access to the state of the blockchain.

You can think of the state as being a collection of key-value stores. Each module is responsible for its own store. Changes to the store are triggered by transactions signed and broadcasted by users. Each transaction contains Cosmos SDK messages (not to be confused with proto `message`). When a transaction is processsed, each message gets routed to its module. A module has message handlers that process messages. Processing a message can trigger changes in the state.

### Handling Messages

A Cosmos SDK message contains information that can trigger changes in the state of a blockchain.

To create a message type and its handler, use the `message` command:

```go
starport message createPost title body
```

The `message` command accepts message name (`createPost`) and a list of fields (`title` and `body`) as arguments.

The `message` command has created and modified several files:

- modified `proto/chain/tx.proto`
- modified `x/blog/handler.go`
- created `x/blog/keeper/msg_server_createPost.go`
- modified `x/blog/client/cli/tx.go`
- created `x/blog/client/cli/txCreatePost.go`
- created `x/blog/types/message_createPost.go`
- modified `x/blog/types/codec.go`

As always, we start with a proto file. Inside `proto/chain/tx.proto`:

```go
message MsgCreatePost {
  string creator = 1;
  string title = 2;
  string body = 3;
}

message MsgCreatePostResponse {
  uint64 id = 1;
}
```

First, we define a Cosmos SDK message type with proto `message`. The `MsgCreatePost` has three fields: creator, title and body. Since the purpose of `MsgCreatePost` is to create new posts in the store, the only thing it needs to return is an ID of a created post. `CreatePost` `rpc` is added to the `Msg` `service`:

```go
service Msg {
  rpc CreatePost(MsgCreatePost) returns (MsgCreatePostResponse);
}
```

Next, let's look into `x/blog/handler.go`. Starport has added a `case` to the `switch` statement inside the `NewHandler` function. This switch statement is responsible for routing messages and calling specific keeper methods based on the type of the message

```go
func NewHandler(k keeper.Keeper) sdk.Handler {
  //...
  return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
    //...
    switch msg := msg.(type) {
    case *types.MsgCreatePost:
      res, err := msgServer.CreatePost(sdk.WrapSDKContext(ctx), msg)
      return sdk.WrapServiceResult(ctx, res, err)
    //...
    }
  }
}
```

`case *types.MsgCreatePost` handles messages of type `MsgCreatePost`, calls `CreatePost` method and returns back the response.

Every module has a handler function like this to process messages and call keeper methods.

### Processing Messages

In the newly scaffolded file `x/blog/keeper/msg_server_createPost.go` we can see a placeholder implementation of `CreatePost`. Right now it does nothing and returns an empty response. For our blog chain we want the contents of the message (title and body) to be written to the state as a new post. To do so we need to do two things: create a variable of type `Post` with title and body from the message and append this `Post` to the store.

```go
func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
  // Get the context
  ctx := sdk.UnwrapSDKContext(goCtx)
  var post = types.Post{
    Creator: msg.Creator,
    Title:   msg.Title,
    Body:    msg.Body,
  }
  // Add a post to the store and get back the ID
  id := k.AppendPost(ctx, post)
  // Return the ID of the post
  return &types.MsgCreatePostResponse{Id: id}, nil
}
```

Now we need to define both `Post` type and `AppendPost` keeper method.

The `Post` type can be defined in a proto file and Starport (with the help of `protoc`) will take care of generating required Go files.

Create a new file `proto/blog/post.proto` and define the `Post` `message`:

```go
syntax = "proto3";
package alice.blog.blog;
option go_package = "github.com/alice/blog/x/blog/types";

message Post {
  string creator = 1;
  uint64 id = 2;
  string title = 3; 
  string body = 4; 
}
```

The contents of `post.proto` are fairly standard. We define a package name (that is used to identify messages, among other things), specify in which Go package new files should be generated, and finally define `message Post`. Now, after we build and start our chain with Starport, the `Post` type will be available.

### Writing Data to the Store

The next step is to define `AppendPost` keeper method. Let's create a new file `x/blog/keeper/post.go` and start thinking about the logic of the function.

To implement `AppendPost` we first need to understand how does the store works. We can think of a store as a key-value database, where keys are lexicographically ordered. We can loop through keys and `Get` and `Set` values based on keys. To distinguish between different types of data that a module can keep in its store, we use prefixes like `product-` or `post-`.

We want to keep a list of posts in what is essentially a key-value store, which means we need to keep track of the index of the posts we insert ourselves. Since we'll be keeping both post values and post count (index) in the store, let's use different prefixes: `Post-value-` and `Post-count-`, respectively. Add the following to `x/blog/types/keys.go`:

```go
const (
  PostKey      = "Post-value-"
  PostCountKey = "Post-count-"
)
```

`AppendPost` given a `Post` should do four things: get the number of posts in the store (count), add a post by using the count as an ID, increment the count, and return the count.

Let's draft the `AppendPost` function in `x/blog/keeper/post.go`:

```go
// func (k Keeper) AppendPost() uint64 {
// 	 count := k.GetPostCount()
// 	 store.Set()
// 	 k.SetPostCount()
// 	 return count
// }
```

Let's implement `GetPostCount` first.

```go
func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
  // Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostCountKey))
  // Convert the PostCountKey to bytes
  byteKey := []byte(types.PostCountKey)
  // Get the value of the count
  bz := store.Get(byteKey)
  // Return zero if the count value is not found (for example, it's the first post)
  if bz == nil {
    return 0
  }
  // Convert the count into a uint64
  count, err := strconv.ParseUint(string(bz), 10, 64)
  // Panic if the count cannot be converted to uint64
  if err != nil {
    panic("cannot decode count")
  }
  return count
}
```

Now that `GetPostCount` returns the correct number of posts in the store, let's implement `SetPostCount`:

```go
func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
  // Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostCountKey))
  // Convert the PostCountKey to bytes
  byteKey := []byte(types.PostCountKey)
  // Convert count from uint64 to string and get bytes
  bz := []byte(strconv.FormatUint(count, 10))
  // Set the value of Post-count- to count
  store.Set(byteKey, bz)
}
```

Now that we've implemented functions for getting the number of posts and setting the post count, we can implement the logic behind `AppendPost`:

```go
package keeper

import (
  "encoding/binary"
  "github.com/alice/blog/x/blog/types"
  "github.com/cosmos/cosmos-sdk/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "strconv"
)

func (k Keeper) AppendPost(ctx sdk.Context, post types.Post) uint64 {
  // Get the current number of posts in the store
  count := k.GetPostCount(ctx)
  // Assign an ID to the post based on the number of posts in the store
  post.Id = count
  // Get the store
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostKey))
  // Convert the post ID into bytes
  byteKey := make([]byte, 8)
  binary.BigEndian.PutUint64(byteKey, post.Id)
  // Marshal the post into bytes
  appendedValue := k.cdc.MustMarshalBinaryBare(&post)
  // Insert the post bytes using post ID as a key
  store.Set(byteKey, appendedValue)
  // Update the post count
  k.SetPostCount(ctx, count+1)
  return count
}
```

We've implemented all the code necessary to create new posts and store them on chain. Now, when a transaction containing a message of type `MsgCreatePost` is broadcasted, Cosmos SDK will route the message to our blog module. `x/blog/handler.go` calls `k.CreatePost`, which in turn calls `AppendPost`. `AppendPost` gets the number of posts from the store, adds a post using the count as an ID, increments the count, and returns the ID.

Let's try it out! If the chain is yet not started, run `starport serve`. Create a post:

```
blogd tx blog createPost foo bar --from alice
```

```
"body":{"messages":[{"@type":"/alice.blog.blog.MsgCreatePost","creator":"cosmos1dad8xvsj3dse928r52yayygghwvsggvzlm730p","title":"foo","body":"bar"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
{"height":"6861","txhash":"6086372860704F5F88F4D0A3CF23523CF6DAD2F637E4068B92582E3BB13800DA","codespace":"","code":0,"data":"0A100A0A437265617465506F737412020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreatePost\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44674","tx":null,"timestamp":""}
```

Now that we've added the functionality to create posts and broadcast them to our chain, let's add querying.

## Displaying Posts

There are two components responsible for querying data: `rpc` inside `service Query` in a proto file (that defines data types and specifies the HTTP API endpoint) and a keeper method that performs the querying from the key-value store.

Let's first review the services and messages in `x/blog/query.proto`. `Posts` `rpc` accepts an empty request and returns an object with two fields: title and body. We would like for it to return a list of posts, instead. The list of posts can be long, so let's also add pagination. For pagination, both request and response should include a page number: we want to be able to request a particular page and we need to know what page has been returned.

`x/blog/query.proto`:

```go
// Import the Post message
import "chain/post.proto";

message QueryPostsRequest {
  // Adding pagination to request
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryPostsResponse {
  // Returning a list of posts
  repeated Post Post = 1;
  // Adding pagination to response
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

Once the types are defined in proto files, we can implement post querying logic. In `grpc_query_posts.go`:

```go
func (k Keeper) Posts(c context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
  // Throw an error if request is nil
  if req == nil {
    return nil, status.Error(codes.InvalidArgument, "invalid request")
  }
  // Define a variable that will store a list of posts
  var posts []*types.Post
  // Get context with the information about the environment
  ctx := sdk.UnwrapSDKContext(c)
  // Get the key-value module store using the store key (in our case store key is "chain")
  store := ctx.KVStore(k.storeKey)
  // Get the part of the store that keeps posts (using post key, which is "Post-value-")
  postStore := prefix.NewStore(store, []byte(types.PostKey))
  // Paginate the posts store based on PageRequest
  pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
    var post types.Post
    if err := k.cdc.UnmarshalBinaryBare(value, &post); err != nil {
      return err
    }
    posts = append(posts, &post)
    return nil
  })
  // Throw an error if pagination failed
  if err != nil {
    return nil, status.Error(codes.Internal, err.Error())
  }
  // Return a struct containing a list of posts and pagination info
  return &types.QueryPostsResponse{Post: posts, Pagination: pageRes}, nil
}
```

## Using CLI to Create and Display Posts

Having implemened logic for both creating and querying posts we can use the node's binary to interact with our chain. To create a post:

```
blogd tx blog createPost foo bar --from alice
```

```
{"body":{"messages":[{"@type":"/alice.chain.chain.MsgCreatePost","creator":"cosmos1c9zy9aajk9fs2f8ygtz4pm22r3rxmg597vw2n3","title":"foo","body":"bar"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
{"height":"2828","txhash":"E04A712E65B0F6F30F5DC291A6552B69F6CB3F77761F28AFFF8EAA535EC4C589","codespace":"","code":0,"data":"0A100A0A437265617465506F737412020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreatePost\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44674","tx":null,"timestamp":""}
```

To query the list of all posts on chain:

```
blogd q blog posts
```

```
Post:
- body: bar
  creator: cosmos1c9zy9aajk9fs2f8ygtz4pm22r3rxmg597vw2n3
  id: "0"
  title: foo
pagination:
  next_key: null
  total: "1"
```
