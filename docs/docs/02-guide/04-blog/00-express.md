---
description: Explore the essentials of module development while creating a dynamic blogging platform on your blockchain, where users can seamlessly submit and access blog posts, gaining practical experience in decentralized application functionalities.
title: Blog tutorial
---

# Build a Blog on a Blockchain with Ignite CLI

## Introduction

This tutorial guides you through creating a blog application as a Cosmos SDK blockchain using Ignite CLI. You'll learn to set up types, messages, queries, and write logic for creating, reading, updating, and deleting blog posts.

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
This step creates a Post type with title, body (both strings), creator (string), and id (unsigned integer).

## Implementing CRUD operations

**Creating Posts**

1. **Scaffold Create Message**
```bash
ignite scaffold message create-post title body --response id:uint
```
This message allows users to create posts with a title and body.

2. **Append Posts to the Store:**
Create the file `x/blog/keeper/post.go`.

Implement `AppendPost` and following functions in `x/blog/keeper/post.go` to add posts to the store.

```go title="x/blog/keeper/post.go"
package keeper

import (
	"encoding/binary"

	"blog/x/blog/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AppendPost(ctx sdk.Context, post types.Post) uint64 {
	count := k.GetPostCount(ctx)
	post.Id = count
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	appendedValue := k.cdc.MustMarshal(&post)
	store.Set(GetPostIDBytes(post.Id), appendedValue)
	k.SetPostCount(ctx, count+1)
	return count
}

func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
    byteKey := types.KeyPrefix(types.PostCountKey)
    bz := store.Get(byteKey)
    if bz == nil {
        return 0
    }
    return binary.BigEndi
    
    an.Uint64(bz)
}

func GetPostIDBytes(id uint64) []byte {
    bz := make([]byte, 8)
    binary.BigEndian.PutUint64(bz, id)
    return bz
}

func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
    byteKey := types.KeyPrefix(types.PostCountKey)
    bz := make([]byte, 8)
    binary.BigEndian.PutUint64(bz, count)
    store.Set(byteKey, bz)
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

Refine the `UpdatePost` function in `x/blog/keeper/msg_server_update_post.go`

```go title="x/blog/keeper/msg_server_update_post.go"
package keeper

import (
	"context"
	"fmt"

	"blog/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

Add the according logic to `x/blog/keeper/msg_server_delete_post`:

```go title="x/blog/keeper/msg_server_delete_post.go"
package keeper

import (
	"context"
	"fmt"

	"blog/x/blog/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

```bash
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

	"blog/x/blog/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	"blog/x/blog/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

Add a `repeated` keyword to return a list of posts and include the option
`[(gogoproto.nullable) = false]` to generate the field without a pointer.

```proto title="proto/blog/blog/query.proto"
message QueryListPostResponse {
  // highlight-next-line
  repeated Post post = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

**Interacting with the Blog**

1. **Create a Post:**

```bash
blogd tx blog create-post "Hello" "World" --from alice
```

2. **View a Post:**

```bash
blogd q blog show-post 0
```

3. **List All Posts:**

```bash
blogd q blog list-post
````

4. **Update a Post:**

```bash
blogd tx blog update-post "Hello" "Cosmos" 0 --from alice
```

5. **Delete a Post:**

```bash
blogd tx blog delete-post 0 --from alice
```

**Summary**

Congratulations on completing the Blog tutorial! You've successfully built a functional blockchain application using Ignite and Cosmos SDK. This tutorial equipped you with the skills to generate code for key blockchain operations and implement business-specific logic in a blockchain context. Continue developing your skills and expanding your blockchain applications with the next tutorials.
