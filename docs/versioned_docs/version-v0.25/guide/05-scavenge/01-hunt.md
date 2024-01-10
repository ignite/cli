---
sidebar_position: 1
slug: /guide/scavenge

---

# Scavenger hunt game

In this tutorial, you will build a blockchain for a scavenger hunt game and learn how to:

* Implement custom logic in the CLI commands
* Use an escrow account to store tokens

This tutorial was first presented as a workshop at GODays 2020 Berlin by [Billy Rennekamp](https://twitter.com/billyrennekamp).

This session aims to get you thinking about what is possible when developing applications that have access to **digital scarcity as a primitive**. The easiest way to think of scarcity is as money; If money grew on trees it would stop being _scarce_ and stop having value. 

Although a long history of software deals with money, the representation of money has not been a first-class citizen in the programming environment. Instead, money has historically been represented as a number or a float. It has been left up to a third party merchant service or process of exchange to swap the _representation_ of money for actual cash. If money were a primitive in a software environment, it would allow for **real economies to exist within games and applications**. Money as a primitive takes one more step in erasing the line between games, life, and play.

This tutorial uses the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework that makes it possible to build **deterministic state machines**. A state machine is simply an application that has a state and explicit functions for updating that state. You can think of a light bulb and a light switch as a kind of state machine: the state of the "application" is either `light on` or `light off`. There is one function in this state machine: `flip switch`. Every time you trigger `flip switch`, the state of the application goes from `light on` to `light off` or vice versa.

## Deterministic state machine

A **deterministic** state machine is a state machine in which an accumulation of actions, taken together and replayed, have the same outcome. So if you were to take all the `switch on` and `switch off` actions of the entire month of January for some room and replay them in August, you have the same final state of `light on` or `light off`. Nothing about the metaphorical months of January or August changes the outcome. Of course, a _real_ room is not deterministic if things like power shortages or maintenance occurred during those months.

A strong feature of deterministic state machines lets you  track changes with **cryptographic hashes** of the state, similar to version control systems like `git`. If there is agreement about the hash of a certain state, it is unnecessary to replay every action from genesis to ensure that two repos are in sync with each other. These properties are useful when dealing with software that is run by many different people in many different situations, just like git.

Another nice property of cryptographically hashing state is the system of **reliable dependencies**. For example, a developer can build software that uses your library and references a specific state in your software. That way if your code changes in a way that breaks code in a specific state, developers are not required to use your new version but can continue to use the referenced version. This same property of knowing exactly what the state of a system (as well as all the ways that state can update) makes it possible to have the necessary assurances that allow for digital scarcity within an application. _If I say there is only one of some thing within a state machine and you know that there is no way for that state machine to create more than one, you can rely on there always being only one._

You might have guessed by now that we're talking about **blockchains**. Blockchains are deterministic state machines that have very specific rules about how state is updated. Blockchains checkpoint state with cryptographic hashes and use asymmetric cryptography to handle **access control**. There are different ways that different blockchains decide who can make a checkpoint of state. These entities are called **validators**. On blockchains like Bitcoin or Ethereum, validators are chosen by an electricity-intensive process called Proof-of-Work (PoW) in tandem with something called the longest chain rule or the Nakamoto consensus. Nakamoto solved the permissionless consensus problem with a remarkably simple but powerful scheme that uses only basic cryptographic primitives (hash functions and digital signatures).

## Proof-of-Stake (PoS)

The state machine you build with this tutorial uses the energy-efficient Proof-of-Stake (PoS) consensus that can consist of one or many validators, either trusted or byzantine. When a system handles _real_ scarcity, the integrity of that system becomes very important. One way to ensure integrity is by sharing the responsibility of maintaining the integrity with a large group of independently motivated participants as validators.

So, now that you know a little more about **why** you might build an app like this, start to dive into the game itself.
