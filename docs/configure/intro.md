---
order: 1
description: Overview of configuration for Starport-launched blockchains
---

# Configuration

A `config.yml` file is created in the blockchain directory for every new blockchain that is launched with Starport. The `config.yml` file describes the development environment for your blockchain.

The generated `config.yml` defines the accounts and a validator of the blockchain.

## config.yml

The generated `config.yml` looks like:

```yml
accounts:
 - name: alice
   coins: ["1000token", "100000000stake"]
 - name: bob
   coins: ["500token"]
validator:
 name: user1
 staked: "100000000stake"
```

## Bootstrap your Blockchain with Changes

To bootstrap your blockchain with changes, you can change values for existing parameters in the generated `config.yml` and add parameters that are not provided by default.

For parameter details, see [config.yml Reference](https://docs.starport.network/configure/reference.html).

## Changes to `config.yml`

When changes to `config.yml` are saved, the blockchain automatically restarts to read the new values. The state of the blockchain is reset.
