---
order: 4
---

# Messages

Messages are a great place to start when building a module because they define the actions that your application can make. Think of all the scenarios where a user would be able to update the state of the application in any way. These should be boiled down into basic interactions, similar to **CRUD** (Create, Read, Update, Delete).

The Scavenge module will have 3 messages:

* Submit scavenge
* Commit solution
* Reveal solution

## Submit Scavenge Message

Submit scavenge message should contain all the necessary information when creating a scavenge:

* Description - what is the question to be solved or description of the challenge.
* Solution hash - the scrambled solution.
* Reward - this is the bounty that is awarded to whoever submits the answer first.

Use the `starport scaffold message` command to scaffold a new Cosmos SDK message for your module. The command accepts message name as the first argument and a list of fields. By default, a message is scaffolded in a module with a name that matches the name of the project, in our case `scavenge` (this behaviour can be overwritten by using a flag).

```
starport scaffold message submit-scavenge solutionHash description reward
```

The command has created and modified several files.

* `proto/scavenge/tx.proto`: `MsgSubmitScavenge` and `MsgSubmitScavengeResponse` proto messages are added and a `SubmitScavenge` RPC is registered in the `Msg` service.
* `x/scavenge/types/message_submit_scavenge.go`: methods are defined to satisfy `Msg` interface.
* `x/scavenge/handler.go`: `MsgSubmitScavenge` message is registered in the module message handler.
* `x/scavenge/keeper/msg_server_submit_scavenge.go`: `SubmitScavenge` keeper method is defined
* `x/scavenge/client/cli/tx_submit_scavenge.go`: CLI command added to brodcast a transaction with a message.
* `x/scavenge/client/cli/tx.go`: CLI command is registered.
* `x/scavenge/types/codec.go`: codecs are registered.

In `x/scavenge/types/message_submit_scavenge.go` you can notice that the message follows the `sdk.Msg` interface. The message `struct` contains all the necessary information when creating a new scavenge: `Description`, `SolutionHash`, `Reward`, and `Creator` (which was added automatically).

The `Msg` interface requires some other methods be set, like validating the content of the `struct`, and confirming the msg was signed and submitted by the Creator.

Now that one can submit a scavenge the only other essential action is to be able to solve it. This should be broken into two separate actions as described before: `MsgCommitSolution` and `MsgRevealSolution`.

## Commit Solution Message

Commit solution message needs to contain the following fields:

* Solution hash - the scrambled solution.
* Solution scavenger hash - this is the hash of the combination of the solution and the person who solved it.

```
starport scaffold message commit-solution solutionHash solutionScavengerHash
```

As you're using the same `starport scaffold message` command the set of modified and created files are the same.

## Reveal Solution Message

Reveal solution message needs only one field:

* Solution - this is the plain text version of the solution.

```
starport scaffold message reveal-solution solution
```

Information about the scavenger (creator of the message is automatically included) and solution hash can be deterministically derived from the solution string.
