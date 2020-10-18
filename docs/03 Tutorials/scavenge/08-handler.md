---
order: 8
---

# Handler

In order for a **Message** to reach a **Keeper**, it has to go through a **Handler**. This is where logic can be applied to either allow or deny a `Message` to succeed. It's also where logic as to exactly how the state should change within the Keeper should take place. If you're familiar with [Model View Controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller) (MVC) architecture, the `Keeper` is a bit like the **Model** and the `Handler` is a bit like the **Controller**. If you're familiar with [React/Redux](<https://en.wikipedia.org/wiki/React_(web_framework)>) or [Vue/Vuex](https://en.wikipedia.org/wiki/Vue.js) architecture, the `Keeper` is a bit like the **Reducer/Store** and the `Handler` is a bit like **Actions**.

Our Handler will go in `./x/scavenge/handler.go` and will follow the suggestions outlined in the boilerplate. We will create handler functions in separate files named `handler<Action.go` for each of our three `Message` types, `MsgCreateScavenge`, `MsgCommitSolution` and `MsgRevealSolution`.

Running the `starport type` command should have already added a `handlerMsgCreateScavenge.go` and a `handlerMsgCreateCommit.go` file. In essence, you can rename `handlerMsgCreateCommit` to `handlerMsgCommitSolution`, and copy the template for `HandlerMsgRevealSolution`. We'll be modifying the files to look as follows:

#### `handlerMsgCreateScavenge.go`
<<< @/scavenge/scavenge/x/scavenge/handlerMsgCreateScavenge.go

#### `handlerMsgCommitSolution.go`
<<< @/scavenge/scavenge/x/scavenge/handlerMsgCommitSolution.go

#### `handlerMsgRevealSolution.go`
<<< @/scavenge/scavenge/x/scavenge/handlerMsgRevealSolution.go

After we're done, we need to register these functions inside our main handler:

#### `handler.go`
<<< @/scavenge/scavenge/x/scavenge/handler.go

## moduleAcct

You might notice the use of `moduleAcct` within the `handleMsgCreateScavenge` and `handleMsgRevealSolution` handler functions. This account is not controlled by a public key pair, but is a reference to an account that is owned by this actual module. It is used to hold the bounty reward that is attached to a scavenge until that scavenge has been solved, at which point the bounty is paid to the account who solved the scavenge.

## Events

At the end of each handler is an `EventManager` which will create logs within the transaction that reveals information about what occurred during the handling of this message. This is useful for client side software that wants to know exactly what happened as a result of this state transition. These Events use a series of pre-defined types that can be found in `./x/scavenge/types/events.go` and look as follows:

<<< @/scavenge/scavenge/x/scavenge/types/events.go

Now that we have all the necessary pieces for updating state (`Message`, `Handler`, `Keeper`) we might want to consider ways in which we can _query_ state. This is typically done via a REST endpoint and/or a CLI. Both of those clients interact with part of the app which queries state, called the `Querier`.
