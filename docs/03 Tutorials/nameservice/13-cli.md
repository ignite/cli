---
order: 14
---

# Nameservice Module CLI

The Cosmos SDK uses the [`cobra`](https://github.com/spf13/cobra) library for CLI interactions. This library makes it easy for each module to expose its own commands. The `type` command should have scaffolded the following files for us -

- `./x/nameservice/client/cli/queryWhois.go`
- `./x/nameservice/client/cli/txWhois.go`

## Queries

Start in `queryWhois.go`. Here, define `cobra.Command`s for each of your modules `Queriers` (`resolve`, and `whois`):

<<< @/nameservice/nameservice/x/nameservice/client/cli/queryWhois.go

Notes on the above code:

- The CLI introduces a new `context`: [`CLIContext`](https://godoc.org/github.com/cosmos/cosmos-sdk/client/context#CLIContext). It carries data about user input and application configuration that are needed for CLI interactions.
- The `path` required for the `cliCtx.QueryWithData()` function maps directly to the names in your query router.
  - The first part of the path is used to differentiate the types of queries possible to SDK applications: `custom` is for `Queriers`.
  - The second piece (`nameservice`) is the name of the module to route the query to.
  - Finally there is the specific querier in the module that will be called.
  - In this example the fourth piece is the query. This works because the query parameter is a simple string. To enable more complex query inputs you need to use the second argument of the [`.QueryWithData()`](https://godoc.org/github.com/cosmos/cosmos-sdk/client/context#CLIContext.QueryWithData) function to pass in `data`. For an example of this see the [queriers in the Staking module](https://github.com/cosmos/cosmos-sdk/blob/5af6bd77aa6c0e8facc936947a3365416892e44d/x/staking/keeper/querier.go).

## Transactions

Now that the query interactions are defined, it is time to move on to transaction generation in `txWhois.go`:

> _*NOTE*_: Your application needs to import the code you just wrote. Here the import path is set to this repository (`github.com/cosmos/sdk-tutorials/nameservice/x/nameservice`). If you are following along in your own repo you will need to change the import path to reflect that (`github.com/{ .Username }/{ .Project.Repo }/x/nameservice`).

<<< @/nameservice/nameservice/x/nameservice/client/cli/txWhois.go

We also need to add the commands to our `tx` command:

<<< @/nameservice/nameservice/x/nameservice/client/cli/tx.go

Notes on the above code:

- The `authcmd` package is used here. [The godocs have more information on usage](https://godoc.org/github.com/cosmos/cosmos-sdk/x/auth/client/cli#GetAccountDecoder). It provides access to accounts controlled by the CLI and facilitates signing.

### Now your ready to define the routes that the REST client will use to communicate with your module
