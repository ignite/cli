---
sidebar_position: 1
description: Information about the generated Typescript client code.
---

# Typescript code generation

The `ignite generate ts-client` command generates a Typescript client for your blockchain project.

## Client code generation

A TypeScript (TS) client is automatically generated for your blockchain for custom and standard Cosmos SDK modules.

To enable client code generation, add the `client` entries to `config.yml`:

```yaml
client:
  typescript:
    path: "ts-client"
```

A TS client is generated in the `ts-client` directory.

## Client code regeneration

By default, the filesystem is watched and the clients are regenerated automatically. Clients for standard Cosmos SDK modules are generated after you scaffold a blockchain.

To regenerate all clients for custom and standard Cosmos SDK modules, run this command:

```bash
ignite generate ts-client
```

## Preventing client code regeneration	

To prevent regenerating the client, remove the `client:typescript` property from `config.yml`.	

## Setup

The best way to get started building with the TypeScript client is by using a [Vite](https://vitejs.dev) boilerplate. Vite provides boilerplates for vanilla TS projects as well as react, vue, lit, svelte and preact frameworks.
You can find additional information at the [Vite Getting Started guide](https://vitejs.dev/guide).

You will also need to polyfill the client's dependencies. The following is an example of setting up a vanilla TS project with the necessary polyfills.

```bash
npm create vite@latest my-frontend-app -- --template vanilla-ts
npm install --save-dev @esbuild-plugins/node-globals-polyfill @rollup/plugin-node-resolve
```

You must then create the necessary `vite.config.ts` file.

```typescript
import { nodeResolve } from '@rollup/plugin-node-resolve'
import { NodeGlobalsPolyfillPlugin } from '@esbuild-plugins/node-globals-polyfill'
import { defineConfig } from 'vite'

export default defineConfig({
  
	plugins: [nodeResolve()],

	optimizeDeps: {
		esbuildOptions: {
			define: {
				global: 'globalThis',
			},
			plugins: [
				NodeGlobalsPolyfillPlugin({
					buffer:true
				}),
			],
		},
	}
})
```

You are then ready to use the generated client code inside this project directly or by publishing the client and installing it as any other npm package.

## Usage

The code generated in `ts-client` comes with a `package.json` file ready to publish which you can modify to suit your needs.

The client is based on a modular architecture where you can configure a client class to support the modules you need and instantiate it.

By default, the generated client exports a client class that includes all the Cosmos SDK, custom and 3rd party modules in use in your project.

To instantiate the client you need to provide environment information (endpoints and chain prefix) and an optional wallet (implementing the CosmJS OfflineSigner interface).

For example, to connect to a local chain instance running under the Ignite CLI defaults, using a CosmJS wallet:

```typescript
import { Client } from '<path-to-ts-client>';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic = "surround miss nominee dream gap cross assault thank captain prosper drop duty group candy wealth weather scale put";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	},
	wallet
);
```

The resulting client instance contains namespaces for each module, each with a `query` and `tx` namespace containing the module's relevant querying and transacting methods with full type and auto-completion support. 

e.g.

```typescript
const balances = await client.CosmosBankV1Beta1.query.queryAllBalances('cosmos1qqqsyqcyq5rqwzqfys8f67');
```

And for transactions:

```typescript
const tx_result = await client.CosmosBankV1Beta1.tx.sendMsgSend(
	{ 
		value: {
			amount: [
				{
					amount: '200',
					denom: 'token',
				},
			],
			fromAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67',
			toAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67'
		},
		fee,
		memo
	}
);
```

If you prefer, you can construct a lighter client using only the modules you are interested in by importing the generic client class and expanding it with the modules you need:

```typescript
import { IgniteClient } from '<path-to-ts-client>/client';
import { Module as CosmosBankV1Beta1 } from '<path-to-ts-client>/cosmos.bank.v1beta1'
import { Module as CosmosStakingV1Beta1 } from '<path-to-ts-client>/cosmos.staking.v1beta1'
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic = "surround miss nominee dream gap cross assault thank captain prosper drop duty group candy wealth weather scale put";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);
const CustomClient = IgniteClient.plugin([CosmosBankV1Beta1, CosmosStakingV1Beta1]);

const client = new CustomClient({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	},
	wallet
);
```

You can also construct TX messages separately and send them in a single TX using a global signing client like so:

```typescript
const msg1 = await client.CosmosBankV1Beta1.tx.msgSend(
	{ 
		value: {
			amount: [
				{
					amount: '200',
					denom: 'token',
				},
			],
			fromAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67',
			toAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67'
		}
	}
);
const msg2 = await client.CosmosBankV1Beta1.tx.msgSend(
	{ 
		value: {
			amount: [
				{
					amount: '200',
					denom: 'token',
				},
			],
			fromAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67',
			toAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67'
		},
	}
);
const tx_result = await client.signAndBroadcast([msg1,msg2], fee, memo);
```

Finally, for additional ease-of-use, apart from the modular client mentioned above, each generated module is usable on its own in a stripped-down way by exposing a separate txClient and queryClient.

e.g.

```typescript
import { queryClient } from '<path-to-ts-client>/cosmos.bank.v1beta1';

const client = queryClient({ addr: 'http://localhost:1317' });
const balances = await client.queryAllBalances('cosmos1qqqsyqcyq5rqwzqfys8f67');
```

and

```typescript
import { txClient } from '<path-to-ts-client>/cosmos.bank.v1beta1';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic = "surround miss nominee dream gap cross assault thank captain prosper drop duty group candy wealth weather scale put";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = txClient({
	signer: wallet,
	prefix: 'cosmos',
	addr: 'http://localhost:26657'
});

const tx_result = await client.sendMsgSend(
	{ 
		value: {
			amount: [
				{
					amount: '200',
					denom: 'token',
				},
			],
			fromAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67',
			toAddress: 'cosmos1qqqsyqcyq5rqwzqfys8f67'
		},
		fee,
		memo
	}
);
```

## Usage with Keplr

Normally, Keplr provides a wallet object implementing the OfflineSigner interface so you can simply replace the wallet argument in client instantiation with it like so:


```typescript
import { Client } from '<path-to-ts-client>';

const chainId = 'mychain-1'
const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	},
	window.keplr.getOfflineSigner(chainId)
);
```

The problem is that for a new Ignite CLI scaffolded chain, Keplr has no knowledge of it thus requiring an initial call to [`experimentalSuggestChain()`](https://docs.keplr.app/api/suggest-chain.html) method to add the chain information to the user's Keplr instance.

The generated client makes this easier by offering a `useKeplr()` method that autodiscovers the chain information and sets it up for you. Thus you can instantiate the client without a wallet and then call `useKeplr()` to enable transacting via Keplr like so:

```typescript
import { Client } from '<path-to-ts-client>';

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr();
```

`useKeplr()` optionally accepts an object argument that contains one or more of the same keys as the `ChainInfo` type argument of `experimentalSuggestChain()` allowing you to override the auto-discovered values.

For example, the default chain name and token precision (which are not recorded on-chain) are set to `<chainId> Network` and `0` while the ticker for the denom is set to the denom name in uppercase. If you wanted to override these, you could do something like:


```typescript
import { Client } from '<path-to-ts-client>';

const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr({ chainName: 'My Great Chain', stakeCurrency : { coinDenom: 'TOKEN', coinMinimalDenom: 'utoken', coinDecimals: '6' } });
```

## Wallet switching

The client also allows you to switch out the wallet for a different one on an already instantiated client like so:

```typescript
import { Client } from '<path-to-ts-client>';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic = "surround miss nominee dream gap cross assault thank captain prosper drop duty group candy wealth weather scale put";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);


const client = new Client({ 
		apiURL: "http://localhost:1317",
		rpcURL: "http://localhost:26657",
		prefix: "cosmos"
	}
);
await client.useKeplr();

// transact using Keplr Wallet

client.useSigner(wallet);

//transact using CosmJS wallet
```
