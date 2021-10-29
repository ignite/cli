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
