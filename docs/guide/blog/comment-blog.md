---
description: Write a query that returns a blog post by ID with associated comments.
order: 2
---

# Add associated comments to a blog post

In this tutorial, you create a new message to add comments to a blog post.  

By completing this tutorial, you will learn about:

* Scaffolding a new `list` with proto functions and keeper functions
* Adding comments to existing blog posts
* Querying for blog posts that have associated comments
* Deleting comments from a blog post
* Implementing logic for writing comments to the blockchain

**Note:** For this tutorial, adding comments is available only to blog posts that are no older than 100 blocks. The 100 block value has been hard coded for rapid testing. You can increase the block count to a larger number to achieve a longer time period before commenting is disabled.

## Prerequisites

This tutorial is an extension of and requires completion of the [Module Basics: Build a Blog](index.md) tutorial. 

## Core concepts 

This tutorial relies on the `blog` blockchain that you built in the `Build a Blog Tutorial.`

## Fetch functions using list command

To get the useful functions for this tutorial, you use the `ignite scaffold list NAME [field]... [flags]` command. Make sure to familiarize yourself with the command.

1. Navigate to the `blog` directory that you created in the [Build a blog](index.md) tutorial.

2. To create the source code files to add CRUD (create, read, update, and delete) functionality for data stored as an array, run:

```bash
ignite scaffold list comment --no-message creator:string title:string body:string postID:uint createdAt:int 
```

The `--no-message` flag disables CRUD interaction messages scaffolding because you will write your own messages.

The command output shows the files that were created and modified:

```text
create proto/blog/comment.proto
modify proto/blog/genesis.proto
modify proto/blog/query.proto
modify vue/src/views/Types.vue
modify x/blog/client/cli/query.go
create x/blog/client/cli/query_comment.go
create x/blog/client/cli/query_comment_test.go
modify x/blog/genesis.go
modify x/blog/genesis_test.go
create x/blog/keeper/comment.go
create x/blog/keeper/comment_test.go
create x/blog/keeper/grpc_query_comment.go
create x/blog/keeper/grpc_query_comment_test.go
modify x/blog/module.go
modify x/blog/types/genesis.go
modify x/blog/types/genesis_test.go
modify x/blog/types/keys.go

🎉 comment added.
```

Make a small modification in proto/blog/comment.proto to change createdAt to int64:

```go
message Comment {
  uint64 id = 1;
  string creator = 2; 
  string title = 3; 
  string body = 4; 
  uint64 postID = 5; 
  int64 createdAt = 6;  
}
```


## Add a comment to a post

To create a new message that adds a comment to the existing post, run:

```bash
ignite scaffold message create-comment postID:uint title body
```

The `ignite scaffold message` command accepts `postID` and a list of fields as arguments. The fields are `title` and `body`.

Here, `postID` is the reference to previously created blog post.

The `message` command has created and modified several files:

```text
modify proto/blog/tx.proto
modify x/blog/client/cli/tx.go
create x/blog/client/cli/tx_create_comment.go
modify x/blog/handler.go
create x/blog/keeper/msg_server_create_comment.go
modify x/blog/types/codec.go
create x/blog/types/message_create_comment.go
create x/blog/types/message_create_comment_test.go

🎉 Created a message `create-comment`.
```

As always, start your development with a proto file. 

In the `proto/blog/tx.proto` file, edit `MsgCreateComment` to:

* Add `id`
* Define the `id` for `message MsgCreateCommentResponse`:

```go
message MsgCreateComment {
  string creator = 1;
  uint64 postID = 2;
  string title = 3;
  string body = 4;
  uint64 id = 5;
}

message MsgCreateCommentResponse {
  uint64 id = 1;
}
```

 You see in the `proto/blog/tx.proto` file that the `MsgCreateComment` has five fields: creator, title, body, PostID, and id. Since the purpose of the `MsgCreateComment` message is to create new comments in the store, the only thing the message needs to return is an ID of a created comments. The `CreateComment` rpc was already added to the `Msg` service:

```go
rpc CreateComment(MsgCreateComment) returns (MsgCreateCommentResponse);
```

Now, add the `id` field to `MsgCreatePost`: 

```go
message MsgCreatePost {
  string creator = 1;
  string title = 2;
  string body = 3;
  uint64 id = 4;
}
```

## Process messages

In the newly scaffolded `x/blog/keeper/msg_server_create_comment.go` file, you can see a placeholder implementation of the `CreateComment` function (marked with `//TODO`). Right now it does nothing and returns an empty response. For your blog chain, you want the contents of the message (title and body) to be written to the state as a new comment.

