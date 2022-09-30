---
sidebar_position: 4
---

# Messages

Messages are a great place to start when building a module because messages define your application actions. Think of all the scenarios where a user would be able to update the state of the application in any way. These scenarios are the basic interactions, similar to CRUD (create, read, update, and delete) operations. Messages are objects whose end-goal is to trigger state-transitions.

For the scavenger hunt game, the scavenge module requires 3 messages:

* Submit scavenge
* Commit solution
* Reveal solution

## Submit scavenge message

The submit scavenge message must contain all the information that is required to create a scavenge:

* Description - the question to be solved or description of the challenge.
* Solution hash - the scrambled solution.
* Reward - the bounty that is awarded to whoever submits the answer first.

Use the `ignite scaffold message` command to scaffold a new Cosmos SDK message for your module. The command accepts the message name as the first argument and a list of fields. By default, a message is scaffolded in a module with a name that matches the name of the project, in our case `scavenge`. You can use a flag to overwrite this default naming behavior.

```bash
ignite scaffold message submit-scavenge solutionHash description reward
```

The command creates and modifies several files:

```
modify app/app.go
create proto/scavenge/genesis.proto
create proto/scavenge/params.proto
create proto/scavenge/query.proto
create proto/scavenge/tx.proto
create testutil/keeper/scavenge.go
create x/scavenge/client/cli/query.go
create x/scavenge/client/cli/query_params.go
create x/scavenge/client/cli/tx.go
create x/scavenge/genesis.go
create x/scavenge/genesis_test.go
create x/scavenge/keeper/grpc_query.go
create x/scavenge/keeper/grpc_query_params.go
create x/scavenge/keeper/grpc_query_params_test.go
create x/scavenge/keeper/keeper.go
create x/scavenge/keeper/msg_server.go
create x/scavenge/keeper/msg_server_test.go
create x/scavenge/keeper/params.go
create x/scavenge/keeper/params_test.go
create x/scavenge/module.go
create x/scavenge/module_simulation.go
create x/scavenge/simulation/simap.go
create x/scavenge/types/codec.go
create x/scavenge/types/errors.go
create x/scavenge/types/expected_keepers.go
create x/scavenge/types/genesis.go
create x/scavenge/types/genesis_test.go
create x/scavenge/types/keys.go
create x/scavenge/types/params.go
create x/scavenge/types/types.go

🎉 Module created scavenge.
```

The `scaffold message` command does all of these code updates for you:

* `proto/scavenge/tx.proto`

  * Adds `MsgSubmitScavenge` and `MsgSubmitScavengeResponse` proto messages
  * Registers a `SubmitScavenge` RPC in the `Msg` service

* `x/scavenge/types/message_submit_scavenge.go`

  * Defines methods to satisfy `Msg` interface

* `x/scavenge/keeper/msg_server_submit_scavenge.go`

  * Defines the `SubmitScavenge` keeper method

* `x/scavenge/client/cli/tx_submit_scavenge.go`

  * Adds CLI command to broadcast a transaction with a message

* `x/scavenge/client/cli/tx.go`

  * Registers the CLI command

* `x/scavenge/types/codec.go`

  * Registers the codecs

In `x/scavenge/types/message_submit_scavenge.go`, you can notice that the message follows the `sdk.Msg` interface. The message `struct` automatically contains the information required to create a new scavenge:

```go
func NewMsgSubmitScavenge(creator string, solutionHash string, description string, reward string) *MsgSubmitScavenge {
	return &MsgSubmitScavenge{
		Creator:      creator,
		SolutionHash: solutionHash,
		Description:  description,
		Reward:       reward,
	}
}
```

The `Msg` interface requires some other methods be set, like validating the content of the `struct` and confirming the message was signed and submitted by the creator.

Now that a user can submit a scavenge, the only other essential action is to be able to solve the scavenge. As described earlier to prevent front running, use two separate actions, `MsgCommitSolution` and `MsgRevealSolution`.

## Commit solution message

The commit solution message requires the following fields:

* Solution hash - the scrambled solution
* Solution scavenger hash - the hash of the combination of the solution and the person who solved it

```bash
ignite scaffold message commit-solution solutionHash solutionScavengerHash
```

Because you're using the same `ignite scaffold message` command, the set of modified and created files is the same:
```
modify proto/scavenge/tx.proto
modify x/scavenge/client/cli/tx.go
create x/scavenge/client/cli/tx_commit_solution.go
create x/scavenge/keeper/msg_server_commit_solution.go
modify x/scavenge/module_simulation.go
create x/scavenge/simulation/commit_solution.go
modify x/scavenge/types/codec.go
create x/scavenge/types/message_commit_solution.go
create x/scavenge/types/message_commit_solution_test.go

🎉 Created a message `commit-solution`.
```

## Reveal solution message

The reveal solution message requires only one field:

* Solution - the plain text version of the solution

```bash
ignite scaffold message reveal-solution solution
```

Again, because you're using the same `ignite scaffold message` command, the set of modified and created files is the same for the `reveal-solution` message.

Information about the scavenger (the creator of the message is automatically included) and the solution hash can be deterministically derived from the solution string.

## Save changes

Now is a good time to store your project in a git commit:

```bash
git add .
git commit -m "add scavenge messages"
```
