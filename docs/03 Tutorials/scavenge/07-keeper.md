---
order: 7
---

# Keeper

After using the `starport` command you should have a boilerplate `Keeper` at `./x/scavenge/keeper/keeper.go`. It contains a basic keeper with references to basic functions like `Set`, `Get` and `Delete`.

Our keeper stores all our data for our module. Sometimes a module will import the keeper of another module. This will allow state to be shared and modified across modules. Since we are dealing with coins in our module as bounty rewards, we will need to access the `bank` module's keeper (which we call `CoinKeeper`). Look at our completed `Keeper` files and you can see where the `bank` keeper is referenced and how `Set`, `Get` and `Delete` are expanded:

#### `keeper/keeper.go`
<<< @/scavenge/scavenge/x/scavenge/keeper/keeper.go

#### `keeper/scavenge.go`
<<< @/scavenge/scavenge/x/scavenge/keeper/scavenge.go

#### `keeper/commit.go`
<<< @/scavenge/scavenge/x/scavenge/keeper/commit.go

## Commits and Scavenges

You may notice reference to `types.Commit` and `types.Scavenge` throughout the `Keeper`. These are new structs defined in `./x/scavenge/types/type<Type>.go` that contin all necessary information about different scavenge challenges, and different commited solutions to those challenges. They appear similar to the `Msg` types we saw earlier because they contain similar information. We will be making some modifications to the scaffolded files.

In the `TypeScavenge.go` file, we need to delete the `ID` field, since we're going to be using the `SolutionHash` as the key. We also need to update `Reward` to `sdk.Coins`, as well as `Scavenger` to `sdk.AccAddress`, so we can make the payout once the scavenge is solved.

Once this is done, your struct in `scavenge/x/scavenge/types/TypeScavenge.go` should look like this -

<<< @/scavenge/scavenge/x/scavenge/types/TypeScavenge.go

For `TypeCommit.go`, we need to delete the `ID` field, and rename the `Creator` field to `Scavenger`.

<<< @/scavenge/scavenge/x/scavenge/types/TypeCommit.go

You can imagine that an unsolved `Scavenge` would contain a `nil` value for the fields `Solution` and `Scavenger` before they are solved. You might also notice that each type has the `String` method. This allows us to render the struct as a string for rendering.

## Prefixes

You may notice the use of `types.ScavengePrefix` and `types.CommitPrefix`. These are defined in a file called `./x/scavenge/types/key.go` and help us keep our `Keeper` organized. The `Keeper` is really just a key value store. That means that, similar to an `Object` in javascript, all values are referenced under a key. To access a value, you need to know the key under which it is stored. This is a bit like a unique identifier (UID).

When storing a `Scavenge` we use the key of the `SolutionHash` as a unique ID, for a `Commit` we use the key of the `SolutionScavengeHash`. However since we are storing these two data types in the same location, we may want to distinguish between the types of hashes we use as keys. We can do this by adding prefixes to the hashes that allow us to recognize which is which. For `Scavenge` we add the prefix `scavenge-` and for `Commit` we add the prefix `commit-`. You should add these to your `key.go` file so it looks as follows:

<<< @/scavenge/scavenge/x/scavenge/types/key.go

## Iterators

Sometimes you will want to access a `Commit` or a `Scavenge` directly by their key. That's why we have the methods `GetCommit` and `GetScavenge`. However, sometimes you will want to get every `Scavenge` at once or every `Commit` at once. To do this we use an **Iterator** called `KVStorePrefixIterator`. This utility comes from the `sdk` and iterates over a key store. If you provide a prefix, it will only iterate over the keys that contain that prefix. Since we have prefixes defined for our `Scavenge` and our `Commit` we can use them here to only return our desired data types.

---

Now that you've seen the `Keeper` where every `Commit` and `Scavenge` are stored, we need to connect the messages to the this storage. This process is called _handling_ the messages and is done inside the `Handler`.
