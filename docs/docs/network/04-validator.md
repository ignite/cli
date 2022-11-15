---
sidebar_position: 4
description: Ignite Network commands for validators.
---

# Validator Guide

Validators join as genesis validators for chain launches on Ignite Chain.

---

## List all published chains

Validators can list and explore published chains to be launched on Ignite.

```
ignite n chain list
```

**Output**

```
Launch Id  Chain Id  Source                              Phase

3   example-1   https://github.com/ignite/example   coordinating
2   spn-10      https://github.com/tendermint/spn   launched
1   example-20  https://github.com/tendermint/spn   launching
```

- `Launch ID` is the unique identifier of the chain on Ignite. This is the ID used to interact with the chain launch.
- `Chain ID` represents the identifer of the chain network once it will be launched. It should be a unique identifier in
  practice but doesn't need to be unique on Ignite.
- `Source` is the repository URL of the project.
- `Phase` is the current phase of the chain launch. A chain can have 3 different phases:
  - `coordinating`: means the chain is open to receive requests from validators
  - `launching`: means the chain no longer receives requests but it hasn't been launched yet
  - `launched`: means the chain network has been launched

---

## Request network participation

When the chain is in the coordination phase, validators can request to be a genesis validator for the chain.
Ignite CLI supports an automatic workflow that can setup a node for the validator and a workflow for advanced users with
a specific setup for their node.

### Simple Flow

`ignite` can handle validator setup automatically. Initialize the node and generate a gentx file with default values:

```
ignite n chain init 3
```

**Output**

```
✔ Source code fetched
✔ Blockchain set up
✔ Blockchain initialized
✔ Genesis initialized
? Staking amount 95000000stake
? Commission rate 0.10
? Commission max rate 0.20
? Commission max change rate 0.01
⋆ Gentx generated: /Users/lucas/spn/3/config/gentx/gentx.json
```

Now, create and broadcast a request to join a chain as a validator:

```
ignite n chain join 3 --amount 100000000stake
```

The join command accepts a `--amount` flag with a comma-separated list of tokens. If the flag is provided, the
command will broadcast a request to add the validator’s address as an account to the genesis with the specific amount.

**Output**

```
? Peer's address 192.168.0.1:26656
✔ Source code fetched
✔ Blockchain set up
✔ Account added to the network by the coordinator!
✔ Validator added to the network by the coordinator!
```

---

### Advanced Flow

Using a more advanced setup (e.g. custom `gentx`), validators must provide an additional flag to their command
to point to the custom file:

```
ignite n chain join 3 --amount 100000000stake --gentx ~/chain/config/gentx/gentx.json
```

---

## Launch the network

### Simple Flow

Generate the final genesis and config of the node:

```
ignite n chain prepare 3
```

**Output**

```
✔ Source code fetched
✔ Blockchain set up
✔ Chain's binary built
✔ Genesis initialized
✔ Genesis built
✔ Chain is prepared for launch
```

Next, start the node:

```
exampled start --home ~/spn/3
```

---

### Advanced Flow

Fetch the final genesis for the chain:

```
ignite n chain show genesis 3
```

**Output**

```
✔ Source code fetched
✔ Blockchain set up
✔ Blockchain initialized
✔ Genesis initialized
✔ Genesis built
⋆ Genesis generated: ./genesis.json
```

Next, fetch the persistent peer list:

```
ignite n chain show peers 3
```

**Output**

```
⋆ Peer list generated: ./peers.txt
```

The fetched genesis file and peer list can be used for a manual node setup.
