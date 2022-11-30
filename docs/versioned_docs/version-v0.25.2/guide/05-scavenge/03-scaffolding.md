---
sidebar_position: 3
---

# Scaffold the scavenge chain

Scaffold a new Cosmos SDK blockchain using the `ignite scaffold chain` command.

By default a chain is scaffolded with a new empty Cosmos SDK module. Use the `--no-module` flag to skip module scaffolding.

```bash
ignite scaffold chain scavenge --no-module
```

This command creates a new `scavenge` directory with a brand new Cosmos SDK blockchain. This blockchain doesn't have any application-specific logic yet, but it imports standard Cosmos SDK modules, such as `auth`, `bank`, `mint`, and others.

Change the current directory to `scavenge`:

```bash
cd scavenge
```

Inside the project directory, you can execute other Ignite CLI commands to start a blockchain node, scaffold modules, messages, types, generate code, and much more.

In a Cosmos SDK blockchain, implement application-specific logic in separate modules. Using modules keeps code easy to understand and reuse.

## Scaffold the scavenge module

Scaffold a new module called `scavenge`. Based on the game design, the `scavenge` module sends tokens between participants. 

- Implement sending tokens in the standard `bank` module.
- Use the optional `--dep` flag to specify the `bank` module.

```bash
ignite scaffold module scavenge --dep bank
```

This command creates the `x/scavenge` directory and imports the scavenge module into the blockchain in the `app/app.go` directory.

## Save changes

Before you go to the next step, you can store your project in a git commit:

```bash
git add .
git commit -m "scaffold scavenge chain and module"
```

