---
order: 5
---

# Mint and Burn Voucher

In this chapter you will learn more about vouchers and how the implementation mints voucher or locks native token from a blockchain.

## Create the SafeBurn Function to Burn Vouchers or Lock Tokens

`SafeBurn` burns tokens if they are IBC vouchers (have an `ibc/` prefix) and locks tokens if they are native to the chain.

Create a new file in the `dex/keeper` called `mint.go`

```go
// x/dex/keeper/mint.go
package keeper

import (
  "fmt"
  sdk "github.com/cosmos/cosmos-sdk/types"
  ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
  "github.com/username/interchange/x/dex/types"
  "strings"
)

// isIBCToken checks if the token came from the IBC module
func isIBCToken(denom string) bool {
  return strings.HasPrefix(denom, "ibc/")
}

func (k Keeper) SafeBurn(ctx sdk.Context, port string, channel string, sender sdk.AccAddress, denom string, amount int32) error {
  if isIBCToken(denom) {
    // Burn the tokens
    if err := k.BurnTokens(ctx, sender, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
      return err
    }
  } else {
    // Lock the tokens
    if err := k.LockTokens(ctx, port, channel, sender, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
      return err
    }
  }
  return nil
}
```

Implement the `BurnTokens` keeper method.

```go
// x/dex/keeper/mint.go
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

```go
// x/dex/keeper/mint.go
import (
  ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

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

`BurnTokens` and `LockTokens` use `SendCoinsFromAccountToModule`, `BurnCoins`, and `SendCoins` keeper methods of the `bank` module. To start using these function from the `dex` module, first add them to the `BankKeeper` interface.

```go
// x/dex/types/expected_keeper.go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
  SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
  BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
  SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## SaveVoucherDenom

`SaveVoucherDenom` saves the voucher denom to be able to convert it back later.

Create a new `denom.go` file in the `keeper` directory.

```go
// x/dex/keeper/denom.go
package keeper

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

Finally, last function we need to implement is `VoucherDenom`. `VoucherDenom` returns the voucher of the denom from the port ID and channel ID.

```go
// x/dex/keeper/denom.go
import (
  sdk "github.com/cosmos/cosmos-sdk/types"
  ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
  "github.com/cosmonaut/interchange/x/dex/types"
)

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

### Implement a OriginalDenom Function

`OriginalDenom` returns back the original denom of the voucher. False is returned if the port ID and channel ID provided are not the origins of the voucher

```go
// x/dex/keeper/denom.go
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
If a token is an IBC token (has an `ibc/` prefix) `SafeMint` mints IBC tokens with `MintTokens`, otherwise, it unlocks native tokens with `UnlockTokens`.
```go
// x/dex/keeper/mint.go
func (k Keeper) SafeMint(ctx sdk.Context, port string, channel string, receiver sdk.AccAddress, denom string, amount int32) error {
	if isIBCToken(denom) {
		// Mint IBC tokens
		if err := k.MintTokens(ctx, receiver, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
			return err
		}
	} else {
		// Unlock native tokens
		if err := k.UnlockTokens(
			ctx,
			port,
			channel,
			receiver,
			sdk.NewCoin(denom, sdk.NewInt(int64(amount))),
		); err != nil {
			return err
		}
	}
	return nil
}
```

#### Implement a `MintTokens` Function

```go
// x/dex/keeper/mint.go
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

Finally, add the function to unlock tokens when they are sent back to the native blockchain.

```go
// x/dex/keeper/mint.go
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

`MintTokens` uses two keeper methods from the `bank` module: `MintCoins` and `SendCoinsFromModuleToAccount`. Import them by adding their signatures to the `BankKeeper` interface.

```go
// x/dex/types/expected_keeper.go
type BankKeeper interface {
  // ...
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```
