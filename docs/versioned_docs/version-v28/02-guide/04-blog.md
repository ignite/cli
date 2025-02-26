---
description: Explore the essentials of module development while creating a dynamic blogging platform on your blockchain, where users can seamlessly submit and access blog posts, gaining practical experience in decentralized application functionalities.
title: Blog tutorial
---

# Build a Blog on a Blockchain with Ignite CLI

## Introduction

This tutorial guides you through creating a blog application as a Cosmos SDK blockchain using Ignite CLI. You'll learn how to set up types, messages, queries, and write logic for creating, reading, updating, and deleting blog posts.

## Creating the Blog Blockchain

1. **Initialize the Blockchain:**

```bash
ignite scaffold chain blog
cd blog
```

2. **Define the Post Type:**

```bash
ignite scaffold type post title body creator id:uint
```
This step creates a Post type with title (string), body (string), creator (string), and id (unsigned integer) fields.

## Implementing CRUD operations

**Creating Posts**

1. **Scaffold Create Message**

```bash
ignite scaffold message create-post title body --response id:uint
```

This message allows users to create posts with a title and body.

2. **Append Posts to the Store:**

Create the file `x/blog/keeper/post.go`.

Implement `AppendPost` and the following functions in `x/blog/keeper/post.go` to add posts to the store.

```go title="x/blog/keeper/post.go"
package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"blog/x/blog/types"
)

func (k Keeper) AppendPost(ctx sdk.Context, post types.Post) uint64 {
	count := k.GetPostCount(ctx)
	post.Id = count
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PostKey))
	appendedValue := k.cdc.MustMarshal(&post)
	store.Set(GetPostIDBytes(post.Id), appendedValue)
	k.SetPostCount(ctx, count+1)
	return count
}

func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PostCountKey)
	bz := store.Get(byteKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func GetPostIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PostCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

func (k Keeper) GetPost(ctx sdk.Context, id uint64) (val types.Post, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PostKey))
	b := store.Get(GetPostIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
```

3. **Add Post key prefix:**

Add the `PostKey` and `PostCountKey` functions to the `x/blog/types/keys.go` file:

```go title="x/blog/types/keys.go"
	// PostKey is used to uniquely identify posts within the system.
	// It will be used as the beginning of the key for each post, followed by their unique ID
	PostKey = "Post/value/"

	// This key will be used to keep track of the ID of the latest post added to the store.
	PostCountKey = "Post/count/"
```

4. **Update Create Post:**

Update the `x/blog/keeper/msg_server_create_post.go` file with the `CreatePost` function:

```go title="x/blog/keeper/msg_server_create_post.go"
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"blog/x/blog/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var post = types.Post{
		Creator: msg.Creator,
		Title:   msg.Title,
		Body:    msg.Body,
	}
	id := k.AppendPost(
		ctx,
		post,
	)
	return &types.MsgCreatePostResponse{
		Id: id,
	}, nil
}
```

**Updating Posts**

1. **Scaffold Update Message:**

```bash
ignite scaffold message update-post title body id:uint
```

This command allows for updating existing posts specified by their ID.

2. **Update Logic**

Implement `SetPost` in `x/blog/keeper/post.go` for updating posts in the store.

```go title="x/blog/keeper/post.go"
func (k Keeper) SetPost(ctx sdk.Context, post types.Post) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PostKey))
	b := k.cdc.MustMarshal(&post)
	store.Set(GetPostIDBytes(post.Id), b)
}
```

Refine the `UpdatePost` function in `x/blog/keeper/msg_server_update_post.go`.

```go title="x/blog/keeper/msg_server_update_post.go"
package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blog/x/blog/types"
)

func (k msgServer) UpdatePost(goCtx context.Context, msg *types.MsgUpdatePost) (*types.MsgUpdatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var post = types.Post{
		Creator: msg.Creator,
		Id:      msg.Id,
		Title:   msg.Title,
		Body:    msg.Body,
	}
	val, found := k.GetPost(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	k.SetPost(ctx, post)
	return &types.MsgUpdatePostResponse{}, nil
}
```

**Deleting Posts**

1. **Scaffold Delete Message:**

```bash
ignite scaffold message delete-post id:uint
```

This command enables the deletion of posts by their ID.

2. **Delete Logic:**
 
Implement RemovePost in `x/blog/keeper/post.go` to delete posts from the store.

```go title="x/blog/keeper/post.go"
func (k Keeper) RemovePost(ctx sdk.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PostKey))
	store.Delete(GetPostIDBytes(id))
}
```

Add the according logic to `x/blog/keeper/msg_server_delete_post`.

```go title="x/blog/keeper/msg_server_delete_post.go"
package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"blog/x/blog/types"
)

func (k msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	val, found := k.GetPost(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	k.RemovePost(ctx, msg.Id)
	return &types.MsgDeletePostResponse{}, nil
}
```

**Reading Posts**

1. **Scaffold Query Messages:**

```bash title="proto/blog/blog/query.proto"
ignite scaffold query show-post id:uint --response post:Post
ignite scaffold query list-post --response post:Post --paginated
```

These queries allow for retrieving a single post by ID and listing all posts with pagination.

2. **Query Implementation:**

Implement `ShowPost` in `x/blog/keeper/query_show_post.go`.

```go title="x/blog/keeper/query_show_post.go"
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog/x/blog/types"
)

func (k Keeper) ShowPost(goCtx context.Context, req *types.QueryShowPostRequest) (*types.QueryShowPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	post, found := k.GetPost(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryShowPostResponse{Post: &post}, nil
}
```

Implement `ListPost` in `x/blog/keeper/query_list_post.go`.

```go title="x/blog/keeper/query_list_post.go"
package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog/x/blog/types"
)

func (k Keeper) ListPost(ctx context.Context, req *types.QueryListPostRequest) (*types.QueryListPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PostKey))

	var posts []types.Post
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var post types.Post
		if err := k.cdc.Unmarshal(value, &post); err != nil {
			return err
		}

		posts = append(posts, post)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListPostResponse{Post: posts, Pagination: pageRes}, nil
}
```

3. **Proto Implementation:** 

Add a `repeated` keyword to return a list of posts in `QueryListPostResponse` and include the option
`[(gogoproto.nullable) = false]` in `QueryShowPostResponse` and `QueryListPostResponse` to generate the field without a pointer.

```proto title="proto/blog/blog/query.proto"
message QueryShowPostResponse {
  Post post = 1 [(gogoproto.nullable) = false];
}

message QueryListPostResponse {
  // highlight-next-line
  repeated Post post = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

Build the blockchain:

```
ignite chain build
```

Start the blockchain:

```
ignite chain serve
```

**Interacting with the Blog**

1. **Create a Post:**

```bash
blogd tx blog create-post hello world --from alice --chain-id blog
```

2. **View a Post:**

```bash
blogd q blog show-post 0
```

3. **List All Posts:**

```bash
blogd q blog list-post
```

4. **Update a Post:**

```bash
blogd tx blog update-post "Hello" "Cosmos" 0 --from alice --chain-id blog
```

5. **Delete a Post:**

```bash
blogd tx blog delete-post 0 --from alice  --chain-id blog
```

## Conclusion

Congratulations on completing the Blog tutorial! You've successfully built a functional blockchain application using Ignite and Cosmos SDK. This tutorial equipped you with the skills to generate code for key blockchain operations and implement business-specific logic in a blockchain context. Continue developing your skills and expanding your blockchain applications with the next tutorials.
