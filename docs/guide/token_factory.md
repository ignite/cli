---
order: 6
title: "Advanced Module: Token Factory"
description: Build a token factory module. Mint new native tokens to your blockchain.
---

# Token Factory

In this tutorial, you build a token factory module.

A token factory module is used to create native denoms on your blockchain.

Unique and scarce digital assets are key promises that blockchains deliver. For example, on Ethereum, the standard of an ERC20 token has seen a big popularity in the crypto scene. By creating native denom on your blockchain, you learn skills to manage one of the core benefits of blockchain technology. 

**You will learn how to:**

* Create a module
* Remove the delete function of a CRUD operation to prevent deletion of an initialized denom 
* Embed the logic for denom creation
* Work with the client, types, keeper, expected keeper, and handlers to apply the Token Factory application logic

**Note:** The code in this tutorial is written specifically for this learning experience and is intended only for educational purposes. This tutorial code is not intended to be used in production.

## Module Design

The Token Factory module allows you to create native denoms on your blockchain at will. A denom is the name of a token that can be used for all purposes with Starport and in the Cosmos ecosystem. To learn more, see [Denom](../kb/denom.md).

Basically, denoms describe token on a blockchain. Denom and token are also to referred as Coins, see [ADR 024: Coin Metadata](https://docs.cosmos.network/master/architecture/adr-024-coin-metadata.html).

A denom in this module always has an owner. An owner is allowed to issue new tokens, change the denoms name, and transfer the ownership to a different account.

The denom has a name (`denom`) and a property (`ticker`).

- The exponential of the denom is held in the `precision` property that defines how many decimal places the denom has.
- To describe the circulating supply, the token has the parameters `maxSupply` and `supply` as current supply. 
- The `canChangeMaxSupply` boolean parameter defines if a token can have an increasing `maxSupply` after issuance.
- The denom has a `description` and a `url` that contains details about the token.

The resulting proto definition looks like:

```proto
message Denom {
  string denom = 1; 
  string description = 2; 
  string ticker = 3; 
  int32 precision = 4; 
  string url = 5; 
  int32 maxSupply = 6; 
  int32 supply = 7; 
  bool canChangeMaxSupply = 8; 
  string owner = 9;
}
```

To bring these tokens into existence, you require functions to:

* Issue new token
* Change the ownership of token
* Track all tokens in existence

Get started by scaffolding the blockchain and the module of the token factory.

## Scaffold a Blockchain

Scaffold a new `tokenfactory` blockchain and use the `--no-module` flag because you want to add the token factory module with certain dependencies:

```bash
starport scaffold chain github.com/cosmonaut/tokenfactory --no-module
```

Change directory to the new scaffolded blockchain

```bash
cd tokenfactory
```

## Scaffold a Module

Next, scaffold a new module with dependencies on the Cosmos SDK [bank](https://docs.cosmos.network/master/modules/bank/) and [auth](https://docs.cosmos.network/master/modules/auth/) modules:

```bash
starport scaffold module tokenfactory --dep account,bank
```

The `--dep` flag is for `dependencies` so that dependencies on the `auth` (account access) and `bank` modules are wired into the right places.

To scaffold the CRUD operations for a denom in the token factory, use a Starport `map` for data stored as key-value pairs that define the data format:

```bash
starport scaffold map Denom description:string ticker:string precision:int url:string maxSupply:int supply:int canChangeMaxSupply:bool --signer owner --index denom --module tokenfactory
```

Check the `proto/tokenfactory/denom.proto` file to see the result.

Congratulations, you have scaffolded an entire CRUD application. While the purpose of the token factory is to create denoms, you want to prevent an initialized denom from being deleted. You remove the delete function for denoms in the following steps.

After scaffolding the denom map, it is a good time to make a first git commit:

```bash
git add .
git commit -m "Add token factory module and denom map"
```
You can come back to this step in case something goes wrong with the following steps.

## Remove Delete Messages Functionality

Since a created denom is subsequently handled by the `bank` module like any other native denom, the denom must not be deletable. To prevent the denom from being deleted, remove all references to the delete action of the scaffolded CRUD type.

In order to remove the functionality to delete token, you must remove these delete functions from proto, client, keeper, and handler as shown in the next sections:

### Proto

In the proto file `proto/tokenfactory/tx.proto`, remove this code:

```proto
rpc DeleteDenom(MsgDeleteDenom) returns (MsgDeleteDenomResponse);
```

Then, remove the `MsgDeleteDenom` and `MsgDeleteDenomResponse` messages.

### Client

Navigate to the client in `x/tokenfactory/client` and make these changes:

- In the `x/tokenfactory/client/cli/tx_denom_test.go` file, remove the entire `TestDeleteDenom()` function.
- In the `x/tokenfactory/client/cli/tx_denom.go` file, remove the entire `CmdDeleteDenom()` function.
- In the `x/tokenfactory/client/cli/tx.go` file, remove the line that has the delete command:

    ```go
    cmd.AddCommand(CmdDeleteDenom())
    ```

### Keeper

In the `keeper` directory, there are a few files contain the delete denom functionality:

- In the `x/tokenfactory/keeper/denom_test.go` file, remove the `TestDenomRemove()` function.
- In the `x/tokenfactory/keeper/denom.go` file, remove the entire `RemoveDenom()` function.
- In the `x/tokenfactory/keeper/msg_server_denom_test.go` file, remove the `TestDenomMsgServerDelete()` function.
- In the `x/tokenfactory/keeper/msg_server_denom.go` file, remove the `DeleteDenom()` function.

### Types

The `types` directory defines functions and validations that describe the format of the blockchain data. You must remove the delete denom functionality from the codec, the message denom test, and the message denom file.

- Start with the codec in `x/tokenfactory/types/codec.go`, and remove the codec and interface registrations for `MsgDeleteDenom`.
- In the `x/tokenfactory/types/messages_denom_test.go` test, remove the `TestMsgDeleteDenom_ValidateBasic()` function.
- In the message denom `x/tokenfactory/types/messages_denom.go` file, remove the entire part that references `MsgDeleteDenom()`.

### Handler

In the handler, update the switch file for all of the messages.

- Open the `x/tokenfactory/handler.go` file and remove `MsgDeleteDenom` case from `NewHandler` function.

Good job, this step finishes all of the updates required to remove the delete denom functionality.

In the next chapter, you implement the custom logic for the token factory.

This is a good time to make another git commit, before moving to the application logic:

```bash
git add .
git commit -m "Remove the delete denom functionality"
```

## Add Application Logic

After removing deletion of denoms, now is the time that you dedicate to the logic of the token factory.

### Proto

Define the format of a new token denom in `proto/tokenfactory/tx.proto`.

For the `MsgCreateDenom` message:

- Remove `int32 supply = 8;`
- Change the field order accordingly, so `canChangeMaxSupply` changes from 9 to 8

These changes result in the following `MsgCreateDenom` message:

```proto
message MsgCreateDenom {
  string owner = 1;
  string denom = 2;
  string description = 3;
  string ticker = 4;
  int32 precision = 5;
  string url = 6;
  int32 maxSupply = 7;
  bool canChangeMaxSupply = 8;
}
```

For the `MsgUpdateDenom` message:

- Remove `string ticker = 4;`, `int32 precision = 5;`, and `int32 supply = 8;`
- Change the field order for the rest of the fields appropriately

These changes result in the following `MsgUpdateDenom` message:

```proto
message MsgUpdateDenom {
  string owner = 1;
  string denom = 2;
  string description = 3;
  string url = 4;
  int32 maxSupply = 5;
  bool canChangeMaxSupply = 6;
}
```

### Client

Application logic for the client is in the `x/tokenfactory/client/cli/tx_denom.go` file. 

For the  `CmdCreateDenom()` message:

- Change the number of args to 7 from 8 in
- Remove references to the supply argument
- Reorder the args accordingly
- Change the usage descriptions

The changes look like:

```go
func CmdCreateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [denom] [description] [ticker] [precision] [url] [max-supply] [can-change-max-supply]",
		Short: "Create a new Denom",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexDenom := args[0]

			// Get value arguments
			argDescription := args[1]
			argTicker := args[2]
			argPrecision, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}
			argUrl := args[4]
			argMaxSupply, err := cast.ToInt32E(args[5])
			if err != nil {
				return err
			}
			argCanChangeMaxSupply, err := cast.ToBoolE(args[6])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDenom(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argDescription,
				argTicker,
				argPrecision,
				argUrl,
				argMaxSupply,
				argCanChangeMaxSupply,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```

For the `CmdUpdateDenom()` message:

- Change the number of args from 8 to 5
- Remove references to the supply, precision, and ticker arguments
- Reordering args accordingly
- Change the usage descriptions

The changes look like:

```go
func CmdUpdateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom [denom] [description] [url] [max-supply] [can-change-max-supply]",
		Short: "Update a Denom",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexDenom := args[0]

			// Get value arguments
			argDescription := args[1]
			argUrl := args[2]
			argMaxSupply, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}
			argCanChangeMaxSupply, err := cast.ToBoolE(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argDescription,
				argUrl,
				argMaxSupply,
				argCanChangeMaxSupply,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```

In the `x/tokenfactory/client/cli/tx_denom_test.go` file, adjust tests to match the changes you just made.

### Types

When creating new denoms, the denom does not have an initial supply. The supply is updated only when tokens are minted, and is based on the amount that is minted.

In `x/tokenfactory/types/messages_denom.go`:

- Remove the `supply` parameter from `NewMsgCreateDenom` 

- A few modifications are also required in `NewMsgUpdateDenom` to ensure that parameters that cannot be changed, remove `ticker`, `precision`, and `supply` from the function

Before you start implementing the custom logic for creating and updating denoms, add some basic validation to the inputs in `x/tokenfactory/types/messages_denom.go`:

- Restrict ticker to between 3 and 10 chars
- Define `maxSupply` to be greater than 0 

```go
func (msg *MsgCreateDenom) ValidateBasic() error {
    _, err := sdk.AccAddressFromBech32(msg.Owner)
    if err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
    }

    tickerLength := len(msg.Ticker)
    if tickerLength < 3 {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Ticker length must be at least 3 chars long")
    }
    if tickerLength > 10 {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Ticker length must be 10 chars long maximum")
    }
    if msg.MaxSupply == 0 {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Max Supply must be greater than 0")
    }
    
    return nil
}
```

Modify the `ValidateBasic()` function in the `MsgUpdateDenom` message to define error messages for invalid owner address and minimum max supply:

```go
func (msg *MsgUpdateDenom) ValidateBasic() error {
    _, err := sdk.AccAddressFromBech32(msg.Owner)
    if err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
    }
    if msg.MaxSupply == 0 {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Max Supply must be greater than 0")
    }
    return nil
}
```

### Keeper

Define the business logic in the keeper. These changes are where you make changes to the database and actually write to the key-value store.

In `x/tokenfactory/keeper/msg_server_denom.go`:

- Modify the `CreateDenom()` function to define the logic for creating unique new denoms

```go
func (k msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value already exists
    _, isFound := k.GetDenom(
        ctx,
        msg.Denom,
    )
    if isFound {
        return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Denom already exists")
    }
    var denom = types.Denom{
        Owner:              msg.Owner,
        Denom:              msg.Denom,
        Description:        msg.Description,
        Ticker:             msg.Ticker,
        Precision:          msg.Precision,
        Url:                msg.Url,
        MaxSupply:          msg.MaxSupply,
        Supply:             0,
        CanChangeMaxSupply: msg.CanChangeMaxSupply,
    }

    k.SetDenom(
        ctx,
        denom,
    )
    return &types.MsgCreateDenomResponse{}, nil
}
```

- Modify the `UpdateDenom()` function to check whether the owner is correct and check if the owner is allowed to change the max supply:

```go
func (k msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetDenom(
        ctx,
        msg.Denom,
    )
    if !isFound {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the the msg owner is the same as the current owner
    if msg.Owner != valFound.Owner {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    if !valFound.CanChangeMaxSupply && valFound.MaxSupply != msg.MaxSupply {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "cannot change maxsupply")
    }
    if !valFound.CanChangeMaxSupply && msg.CanChangeMaxSupply {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot revert change maxsupply flag")
    }
    var denom = types.Denom{
        Owner:              msg.Owner,
        Denom:              msg.Denom,
        Description:        msg.Description,
        Ticker:             valFound.Ticker,
        Precision:          valFound.Precision,
        Url:                msg.Url,
        MaxSupply:          msg.MaxSupply,
        Supply:             valFound.Supply,
        CanChangeMaxSupply: msg.CanChangeMaxSupply,
    }

    k.SetDenom(ctx, denom)

    return &types.MsgUpdateDenomResponse{}, nil
}
```

### Expected Keepers

As you learned in previous tutorials, when you work closely with other modules, you must define the functions you want to use from the other modules in the `expected_keepers.go` file. 

Initially, you scaffolded the module with dependencies on `auth` and `bank` module. Here you can define which functions of these modules can be accessed by your module.

Use the following code in `x/tokenfactory/types/expected_keepers.go`:

```go
package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
    GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
    GetModuleAddress(name string) sdk.AccAddress
    GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

type BankKeeper interface {
    SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
    MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}
```

Before you scaffold new messages and move on to minting and sending tokens, this is a good time to make another git commit:

```bash
git add .
git commit -m "Add Token Factory Create and Update logic"
```

You can see the work you have done with:

```bash
git log
```

Be proud! The command output shows what you have already accomplished:

```bash
Add token factory Create and Update logic
...
Remove the delete denom functionality
...
Add token factory module and denom map
...
Initialized with Starport
```

The last part is to scaffold and add logic for minting-and-sending and updating token.

### Scaffold Messages

Everything is in place for you to create a new denom. First, you need to scaffold two new messages to complete the token factory functionality: 

- `MintAndSendTokens`
- `UpdateOwner`

The `MintAndSendTokens` message requires this input:

- The `denom`
- An amount to mint `amount` 
- Where the tokens are minted `recipient`

```bash
starport scaffold message MintAndSendTokens denom:string amount:int recipient:string --module tokenfactory --signer owner
```

The `UpdateOwner` message requires this input:

- The `denom`
- The new owner `newOwner`

```bash
starport scaffold message UpdateOwner denom:string newOwner:string --module tokenfactory --signer owner
```

The two new messages are now available in the newly created file `x/tokenfactory/keeper/msg_server_mint_and_send_tokens.go`. 

- Add the details of the logic to mint new tokens:

```go
package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/cosmonaut/tokenfactory/x/tokenfactory/types"
)

func (k msgServer) MintAndSendTokens(goCtx context.Context, msg *types.MsgMintAndSendTokens) (*types.MsgMintAndSendTokensResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetDenom(
        ctx,
        msg.Denom,
    )
    if !isFound {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "denom does not exist")
    }

    // Checks if the the msg owner is the same as the current owner
    if msg.Owner != valFound.Owner {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    if valFound.Supply+msg.Amount > valFound.MaxSupply {
        return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Cannot mint more than Max Supply")
    }
    moduleAcct := k.accountKeeper.GetModuleAddress(types.ModuleName)

    recipientAddress, err := sdk.AccAddressFromBech32(msg.Recipient)
    if err != nil {
        return nil, err
    }

    var mintCoins sdk.Coins

    mintCoins = mintCoins.Add(sdk.NewCoin(msg.Denom, sdk.NewInt(int64(msg.Amount))))
    if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
        return nil, err
    }
    if err := k.bankKeeper.SendCoins(ctx, moduleAcct, recipientAddress, mintCoins); err != nil {
        return nil, err
    }

    var denom = types.Denom{
        Owner:              valFound.Owner,
        Denom:              valFound.Denom,
        Description:        valFound.Description,
        MaxSupply:          valFound.MaxSupply,
        Supply:             valFound.Supply + msg.Amount,
        Precision:          valFound.Precision,
        Ticker:             valFound.Ticker,
        Url:                valFound.Url,
        CanChangeMaxSupply: valFound.CanChangeMaxSupply,
    }

    k.SetDenom(
        ctx,
        denom,
    )
    return &types.MsgMintAndSendTokensResponse{}, nil
}
```

In the `x/tokenfactory/keeper/msg_server_update_owner.go` file, make updates to allow changing an owner of the denom. 

- Update the appropriate fields of the denom
- Add a check for the existence of the denom
- Add a check for and the right owner
- Then you can save the updated demon in the keeper

```go
package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/cosmonaut/tokenfactory/x/tokenfactory/types"
)

func (k msgServer) UpdateOwner(goCtx context.Context, msg *types.MsgUpdateOwner) (*types.MsgUpdateOwnerResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetDenom(
        ctx,
        msg.Denom,
    )
    if !isFound {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "denom does not exist")
    }

    // Checks if the the msg owner is the same as the current owner
    if msg.Owner != valFound.Owner {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    var denom = types.Denom{
        Owner:              msg.NewOwner,
        Denom:              msg.Denom,
        Description:        valFound.Description,
        MaxSupply:          valFound.MaxSupply,
        Supply:             valFound.Supply,
        Precision:          valFound.Precision,
        Ticker:             valFound.Ticker,
        Url:                valFound.Url,
        CanChangeMaxSupply: valFound.CanChangeMaxSupply,
    }

    k.SetDenom(
        ctx,
        denom,
    )

    return &types.MsgUpdateOwnerResponse{}, nil
}
```

After adding minting and sending functionality, this is another a good time to add another git commit:

```bash
git add .
git commit -m "Add minting and sending"
```

## Walk Through

You can now test the token factory module. 

First, build and start the chain:

```bash
starport chain serve
```

Leave this terminal window open.

### Create Denom

After the chain starts, in a different terminal, run:

```bash
tokenfactoryd tx tokenfactory create-denom ustarport "My denom" STARPORT 6 "some/url" 1000000000 true --from alice
```

Confirm the transaction.

### Query Denom

From here, you can query the denoms to see your newly created denom:

```bash
tokenfactoryd query tokenfactory list-denom
```

### Update Denom

To test the update denom functionality:

- Change the max supply to 2,000,000,000 
- Change the description and URL fields
- Lock down the max supply

```bash
tokenfactoryd tx tokenfactory update-denom ustarport "Starport" "newurl" 2000000000 false --from alice
```

Run the query for denoms again to see the changes taking effect:

```bash
tokenfactoryd query tokenfactory list-denom
```

### Mint-And-Send-Tokens

Mint some `ustarport` tokens and send them to a different address:

```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens ustarport 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

### Query Balance

Query the user balance from the `alice` account to see your newly created native denom:

```bash
tokenfactoryd query bank balances cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t
```

Query for the denom list again, to see the updated supply:

```bash
tokenfactoryd query tokenfactory list-denom
```

### Transfer Ownership

The next function to test is transfer ownership of a denom to a different account.

You created the `update-owner` command for this transfer:

```bash
tokenfactoryd tx tokenfactory update-owner ustarport cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

Query for the denom list again to see the updated ownership:

```bash
tokenfactoryd query tokenfactory list-denom
```

### Confirm Minting Restrictions

To confirm that the `alice` is no longer able to mint and send tokens and is no longer the owner, test with the command:

```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens ustarport 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

## Congratulations

Celebrate! You have built a token factory module. 

By completing this advanced tutorial, you successfully performed critical steps to:

- Add other modules to the keeper and use their functions
- Delete one of the CRUD operations when it is not required for your blockchain
- Scaffolded a module and messages
  
In the next chapters, you learn more about IBC. 

If you are an enthusiastic developer, a good challenge is to add IBC functionality to this module.
