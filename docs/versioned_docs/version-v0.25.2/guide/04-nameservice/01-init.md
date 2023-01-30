---
sidebar_position: 1
description: Scaffold a blockchain and create a nameservice module.
---

# Scaffold the Nameservice Module

Scaffold a blockchain and create a `nameservice` module for the nameservice app. Remember, the goal of the nameservice app is to send tokens between participants so that end users can buy names and set a value to the names.

## Create a Blockchain

Scaffold a new Cosmos SDK blockchain using the `ignite scaffold chain` command. The [ignite scaffold chain](../../../references/cli#ignite-scaffold-chain) command accepts one argument: the Go module path that is used for the project.

By default, a chain is scaffolded with a new empty Cosmos SDK module. You want to create the nameservice module without scaffolding a module, so use the `--no-module` flag:

```bash
ignite scaffold chain nameservice --no-module
```

This command created a new directory `nameservice` with a brand new Cosmos SDK blockchain. This blockchain doesn't have any application-specific logic yet, but it imports standard Cosmos SDK modules, such as `auth`, `bank`, `mint`, and others.

Change the current directory to `nameservice`:

```bash
cd nameservice
```

Inside the `nameservice` project directory you can execute other Ignite CLI commands to start a blockchain node, scaffold modules, messages, types, generate code, and much more.

## Create the Module

Scaffold a new module called `nameservice`. By design, the `nameservice` module must send tokens between participants. The send tokens functionality is implemented in the standard `bank` module.

To specify `bank` as a dependency, use the optional `--dep` flag:

```bash
ignite scaffold module nameservice --dep bank
```

## Results

The Ignite CLI scaffold command has done all of the work for you!

- The `nameservice` module was created in the `x/nameservice` directory.
- The `nameservice` module was imported into the blockchain in the `app/app.go` file.

Now, define the actions your app can make with messages.
