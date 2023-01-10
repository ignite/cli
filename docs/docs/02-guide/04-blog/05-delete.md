# Deleting posts

In this chapter, we will be focusing on the process of handling a "delete post"
message.

## Removing posts

```go title="x/blog/keeper/post.go"
func (k Keeper) RemovePost(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	store.Delete(GetPostIDBytes(id))
}
```

`RemovePost` function takes in two arguments: a context object `ctx` and an
unsigned integer `id`. The function removes a post from a key-value store by
deleting the key-value pair associated with the given `id`. The key-value store
is accessed using the `store` variable, which is created by using the `prefix`
package to create a new store using the context's key-value store and a prefix
based on the `PostKey` constant. The `Delete` method is then called on the
`store` object, using the `GetPostIDBytes` function to convert the `id` to a
byte slice as the key to delete.

## Deleting posts

```go title="x/blog/keeper/msg_server_delete_post.go"
package keeper

import (
	"context"
	"fmt"

	"blog/x/blog/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	val, found := k.GetPost(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	k.RemovePost(ctx, msg.Id)
	return &types.MsgDeletePostResponse{}, nil
}
```

`DeletePost` takes in two arguments: a context `goCtx` of type `context.Context`
and a pointer to a message of type `*types.MsgDeletePost`. The function returns
a pointer to a message of type `*types.MsgDeletePostResponse` and an `error`.

Inside the function, the context is unwrapped using the `sdk.UnwrapSDKContext`
function and the value of the post with the ID specified in the message is
retrieved using the `GetPost` function. If the post is not found, an error is
returned using the `sdkerrors.Wrap` function. If the creator of the message does
not match the creator of the post, another error is returned. If both of these
checks pass, the `RemovePost` function is called with the context and the ID of
the post to delete the post. Finally, the function returns a response message
with no data and a `nil` error.

In short, `DeletePost` handles a request to delete a post, ensuring that the
requester is the creator of the post before deleting it.

## Summary

Congratulations on completing the implementation of the `RemovePost` and
`DeletePost` methods in the keeper package! These methods provide functionality
for removing a post from a store and handling a request to delete a post,
respectively.