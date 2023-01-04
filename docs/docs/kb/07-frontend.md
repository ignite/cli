---
description: Details on the Vue frontend app created by Ignite CLI.
sidebar_position: 7
---

# Frontend overview

A Vue frontend app is created in the `vue` directory when running `ignite scaffold vue`.

To start the frontend app run `npm i && npm run dev` in the `vue` directory.

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages.
For details, see the [monorepo for Ignite CLI front-end development](https://github.com/ignite/web).

## Client code generation

To configure client code generation add the `client` settings to `config.yml`:

```yaml
client:
  typescript:
    path: ts-client
  vuex:
    path: vue/src/store
```

When using this configuration a TS client is generated in the `ts-client` directory (see: [TypeScript client information](/clients/typescript))
and Vuex store modules making use of this client are generated in the `vue/src/store` directory.

## Client code generation

To generate all clients for custom and standard Cosmos SDK modules, run this command:

```bash
ignite generate vuex
```

(Note: this command also runs the typescript client generation and you do not need to run `ignite generate ts-client` separately.)

## Setting the Correct Address Prefix

It is necessary to set the correct address prefix in order for the Vue app to properly interact with a Cosmos chain. The address prefix is used to identify the chain that the app is connected to, and must match the prefix used by the chain.

There are two ways to set the address prefix in a Vue app:

1. Using the `VITE_ADDRESS_PREFIX` environment variable:

You can set the `VITE_ADDRESS_PREFIX` environment variable to the correct address prefix for your chain. This will override the default prefix used by the app.

To set the `VITE_ADDRESS_PREFIX` environment variable, you can use the following command:

```bash
export VITE_ADDRESS_PREFIX=your-prefix
```

Replace `your-prefix` with the actual address prefix for your chain.

2. Replacing the fallback value of the `prefix` variable in the file `./vue/src/env.ts`:

Alternatively, you can manually set the correct address prefix by replacing the fallback value of the `prefix` variable in the file `./vue/src/env.ts`.

To do this, open the file `./vue/src/env.ts` and find the following line:

```js
const prefix = process.env.VITE_ADDRESS_PREFIX || 'your-prefix';
```

Replace `your-prefix` with the actual address prefix for your chain.

Save the file and restart the Vue app to apply the changes.