You need to do the following things:

* Create a variable of type `Comment` with title and body from the message
* Check if the the comment posted for the respective blog id exists and comment is not older than 100 blocks
* Append this `Comment` to the store

```go
import (
    //...
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateComment(goCtx context.Context, msg *types.MsgCreateComment) (*types.MsgCreateCommentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	post := k.GetPost(ctx, msg.PostID)
	postId := post.Id
	
	// Check if the Post Exists for which a comment is being created
	if msg.PostID == 0 {
		return nil, sdkerrors.Wrapf(types.ErrID, "Post Blog Id does not exist for which comment with Blog Id %d was made", msg.PostID)
	}

	// Create variable of type comment
	var comment = types.Comment{
		Creator:   msg.Creator,
		Id:        msg.Id,
		Body:      msg.Body,
		Title:     msg.Title,
		PostID:    msg.PostID,
		CreatedAt: ctx.BlockHeight(),
	}
	
	// Check if the comment is older than the Post. If more than 100 blocks, then return error.
	if comment.CreatedAt > post.CreatedAt+100 {
		return nil, sdkerrors.Wrapf(types.ErrCommentOld, "Comment created at %d is older than post created at %d", comment.CreatedAt, post.CreatedAt)
	}

	id := k.AppendComment(ctx, comment)
	return &types.MsgCreateCommentResponse{Id: id}, nil
}
```

When the Comment validity is checked, it throws 2 error messages - `ErrID` and `ErrCommendOld`. You can define the error messages by navigating to `x/blog/types/errors.go` and replacing the current values in 'var' with:

```go
//...
var (
	ErrCommentOld = sdkerrors.Register(ModuleName, 1300, "")
	ErrID = sdkerrors.Register(ModuleName, 1400, "")
)
```


In the existing `x/blog/keeper/msg_server_create_post.go` file, you need to make a modification to add `createdAt`

```go
func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	// Get the context
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Create variable of type Post
	var post = types.Post{
		Creator:   msg.Creator,
		Id:        msg.Id,
		Title:     msg.Title,
		Body:      msg.Body,
		CreatedAt: ctx.BlockHeight(),
	}
	// Add a post to the store and get back the ID
	id := k.AppendPost(ctx, post)
	// Return the ID of the post
	return &types.MsgCreatePostResponse{Id: id}, nil
}
```

## Write data to the store

When you define the `Comment` type in a proto file, Ignite CLI (with the help of `protoc`) takes care of generating the required Go files.

Inside the `proto/blog/comment.proto` file, you can observe, Ignite CLI has already added the required fields inside the `Comment` message.

The contents of the `comment.proto` file are fairly standard and similar to `post.proto`. The file defines a package name that is used to identify messages, among other things, specifies the Go package where new files are generated, and finally defines `message Comment`. 

Each file save triggers an automatic rebuild.  Now, after you build and start your chain with Ignite CLI, the `Comment` type is available.

Also, make a small modification in `proto/blog/post.proto` to add `createdAt`:

```go
//...

message Post {
 //...
  int64 createdAt = 5;
}
```

### Define keeper methods

The function `ignite scaffold list comment --no-message` has fetched all of the required functions for keeper. 

Inside `x/blog/types/keys.go` file, you can see that the `comment-value` and `comment-count` keys are added.

## Write data to the store

In `x/blog/keeper/post.go`, add a new function to get the post:

```go
import (
	"encoding/binary"
	"github.com/username/blog/x/blog/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//...

func (k Keeper) GetPost(ctx sdk.Context, id uint64) (post types.Post) {
	// Get the store using storeKey (which is "blog") and PostKey (which is "Post-")
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostKey))
	// Convert the post ID into bytes
	byteKey := make([]byte, 8)
	binary.BigEndian.PutUint64(byteKey, id)
	// Get the post bytes using post ID as a key
	bz := store.Get(byteKey)
	// Unmarshal the post bytes into the post object
	k.cdc.MustUnmarshal(bz, &post)
	return post
}
```

You have manually added the functions to `x/blog/keeper/post.go`. 

When you ran the `ignite scaffold list comment --no-message` command, these functions are automatically implemented in `x/blog/keeper/comment.go`:

- `GetCommentCount`
- `SetCommentCount`
- `AppendCommentCount`

By following these steps, you have implemented all of the code required to create comments and store them on-chain. Now, when a transaction that contains a message of type `MsgCreateComment` is broadcast, the message is routed to your blog module.

