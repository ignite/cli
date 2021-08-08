---
order: 2
---

# Scaffolding

## Creating a blockchain

Scaffold a new Cosmos SDK blockchain using the `starport scaffold chain` command. The command accepts one argument: the Go module path that will be used for the project.

By default a chain is scaffolded with a new empty Cosmos SDK module. Use the `--no-module` flag to skip module scaffolding.

```
starport scaffold chain github.com/cosmonaut/nameservice --no-module
```

This command created a new directory `nameservice` with a brand new Cosmos SDK blockchain. This blockchain doesn't have any application-specific logic yet, but it imports standard Cosmos SDK modules, such as `auth`, `bank`, `mint` and others.

Change the current directory to `nameservice`:

```
cd nameservice
```

Inside the project directory you can execute other Starport commands to start a blockchain node, scaffold modules, messages, types, generate code, and much more.

## Creating a module

In a Cosmos SDK blockchain application-specific logic is implemented in separate modules. This keeps code easy to understand and reuse.

Scaffold a new module called `nameservice`. Based on our design the `nameservice` module will be sending tokens between participants. Sending tokens functionality is implemented in the standard `bank` module. Specify `bank` as a dependency using the optional `--dep` flag.

```
starport scaffold module nameservice --dep bank
```

A module has been created in the `x/nameservice` directory and imported into the blockchain in `app/app.go`.