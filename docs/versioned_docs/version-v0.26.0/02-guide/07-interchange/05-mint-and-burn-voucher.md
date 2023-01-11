---
order: 5
description: Mint vouchers and lock and unlock native token from a blockchain.
---

# Mint and Burn Vouchers

In this chapter, you learn about vouchers. The `dex` module implementation mints vouchers and locks and unlocks native
token from a blockchain.

There is a lot to learn from this `dex` module implementation:

- You work with the `bank` keeper and use several methods it offers.
- You interact with another module and use the module account to lock tokens.

This implementation can teach you how to use various interactions with module accounts or minting, locking or burning
tokens.

## Create the SafeBurn Function to Burn Vouchers or Lock Tokens

The `SafeBurn` function burns tokens if they are IBC vouchers (have an `ibc/` prefix) and locks tokens if they are
native to the chain.

Create a new `x/dex/keeper/mint.go` file:

```go
// x/dex/keeper/mint.go

package keeper

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"interchange/x/dex/types"
)

// isIBCToken checks if the token came from the IBC module
// Each IBC token starts with an ibc/ denom, the check is rather simple
func isIBCToken(denom string) bool {
	return strings.HasPrefix(denom, "ibc/")
}

func (k Keeper) SafeBurn(ctx sdk.Context, port string, channel string, sender sdk.AccAddress, denom string, amount int32) error {
	if isIBCToken(denom) {
		// Burn the tokens
		if err := k.BurnTokens(ctx, sender, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
			return err
		}
	} else {
		// Lock the tokens
		if err := k.LockTokens(ctx, port, channel, sender, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
			return err
		}
	}

	return nil
}
```

If the token comes from another blockchain as an IBC token, the burning method actually burns those IBC tokens on one
chain and unlocks them on the other chain. The native token are locked away.

Now, implement the `BurnTokens` keeper method as used in the previous function. The `bankKeeper` has a useful function
for this:

```go
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) BurnTokens(ctx sdk.Context, sender sdk.AccAddress, tokens sdk.Coin) error {
	// transfer the coins to the module account and burn them
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(tokens)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		// NOTE: should not happen as the module account was
		// retrieved on the step above and it has enough balance
		// to burn.
		panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
	}

	return nil
}
```

Implement the `LockTokens` keeper method.

To lock token from a native chain, you can send the native token to the Escrow Address:

```go
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) LockTokens(ctx sdk.Context, sourcePort string, sourceChannel string, sender sdk.AccAddress, tokens sdk.Coin) error {
	// create the escrow address for the tokens
	escrowAddress := ibctransfertypes.GetEscrowAddress(sourcePort, sourceChannel)

	// escrow source tokens. It fails if balance insufficient
	if err := k.bankKeeper.SendCoins(
		ctx, sender, escrowAddress, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}

	return nil
}
```

`BurnTokens` and `LockTokens` use `SendCoinsFromAccountToModule`, `BurnCoins`, and `SendCoins` keeper methods of the
`bank` module.

To start using these function from the `dex` module, first add them to the `BankKeeper` interface in the
`x/dex/types/expected_keepers.go` file.

```go
// x/dex/types/expected_keepers.go

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	//...
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## SaveVoucherDenom

The `SaveVoucherDenom` function saves the voucher denom to be able to convert it back later.

Create a new `x/dex/keeper/denom.go` file:

```go
// x/dex/keeper/denom.go

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"interchange/x/dex/types"
)

func (k Keeper) SaveVoucherDenom(ctx sdk.Context, port string, channel string, denom string) {
	voucher := VoucherDenom(port, channel, denom)

	// Store the origin denom
	_, saved := k.GetDenomTrace(ctx, voucher)
	if !saved {
		k.SetDenomTrace(ctx, types.DenomTrace{
			Index:   voucher,
			Port:    port,
			Channel: channel,
			Origin:  denom,
		})
	}
}
```

Finally, the last function to implement is the `VoucherDenom` function that returns the voucher of the denom from the
port ID and channel ID:

```go
// x/dex/keeper/denom.go

package keeper

// ...

func VoucherDenom(port string, channel string, denom string) string {
	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := ibctransfertypes.GetDenomPrefix(port, channel)

	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom

	// construct the denomination trace from the full raw denomination
	denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)
	voucher := denomTrace.IBCDenom()
	return voucher[:16]
}
```

### Implement an OriginalDenom Function

The `OriginalDenom` function returns back the original denom of the voucher.

False is returned if the port ID and channel ID provided are not the origins of the voucher:

```go
// x/dex/keeper/denom.go

package keeper

// ...

func (k Keeper) OriginalDenom(ctx sdk.Context, port string, channel string, voucher string) (string, bool) {
	trace, exist := k.GetDenomTrace(ctx, voucher)
	if exist {
		// Check if original port and channel
		if trace.Port == port && trace.Channel == channel {
			return trace.Origin, true
		}
	}

	// Not the original chain
	return "", false
}
```

### Implement a SafeMint Function

If a token is an IBC token (has an `ibc/` prefix), the  `SafeMint` function mints IBC token with `MintTokens`.
Otherwise, it unlocks native token with `UnlockTokens`.

Go back to the `x/dex/keeper/mint.go` file and add the following code:

```go
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) SafeMint(ctx sdk.Context, port string, channel string, receiver sdk.AccAddress, denom string, amount int32) error {
	if isIBCToken(denom) {
		// Mint IBC tokens
		if err := k.MintTokens(ctx, receiver, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
			return err
		}
	} else {
		// Unlock native tokens
		if err := k.UnlockTokens(
			ctx,
			port,
			channel,
			receiver,
			sdk.NewCoin(denom, sdkmath.NewInt(int64(amount))),
		); err != nil {
			return err
		}
	}

	return nil
}
```

#### Implement a `MintTokens` Function

You can use the `bankKeeper` function again to MintCoins. These token will then be sent to the receiver account:

```go
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) MintTokens(ctx sdk.Context, receiver sdk.AccAddress, tokens sdk.Coin) error {
	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}

	// send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiver, sdk.NewCoins(tokens),
	); err != nil {
		panic(fmt.Sprintf("unable to send coins from module to account despite previously minting coins to module account: %v", err))
	}

	return nil
}
```

Finally, add the function to unlock token after they are sent back to the native blockchain:

```go
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) UnlockTokens(ctx sdk.Context, sourcePort string, sourceChannel string, receiver sdk.AccAddress, tokens sdk.Coin) error {
	// create the escrow address for the tokens
	escrowAddress := ibctransfertypes.GetEscrowAddress(sourcePort, sourceChannel)

	// escrow source tokens. It fails if balance insufficient
	if err := k.bankKeeper.SendCoins(
		ctx, escrowAddress, receiver, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}

	return nil
}
```

The `MintTokens` function uses two keeper methods from the `bank` module: `MintCoins` and `SendCoinsFromModuleToAccount`
.
To import these methods, add their signatures to the `BankKeeper` interface in the `x/dex/types/expected_keepers.go`
file:

```go
// x/dex/types/expected_keepers.go

package types

// ...

type BankKeeper interface {
	// ...
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## Summary

You finished the mint and burn voucher logic.

It is a good time to make another git commit to save the state of your work:

```bash
git add .
git commit -m "Add Mint and Burn Voucher"
```

In the next chapter, you look into creating sell orders.
