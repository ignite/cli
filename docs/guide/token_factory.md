---
order: 6
description: Step-by-step guidance to build a token factory module. Mint new native tokens to your blockchain.
---

# Token Factory

In this tutorial you will learn how to build a token factory.

Unique and scarce digital assets are one of the key promises that blockchains deliver. On Ethereum the standard of a ERC20 token has seen a big popularity in the crypto scene.

In this tutorial you will learn how to create a module with a logic of creating and maintaing tokens on the Cosmos SDK.

Be aware, this tutorial is for learning purposes only.

## Module Design

With this module you will be able to create new denoms on your blockchain at will. Learn what [denoms](../kb/denom.md) are and how they are used in the Cosmos Ecosystem.

A denom in this module will always have an owner, who is allowed to issue new tokens, change the denoms name or transfer the ownership to a different account.

The denom has a name `base_denom` and a `ticker` property.
The exponential of the denom is hold in the `precision` property, which defines how many decimal places the denom has.
In order to describe the circulating supply of the token, it has the parameters `maxSupply` and `supply` as current supply. The `canChangeMaxSupply` boolean parameter defines if a token can have an increasing `maxSupply` or not.
Furthermore, the denom has a `description` and a `url` for further information about the token.

The resulting proto definition should look as follows:

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

In this tutorial you will learn how to bring these tokens into existence, you will need functions to issue new token, change the ownership of token and track all tokens in existence.

Get started with scaffolding the blockchain and the module of the `Token Factory`.

## Scaffold a new module

Scaffold a new `tokenfactory` module, use the `--no-module` flag because in the next step you want to add a module with certain dependencies.

```bash
starport scaffold chain github.com/cosmonaut/tokenfactory --no-module
```

Change directory to the new scaffolded blockchain

```bash
cd tokenfactory
```

Next, scaffold a new module with bank and account access

```bash
starport scaffold module tokenfactory --dep account,bank
```

The `--dep` flag is for `dependencies` and Starport wires the dependencies into the right places.

To scaffold the above mentioned data format for a denom in the Token Factory, use a Starport `map`.

```bash
starport scaffold map Denom description:string ticker:string precision:int url:string maxSupply:int supply:int canChangeMaxSupply:bool --signer owner --index denom --module tokenfactory
```

Check the `proto/tokenfactory/denom.proto` file to see the result.

Starport has scaffolded a whole CRUD application.

While the Token Factory is there to create denoms, a once initialized denom should not be deletable. The delete function for denoms is something you should remove in the following steps.

After scaffolding the Denom map, it is a good time to make a first git commit, so you can come back to this step in case something goes wrong with the following steps.

```bash
git add .
git commit -m "Add token factory module and denom map"
```

## Remove Delete Messages

Since a created denom is subsequently handled by the Bank module like any other native denom, it should not be deletable. Hence, remove all references to the Delete action of the scaffolded CRUD type.

In order to remove the functionality to delete token, you will need to remove these functions from the `proto`, the `client`, the `keeper` and the `handler`.

### Proto

In the proto file `proto/tokenfactory/tx.proto` remove the part

```
rpc DeleteDenom(MsgDeleteDenom) returns (MsgDeleteDenomResponse);
```

from the tx.proto service. 
Then, the `MsgDeleteDenom` and `MsgDeleteDenomResponse` messages.

### Client

Navigate to the client in `x/tokenfactory/client`.

First, in the `x/tokenfactory/client/cli/tx_denom_test.go` file, remove the entire `TestDeleteDenom()` function.

In the `x/tokenfactory/client/cli/tx_denom.go` file, remove the entire `CmdDeleteDenom()` function.

In the `x/tokenfactory/client/cli/tx.go` file, remove the line that adds the delete command.

```go
cmd.AddCommand(CmdDeleteDenom())
```

### Keeper

In the keeper there are a few references that you need to take care of in order to remove the delete denom functionlity.

Navigate to the file at `x/tokenfactory/keeper/denom_test.go`.

Remove the `TestDenomRemove` function.

In the `x/tokenfactory/keeper/denom.go` file, remove the entire `RemoveDenom` function.

In the `x/tokenfactory/keeper/msg_server_denom_test.go` file, remove the `TestDenomMsgServerDelete` function.

And finally `x/tokenfactory/keeper/msg_server_denom.go` file, remove the `DeleteDenom` function.

### Types

The types directory defines useful functions and validations that describe the format of the blockchain data. We will have to remove the delete denom functionality from the codec, the message denom test and message denom file.

Start with the codec in `x/tokenfactory/types/codec.go`.

Remove the codec and interface registrations for `MsgDeleteDenom`.

There is a test written in `x/tokenfactory/types/messages_denom_test.go`.

Remove the `TestMsgDeleteDenom_ValidateBasic` function

Lastly in the message denom file at `x/tokenfactory/types/messages_denom.go`.

Remove the entire part referring to `MsgDeleteDenom`.

### Handler

In the handler we have the switch file for all the messages.

