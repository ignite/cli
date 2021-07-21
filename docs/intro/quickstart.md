---
order: 3
description: Use commands to build, start, and add features to your first blockchain.
---

# Quickstart

Use these three commands to build, start, and add features to your first blockchain. Before moving forward, make sure to have Starport locally installed, visit the [Installation guide](https://docs.starport.network/intro/install.html) for more information.

## Create a blockchain

To create a blockchain:

```
starport scaffold chain github.com/username/myapp && cd myapp
```

The `app` command creates the blockchain directory `myapp` and scaffolds a [Cosmos SDK](https://docs.cosmos.network/) blockchain.

## Run a blockchain

To run a blockchain in your development environment:

```
starport chain serve
```

The `serve` command installs dependencies, builds, initializes, and starts the blockchain.

## Add Features

To add a custom type with create, read, update, and delete (CRUD) functionality:

```
starport scaffold type post title body
```

The `type` command scaffolds functionality a custom type.

## Learn more

- [Configure a Blockchain](../configure/index.md)
- [Run a Blockchain](../run/index.md)
- [Starport repo in GitHub](https://github.com/tendermint/starport)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Starport Channel on Discord](https://discord.com/channels/669268347736686612/737461683588431924)

## Next steps

To learn about building Cosmos SDK blockchains with Starport, try a beginner-friendly [IBC Hello world tutorial](https://tutorials.cosmos.network/hello-world/tutorial/).