- `x/blog/handler.go` calls `k.CreateComment` which in turn calls `AppendComment`.
- `AppendComment` gets the number of comments from the store, adds a comment using the count as an ID, increments the count, and returns the ID.

## Create the delete-comment message

To create a message, use the `message` command:

```bash
ignite scaffold message delete-comment commentID:uint postID:uint 
```

The `message` commands accepts `commentID` and `postID` as arguments.

Here, `commentID` and `postID` are the references to previously created comment and blog post.

```bash
The `message` command has created and modified several files:

modify proto/blog/tx.proto
modify x/blog/client/cli/tx.go
create x/blog/client/cli/tx_delete_comment.go
modify x/blog/handler.go
create x/blog/keeper/msg_server_delete_comment.go
modify x/blog/types/codec.go
create x/blog/types/message_delete_comment.go
create x/blog/types/message_delete_comment_test.go
```

As always, start your development with a proto file. 

In the `proto/blog/tx.proto` file, edit `MsgDeleteComment` to:

* Add `id`
* Define the `id` for `message MsgDeleteCommentResponse`:

```go
message MsgDeleteComment {
  string creator = 1;
  uint64 commentID = 2;
  uint64 postID = 3;
  uint64 id = 4;
}

message MsgDeleteCommentResponse {
  uint64 id = 1;
}
```

## Process messages

In the newly scaffolded `x/blog/keeper/msg_server_delete_comment.go` file, you can see a placeholder implementation of the `DeleteComment` function. Right now it does nothing and returns an empty response. 

For your blog chain, you want to delete the contents of the comment. Add the code to:

- Check if the post Id exists to see which comment was deleted.
- Delete the comment from the store.

```go
package keeper

 import (
	"context"
	"encoding/binary"

	"github.com/username/blog/x/blog/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
 )

 func (k msgServer) DeleteComment(goCtx context.Context, msg *types.MsgDeleteComment) (*types.MsgDeleteCommentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	comment, exist := k.GetComment(ctx, msg.Id)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrID, "Comment doesnt exist")
	}

	if msg.PostID != comment.PostID {
		return nil, sdkerrors.Wrapf(types.ErrID, "Post Blog Id does not exist for which comment with Blog Id %d was made", msg.PostID)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, comment.Id)
	store.Delete(bz)

	return &types.MsgDeleteCommentResponse{}, nil
}

```

## Display posts

Implement logic to query existing posts:

```bash
ignite scaffold query comments id:uint --response title,body
```

Also in `proto/blog/query.proto`, make these updates:

```go
import "blog/post.proto";

message QueryCommentsRequest {
	uint64 id = 1;
    // Adding pagination to request
    cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
//...
message QueryCommentsResponse {
  Post Post = 1;
  	// Returning a list of comments
  repeated Comment Comment = 2;
    // Adding pagination to response
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

After the types are defined in proto files, you can implement post querying logic in `x/blog/keeper/grpc_query_comments.go` by registering the `Comments` function:

```go

