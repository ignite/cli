# Quickstart

Use these three commands to build, start, and define your first blockchain.

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

## Define

To add a custom type with create, read, update, and delete (CRUD) functionality:

```
starport type post title body
```

The `type` command scaffolds functionality a custom type.

<!-- link to starport tutorial for 2 blockchains? -->
