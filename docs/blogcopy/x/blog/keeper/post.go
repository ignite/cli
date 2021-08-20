package keeper

import (
  "encoding/binary"
  "github.com/cosmonaut/blog/x/blog/types"
  "github.com/cosmos/cosmos-sdk/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "strconv"
)


type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
  // Get the context
  ctx := sdk.UnwrapSDKContext(goCtx)
  var post = types.Post{
    Creator: msg.Creator,
    Title:   msg.Title,
    Body:    msg.Body,
  }
  // Add a post to the store and get back the ID
  id := k.AppendPost(ctx, post)
  // Return the ID of the post
  return &types.MsgCreatePostResponse{Id: id}, nil
}

func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
  // Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostCountKey))
  // Convert the PostCountKey to bytes
  byteKey := []byte(types.PostCountKey)
  // Get the value of the count
  bz := store.Get(byteKey)
  // Return zero if the count value is not found (for example, it's the first post)
  if bz == nil {
    return 0
  }
  // Convert the count into a uint64
  count, err := strconv.ParseUint(string(bz), 10, 64)
  // Panic if the count cannot be converted to uint64
  if err != nil {
    panic("cannot decode count")
  }
  return count
}

func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
  // Get the store using storeKey (which is "blog") and PostCountKey (which is "Post-count-")
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostCountKey))
  // Convert the PostCountKey to bytes
  byteKey := []byte(types.PostCountKey)
  // Convert count from uint64 to string and get bytes
  bz := []byte(strconv.FormatUint(count, 10))
  // Set the value of Post-count- to count
  store.Set(byteKey, bz)
}

func (k Keeper) AppendPost(ctx sdk.Context, post types.Post) uint64 {
  // Get the current number of posts in the store
  count := k.GetPostCount(ctx)
  // Assign an ID to the post based on the number of posts in the store
  post.Id = count
  // Get the store
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PostKey))
  // Convert the post ID into bytes
  byteKey := make([]byte, 8)
  binary.BigEndian.PutUint64(byteKey, post.Id)
  // Marshal the post into bytes
  appendedValue := k.cdc.MustMarshalBinaryBare(&post)
  // Insert the post bytes using post ID as a key
  store.Set(byteKey, appendedValue)
  // Update the post count
  k.SetPostCount(ctx, count+1)
  return count
}

var _ types.MsgServer = msgServer{}
