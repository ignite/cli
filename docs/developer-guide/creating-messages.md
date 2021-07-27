---
order: 7
description: Creating messages
---

# Creating Posts

We discussed how to modify proto files to define a new API endpoint and modify a keeper query function to return static data back to the user. Of course, a keeper can do more than return a string of data. Its purpose is to manage access to the state of the blockchain.

You can think of the state as being a collection of key-value stores. Each module is responsible for its own store. Changes to the store are triggered by transactions signed and broadcasted by users. Each transaction contains Cosmos SDK messages (not to be confused with proto `message`). When a transaction is processsed, each message gets routed to its module. A module has message handlers that process messages. Processing a message can trigger changes in the state.

## Handling Messages

A Cosmos SDK message contains information that can trigger changes in the state of a blockchain.

To create a message type and its handler, use the `message` command:

```go
starport scaffold message createPost title body
```

The `message` command accepts message name (`createPost`) and a list of fields (`title` and `body`) as arguments.

The `message` command has created and modified several files:

- modified `proto/blog/tx.proto`
- modified `x/blog/handler.go`
- created `x/blog/keeper/msg_server_createPost.go`
- modified `x/blog/client/cli/tx.go`
- created `x/blog/client/cli/txCreatePost.go`
- created `x/blog/types/message_createPost.go`
- modified `x/blog/types/codec.go`

As always, we start with a proto file. Inside `proto/blog/tx.proto`:

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

`case *types.MsgCreatePost` handles messages of type `MsgCreatePost`, calls `CreatePost` method, and returns back the response.

Every module has a handler function like this to process messages and call keeper methods.
