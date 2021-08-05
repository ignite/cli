---
order: 4
---

# Keeper

The main core of a Cosmos SDK module is a piece called the `Keeper`. It is what handles interaction with the store, has references to other keepers for cross-module interactions, and contains most of the core functionality of a module.

## Buy Name

```go
// x/nameservice/keeper/msg_server_buy_name.go
func (k msgServer) BuyName(goCtx context.Context, msg *types.MsgBuyName) (*types.MsgBuyNameResponse, error) {
  ctx := sdk.UnwrapSDKContext(goCtx)
  // Try getting a name from the store
  whois, isFound := k.GetWhois(ctx, msg.Name)
  // Set the price at which the name has to be bought if it didn't have an owner before
  minPrice := sdk.Coins{sdk.NewInt64Coin("token", 10)}
  // Convert price and bid strings to sdk.Coins
  price, _ := sdk.ParseCoinsNormalized(whois.Price)
  bid, _ := sdk.ParseCoinsNormalized(msg.Bid)
  // Convert owner and buyer address strings to sdk.AccAddress
  owner, _ := sdk.AccAddressFromBech32(whois.Creator)
  buyer, _ := sdk.AccAddressFromBech32(msg.Creator)
  // If a name is found in store
  if isFound {
    // If the current price is higher than the bid
    if price.IsAllGT(bid) {
      // Throw an error
      return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is not high enough")
    }
    // Otherwise (when the bid is higher), send tokens from the buyer to the owner
    k.bankKeeper.SendCoins(ctx, buyer, owner, bid)
  } else { // If the name is not found in the store
    // If the minimum price is higher than the bid
    if minPrice.IsAllGT(bid) {
      // Throw an error
      return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is less than min amount")
    }
    // Otherwise (when the bid is higher), burn tokens from the buyer's account (as a payment for the name)
    k.bankKeeper.SubtractCoins(ctx, buyer, bid)
  }
  // Create an updated whois record
  newWhois := types.Whois{
    Index:   msg.Name,
    Name:    msg.Name,
    Value:   whois.Value,
    Price:   bid.String(),
    Creator: buyer.String(),
  }
  // Write whois information to the store
  k.SetWhois(ctx, newWhois)
  return &types.MsgBuyNameResponse{}, nil
}
```

`BuyName` uses `SendCoins` and `SubtractCoins` methods from the `bank` module. In the beginning when scaffolding a module you used `--dep bank` to specify a dependency between the `nameservice` and `bank` modules. This created an `expected_keepers.go` file with a `BankKeeper` interface. Add `SendCoins` and `SubtractCoins` to be able to use it in the keeper methods of the `nameservice` module.

```go
// x/nameservice/types/expected_keepers.go
type BankKeeper interface {
  SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
  SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## Set Name

```go
// x/nameservice/keeper/msg_server_set_name.go
func (k msgServer) SetName(goCtx context.Context, msg *types.MsgSetName) (*types.MsgSetNameResponse, error) {
  ctx := sdk.UnwrapSDKContext(goCtx)
  // Try getting name information from the store
  whois, _ := k.GetWhois(ctx, msg.Name)
  // If the message sender address doesn't match the name owner, throw an error
  if !(msg.Creator == whois.Creator) {
    return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
  }
  // Otherwise, create an updated whois record
  newWhois := types.Whois{
    Index:   msg.Name,
    Name:    msg.Name,
    Value:   msg.Value,
    Creator: whois.Creator,
    Price:   whois.Price,
  }
  // Write whois information to the store
  k.SetWhois(ctx, newWhois)
  return &types.MsgSetNameResponse{}, nil
}
```

## Delete Name

```go
// x/nameservice/keeper/msg_server_delete_name.go
func (k msgServer) DeleteName(goCtx context.Context, msg *types.MsgDeleteName) (*types.MsgDeleteNameResponse, error) {
  ctx := sdk.UnwrapSDKContext(goCtx)
  // Try getting name information from the store
  whois, isFound := k.GetWhois(ctx, msg.Name)
  // If a name is not found, throw an error
  if !isFound {
    return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name doesn't exist")
  }
  // If the message sender address doesn't match the name owner, throw an error
  if !(whois.Creator == msg.Creator) {
    return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
  }
  // Otherwise, remove the name information from the store
  k.RemoveWhois(ctx, msg.Name)
  return &types.MsgDeleteNameResponse{}, nil
}
```