Open the `x/tokenfactory/handler.go` file and remove `MsgDeleteDenom` case from `NewHandler` function.

This finishes up removing the delete denom functionality.

In the next Chapter, you will implement the custom logic for the Token Factory.

## Add Application Logic

After removing deletion of denoms now is the time to dedicate to the logic of the Token Factory.

### Proto

Define the format of a new token denom in `proto/tokenfactory/tx.proto`.

Remove `int32 supply = 8;` from MsgCreateDenom and change the field order acccordingly so `canChangeMaxSupply` becomes 8 from 9

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

Remove `string ticker = 4;` , `int32 precision = 5;`, `int32 supply = 8;` from MsgUpdateDenom and change field order for the rest of the fields appropriately

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

`x/tokenfactory/client/cli/tx_denom.go`
Change the number of args to 7 from 8 in `CmdCreateDenom()`and remove references to the supply argument, reordering args accordingly. Also change the usage descriptions.

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

Change the number of args to 5 from 8 in `CmdUpdateDenom()`and remove references to the supply, precision and ticker arguments, reordering args accordingly. Also change the usage descriptions.

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

`x/tokenfactory/client/cli/tx_denom_test.go`
Adjust tests to match changes just made

### Types

When creating new denoms, the denom will not have an initial supply. The supply is only updated when minting tokens based on the amount minted.

Remove `supply` parameter from `NewMsgCreateDenom` in `x/tokenfactory/types/messages_denom.go`.

A few modifications needs to be done to `NewMsgUpdateDenom`in `x/tokenfactory/types/messages_denom.go`.
Remove `ticker`, `precision` and `supply` from the function as these are parameters that cannot be changed.

Before you start implementing the custom logic for creating and updating denoms, add some basic validation to the inputs. You can restrict ticker to between 3 and 10 chars and also we want maxSupply to be greater than 0 in `x/tokenfactory/types/messages_denom.go`

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

Modify MsgUpdateDenom's `ValidateBasic()` function like so:

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

Define the business logic in the keeper. This is the place where you make changes to the database and actually write to the Key/Value Store.

In `x/tokenfactory/keeper/msg_server_denom.go` modify the `CreateDenom()` function like so:

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

Modify the `UpdateDenom()` function.

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

As in previous tutorials, when working closely with other modules, you need to define the functions you want to use from other modules in the `expected_keepers.go` file. Initially, you already scaffolded the module with dependencies on `auth` and `bank` module. Here you can define which functions of these modules can be accessed by your module.

Use the following code in `x/tokenfactory/types/expected_keepers.go`

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

### Scaffold new messages

Everything is in place, scaffold two additional messages to complete the Token Factory's functionality: a `MintAndSendTokens` message and an `UpdateOwner` message

```bash
starport scaffold message MintAndSendTokens denom:string amount:int recipient:string --module tokenfactory --signer owner
```

```bash
starport scaffold message UpdateOwner denom:string newOwner:string --module tokenfactory --signer owner
```

Modify `x/tokenfactory/keeper/msg_server_mint_and_send_tokens.go` like so:

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

Modify `x/tokenfactory/keeper/msg_server_update_owner.go` like so:

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

## Walkthrough

You can now test the Token Factory.

First build and start the chain with 

```bash
starport chain serve
```

Once the chain starts, in a different terminal, run 

```bash
tokenfactoryd tx tokenfactory create-denom ustarport "My denom" STARPORT 6 "someurl" 1000000000 true --from alice
```

and confirm the transaction.

From here you can query the denoms to see your newly created denom.

```bash
tokenfactoryd query tokenfactory list-denom
```

To test the update denom functionality, change the max supply to 2000000000 and the description and URL fields as well as locking down the max supply by running:

```bash
tokenfactoryd tx tokenfactory update-denom ustarport "Starport" "newurl" 2000000000 false --from alice
```

Run the query for denoms again to see the changes taking effect

```bash
tokenfactoryd query tokenfactory list-denom
```

Mint some ustarport tokens and send them to a different address.

```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens ustarport 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

Query the user balance from your `alice` user to see your newly created native denom.

```bash
tokenfactoryd query bank balances cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t
```

Query for the denom list again, to see the updated supply

```bash
tokenfactoryd query tokenfactory list-denom
```

The next function to test is transfer ownership of a denom to a different account.
The command for this that you created is:

```bash
tokenfactoryd tx tokenfactory update-owner ustarport cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

Query for the denom list again, to see the updated ownership

```bash
tokenfactoryd query tokenfactory list-denom
```

To confirm that alice may no longer mint and send tokens and is no longer the owner, test with the command:

```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens ustarport 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

## Congratulations

You have built a Token Factory module. You have learned how to add other modules to the keeper and use their functions. To delete one of the CRUD operations in case they are not needed for your blockchain. You have scaffolded a module and messages.
In the next chapters you will learn more about IBC. One challenge for the enthusiastic reader might be to add IBC functionlity to this module.
