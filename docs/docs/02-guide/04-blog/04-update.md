# Updating posts

In this chapter, we will be focusing on the process of handling an "update post"
message.

To update a post, you need to retrieve the specific post from the store using
the "Get" operation, modify the values, and then write the updated post back to
the store using the "Set" operation.

Let's first implement a getter and a setter logic.

## Getting posts

Implement the `GetPost` keeper method in `post.go`:

```go title="x/blog/keeper/post.go"
func (k Keeper) GetPost(ctx sdk.Context, id uint64) (val types.Post, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	b := store.Get(GetPostIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
```

`GetPost` takes in two arguments: a context `ctx` and an `id` of type `uint64`
representing the ID of the post to be retrieved. It returns a `types.Post`
struct containing the values of the post, and a boolean value indicating whether
the post was found in the database.

The function first creates a `store` using the `prefix.NewStore` method, passing
in the key-value store from the context and the `types.KeyPrefix` function
applied to the `types.PostKey` constant as arguments. It then attempts to
retrieve the post from the store using the `store.Get` method, passing in the ID
of the post as a byte slice. If the post is not found in the store, it returns
an empty `types.Post` struct and a boolean value of false.

If the post is found in the store, the function unmarshals the retrieved byte
slice into a `types.Post` struct using the `cdc.MustUnmarshal` method, passing
in a pointer to the val variable as an argument. It then returns the val struct
and a boolean value of true to indicate that the post was found in the database.

## Setting posts

Implement the `SetPost` keeper method in `post.go`:

```go title="x/blog/keeper/post.go"
func (k Keeper) SetPost(ctx sdk.Context, post types.Post) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	b := k.cdc.MustMarshal(&post)
	store.Set(GetPostIDBytes(post.Id), b)
}
```

`SetPost` takes in two arguments: a context `ctx` and a `types.Post` struct
containing the updated values for the post. The function does not return
anything.

The function first creates a store using the `prefix.NewStore` method, passing
in the key-value store from the context and the `types.KeyPrefix` function
applied to the `types.PostKey` constant as arguments. It then marshals the
updated post struct into a byte slice using the `cdc.MustMarshal` method,
passing in a pointer to the post struct as an argument. Finally, it updates the
post in the store using the `store.Set` method, passing in the ID of the post as
a byte slice and the marshaled post struct as arguments.


## Update posts

```go title="x/blog/keeper/msg_server_update_post.go"
package keeper

import (
	"context"
	"fmt"

	"blog/x/blog/types"

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
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	k.SetPost(ctx, post)
	return &types.MsgUpdatePostResponse{}, nil
}
```

`UpdatePost` takes in a context and a message `MsgUpdatePost` as input, and
returns a response `MsgUpdatePostResponse` and an `error`. The function first
retrieves the current values of the post from the database using the provided
`msg.Id`, and checks if the post exists and if the `msg.Creator` is the same as
the current owner of the post. If either of these checks fail, it returns an
error. If both checks pass, it updates the post in the database with the new
values provided in `msg`, and returns a response without an error.

## Summary

Well done! You have successfully implemented a number of important methods for
managing posts within a store.

The `GetPost` method allows you to retrieve a specific post from the store based
on its unique identification number, or post ID. This can be useful for
displaying a specific post to a user, or for updating it.

The `SetPost` method enables you to update an existing post in the store. This
can be useful for correcting mistakes or updating the content of a post as new
information becomes available.

Finally, you implemented the `UpdatePost` method, which is called whenever the
blockchain processes a message requesting an update to a post.