---
order: 7
---

# Keeper

## Create Scavenge

Create scavenge method should do the following:

* Check that a scavenge with a given solution hash doesn't exist
* Send tokens from scavenge creator account to a module account
* Write the scavenge to the store

```go
// x/scavenge/keeper/msg_server_submit_scavenge.go
func (k msgServer) SubmitScavenge(goCtx context.Context, msg *types.MsgSubmitScavenge) (*types.MsgSubmitScavengeResponse, error) {
  // get context that contains information about the environment, such as block height
	ctx := sdk.UnwrapSDKContext(goCtx)
  // create a new scavenge from the data in the MsgSubmitScavenge message
	var scavenge = types.Scavenge{
		Index:        msg.SolutionHash,
		Creator:      msg.Creator,
		Description:  msg.Description,
		SolutionHash: msg.SolutionHash,
		Reward:       msg.Reward,
	}
  // try getting a scavenge from the store using the solution hash as the key
	_, isFound := k.GetScavenge(ctx, scavenge.SolutionHash)
  // return an error if a scavenge already exists in the store
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists")
	}
  // get address of the Scavenge module account
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
  // convert the message creator address from a string into sdk.AccAddress
	scavenger, err := sdk.AccAddressFromBech32(scavenge.Creator)
	if err != nil {
		panic(err)
	}
  // convert tokens from string into sdk.Coins
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}
  // send tokens from the scavenge creator to the module account
	sdkError := k.bankKeeper.SendCoins(ctx, scavenger, moduleAcct, reward)
	if sdkError != nil {
		return nil, sdkError
	}
  // write the scavenge to the store
	k.SetScavenge(ctx, scavenge)
	return &types.MsgSubmitScavengeResponse{}, nil
}
```

Notice the use of `moduleAcct`. This account is not controlled by a public key pair, but is a reference to an account that is owned by this actual module. It is used to hold the bounty reward that is attached to a scavenge until that scavenge has been solved, at which point the bounty is paid to the account who solved the scavenge.

`SubmitScavenge` uses `SendCoins` method from the `bank` module. In the beginning when scaffolding a module you used `--dep bank` to specify a dependency between the `scavenge` and `bank` modules. This created an `expected_keepers.go` file with a `BankKeeper` interface. Add `SendCoins` to be able to use it in the keeper methods of the `scavenge` module.

```go
// x/scavenge/types/expected_keepers.go
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## Commit Solution

Commit solution method should do the following:

* Check that commit with a given hash doesn't exist in the store
* Write a new commit to the store

```go
// x/scavenge/keeper/msg_server_commit_solution.go
func (k msgServer) CommitSolution(goCtx context.Context, msg *types.MsgCommitSolution) (*types.MsgCommitSolutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
  // create a new commit from the information in the MsgCommitSolution message
	var commit = types.Commit{
		Index:                 msg.SolutionScavengerHash,
		Creator:               msg.Creator,
		SolutionHash:          msg.SolutionHash,
		SolutionScavengerHash: msg.SolutionScavengerHash,
	}
  // try getting a commit from the store using the solution+scavenger hash as the key
	_, isFound := k.GetCommit(ctx, commit.SolutionScavengerHash)
	// return an error if a commit already exists in the store
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash already exists")
	}
  // write commit to the store
	k.SetCommit(ctx, commit)
	return &types.MsgCommitSolutionResponse{}, nil
}
```

## Reveal Solution

Reveal solution method should do the following:

* Check that a commit with a given hash exists in the store
* Check that a scavenge with a given solution hash exists in the store
* Check that the scavenge hasn't already been solved
* Send tokens from the module account to the account that revealed the correct anwer
* Write the updated scavenge to the store

```go
// x/scavenge/keeper/msg_server_reveal_solution.go
func (k msgServer) RevealSolution(goCtx context.Context, msg *types.MsgRevealSolution) (*types.MsgRevealSolutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
  // concatenate a solution and a scavenger address and convert it to bytes
	var solutionScavengerBytes = []byte(msg.Solution + msg.Creator)
  // find the hash of solution and address
	var solutionScavengerHash = sha256.Sum256(solutionScavengerBytes)
  // convert the hash to a string
	var solutionScavengerHashString = hex.EncodeToString(solutionScavengerHash[:])
  // try getting a commit using the the hash of solution and address
	_, isFound := k.GetCommit(ctx, solutionScavengerHashString)
  // return an error if a commit doesn't exist
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash doesn't exists")
	}
  // find a hash of the solution
	var solutionHash = sha256.Sum256([]byte(msg.Solution))
  // encode the solution hash to string
	var solutionHashString = hex.EncodeToString(solutionHash[:])
	var scavenge types.Scavenge
  // get a scavenge from the stre using the solution hash
	scavenge, isFound = k.GetScavenge(ctx, solutionHashString)
  // return an error if the solution doesn't exist
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash doesn't exists")
	}
  // check that the scavenger property contains a valid address
	_, err := sdk.AccAddressFromBech32(scavenge.Scavenger)
  // return an error if a scavenge has already been solved
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge has already been solved")
	}
  // save the scavebger address to the scavenge
	scavenge.Scavenger = msg.Creator
  // save the correct solution to the scavenge
	scavenge.Solution = msg.Solution
  // get address of the module account
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
  // convert scavenger address from string to sdk.AccAddress
	scavenger, err := sdk.AccAddressFromBech32(scavenge.Scavenger)
	if err != nil {
		panic(err)
	}
  // parse tokens from a string to sdk.Coins
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}
  // send tokens from a module account to the scavenger
	sdkError := k.bankKeeper.SendCoins(ctx, moduleAcct, scavenger, reward)
	if sdkError != nil {
		return nil, sdkError
	}
  // save the udpated scavenge to the store
	k.SetScavenge(ctx, scavenge)
	return &types.MsgRevealSolutionResponse{}, nil
}
```