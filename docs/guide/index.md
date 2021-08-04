---
title: Introduction
order: 0
parent:
  order: 0
  title: Developer Guide
---

# Introduction

By following this guide you will learn how to:

* Install Starport CLI on your local machine
* Create a new blockchain and start a node locally for development
* Make your blockchain say "Hello, World!"
  * Scaffold a Cosmos SDK query
  * Modify a keeper method to return a static string
  * Use the blockchain's CLI to make a query
* Write and read blog posts to your chain in the Blog tutorial
  * Scaffold a Cosmos SDK message
  * Define new types in protocol buffer files
  * Write keeper methods to write data to the store
  * Read data from the store and return it as a result a query
  * Use the blockchain's CLI to broadcast transactions
* Build a blockchain for buying and selling names in the Nameservice tutorial
  * Scaffold a `map` without messages
  * Use other module's methods in your custom module
  * Send tokens between addresses
* Build a guessing game with rewards in the Scavenge tutorial
  * Use an escrow account to store tokens
* Use the Inter-Blockchain Communication (IBC) protocol
  * Scaffold an IBC-enabled module
  * Send and receive IBC packets
  * Configure and run a built-in IBC relayer
* Build a decentralized order-book token exchange
  * Build an advanced IBC-enabled module