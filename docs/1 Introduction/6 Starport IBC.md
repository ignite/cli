# Starport IBC

Inter-Blockchain Communication Protocol in short [IBC](https://github.com/cosmos/ics#ibc-quick-references), enables communication of blockchains. IBC allows a reliable and secure connection between two chains which can be used for transferring of tokens, multi-chain smart contracts, atomic swaps or data and code sharding of any kind.

In order to communicate between chains, bootstrapping two blockchains with starport will give us a picture of how it is intended to work and what happens on each of the blockchain applications in the communication process. In this tutorial, we will be creating two blockchains, connecting those blockchains and transferring tokens via IBC.

## Scaffolding chain `foo`

To start using IBC with Starport open up a [web-based development environment](https://gitpod.io/#https://github.com/tendermint/starport/), then scaffold and launch a Stargate chain:

```
starport app github.com/foo/foo --sdk-version stargate

cd foo

starport serve
```

You now have a blockchain `foo` running, but it's not connected to anything yet.

## Scaffolding chain `bar`

To connect this blockchain to another one, open up one more [web-based development environment](https://gitpod.io/#https://github.com/tendermint/starport/) instance and follow the steps above to scaffold and launch another chain (let's call it `bar`).

To connect our blockchains, we will be using the [relayer](https://github.com/cosmos/ics/tree/master/spec/ics-018-relayer-algorithms). The relayer is our "physical" connection between the two blockchains. It is responsible for monitoring both blockchains, relaying data between them, construct appropriate diagrams and execute them accordingly on both blockchains.
Once the chain is running, you will see a "Relayer info" string in the terminal output (your string will be different):

```
✨ Relayer info: eyJDaGFpbklEIjoiYmFyIiwiTW5lbW9uaWMiOiJmcm9zdCByYXpvciBoYWxmIGxhdW5kcnkgcHJvZml0IHdpc2UgdG9uZSBibHVzaCBzdXJnZSBrZWVwIHRvZ2V0aGVyIHNsaWNlIHlvdXRoIHRydXRoIGVubGlzdCBjdXBib2FyZCBhYnNvcmIgc2VlZCBzZXJpZXMgZG91YmxlIHZpbGxhZ2UgdG9uZ3VlIGZsYXNoIGdvcmlsbGEiLCJSUENBZGRyZXNzIjoiaHR0cHM6Ly8yNjY1Ny1jNzllNDk2ZC1kZDk4LTQ4MWQtOTlmZi1jZGQ4OTA2NWQ4MWIud3MtZXUwMS5naXRwb2QuaW86NDQzIn0
```

This is a `base64` encoded JSON that contains information about the chain ID, a relayer account mnemonic and an RPC URL.

## Connecting `foo` with `bar`

To connect these chains together, copy the relayer info of chain `bar`, switch to the terminal of chain `foo` and run the following command (use your own string):

```
starport chain add eyJDaGFpbklEIjoiYmFyIiwiTW5lbW9uaWMiOiJmcm9zdCByYXpvciBoYWxmIGxhdW5kcnkgcHJvZml0IHdpc2UgdG9uZSBibHVzaCBzdXJnZSBrZWVwIHRvZ2V0aGVyIHNsaWNlIHlvdXRoIHRydXRoIGVubGlzdCBjdXBib2FyZCBhYnNvcmIgc2VlZCBzZXJpZXMgZG91YmxlIHZpbGxhZ2UgdG9uZ3VlIGZsYXNoIGdvcmlsbGEiLCJSUENBZGRyZXNzIjoiaHR0cHM6Ly8yNjY1Ny1jNzllNDk2ZC1kZDk4LTQ4MWQtOTlmZi1jZGQ4OTA2NWQ4MWIud3MtZXUwMS5naXRwb2QuaW86NDQzIn0
```

Chain `foo` will now restart and you should see information about two being connected:

```
Detected chains, linking them...
Linked foo <--> bar
```

The two chains are now connected via IBC and you have successfully created a relayer.

## Sending tokens from `foo` to `bar`

Once the chains are connected, you can use a [relayer](https://github.com/ovrclk/relayer) CLI `rly` to create an IBC token send transaction:

```
rly tx transfer foo bar 5token $(rly chains address bar)
```

After a transaction is successfully created, you can now relay it to a connected chain:

```
rly tx relay foo-bar
```

## Checking token balances on chain `bar`

To verify that an IBC transaction was relayed correctly, let's check the balances of our relayer account:

```
bard q bank balances $(bard keys show bar -a --keyring-backend test)
```

This command will output token balances for the relayer account and you should see 5 token transferred with IBC.

## Configuration

Inside your chain's project directory you will see `secret.yml`. This file contains information about the local chain's relayer account (under `accounts` property) and relayer accounts of connected chains (under `relayer` property).

Once the chain is launched with `starport serve`, Starport uses information from `secret.yml` to create a relayer config in `~/.relayer/`. Every time the chain is restarted relayer config is reset, and connections are re-established.
