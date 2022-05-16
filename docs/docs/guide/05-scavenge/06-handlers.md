---
sidebar_position: 6
---

# Handlers

For a message to reach a keeper, it has to go through a handler. A handler is where you can apply logic to  allow or deny a message to succeed.

* If you're familiar with the [Model-view-controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller) (MVC) software architecture, the keeper is a bit like the model, and the handler is a bit like the controller. 
* If you're familiar with [React](<https://en.wikipedia.org/wiki/React_(web_framework)>) or [Vue](https://en.wikipedia.org/wiki/Vue.js) architecture, the keeper is a bit like the reducer store and the handler is a bit like actions.

The module-wide message handler is defined in `x/scavenge/handler.go`. Three message types were automatically added to the handler:

* `MsgSubmitScavenge`
* `MsgCommitSolution`
* `MsgRevealSolution`

Each message, when handled, calls the appropriate keeper method that is responsible for committing changes to the store.
