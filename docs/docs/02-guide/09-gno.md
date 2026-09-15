---
description: Build gno.land smart contracts with Ignite CLI.
---

# Gno.land smart contracts

Ignite CLI is the developer tooling for [gno.land](https://gno.land). Gno (Gnolang) is a deterministic
variant of Go: smart contracts ("realms") are plain Go-like code, with
package-level state persisted on-chain.

This guide walks through the full development loop.

## Scaffold a realm

A **realm** (`gno.land/r/...`) is a stateful smart contract. A **package**
(`gno.land/p/...`) is a stateless library.

```sh
ignite scaffold realm counter
cd counter
```

This creates a `gnomod.toml`, a working realm with a counter, its tests, and a
`Render` function for the gno.land web. Use a full path for nested modules:

```sh
ignite scaffold realm gno.land/r/demo/counter
ignite scaffold package utils
```

## Start a dev chain

```sh
ignite chain serve
```

This runs a local, in-memory gno.land chain:

- the well-known dev account `test1` is pre-funded (its mnemonic is printed),
- your realm is deployed at genesis,
- every `.gno` file change reloads the chain,
- the RPC endpoint listens on `tcp://127.0.0.1:26657` (override with
  `--remote`).

The gno stdlibs are embedded in the binary; no gno installation is required.

## Interact with the chain

```sh
# invoke a realm function
ignite chain call gno.land/r/counter Increment

# evaluate a read-only expression
ignite chain query "gno.land/r/counter.Get()"

# send coins (by address or key name)
ignite chain send alice 10000000ugnot
```

## Deploy

```sh
ignite chain deploy               # deploy the package in the current dir
ignite chain deploy path/to/pkg   # or point at a package directory
```

The module path comes from `gnomod.toml`. Point at any chain with `--remote`
and sign with `--from <key>` (defaults to the dev account `test1` on the local
dev chain).

## Manage accounts

Accounts live in the gno keybase (`~/.config/gno`) and are shared with
`gnokey`:

```sh
ignite account create alice
ignite account list
ignite account export alice --output alice.asc
ignite account import alice alice.asc
```

New accounts hold no coins; fund them from the dev account with
`ignite chain send`.

## Cosmos SDK chains

The Cosmos SDK tooling for sovereign blockchains lives under
`ignite cosmos` (see `ignite cosmos --help`).
