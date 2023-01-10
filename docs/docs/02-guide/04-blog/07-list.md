# List posts

In this chapter, you will develop a feature that enables users to retrieve all
of the blog posts stored on your blockchain application. The feature will allow
users to perform a query and receive a paginated response, which means that the
output will be divided into smaller chunks or "pages" of data. This will allow
users to more easily navigate and browse through the list of posts, as they will
be able to view a specific number of posts at a time rather than having to
scroll through a potentially lengthy list all at once.

## List posts

Let's implement the `ListPost` keeper method that will be called when a user
makes a query to the blockchain application, requesting a paginated list of all
the posts stored on chain.

```go title="x/blog/keeper/query_list_post.go"
package keeper

import (
	"context"

	"blog/x/blog/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListPost(goCtx context.Context, req *types.QueryListPostRequest) (*types.QueryListPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var posts []types.Post
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	postStore := prefix.NewStore(store, types.KeyPrefix(types.PostKey))

	pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
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

`ListPost` takes in two arguments: a context object and a request object of type
`QueryListPostRequest`. It returns a response object of type
`QueryListPostResponse` and an error.

The function first checks if the request object is `nil` and returns an error
with a `InvalidArgument` code if it is. It then initializes an empty slice of
`Post` objects and unwraps the context object.

It retrieves a key-value store from the context using the `storeKey` field of
the keeper struct and creates a new store using a prefix of the `PostKey`. It
then calls the `Paginate` function from the `query` package on the store and the
pagination information in the request object. The function passed as an argument
to Paginate iterates over the key-value pairs in the store and unmarshals the
values into `Post` objects, which are then appended to the `posts` slice.

If an error occurs during pagination, the function returns an `Internal error`
with the error message. Otherwise, it returns a `QueryListPostResponse` object
with the list of posts and pagination information.

## Modify `QueryListPostResponse`

Add a `repeated` keyword to return a list of posts and include the option
`[(gogoproto.nullable) = false]` to generate the field without a pointer.

```proto title="proto/blog/blog/query.proto"
message QueryListPostResponse {
  // highlight-next-line
  repeated Post post = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

Run the command to generate Go files from proto:

```
ignite generate proto-go
```