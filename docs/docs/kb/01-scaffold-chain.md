---
sidebar_position: 1
description: High-level overview of a new Cosmos SDK blockchain project built with Ignite CLI.
---

# Scaffold a chain

The `ignite scaffold chain` command scaffolds a new Cosmos SDK blockchain project.

## Build a blockchain app

To build the planet application:

```bash
ignite scaffold chain planet
```

## Directory structure

The `ignite scaffold chain planet` command creates a directory called `planet` that contains all the files for your project and initializes a local git repository. The `planet` argument is a string that is used for the Go module path. The repository name (`planet`, in this case) is used as the project's name.

The project directory structure:

- `app`: files that wire the blockchain together
- `cmd`: binary for the blockchain node
- `docs`: static `openapi.yml` API doc for the blockchain node
- `proto`: protocol buffer files for custom modules
- `x`: modules
- `vue`: scaffolded web application (optional)
- `config.yml`: configuration file

### Application-specific logic

Most of the logic of your application-specific blockchain is written in custom modules. Each module effectively encapsulates an independent piece of functionality. Following the Cosmos SDK convention, custom modules are stored inside the `x` directory. By default, `ignite scaffold chain` scaffolds a module with a name that matches the name of the project. In this example, the module name is `x/planet`.

### Proto files

Every Cosmos SDK module has protocol buffer files that define data structures, messages, queries, RPCs, and so on. The `proto` directory contains a directory with proto files for each custom module in the `x` directory.

### Global settings

Global changes to your blockchain are defined in files inside the `app` directory. These changes include importing third-party modules, defining relationships between modules, and configuring blockchain-wide settings.

### Configuration

The `config.yml` file contains configuration options that Ignite CLI uses to build, initialize, and start your blockchain node in development.  

## Address prefix

Account addresses on Cosmos SDK-based blockchains have string prefixes. For example, the Cosmos Hub blockchain uses the default `cosmos` prefix, so that addresses look like this: `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

### Change prefix on new blockchains

When you create a new blockchain, pass a prefix as a value to the `--address-prefix` flag:

```bash
ignite scaffold chain planet --address-prefix moonlight
```

Using the `moonlight` prefix, account addresses on your blockchain look like this: `moonlight12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

### Change prefix on existing blockchains

To change the prefix after the blockchain has been scaffolded, modify the `AccountAddressPrefix` in the `app/app.go` file.

## Cosmos SDK version

By default, the `ignite scaffold chain` command creates a Cosmos SDK blockchain using the latest stable version of the Cosmos SDK.