import (
	"context"

	"github.com/username/blog/x/blog/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) Comments(c context.Context, req *types.QueryCommentsRequest) (*types.QueryCommentsResponse, error) {

	// Throw an error if request is nil
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Define a variable that will store a list of posts
	var comments []*types.Comment
	
	// Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get the key-value module store using the store key (in this case store key is "chain")
	store := ctx.KVStore(k.storeKey)
	
	// Get the part of the store that keeps posts (using post key, which is "Post-value-")
	commentStore := prefix.NewStore(store, []byte(types.CommentKey))

	// Get the post by ID 
	post := k.GetPost(ctx, req.Id)

	// Get the post ID
	postID := post.Id
	
	// Paginate the posts store based on PageRequest
	pageRes, err := query.Paginate(commentStore, req.Pagination, func(key []byte, value []byte) error {
		var comment types.Comment
		if err := k.cdc.Unmarshal(value, &comment); err != nil {
			return err
		}

		if comment.PostID == postID {
			comments = append(comments, &comment)	
		}
		
		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of posts and pagination info
	return &types.QueryCommentsResponse{Post: &post, Comment: comments, Pagination: pageRes}, nil
}

```

**Note:** Since gRPC has been already added to module handler in the previous tutorial, you don't need to add it again.

## Create post and comment

Try it out! 

If the chain is yet not started, run `ignite chain serve`.

Create a post:

```bash
blogd tx blog create-post Uno "This is the first post" --from alice
```

As before, you are prompted to confirm the transaction:

```bash
"body":{"messages":[{"@type":"/username.blog.blog.MsgCreatePost","creator":"cosmos1dad8xvsj3dse928r52yayygghwvsggvzlm730p","title":"foo","body":"bar"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

The transaction output:

```bash
{"height":"6861","txhash":"6086372860704F5F88F4D0A3CF23523CF6DAD2F637E4068B92582E3BB13800DA","codespace":"","code":0,"data":"0A100A0A437265617465506F737412020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreatePost\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"CreatePost"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44674","tx":null,"timestamp":""}
```

Create a comment:

```bash
blogd tx blog create-comment 0  Uno "This is the first comment" --from alice
```

```bash
{"body":{"messages":[{"@type":"/username.blog.blog.MsgCreateComment","creator":"cosmos17pvwgu36fu37j8y9gc4pasxsj3p26ghmlqvngd","id":"0","title":"Uno","body":"This is the first comment","PostID":"2","createdAt":"0"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}
```

When prompted, press Enter to confirm the transaction:

```bash
confirm transaction before signing and broadcasting [y/N]: y
```

The command output shows the results of the transaction:

```bash
code: 0
codespace: ""
data: 0A270A252F636F736D6F6E6175742E626C6F672E626C6F672E4D7367437265617465436F6D6D656E74
gas_used: "45891"
gas_wanted: "200000"
height: "118"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: CreateComment
    type: message
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"CreateComment"}]}]}]'
timestamp: ""
tx: null
txhash: 0CAFC113D1C73BC0210EFEA5964EBD2EB530311169FB442C5CBF0B5E92521C41
```

## Display post and comment

```bash
blogd q blog comments 0
```

The results are output:

```bash
Comment:
- body: Let us add random comment
  createdAt: "14094"
  creator: cosmos1g7x9cpj6w0jklshe3se57tlwydx6yfl8ex5g7n
  id: "0"
  postID: "0"
  title: comment
Post:
  body: This is the first post
  createdAt: "14046"
  creator: cosmos1g7x9cpj6w0jklshe3se57tlwydx6yfl8ex5g7n
  id: "0"
  title: Uno
```

## Delete comment

```bash
blogd tx blog delete-comment 0 0 --from alice -y
```

The results are output:

```bash
code: 0
codespace: ""
data: 0A270A252F636F736D6F6E6175742E626C6F672E626C6F672E4D736744656C657465436F6D6D656E74
gas_used: "40125"
gas_wanted: "200000"
height: "253"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: DeleteComment
    type: message
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"DeleteComment"}]}]}]'
timestamp: ""
tx: null
txhash: 0312234CBB9EEA1A59D474496E100AFC5A460A0E60E7D009D3E9417530148A75
```

## Display the post and all associated comments

```bash
blogd q blog comments 0
```

The results are output:

```bash
Comment:
[]
Post:
  body: This is the first post
  createdAt: "14046"
  creator: cosmos1g7x9cpj6w0jklshe3se57tlwydx6yfl8ex5g7n
  id: "0"
  title: Uno
```

## Edge cases

1. Add comment to a nonexistent blog id:

```bash
blogd tx blog create-comment 53 "Edge1"  "This is the 53 comment" --from alice -y
```

The transaction is not able to be completed because the blog id does not exist:

```bash
code: 1400
codespace: blog
data: ""
gas_used: "38151"
gas_wanted: "200000"
height: "1019"
info: ""
logs: []
raw_log: 'failed to execute message; message index: 0: Post Blog Id 53 does not exist for which comment was made: '
timestamp: ""
tx: null
txhash: B99BD295A0B08DF58B9FEC8EB41D467C2F28BD4EC8CDB56FBF30DB728B877ABA
```

1. Add comment to a blog post that is older than 100 blocks:

```bash
blogd tx blog create-comment 0 "Comment" "This is a comment" --from alice -y
```

The transaction is not executed:

```bash
code: 1300
codespace: blog
data: ""
gas_used: "38101"
gas_wanted: "200000"
height: "1191"
info: ""
logs: []
raw_log: 'failed to execute message; message index: 0: Comment created at 1191 is older than post created at 1047:'
timestamp: ""
tx: null
txhash: A87AAD5E2E6A26F9B80796D013139E9A18DB286D9CF769BC6AA6601DD64C6A35
```

## Conclusion

Congratulations. You have added comments to your blog blockchain! 

You have successfully completed these steps:

* Scaffolding a new `list` with proto functions and keeper functions
* Add comments to existing blog post
* Display the blog post by ID with associated comments
* Delete comments from a given blog post
