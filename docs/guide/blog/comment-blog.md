---
description: Add comment on Blog
order: 2
---

# Add comment to Blog post Blockchain 

In this tutorial, you will create a new message module called comment. The module will let you read and write comments to an existing blog post.

You can only add comments to post which are no older than 100 Blocks. 

**Note:** This value has been hard coded to a low number for rapid testing. You can increase it to a greater number to achieve longer period of time before commenting is not allowed.

### Prerequisites:

- This tutorial is an extension of previously written blog tutorial. Make sure you complete that first before proceeding with this tutorial.
- This tutorial also assumes basic knowledge of blog tutorial implementation.
- Make sure you are inside the `blog` directory created in the previous blog tutorial.

## Create a new message called comment

To create a new message module for comment, use the `message` command:

```bash
starport scaffold message create-comment blogID:int title body
```

The `message` commands accepts `blogID` and a list of fields (`title` and `body` as arguments )
Here, `blogID` is the reference to previously created blog posts.

The `message` command has created and modified several files:

----------->>> TODO

As always, start with a proto file. Inside the `proto/blog/tx.proto` file, the `MsgCreateComment` message has been created. Edit the file to define the id for `message MsgCreateCommentResponse`:

```go
message MsgCreateComment {
  string creator = 1;
  uint64 id = 2;
  string title = 3;
  string body = 4;
  uint64 blogID = 5;
  int64 createdAt = 6;
}

message MsgCreateCommentResponse {
  uint64 id = 1;
}
```

First, define a Cosmos SDK message type with proto `message`. The `MsgCreateComment` has three fields: creator, title, body, blogID and createdAt. Since the purpose of the `MsgCreateComment` message is to create new comments in the store, the only thing the message needs to return is an ID of a created comments. The `CreateComment` rpc was already added to the `Msg` service:

```go
  rpc CreateComment(MsgCreateComment) returns (MsgCreateCommentResponse);
```

Next, look at the `x/blog/handler.go` file. Starport has added a `case` to the `switch` statement inside the `NewHandler` function. This switch statement is responsible for routing messages and calling specific keeper methods based on the type of the message. `case *types.MsgCreateComment` has been added along with `case *types.MsgCreatePost`

```go
func NewHandler(k keeper.Keeper) sdk.Handler {
	//...
		switch msg := msg.(type) {
		//...
		case *types.MsgCreateComment:
			res, err := msgServer.CreateComment(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			//...
		}
	}
}
```

The `case *types.MsgCreateComment` statement handles messages of type `MsgCreateComment`, calls the `CreateComment` method, and returns back the response.

Every module has a handler function like this to process messages and call keeper methods.

## Process Messages

In the newly scaffolded `x/blog/keeper/msg_server_create_comment.go` file, you can see a placeholder implementation of the `CreateComment` function. Right now it does nothing and returns an empty response. For your blog chain, you want the contents of the message (title and body) to be written to the state as a new comment.

You need to do two things:

- Create a variable of type `Comment` with title and body from the message
- Check if the the comment posted for the respective blog id exists and comment is not older than 100 blocks.
- Append this `Comment` to the store

```go
func (k msgServer) CreateComment(goCtx context.Context, msg *types.MsgCreateComment) (*types.MsgCreateCommentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	BlockHeight := CreatedAt() // Invoke method to get the block height of Post
	PostID := CreatedId() // Invoke method to get the Post ID

	// Create variable of type comment
	var comment = types.Comment{
		Creator: msg.Creator,
		Id:      msg.Id,
		Body:    msg.Body,
		Title:   msg.Title,
		BlogID:  msg.BlogID,
		CreatedAt: ctx.BlockHeight(),
	}

	// Check if the Post Exists for which a comment is being created
	if comment.BlogID > PostID {
		return nil, sdkerrors.Wrapf(types.ErrID, "Post Blog Id %d does not exist for which comment with Blog Id %d was made", PostID, comment.BlogID)
	}

	BlockHeight = BlockHeight + 10 // Hardcoded value to 100. This can be changed as per requirement.
	
	// Check if the comment is older than the Post. If more than 100 blocks, then return error.
	if comment.CreatedAt > BlockHeight {
		return nil, sdkerrors.Wrapf(types.ErrCommentOld, "Comment created at %d is older than post created at %d", comment.CreatedAt, BlockHeight)
	} else {
		id := k.AppendComment(ctx, comment)
		return &types.MsgCreateCommentResponse{Id: id}, nil
	}
}
```

When Comment's validity is checked, it throws 2 error messages - `ErrID` and `ErrCommendOld`. You can define the error messaged by adding it to errors definition in x/blog/types/errors.go

```go
//...
var (
	ErrCommentOld = sdkerrors.Register(ModuleName, 1300, "")
)

var (
	ErrID = sdkerrors.Register(ModuleName, 1400, "")
)
```

Notice, 2 methods `CreatedAt` and `CreatedId` are being invoked inside `msg_server_create_comment.go`. They are responsible for returning the Block height and latest ID of Post respectively.

Define `CreatedAt` and `CreatedId` method inside `x/blog/keeper/msg_server_create_post.go`

```go
//...

var BlockHeight int64
var BlockId uint64 

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
  //...
  var post = types.Post{
    //...
    CreatedAt: ctx.BlockHeight(),
  }

    // Define the globally declared `BlockHeight` variable
    BlockHeight = post.CreatedAt

    id := k.AppendPost(ctx, post)
    //...

    // After the post has been appended, fetch the latest id in globally declared `BlockId` variable
    BlockId = id
}

// Define two functions to return BlockId and Height respectively.

func CreatedAt() int64 {
	return BlockHeight
}

func CreatedId() uint64 {
	return BlockId
}
```