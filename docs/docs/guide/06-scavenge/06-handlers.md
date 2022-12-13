---
sidebar_position: 6
---

# Handlers

In the previous sections you've added three message types to the project:

* `SubmitScavenge`
* `CommitSolution`
* `RevealSolution`

In the Cosmos SDK messages are registered as RPCs in the `Msg` service in
a protocol buffer file.

```proto title="proto/scavenge/scavenge/tx.proto"
service Msg {
  rpc SubmitScavenge(MsgSubmitScavenge) returns (MsgSubmitScavengeResponse);
  rpc CommitSolution(MsgCommitSolution) returns (MsgCommitSolutionResponse);
  rpc RevealSolution(MsgRevealSolution) returns (MsgRevealSolutionResponse);
}
```

When a message is processed, your blockchain calls an appropriate keeper method.
In the next section you will define logic inside these keeper methods.