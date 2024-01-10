---
description: Information about the generated TypeScript client code.
---

# TypeScript frontend

Ignite offers powerful functionality for generating client-side code for your
blockchain. Think of this as a one-click client SDK generation tailored
specifically for your blockchain.

See `ignite generate ts-client --help` learn more on how to use TypeScript code generation.

## Starting a node

Create a new blockchain with `ignite scaffold chain`. You can use an existing
blockchain project if you have one, instead.

```
ignite scaffold chain example
```

For testing purposes add a new account to `config.yml` with a mnemonic:

```yml title="config.yml"
accounts:
  - name: frank
    coins: ["1000token", "100000000stake"]
    mnemonic: play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint
```

Run a command to generate TypeScript clients for both standard and custom Cosmos
SDK modules:

```
ignite generate ts-client --clear-cache
```

Run a command to start your blockchain node:

```
ignite chain serve -r
```

## Setting up a TypeScript frontend client

The best way to get started building with the TypeScript client is by using 
[Vite](https://vitejs.dev). Vite provides boilerplate code for
vanilla TS projects as well as React, Vue, Lit, Svelte and Preact frameworks.
You can find additional information at the [Vite Getting Started
guide](https://vitejs.dev/guide).

You will also need to [polyfill](https://developer.mozilla.org/en-US/docs/Glossary/Polyfill) the client's dependencies. The following is an
example of setting up a vanilla TS project with the necessary polyfills:

```bash
npm create vite@latest my-frontend-app -- --template vanilla-ts
cd my-frontend-app
npm install --save-dev @esbuild-plugins/node-globals-polyfill @rollup/plugin-node-resolve
```

You must then create the necessary `vite.config.ts` file.

```typescript title="my-frontend-app/vite.config.ts"
import { nodeResolve } from "@rollup/plugin-node-resolve";
import { NodeGlobalsPolyfillPlugin } from "@esbuild-plugins/node-globals-polyfill";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [nodeResolve()],

  optimizeDeps: {
    esbuildOptions: {
      define: {
        global: "globalThis",
      },
      plugins: [
        NodeGlobalsPolyfillPlugin({
          buffer: true,
        }),
      ],
    },
  },
});
```

You are then ready to use the generated client code inside this project directly
or by publishing the client and installing it like any other `npm` package.

After the chain starts, you will see Frank's address is
`cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7`. We'll be using Frank's account
for querying data and broadcasting transactions in the next section.

## Querying

The code generated in `ts-client` comes with a `package.json` file ready to
publish which you can modify to suit your needs. To use`ts-client` install the
required dependencies:

```
cd ts-client
npm install
```

The client is based on a modular architecture where you can configure a client
class to support the modules you need and instantiate it.

By default, the generated client exports a client class that includes all the
Cosmos SDK, custom and 3rd party modules in use in your project.

To instantiate the client you need to provide environment information (endpoints
and chain prefix). For querying that's all you need:

```typescript title="my-frontend-app/src/main.ts"
import { Client } from "../../ts-client";

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  }
);
```

The example above uses `ts-client` from a local directory. If you have published
your `ts-client` on `npm` replace `../../ts-client` with a package name.

The resulting client instance contains namespaces for each module, each with a
`query` and `tx` namespace containing the module's relevant querying and
transacting methods with full type and auto-completion support.

To query for a balance of an address:

```typescript
const balances = await client.CosmosBankV1Beta1.query.queryAllBalances(
  'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7'
);
```

## Broadcasting a transaction

Add signing capabilities to the client by creating a wallet from a mnemonic
(we're using the Frank's mnemonic added to `config.yml` earlier) and passing it
as an optional argument to `Client()`. The wallet implements the CosmJS
OfflineSigner` interface.

```typescript title="my-frontend-app/src/main.ts"
import { Client } from "../../ts-client";
// highlight-start
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  "play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);
// highlight-end

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  },
  // highlight-next-line
  wallet
);
```

Broadcasting a transaction:

```typescript title="my-frontend-app/src/main.ts"
const tx_result = await client.CosmosBankV1Beta1.tx.sendMsgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
```

## Broadcasting a transaction with a custom message

If your chain already has custom messages defined, you can use those. If not,
we'll be using Ignite's scaffolded code as an example. Create a post with CRUD
messages:

```
ignite scaffold list post title body
```

After adding messages to your chain you may need to re-generate the TypeScript
client:

```
ignite generate ts-client --clear-cache
```

Broadcast a transaction containing the custom `MsgCreatePost`:

```typescript title="my-frontend-app/src/main.ts"
import { Client } from "../../ts-client";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  "play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  },
  wallet
);
// highlight-start
const tx_result = await client.ExampleExample.tx.sendMsgCreatePost({
  value: {
    title: 'foo',
    body: 'bar',
    creator: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
// highlight-end
```

## Lightweight client

If you prefer, you can construct a lighter client using only the modules you are
interested in by importing the generic client class and expanding it with the
modules you need:

```typescript title="my-frontend-app/src/main.ts"
// highlight-start
import { IgniteClient } from '../../ts-client/client'
import { Module as CosmosBankV1Beta1 } from '../../ts-client/cosmos.bank.v1beta1'
import { Module as CosmosStakingV1Beta1 } from '../../ts-client/cosmos.staking.v1beta1'
// highlight-end
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)
// highlight-next-line
const Client = IgniteClient.plugin([CosmosBankV1Beta1, CosmosStakingV1Beta1])

const client = new Client(
  {
    apiURL: 'http://localhost:1317',
    rpcURL: 'http://localhost:26657',
    prefix: 'cosmos',
  },
  wallet,
)
```

## Broadcasting a multi-message transaction

You can also construct TX messages separately and send them in a single TX using
a global signing client like so:

```typescript title="my-frontend-app/src/main.ts"
const msg1 = await client.CosmosBankV1Beta1.tx.msgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
})

