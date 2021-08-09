---
order: 3
---

# Scaffolding

Scaffold a new Cosmos SDK blockchain using the `starport scaffold chain` command.

By default a chain is scaffolded with a new empty Cosmos SDK module. Use the `--no-module` flag to skip module scaffolding.

```
starport scaffold chain github.com/cosmonaut/scavenge --no-module
```

This command created a new directory `scavenge` with a brand new Cosmos SDK blockchain. This blockchain doesn't have any application-specific logic yet, but it imports standard Cosmos SDK modules, such as `auth`, `bank`, `mint` and others.

Change the current directory to `scavenge`:

```
cd scavenge
```

Inside the project directory you can execute other Starport commands to start a blockchain node, scaffold modules, messages, types, generate code, and much more.

In a Cosmos SDK blockchain, application-specific logic is implemented in separate modules. Using modules keeps code easy to understand and reuse.

Scaffold a new module called `scavenge`. Based on our design the `scavenge` module will be sending tokens between participants. Sending tokens is implemented in the standard `bank` module. Specify `bank` as a dependency using the optional `--dep` flag.

```
starport scaffold module scavenge --dep bank
```

A module has been created in the `x/scavenge` directory and imported into the blockchain in `app/app.go`.
