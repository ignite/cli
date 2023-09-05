# Show a post

In this chapter, you will implement a feature in your blogging application that
enables users to retrieve individual blog posts by their unique ID. This ID is
assigned to each blog post when it is created and stored on the blockchain. By
adding this querying functionality, users will be able to easily retrieve
specific blog posts by specifying their ID.

## Show post

Let's implement the `ShowPost` keeper method that will be called when a user
makes a query to the blockchain application, specifying the ID of the desired
post.

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

`ShowPost` is a function for retrieving a single post object from the
blockchain's state. It takes in two arguments: a `context.Context` object called
`goCtx` and a pointer to a `types.QueryShowPostRequest` object called `req`. It
returns a pointer to a `types.QueryShowPostResponse` object and an `error`.

The function first checks if the `req` argument is `nil`. If it is, it returns
an `error` with the code `InvalidArgument` and the message "invalid request"
using the `status.Error` function from the `google.golang.org/grpc/status`
package.

If the `req` argument is not `nil`, the function unwraps the `sdk.Context`
object from the `context.Context` object using the `sdk.UnwrapSDKContext`
function. It then retrieves a post object with the specified `Id` from the
blockchain's state using the `GetPost` function, and checks if the post was
found by checking the value of the `found` boolean variable. If the post was not
found, it returns an error with the type `sdkerrors.ErrKeyNotFound`.

If the post was found, the function creates a new `types.QueryShowPostResponse`
object with the retrieved post object as a field, and returns a pointer to this
object and a `nil` error.

## Modify `QueryShowPostResponse`

Include the option `[(gogoproto.nullable) = false]` in the `post` field in the
`QueryShowPostResponse` message to generate the field without a pointer.

```proto title="proto/blog/blog/query.proto"
message QueryShowPostResponse {
  // highlight-next-line
  Post post = 1 [(gogoproto.nullable) = false];
}
```

Run the command to generate Go files from proto:

```
ignite generate proto-go
```
