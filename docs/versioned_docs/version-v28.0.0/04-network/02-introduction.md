---
sidebar_position: 2
description: Introduction to Ignite Network commands.
---

# Ignite Network commands

The `ignite network` commands allow to coordinate the launch of sovereign Cosmos blockchains by interacting with the
Ignite Chain.

To launch a Cosmos blockchain you need someone to be a coordinator and others to be validators. These are just roles,
anyone can be a coordinator or a validator.

- A coordinator publishes information about a chain to be launched on the Ignite blockchain, approves validator requests
  and coordinates the launch.
- Validators send requests to join a chain and start their nodes when a blockchain is ready for launch.

## Launching a chain on Ignite

Launching with the CLI can be as simple as a few short commands with the CLI using `ignite network` command
namespace.

> **NOTE:** `ignite n` can also be used as a shortcut for `ignite network`.

To publish the information about your chain as a coordinator, run the following command (the URL should point to a
repository with a Cosmos SDK chain):

```
ignite network chain publish github.com/ignite/example
```

This command will return the launch identifier you will be using in the following
commands. Let's say this identifier is 42.
Next, ask validators to initialize their nodes and request to join the network.
For a testnet you can use the default values suggested by the
CLI.

```
ignite network chain init 42
ignite network chain join 42 --amount 95000000stake
```

As a coordinator, list all validator requests:

```
ignite network request list 42
```

Approve validator requests:

```
ignite network request approve 42 1,2
```

Once you've approved all validators you need in the validator set, announce that
the chain is ready for launch:

```
ignite network chain launch 42
```

Validators can now prepare their nodes for launch:

```
ignite network chain prepare 42
```

The output of this command will show a command that a validator would use to
launch their node, for example `exampled --home ~/.example`. After enough
validators launch their nodes, a blockchain will be live.

---

The next two sections provide more information on the process of coordinating a chain launch from a coordinator and
participating in a chain launch as a validator.
