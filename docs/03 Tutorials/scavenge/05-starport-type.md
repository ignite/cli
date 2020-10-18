---
order: 5
---


# Scaffolding types

In this section, we'll be explaining how to quickly scaffold types for your application using the `starport type` command.

## `starport type`

Let's run the `starport type` command to generate our `scavenge` type - 

```
starport type scavenge description solutionHash reward solution scavenger
```

From this command, we've made changes to the following files:
- `scavenge/vue/src/store/app.js` - The `scavenge` type in our front-end application
- `scavenge/x/scavenge/client/cli/query.go` - adding the `GetCmdListScavenge` query function to the CLI commands
- `scavenge/x/scavenge/client/cli/queryScavenge.go` - defining the `GetCmdListScavenge` function
- `scavenge/x/scavenge/client/cli/tx.go` - adding `GetCmdCreateScavenge` transaction function to the CLI commands
- `scavenge/x/scavenge/client/cli/txScavenge.go` - defining the `GetCmdCreateScavenge` function
- `scavenge/x/scavenge/client/rest/queryScavenge.go` - defining the `listScavengeHandler` query function
- `scavenge/x/scavenge/client/rest/txScavenge.go` - defining the `createScavengeRequest` type and `createScavengeHandler` function
- `scavenge/x/scavenge/client/rest/rest.go` - adding the `listScavengeHandler` and `createScavengeHandler` function to the CLI commands
- `scavenge/x/scavenge/handler.go` - handle the case where `MsgCreateScavenge` is passed, and handle it
- `scavenge/x/scavenge/handlerMsgCreateScavenge.go` - define `handleMsgCreateScavenge`, which creates the scavenge
- `scavenge/x/scavenge/keeper/querier.go` - Handle the `QueryListScavenge` case to use the `listScavenge` function
- `scavenge/x/scavenge/keeper/scavenge.go` - Define the `CreateScavenge` and `listScavenge` functions
- `scavenge/x/scavenge/types/MsgCreateScavenge.go` - define the  `MsgCreateScavenge` function
- `scavenge/x/scavenge/types/TypeScavenge.go` - define the `Scavenge` type
- `scavenge/x/scavenge/types/key.go` - Adding the `ScavengePrefix` constant
- `scavenge/x/scavenge/types/querier.go` - Adding the `QueryListScavenge` constant

We also want to create a second type, `Commit`, in order to prevent frontrunning of our submitted solutions as mentioned earlier.

```
starport type commit solutionHash solutionScavengerHash
```

Here, `starport` has already done the majority of the work by helping us scaffold the necessary files and functions.

In the next sections, we'll be modifying these to give our appliation the functionality we want, according to the game.
