---
order: 10
---

# CLI
A Command Line Interface (CLI) will help us interact with our app once it is running on a machine somewhere. Each Module has it's own namespace within the CLI that gives it the ability to create and sign Messages destined to be handled by that module. It also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI will let you do things like generate keys for a new account or check the status of an interaction you already had with the application.

The CLI for our module is broken into two files called `tx.go` and `query.go` which are located in `./x/scavenge/client/cli/`. One file is for making transactions that contain messages which will ultimately update our state. The other is for making queries which will give us the ability to read information from our state. Both files utilize the [Cobra](https://github.com/spf13/cobra) library.

## tx.go
The `tx.go` file contains `GetTxCmd` which is a standard method within the Cosmos SDK. It is referenced later in the `module.go` file which describes exactly which attributes a modules has. This makes it easier to incorporate different modules for different reasons at the level of the actual application. After all, we are focusing on a module at this point, but later we will create an application that utilizes this module as well as other modules which are already available within the Cosmos SDK.

Inside `GetTxCmd` we create a new module-specific command and call is `scavenge`. Within this command we add a sub-command for each Message type we've defined: 
* `GetCmdCreateScavenge`
* `GetCmdCommitSolution`
* `GetCmdRevealSolution`


Each function takes parameters from the **Cobra** CLI tool to create a new msg, sign it and submit it to the application to be processed. These functions should go into the `tx.go` and `tx<Type>.go` files, and look as follows.

#### `scavenge/x/scavenge/client/cli/tx.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/tx.go

#### `scavenge/x/scavenge/client/cli/txScavenge.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/txScavenge.go

#### `scavenge/x/scavenge/client/cli/txCommit.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/txCommit.go

### sha256
Note that the `txScavenge` and `txCommit` files file make use of the `sha256` library for hashing our plain text solutions into the scrambled hashes. This activity takes place on the client side so the solutions are never leaked to any public entity which might want to sneak a peak and steal the bounty reward associated with the scavenges. You can also notice that the hashes are converted into hexadecimal representation to make them easy to read as strings (which is how they are ultimately stored in the keeper).

## query.go
The `query.go` file contains similar **Cobra** commands that reserve a new name space for referencing our `scavenge` module. Instead of creating and submitting messages however, the `query.go` and `query<Type>.go` files create queries and return the results in human readable form. The queries it handles are the same we defined in our `querier.go` file earlier:
* `GetCmdListScavenge`
* `GetCmdGetScavenge`
* `GetCmdListCommit`
* `GetCmdGetCommit`

After defining these commands, your files should look like:

#### `query.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/query.go

#### `queryScavenge.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/queryScavenge.go

#### `queryCommit.go`
<<< @/scavenge/scavenge/x/scavenge/client/cli/queryCommit.go

Your application is now complete! In the next step, we can start playing our game.