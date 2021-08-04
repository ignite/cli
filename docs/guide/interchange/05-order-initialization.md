---
order: 6
---

# Initialize the Orders

In this chapter you will initialize the bank keeper into `ibcdex` and update the module in the `app.go` file. 
You will make sure that there can be no new native tokens minted on a blockchain, which will secure the exchange module from any counterparty misbehaviour.

## Add the bank keeper to the module

Add the bank keeper reference to the `ibcdex` module. First, in the `types` directory.

```go
// x/ibcdex/types/expected_keeper.go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
```

Next, in the `keeper` directory, specify the bank so that you can access it in your module. 

```go
// x/ibcdex/keeper/keeper.go
type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		channelKeeper types.ChannelKeeper
		portKeeper    types.PortKeeper
		scopedKeeper  capabilitykeeper.ScopedKeeper
		bankKeeper    types.BankKeeper  // <--
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	bankKeeper types.BankKeeper,  // <--
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
		bankKeeper:    bankKeeper,  // <--
	}
}
```

Lastly, the `app.go` file that describes which modules are used in the blockchain application, add the bank keeper to the `ibcdexKeeper`

```go
// app/app.go
app.ibcdexKeeper = *ibcdexkeeper.NewKeeper(
	appCodec,
	keys[ibcdextypes.StoreKey],
	keys[ibcdextypes.MemStoreKey],
	app.IBCKeeper.ChannelKeeper,
	&app.IBCKeeper.PortKeeper,
	scopedIbcdexKeeper,
	app.BankKeeper,  // <--
)
```

The app will mint and burn token using the bank account. The use this feature, the module must have a _module account_. To enable the _module account_ you declare this permission in the _module account permissions_ structure of the auth module. This is done in `app.go` file.

```go
// app/app.go
maccPerms = map[string][]string{
    authtypes.FeeCollectorName:     nil,
    distrtypes.ModuleName:          nil,
    minttypes.ModuleName:           {authtypes.Minter},
    stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
    stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
    govtypes.ModuleName:            {authtypes.Burner},
    ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
    ibcdextypes.ModuleName:         {authtypes.Minter, authtypes.Burner}, // <--
}
```

## Ensure no native token is minted

Ensure that no new native tokens are minted in case the module is connected to a malicious or faulty counterparty. Rename this functionality into lock and unlock and reuse functions from `ibc-transfer` module with an escrow address to ensure no new token creation.

For non-native tokens, use the _module accounts_ functionality.

To abstract the process of deciding if a token must be locked or burned, create the `SafeBurn` and `SafeMint` functions. Those functions will check if the denom is native to the chain or a voucher.

Consider that any denom starting with `ibc/` is a voucher.

Define utilitary functions.

```go
// x/ibcdex/keeper/mint.go
package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/username/interchange/x/ibcdex/types"
	"strings"
)

// isIBCToken checks if the token came from the IBC module
func isIBCToken(denom string) bool {
	return strings.HasPrefix(denom, "ibc/")
}

func (k Keeper) SafeBurn(
	ctx sdk.Context,
	port string,
	channel string,
	sender sdk.AccAddress,
	denom string,
	amount int32,
) error {
	if isIBCToken(denom) {
		// Burn the tokens
		if err := k.BurnTokens(
			ctx, sender,
			sdk.NewCoin(denom, sdk.NewInt(int64(amount))),
		); err != nil {
			return err
		}
	} else {
		// Lock the token to send
		if err := k.LockTokens(
			ctx,
			port,
			channel,
			sender,
			sdk.NewCoin(denom, sdk.NewInt(int64(amount))),
		); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) SafeMint(
	ctx sdk.Context,
	port string,
	channel string,
	receiver sdk.AccAddress,
	denom string,
	amount int32,
) error {
	if isIBCToken(denom) {
		// Mint IBC tokens
		if err := k.MintTokens(
			ctx,
			receiver,
			sdk.NewCoin(denom, sdk.NewInt(int64(amount))),
		); err != nil {
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

func (k Keeper) BurnTokens(
	ctx sdk.Context,
	sender sdk.AccAddress,
	tokens sdk.Coin,
) error {
	// transfer the coins to the module account and burn them
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		// NOTE: should not happen as the module account was
		// retrieved on the step above and it has enough balace
		// to burn.
		panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
	}

	return nil
}

func (k Keeper) MintTokens(
	ctx sdk.Context,
	receiver sdk.AccAddress,
	tokens sdk.Coin,
) error {
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

func (k Keeper) LockTokens(
	ctx sdk.Context,
	sourcePort string,
	sourceChannel string,
	sender sdk.AccAddress,
	tokens sdk.Coin,
) error {
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

func (k Keeper) UnlockTokens(
	ctx sdk.Context,
	sourcePort string,
	sourceChannel string,
	receiver sdk.AccAddress,
	tokens sdk.Coin,
) error {
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

## Declare the denom

As defined in the [App Design](./01-app-design.md), token sent with IBC will create `voucher` token.
The following rules apply:

- An external token received from a chain has a unique denom, referred to as `voucher`
- When a token sent to a chain is received back, the chain can resolve the voucher and convert it back to the original token denomination

Vouchers are hashes, therefore you must store which original denomination is related to a voucher, do this with an `indexed` type transaction.

For a voucher, store the following information: 
- the source port ID 
- source channel ID 
- the original denom

Define the methods to:

- determine voucher from a port, channel and denom
- save all voucher to resolve into original denoms
- resolve vouchers

```go
// x/ibcdex/keeper/denom.go
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/username/interchange/x/ibcdex/types"
)

// VoucherDenom returns the voucher of the denom from the port ID and channel ID
// and saves the origin in the store
func VoucherDenom(port string, channel string, denom string) string {
	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := ibctransfertypes.GetDenomPrefix(port, channel)

	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom

	// construct the denomination trace from the full raw denomination
	denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)

	voucher := denomTrace.IBCDenom()

	return voucher
}

// SaveVoucherDenom saves the voucher denom to be able to convert it back later
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

// OriginalDenom returns back the original denom of the voucher
// False is returned if the port ID and channel ID provided are not the origins of the voucher
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

In the next chapter, you will learn how to create the IBC packets for buy and sell orders. How to use packets that are received, acknowledged or timed out.