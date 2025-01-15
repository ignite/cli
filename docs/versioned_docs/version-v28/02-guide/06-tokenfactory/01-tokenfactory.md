# Token Factory

## Introduction to Building a Token Factory Module with Ignite CLI

In this tutorial, we will guide you through the process of building a token factory module using the Ignite CLI. This module is a powerful tool for creating native denominations (denoms) on your blockchain, providing you with the capability to issue and manage digital assets natively within your network.

Digital assets, characterized by their uniqueness and scarcity, are fundamental to the value proposition of blockchain technology. A well-known example is the ERC20 standard on Ethereum, which has gained widespread popularity. By learning to create and manage native denoms on your blockchain, you will gain hands-on experience with one of blockchain's key functionalities.

**You will learn how to:**

* Develop a module from scratch.
* Implement a CRUD (Create, Read, Update, Delete) operation while specifically removing the delete functionality to safeguard the integrity of initialized denoms.
* Integrate logic for creating new denoms.
* Engage with various components such as the client, types, keeper, expected keeper, and handlers to effectively implement the Token Factory module.

**Note:** The code provided in this tutorial is tailored for educational purposes. It is not designed for deployment in production environments.

## Understanding the Module Design

The Token Factory module empowers you to create and manage native denoms on your blockchain. In the Cosmos ecosystem and with Ignite CLI, a denom represents the name of a token that is universally usable. To learn more, see [Denom](02-denoms.md).

## What is a Denom?

Denoms are essentially identifiers for tokens on a blockchain, synonymous with terms like 'coin' or 'token'. For an in-depth understanding, refer to the Cosmos SDK's [ADR 024: Coin Metadata](https://docs.cosmos.network/main/build/architecture/adr-024-coin-metadata#context).

A denom in this module always has an owner. An owner is allowed to issue new tokens, change the denoms name, and transfer the ownership to a different account. Learn more about [denoms](02-denoms.md).

In our Token Factory module:

1. Ownership and Control: Each denom is assigned an owner, who has the authority to issue new tokens, rename the denom, and transfer ownership.

2. Properties of a Denom:

    - denom: The unique name of the denom.
    - description: A brief about the denom.
    - ticker: The symbolic representation.
    - precision: Determines the number of decimal places for the denom.
    - url: Provides additional information.
    - maxSupply & supply: Define the total and current circulating supply.
    - canChangeMaxSupply: A boolean indicating if maxSupply can be altered post-issuance.
    - owner: The account holding ownership rights.

3. Proto Definition:

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

4. Core Functionalities:

- Issuing new tokens.
- Transferring ownership of tokens.
- Keeping a ledger of all tokens.

## Chapter 2: Getting Started with Your Token Factory Module

Welcome to the next step in your journey of building a token factory module. In this chapter, we'll walk you through setting up your blockchain and beginning the development of your token factory module.

### Setting up your blockchain

First, we'll scaffold a new blockchain specifically for your token factory. We use the --no-module flag to ensure that we add the token factory module with the required dependencies later. Run the following command in your terminal:

```bash
ignite scaffold chain tokenfactory --no-module
```

This command establishes a new Cosmos SDK blockchain named `tokenfactory` and places it in a directory of the same name. Inside this directory, you'll find a fully functional blockchain ready for further customization.

Now, navigate into your newly created blockchain directory:

```bash
cd tokenfactory
```

### Scaffold Your Token Factory Module

