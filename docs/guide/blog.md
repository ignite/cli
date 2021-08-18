---
title: "Module Basics: Blog"
order: 3
description: Learn module basics by writing and reading blog posts to your chain.
---

# Building a Blog

Learn module basics by building a blockchain app to write and read blog posts. 

By completing this tutorial, you will:

* Write and read blog posts to your chain
* Scaffold a Cosmos SDK message
* Define new types in protocol buffer files
* Write keeper methods to write data to the store
* Read data from the store and return it as a result a query
* Use the blockchain's CLI to broadcast transactions

### Prerequisite 

To complete this tutorial, you will need:

- A supported version of Starport. This tutorial is verified for Starport 0.17.2. See [Install Starport](./install.md). 

## Create Your Blog Chain

First, create a new blockchain. 

Open a terminal and navigate to a directory where you have permissions to create files. To create your Cosmos SDK blockchain, run this command:

```bash
starport scaffold chain github.com/cosmonaut/blog
```

The `blog` directory is created with the default directory structure. 

## High Level Transaction Review 

So far, you have learned how to modify proto files to define a new API endpoint and modify a keeper query function to return static data back to the user. Of course, a keeper can do more than return a string of data. Its purpose is to manage access to the state of the blockchain.

You can think of the state as being a collection of key-value stores. Each module is responsible for its own store. Changes to the store are triggered by transactions that are signed and broadcasted by users. Each transaction contains Cosmos SDK messages (not to be confused with proto `message`). When a transaction is processed, each message gets routed to its module. A module has message handlers that process messages. Processing a message can trigger changes in the state.

## Create Message Types

A Cosmos SDK message contains information that can trigger changes in the state of a blockchain.

First, change into the `blog` directory:

```bash
cd blog
```

To create a message type and its handler, use the `message` command:

```bash
starport scaffold message createPost title body
```

The `message` command accepts message name (`createPost`) and a list of fields (`title` and `body`) as arguments.

The `message` command has created and modified several files:

* modified `proto/blog/tx.proto`
* modified `x/blog/handler.go`
* created `x/blog/keeper/msg_server_createPost.go`
* modified `x/blog/client/cli/tx.go`
* created `x/blog/client/cli/txCreatePost.go`
* created `x/blog/types/message_createPost.go`
* modified `x/blog/types/codec.go`

As always, start with a proto file. Inside the `proto/blog/tx.proto` file, the `MsgCreatePost` message has been created:

```go
message MsgCreatePost {
  string creator = 1;
  string title = 2;
  string body = 3;
}

message MsgCreatePostResponse {
}
```

First, define a Cosmos SDK message type with proto `message`. The `MsgCreatePost` has three fields: creator, title and body. Since the purpose of the `MsgCreatePost` message is to create new posts in the store, the only thing the message needs to return is an ID of a created post. The `CreatePost` rpc was already added to the `Msg` service:

```go
service Msg {
  rpc CreatePost(MsgCreatePost) returns (MsgCreatePostResponse);
}
```

Next, look at the `x/blog/handler.go` file. Starport has added a `case` to the `switch` statement inside the `NewHandler` function. This switch statement is responsible for routing messages and calling specific keeper methods based on the type of the message:

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

The `case *types.MsgCreatePost` statement handles messages of type `MsgCreatePost`, calls the `CreatePost` method, and returns back the response.

Every module has a handler function like this to process messages and call keeper methods.

## Process Messages

In the newly scaffolded `x/blog/keeper/msg_server_create_post.go` file, you can see a placeholder implementation of the `CreatePost`. Right now it does nothing and returns an empty response. For your blog chain, you want the contents of the message (title and body) to be written to the state as a new post. 

You need to do two things: 

- Create a variable of type `Post` with title and body from the message 
- Append this `Post` to the store

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
## Write Data to the Store 

Define the `Post` type and the `AppendPost` keeper method.

When you define the `Post` type in a proto file, Starport (with the help of `protoc`) takes care of generating the required Go files.

Create the `proto/blog/post.proto` file and define the `Post` message:

```go
syntax = "proto3";
package cosmonaut.blog.blog;
option go_package = "github.com/cosmonaut/blog/x/blog/types";

message Post {
  string creator = 1;
  uint64 id = 2;
  string title = 3; 
  string body = 4; 
}
```

The contents of the `post.proto` file are fairly standard. The file defines a package name that is used to identify messages, among other things, specifies the Go package where new files are generated, and finally defines `message Post`. 

Now, after you build and start your chain with Starport, the `Post` type is available.

### Define Keeper Methods 

The next step is to define the `AppendPost` keeper method. Create the `x/blog/keeper/post.go` file and start thinking about the logic of the function.

