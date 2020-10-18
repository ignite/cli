---
order: 9
---

# Querier

In order to query the data of our app we need to make it accessible using our `Querier`. This piece of the app works in tandem with the `Keeper` to access state and return it. The `Querier` is defined in `./x/scavenge/keeper/querier.go`. Our `starport` tool starts us out with some suggestions on how it should look, and similar to our `Handler` we want to handle different queried routes. You could make many different routes within the `Querier` for many different types of queries, but we will just make three:

- `listScavenge` will list all scavenges
- `getScavenge` will get a single scavenge by `solutionHash`
- `listCommit` will list all commits
- `getCommit` will get a single commit by `solutionScavengerHash`

Combined into a switch statement and with each of the functions fleshed out, the file should look as follows:

#### `querier.go`
<<< @/scavenge/scavenge/x/scavenge/keeper/querier.go


## Types

You may notice that we use four different imported types on our initial switch statement. These are defined within our `./x/scavenge/types/querier.go` file as simple strings. That file should look like the following:

<<< @/scavenge/scavenge/x/scavenge/types/querier.go

Our queries are rather simple since we've already outfitted our `Keeper` with all the necessary functions to access state. You can see the iterator being used here as well.

Now that we have all of the basic actions of our module created, we want to make them accessible. We can do this with a CLI client and a REST client. For this tutorial we will be creating a CLI client. If you are interested in what goes into making a REST client, check out the [Nameservice Tutorial](../../nameservice/tutorial/00-intro.md).

Let's take a look at what goes into making a CLI.