Next, we'll scaffold a new module for your token factory. This module will depend on the Cosmos SDK's [bank](https://docs.cosmos.network/main/build/modules/bank#abstract) and [auth](https://docs.cosmos.network/main/build/modules/auth#abstract) modules, which provide essential functionalities like account access and token management. Use the following command:

```bash
ignite scaffold module tokenfactory --dep account,bank
```

The successful execution of this command will be confirmed with a message indicating that the `tokenfactory` module has been created.

### Defining Denom Data Structure

To manage denoms within your token factory, define their structure using an Ignite map. This will store the data as key-value pairs. Run this command:

```bash
ignite scaffold map Denom description:string ticker:string precision:int url:string maxSupply:int supply:int canChangeMaxSupply:bool --signer owner --index denom --module tokenfactory
```

Review the `proto/tokenfactory/tokenfactory/denom.proto` file to see the scaffolding results, which include modifications to various files indicating successful creation of the denom structure.

### Git Commit

After scaffolding your denom map, it's a good practice to save your progress. Use the following commands to make your first Git commit:

```bash
git add .
git commit -m "Add tokenfactory module and denom map"
```

This saves a snapshot of your project, allowing you to revert back if needed.

## Removing Delete Functionality

In a blockchain context, once a denom is created, it's crucial to ensure it remains immutable and cannot be deleted. This immutability is key to maintaining the integrity and trust in the blockchain. Therefore, we'll remove the delete functionality from the scaffolded CRUD operations. Follow these steps:

**Proto Adjustments**

In `proto/tokenfactory/tokenfactory/tx.proto`, remove the `DeleteDenom` RPC method and the associated message types.

**Client Updates**

Navigate to the client in `x/tokenfactory/client` and make these changes:

- Remove `TestDeleteDenom()` from `tx_denom_test.go`.
- Eliminate `CmdDeleteDenom()` from `tx_denom.go`.
- In `tx.go`, delete the line referencing the delete command.

**Keeper Modifications**

In `denom_test.go`, remove `TestDenomRemove()`.
Delete `RemoveDenom()` from `denom.go`.
Exclude `TestDenomMsgServerDelete()` and `DeleteDenom()` functions from `msg_server_denom_test.go` and `msg_server_denom.go`, respectively.

**Types Directory Changes**

- Update `codec.go` to remove references to `MsgDeleteDenom`.
- Remove `TestMsgDeleteDenom_ValidateBasic()` from `messages_denom_test.go`.
- Eliminate all references to `MsgDeleteDenom()` in `messages_denom.go`.

After making these changes, commit your updates:

```bash
git add .
git commit -m "Remove the delete denom functionality"
```

This concludes the second chapter, setting a solid foundation for your token factory module. In the next chapter, we'll delve into implementing the application logic that will bring your token factory to life.

## Chapter 3: Implementing Core Functionality in Your Token Factory

Having disabled the deletion of denoms, we now turn our attention to the heart of the token factory module: defining the structure of new denoms and implementing their creation and update logic.

**Proto Definition Updates**

Start by defining the structure of a new token denom in `proto/tokenfactory/tokenfactory/tx.proto`.

For `MsgCreateDenom`:

- Remove `int32 supply = 8;` and adjust the field order so `canChangeMaxSupply` becomes the 8th field.

Resulting `MsgCreateDenom` message:

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

For `MsgUpdateDenom`:

- Omit `string ticker = 4;`, `int32 precision = 5;`, and `int32 supply = 8;`, and reorder the remaining fields.

Resulting `MsgUpdateDenom` message:

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

### Client Logic

In the `x/tokenfactory/client/cli/tx_denom.go` file, update the client application logic.

**For `CmdCreateDenom`:**

- Adjust the number of arguments from 8 to 7, removing references to the supply argument, and update the usage descriptions.

**For `CmdUpdateDenom()`:**

- Reduce the number of arguments to 5, excluding `supply`, `precision`, and `ticker`, and modify the usage descriptions accordingly.

Also, update the tests in `x/tokenfactory/client/cli/tx_denom_test.go` to reflect these changes.

### Types Updates

When creating new denoms, they initially have no supply. The supply is determined only when tokens are minted.

In `x/tokenfactory/types/messages_denom.go`:

- Remove the `supply` parameter from `NewMsgCreateDenom`.
- Update `NewMsgUpdateDenom` to exclude unchangeable parameters like `ticker`, `precision`, and `supply`.

Implement basic input validation in `x/tokenfactory/types/messages_denom.go`:

- Ensure the ticker length is between 3 and 10 characters.
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

- Set `maxSupply` to be greater than 0.

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

### Keeper Logic

The keeper is where you define the business logic for manipulating the database and writing to the key-value store.

**In `x/tokenfactory/keeper/msg_server_denom.go`:**

- Update `CreateDenom()` to include logic for creating unique denoms. Modify the error message to point to existing denoms. Set `Supply` to `0`.
- Modify `UpdateDenom()` to verify ownership and manage max supply changes.

```go
func (k msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetDenom(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "Denom to update not found")
	}

	// Checks if the msg owner is the same as the current owner
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

`x/tokenfactory/types/expected_keepers.go` is where you define interactions with other modules. Since your module relies on the `auth` and `bank` modules, specify which of their functions your module can access.

Replace the existing code in `expected_keepers.go` with the updated definitions that interface with `auth` and `bank` modules.

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
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
```

### Committing Your Changes

Regular commits are vital for tracking progress and ensuring a stable rollback point if needed. After implementing these changes, use the following commands to commit:

```bash
git add .
git commit -m "Add token factory create and update logic"
```

To review your progress, use `git log` to see the list of commits, illustrating the journey from initialization to the current state of your module.


## Chapter 4: Expanding Functionality with New Messages

In this chapter, we focus on enhancing the token factory module by adding two critical messages: `MintAndSendTokens` and `UpdateOwner`. These functionalities are key to managing tokens within your blockchain.

### Scaffolding New Messages

**MintAndSendTokens:**

This message allows the creation (minting) of new tokens and their allocation to a specified recipient. The necessary inputs are the denom, the amount to mint, and the recipient's address.

Scaffold this message with:

```bash
ignite scaffold message MintAndSendTokens denom:string amount:int recipient:string --module tokenfactory --signer owner
```

**UpdateOwner:**

This message facilitates the transfer of ownership of a denom. It requires the denom name and the new owner's address.

Scaffold this message with:

```bash
ignite scaffold message UpdateOwner denom:string newOwner:string --module tokenfactory --signer owner
```

### Implementing Logic for New Messages

**In the `MintAndSendTokens` Functionality:**

Located in `x/tokenfactory/keeper/msg_server_mint_and_send_tokens.go`, this function encompasses the logic for minting new tokens. Key steps include:

- Verifying the existence and ownership of the denom.
- Ensuring minting does not exceed the maximum supply.
- Minting the specified amount and sending it to the recipient.

**In the `UpdateOwner` Functionality:**

Found in `x/tokenfactory/keeper/msg_server_update_owner.go`, this function allows transferring ownership of a denom. It involves:

- Checking if the denom exists.
- Ensuring that the request comes from the current owner.
- Updating the owner field in the denom's record.

### Keeper Logic

- For `MintAndSendTokens`, add logic to mint new tokens as per the request parameters. This includes checking for maximum supply limits and transferring the minted tokens to the specified recipient.

```go
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"tokenfactory/x/tokenfactory/types"
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

	// Checks if the msg owner is the same as the current owner
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

- For `UpdateOwner`, implement the logic to update the owner of a denom, ensuring that only the current owner can initiate this change.

```go
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"tokenfactory/x/tokenfactory/types"
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

	// Checks if the msg owner is the same as the current owner
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

### Committing Your Changes

After implementing these new functionalities, it's crucial to save your progress. Use the following commands:

```bash
git add .
git commit -m "Add minting and sending functionality"
```

This commit not only tracks your latest changes but also acts as a checkpoint to which you can revert if needed.

## Chapter 5: Walkthrough and Manual Testing of the Token Factory Module

Congratulations on reaching the final stage! It's time to put your token factory module to the test. This walkthrough will guide you through building, starting your chain, and testing the functionalities you've implemented.

### Building and Starting the Chain

First, build and initiate your blockchain:

```bash
ignite chain serve
```

Keep this terminal running as you proceed with the tests.

### Testing Functionalities

**1. Creating a New Denom:**

- In a new terminal, create a denom named uignite with the command:

```bash
tokenfactoryd tx tokenfactory create-denom uignite "My denom" IGNITE 6 "some/url" 1000000000 true --from alice
```

- Confirm the transaction in your blockchain.

**2. Querying the Denom:**

Check the list of denoms to see your new creation:

```bash
tokenfactoryd query tokenfactory list-denom
```

**3. Updating the Denom:**

- Modify the uignite denom:

```bash
tokenfactoryd tx tokenfactory update-denom uignite "Ignite" "newurl" 2000000000 false --from alice
```

- Query the denoms again to observe the changes: 
```bash
tokenfactoryd query tokenfactory list-denom
```

**4. Minting and Sending Tokens:**

- Mint uignite tokens and send them to a recipient:
```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens uignite 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

- Check the recipient’s balance:
```bash
tokenfactoryd query bank balances cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t
```

- Verify the updated supply in denom list:
```bash
tokenfactoryd query tokenfactory list-denom
```

**5. Transferring Ownership:**

- Transfer the ownership of uignite:
```bash
tokenfactoryd tx tokenfactory update-owner uignite cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

- Confirm the ownership change:
```bash
tokenfactoryd query tokenfactory list-denom
```

**6. Confirming Minting Restrictions:**

- Test minting with alice to ensure restrictions apply:

```bash
tokenfactoryd tx tokenfactory mint-and-send-tokens uignite 1200 cosmos16x46rxvtkmgph6jnkqs80tzlzk6wpy6ftrgh6t --from alice
```

## Congratulations!

You've successfully built and tested a token factory module. This advanced tutorial has equipped you with the skills to:

- Integrate other modules and utilize their functionalities.
- Customize CRUD operations to fit your blockchain's needs.
- Scaffold modules and messages effectively.

## Looking Ahead: IBC Functionality

As you progress, the next learning adventure involves exploring IBC (Inter-Blockchain Communication). If you're up for a challenge, consider adding IBC functionality to your token factory module. This will not only enhance your module's capabilities but also deepen your understanding of the Cosmos ecosystem.