To implement `AppendPost` you must first understand how the store works. You can think of a store as a key-value database where keys are lexicographically ordered. You can loop through keys and use `Get` and `Set` to retrieve and set values based on keys. To distinguish between different types of data that a module can keep in its store, you can use prefixes like `product-` or `post-`.

To keep a list of posts in what is essentially a key-value store, you need to keep track of the index of the posts you insert. Since both post values and post count (index) values are kept in the store, you can use different prefixes: `Post-value-` and `Post-count-`. 

Add these prefixes to the `x/blog/types/keys.go` file:

```go
const (
  PostKey      = "Post-value-"
  PostCountKey = "Post-count-"
)
```

When a `Post` message is sent to the `AppendPost` function, four actions occur: 

- Get the number of posts in the store (count)
- Add a post by using the count as an ID
- Increment the count
- Return the count

## Write Data to the Store

Now, after the `import` section, draft the `AppendPost` function in the `x/blog/keeper/post.go` file:

```go
// func (k Keeper) AppendPost() uint64 {
// 	 count := k.GetPostCount()
// 	 store.Set()
// 	 k.SetPostCount()
// 	 return count
// }
```

First, implement `GetPostCount`:

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

Now that `GetPostCount` returns the correct number of posts in the store, implement `SetPostCount`:

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

Now that you have implemented functions for getting the number of posts and setting the post count, you can implement the logic behind `AppendPost`. The `import` section is also shown in this example:

```go
package keeper

import (
  "encoding/binary"
  "github.com/cosmonaut/blog/x/blog/types"
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

By following these steps, you have implemented all of the code required to create new posts and store them on-chain. Now, when a transaction that contains a message of type `MsgCreatePost` is broadcast, the message is routed to your blog module.

- `x/blog/handler.go` calls `k.CreatePost` which in turn calls `AppendPost`. 
- `AppendPost` gets the number of posts from the store, adds a post using the count as an ID, increments the count, and returns the ID.

Try it out! If the chain is yet not started, run `starport chain serve`. 

Create a post:

```bash
blogd tx blog create-post foo bar --from alice
```

```bash
"body":{"messages":[{"@type":"/cosmonaut.blog.blog.MsgCreatePost","creator":"cosmos1dad8xvsj3dse928r52yayygghwvsggvzlm730p","title":"foo","body":"bar"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
{"height":"6861","txhash":"6086372860704F5F88F4D0A3CF23523CF6DAD2F637E4068B92582E3BB13800DA","codespace":"","code":0,"data":"0A100A0A437265617465506F737412020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreatePost\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44674","tx":null,"timestamp":""}
```

Now that you have added the functionality to create posts and broadcast them to our chain, you can add querying.

## Display Posts

```bash
starport scaffold query posts --response title,body
```

Two components are responsible for querying data: 

- An rpc inside `service Query` in a proto file that defines data types and specifies the HTTP API endpoint 
- A keeper method that performs the querying from the key-value store

First, review the services and messages in `proto/blog/query.proto`. The `Posts` rpc accepts an empty request and returns an object with two fields: title and body. Now you can make changes so it can return a list of posts. The list of posts can be long, so add pagination. When pagination is added, the request and response include a page number so you can request a particular page when you know what page has been returned.

In the `proto/blog/query.proto` file:

```go
// Import the Post message
import "blog/post.proto";

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

After the types are defined in proto files, you can implement post querying logic. In `grpc_query_posts.go`:

```go
package keeper

import (
  "context"

  "github.com/cosmonaut/blog/x/blog/types"
  "github.com/cosmos/cosmos-sdk/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/cosmos/cosmos-sdk/types/query"  
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

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

Now that you have implemented logic for creating and querying posts, you can use the node's binary to interact with your chain. Your blog chain binary is `blogd`.

To create a post:

```bash
blogd tx blog create-post foo bar --from alice
```

```bash
{"body":{"messages":[{"@type":"/cosmonaut.blog.blog.MsgCreatePost","creator":"cosmos1c9zy9aajk9fs2f8ygtz4pm22r3rxmg597vw2n3","title":"foo","body":"bar"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
{"height":"2828","txhash":"E04A712E65B0F6F30F5DC291A6552B69F6CB3F77761F28AFFF8EAA535EC4C589","codespace":"","code":0,"data":"0A100A0A437265617465506F737412020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreatePost\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44674","tx":null,"timestamp":""}
```

To query the list of all on-chain posts:

```bash
blogd q blog posts
```

```bash
Post:
- body: bar
  creator: cosmos1c9zy9aajk9fs2f8ygtz4pm22r3rxmg597vw2n3
  id: "0"
  title: foo
pagination:
  next_key: null
  total: "1"
```

## Conclusion 
 
Congratulations. You have built a blog blockchain.
