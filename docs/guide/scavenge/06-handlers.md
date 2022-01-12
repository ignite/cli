---
order: 6
---

# Handlers

In order for a **Message** to reach a **Keeper**, it has to go through a **Handler**. This is where logic can be applied to either allow or deny a `Message` to succeed. If you're familiar with [Model View Controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller) (MVC) architecture, the `Keeper` is a bit like the **Model** and the `Handler` is a bit like the **Controller**. If you're familiar with [React/Redux](<https://en.wikipedia.org/wiki/React_(web_framework)>) or [Vue/Vuex](https://en.wikipedia.org/wiki/Vue.js) architecture, the `Keeper` is a bit like the **Reducer/Store** and the `Handler` is a bit like **Actions**.

Module-wide message handler is defined in `x/scavenge/handler.go`. Three message types were added to the handler:

* `MsgSubmitScavenge`
* `MsgCommitSolution`
* `MsgRevealSolution`

Each message, when handled, calls the appropriate keeper method, responsible for committing changes to the store.
