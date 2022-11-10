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

When using this configuration a TS client is generated in the `ts-client` directory (see: [TypeScript client information](https://docs.ignite.com/clients/typescript))
and Vuex store modules making use of this client are generated in the `vue/src/store` directory.

## Client code generation

To generate all clients for custom and standard Cosmos SDK modules, run this command:

```bash
ignite generate vuex
```

(Note: this command also runs the typescript client generation and you do not need to run `ignite generate ts-client` separately.)
