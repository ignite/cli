# Importing methods from the Bank keeper

In the previous step you have created the `loan` module with `ignite scaffold
module` using `--dep bank`. This command created a new module and added the
`bank` keeper to the `loan` module, which allows you to add and use bank's
keeper methods in loan's keeper methods.

To see the changes made by `--dep bank`, review the following files:
`x/loan/keeper/keeper.go` and `x/loan/module.go`.

Ignite takes care of adding the `bank` keeper, but you still need to tell the
`loan` module which `bank` methods you will be using. You will be using three
methods: `SendCoins`, `SendCoinsFromAccountToModule`, and
`SendCoinsFromModuleToAccount`. You can do that by adding method signatures to
the `BankKeeper` interface:

```go title="x/loan/types/expected_keepers.go"
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// highlight-start
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// highlight-end
}
```