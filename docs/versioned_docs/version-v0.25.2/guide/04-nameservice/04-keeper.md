---
sidebar_position: 4
description: Define keepers for the nameservice module. 
---

# Keeper

> The main core of a Cosmos SDK module is a piece called the keeper. The keeper handles interactions with the store, has references to other [keepers](https://docs.cosmos.network/main/building-modules/keeper.html) for cross-module interactions, and contains most of the core functionality of a module.

## Define Keepers for the Nameservice Module 

Keepers are module-specific. Keeper is part of the Cosmos SDK that is responsible for writing data to the store. Each module uses its own keeper. 

In this section, define the keepers that are required by the nameservice module:

- Buy name
- Set name
- Delete name

## Buy Name

To define the keeper for the buy name transaction, add this code to the `x/nameservice/keeper/msg_server_buy_name.go` file:

```go
// x/nameservice/keeper/msg_server_buy_name.go

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"nameservice/x/nameservice/types"
)

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
	owner, _ := sdk.AccAddressFromBech32(whois.Owner)
	buyer, _ := sdk.AccAddressFromBech32(msg.Creator)

	// If a name is found in store
	if isFound {
		// If the current price is higher than the bid
		if price.IsAllGT(bid) {
			// Throw an error
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is not high enough")
		}

		// Otherwise (when the bid is higher), send tokens from the buyer to the owner
		err := k.bankKeeper.SendCoins(ctx, buyer, owner, bid)
		if err != nil {
			return nil, err
		}
	} else { // If the name is not found in the store
		// If the minimum price is higher than the bid
		if minPrice.IsAllGT(bid) {
			// Throw an error
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is less than min amount")
		}

		// Otherwise (when the bid is higher), send tokens from the buyer's account to the module's account (as a payment for the name)
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, bid)
		if err != nil {
			return nil, err
		}
	}

	// Create an updated whois record
	newWhois := types.Whois{
		Index: msg.Name,
		Name:  msg.Name,
		Value: whois.Value,
		Price: bid.String(),
		Owner: buyer.String(),
	}

	// Write whois information to the store
	k.SetWhois(ctx, newWhois)
	return &types.MsgBuyNameResponse{}, nil
}
```

When you scaffolded the `nameservice` module you used `--dep bank` to specify a dependency between the `nameservice` and `bank` modules. 

This dependency automatically created an `expected_keepers.go` file with a `BankKeeper` interface. 

The `BuyName` transaction uses `SendCoins` and `SendCoinsFromAccountToModule` methods from the `bank` module. 

Edit the `x/nameservice/types/expected_keepers.go` file to add `SendCoins` and `SendCoinsFromAccountToModule` to be able to use it in the keeper methods of the `nameservice` module.

```go
// x/nameservice/types/expected_keepers.go

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## Set Name

To define the keeper for the set name transaction, add this code to the `x/nameservice/keeper/msg_server_set_name.go` file:

```go
// x/nameservice/keeper/msg_server_set_name.go

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"nameservice/x/nameservice/types"
)

func (k msgServer) SetName(goCtx context.Context, msg *types.MsgSetName) (*types.MsgSetNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Try getting name information from the store
	whois, _ := k.GetWhois(ctx, msg.Name)

	// If the message sender address doesn't match the name owner, throw an error
	if !(msg.Creator == whois.Owner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	// Otherwise, create an updated whois record
	newWhois := types.Whois{
		Index: msg.Name,
		Name:  msg.Name,
		Value: msg.Value,
		Owner: whois.Owner,
		Price: whois.Price,
	}

	// Write whois information to the store
	k.SetWhois(ctx, newWhois)
	return &types.MsgSetNameResponse{}, nil
}
```

## Delete Name

To define the keeper for the delete name transaction, add this code to the `x/nameservice/keeper/msg_server_delete_name.go` file:

```go
// x/nameservice/keeper/msg_server_delete_name.go

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"nameservice/x/nameservice/types"
)

func (k msgServer) DeleteName(goCtx context.Context, msg *types.MsgDeleteName) (*types.MsgDeleteNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Try getting name information from the store
	whois, isFound := k.GetWhois(ctx, msg.Name)

	// If a name is not found, throw an error
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name doesn't exist")
	}

	// If the message sender address doesn't match the name owner, throw an error
	if !(whois.Owner == msg.Creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	// Otherwise, remove the name information from the store
	k.RemoveWhois(ctx, msg.Name)
	return &types.MsgDeleteNameResponse{}, nil
}
```
