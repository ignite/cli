---
sidebar_position: 2
---

# Getting started

In this tutorial, we will be using Ignite CLI to create a new blockchain. Ignite
CLI is a command line interface that allows users to quickly and easily create
blockchain networks. By using Ignite CLI, we can quickly create a new blockchain
without having to manually set up all the necessary components.

Once we have created our blockchain with Ignite CLI, we will take a look at the
directory structure and files that were created. This will give us an
understanding of how the blockchain is organized and how the different
components of the blockchain interact with each other.

By the end of this tutorial, you will have a basic understanding of how to use
Ignite CLI to create a new blockchain, and you will have a high-level
understanding of the directory structure and files that make up a blockchain.
This knowledge will be useful as you continue to explore the world of blockchain
development.

## Creating a new blockchain

To create a new blockchain project with Ignite, you will need to run the
following command:

```
ignite scaffold chain example
```

The `ignite scaffold chain` command will create a new blockchain in a new
directory `example`.

The new blockchain is built using the Cosmos SDK framework and imports several
standard modules to provide a range of functionality. These modules include
`staking`, which enables a delegated Proof-of-Stake consensus mechanism, `bank`
for facilitating fungible token transfers between accounts, and `gov` for
on-chain governance. In addition to these modules, the blockchain also imports
other modules from the Cosmos SDK framework.

The `example` directory contains the generated files and directories that make
up the structure of a Cosmos SDK blockchain. This directory includes files for
the chain's configuration, application logic, and tests, among others. It
provides a starting point for developers to quickly set up a new Cosmos SDK
blockchain and build their desired functionality on top of it.

By default, Ignite creates a new empty custom module with the same name as the
blockchain being created (in this case, `example`) in the `x/` directory. This
module doesn't have any functionality by itself, but can serve as a starting
point for building out the features of your application. If you don't want to
create this module, you can use the `--no-module` flag to skip it.

## Directory structure

In order to understand what the Ignite CLI has generated for your project, you
can inspect the contents of the `example/` directory.

The `app/` directory contains the files that connect the different parts of the
blockchain together. The most important file in this directory is `app.go`,
which includes the type definition of the blockchain and functions for creating
and initializing it. This file is responsible for wiring together the various
components of the blockchain and defining how they will interact with each
other.

The `cmd/` directory contains the main package responsible for the command-line
interface (CLI) of the compiled binary. This package defines the commands that
can be run from the CLI and how they should be executed. It is an important part
of the blockchain project as it provides a way for developers and users to
interact with the blockchain and perform various tasks, such as querying the
blockchain state or sending transactions.

The `docs/` directory is used for storing project documentation. By default,
this directory includes an OpenAPI specification file, which is a
machine-readable format for defining the API of a software project. The OpenAPI
specification can be used to automatically generate human-readable documentation
for the project, as well as provide a way for other tools and services to
interact with the API. The `docs/` directory can be used to store any additional
documentation that is relevant to the project.

The `proto/` directory contains protocol buffer files, which are used to
describe the data structure of the blockchain. Protocol buffers are a language-
and platform-neutral mechanism for serializing structured data, and are often
used in the development of distributed systems, such as blockchain networks. The
protocol buffer files in the `proto/` directory define the data structures and
messages that are used by the blockchain, and are used to generate code for
various programming languages that can be used to interact with the blockchain.
In the context of the Cosmos SDK, protocol buffer files are used to define the
specific types of data that can be sent and received by the blockchain, as well
as the specific RPC endpoints that can be used to access the blockchain's
functionality.

The `testutil/` directory contains helper functions that are used for testing.
These functions provide a convenient way to perform common tasks that are needed
when writing tests for the blockchain, such as creating test accounts,
generating transactions, and checking the state of the blockchain. By using the
helper functions in the `testutil/` directory, developers can write tests more
quickly and efficiently, and can ensure that their tests are comprehensive and
effective.

The `x/` directory contains custom Cosmos SDK modules that have been added to
the blockchain. Standard Cosmos SDK modules are pre-built components that
provide common functionality for Cosmos SDK-based blockchains, such as support
for staking and governance. Custom modules, on the other hand, are modules that
have been developed specifically for the blockchain project and provide
project-specific functionality.

The `config.yml` file is a configuration file that can be used to customize the
blockchain during development. This file includes settings that control various
aspects of the blockchain, such as the network's ID, account balances, and the
node parameters.

The `.github` directory contains a GitHub Actions workflow that can be used to
automatically build and release a blockchain binary. GitHub Actions is a tool
that allows developers to automate their software development workflows,
including building, testing, and deploying their projects. The workflow in the
`.github` directory is used to automate the process of building the blockchain
binary and releasing it, which can save time and effort for developers. 

The `readme.md` file is a readme file that provides an overview of the
blockchain project. This file typically includes information such as the
project's name and purpose, as well as instructions on how to build and run the
blockchain. By reading the `readme.md` file, developers and users can quickly
understand the purpose and capabilities of the blockchain project and get
started using it.

## Starting a blockchain node

To start a blockchain node in development, you can run the following command:

```
ignite chain serve
```

The `ignite chain serve` command is used to start a blockchain node in
development mode. It first compiles and installs the binary using the
`ignite chain build` command, then initializes the blockchain's data directory
for a single validator using the `ignite chain init` command. After that, it
starts the node locally and enables automatic code reloading so that changes to
the code can be reflected in the running blockchain without having to restart
the node. This allows for faster development and testing of the blockchain.

Congratulations! 🥳 You have successfully created a brand-new Cosmos blockchain
using the Ignite CLI. This blockchain uses the delegated proof of stake (DPoS)
consensus algorithm, and comes with a set of standard modules for token
transfers, governance, and inflation. Now that you have a basic understanding of
your Cosmos blockchain, it's time to start building custom functionality. In the
following tutorials, you will learn how to build custom modules and add new
features to your blockchain, allowing you to create a unique and powerful
decentralized application.
