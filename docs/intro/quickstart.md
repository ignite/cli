---
order: 2
---

# Quickstart

Use these three commands to build, start, and add features to your first blockchain.

## Create a blockchain

To create a blockchain:

```
starport app github.com/username/myapp && cd myapp
```

The `app` command creates the blockchain directory `myapp` and scaffolds a [Cosmos SDK](https://docs.cosmos.network/) blockchain.

## Run a blockchain

To run a blockchain in your development environment:

```
starport serve
```

The `serve` command installs dependencies, builds, initializes, and starts the blockchain.

## Add Features

To add a custom type with create, read, update, and delete (CRUD) functionality:

```
starport type post title body
```

The `type` command scaffolds functionality a custom type.

## Learn more

- [Configure a Blockchain](./3-Configure-Your-Blockchain.html)
- [Run a Blockchain](./2-Run-a-Blockchain.html)
- [Starport repo in GitHub](https://github.com/tendermint/starport)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Starport Channel on Discord](https://discord.com/channels/669268347736686612/737461683588431924)

## Next steps

To learn about building Cosmos SDK blockchains with Starport, try a beginner-friendly [IBC Hello world tutorial](https://tutorials.cosmos.network/hello-world/tutorial/).