const msg2 = await client.CosmosBankV1Beta1.tx.msgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
})

const tx_result = await client.signAndBroadcast(
  [msg1, msg2],
  {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  '',
)
```

Finally, for additional ease-of-use, apart from the modular client mentioned
above, each generated module is usable on its own in a stripped-down way by
exposing a separate txClient and queryClient.

```typescript title="my-frontend-app/src/main.ts"
import { txClient } from '../../ts-client/cosmos.bank.v1beta1'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)

const client = txClient({
  signer: wallet,
  prefix: 'cosmos',
  addr: 'http://localhost:26657',
})

const tx_result = await client.sendMsgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
```

## Usage with Keplr

Normally, Keplr provides a wallet object implementing the `OfflineSigner`
interface, so you can simply replace the `wallet` argument in client
instantiation with `window.keplr.getOfflineSigner(chainId)`. However, Keplr
requires information about your chain, like chain ID, denoms, fees, etc.
[`experimentalSuggestChain()`](https://docs.keplr.app/api/suggest-chain.html) is
a method Keplr provides to pass this information to the Keplr extension.

The generated client makes this easier by offering a `useKeplr()` method that
automatically discovers the chain information and sets it up for you. Thus, you
can instantiate the client without a wallet and then call `useKeplr()` to enable
transacting via Keplr like so:

```typescript title="my-frontend-app/src/main.ts"
import { Client } from '../../ts-client';

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr();
```

`useKeplr()` optionally accepts an object argument that contains one or more of
the same keys as the `ChainInfo` type argument of `experimentalSuggestChain()`
allowing you to override the auto-discovered values.

For example, the default chain name and token precision (which are not recorded
on-chain) are set to `<chainId> Network` and `0` while the ticker for the denom
is set to the denom name in uppercase. If you want to override these, you can do
something like:

```typescript title="my-frontend-app/src/main.ts"
import { Client } from '../../ts-client';

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr({
  chainName: 'My Great Chain',
  stakeCurrency: {
    coinDenom: 'TOKEN',
    coinMinimalDenom: 'utoken',
    coinDecimals: '6',
  },
})
```

## Wallet switching

The client also allows you to switch out the wallet for a different one on an
already instantiated client like so:

```typescript
import { Client } from '../../ts-client';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr();

// broadcast transactions using the Keplr wallet

client.useSigner(wallet);

// broadcast transactions using the CosmJS wallet
```
