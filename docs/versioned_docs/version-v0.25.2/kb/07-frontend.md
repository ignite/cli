---
description: Details on the Vue frontend app created by Ignite CLI.
sidebar_position: 7
---

# Frontend overview

A Vue frontend app is created in the `vue` directory when a blockchain is scaffolded. To start the frontend app run `npm i && npm run dev` in the `vue` directory.

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite CLI front-end development](https://github.com/ignite/web).

## Client code generation

A TypeScript (TS) client and associated Vuex stores are automatically generated for your blockchain for custom and standard Cosmos SDK modules.

To enable client code generation, add the `client` entries to `config.yml`:

```yaml
client:
  typescript:
    path: "ts-client"
  vuex:
    path: "vue/src/store"
```

A TS client is generated in the `ts-client` directory (see: [TypeScript client information](/clients/typescript)) and Vuex store modules making use of this client are generated in the `vue/src/store` directory.

## Client code regeneration

By default, the filesystem is watched and the clients are regenerated automatically. Clients for standard Cosmos SDK modules are generated after you scaffold a blockchain.

To regenerate all clients for custom and standard Cosmos SDK modules, run this command:

```bash
ignite generate vuex
```

(Note: this command also runs the typescript client generation and you do not need to run `ignite generate ts-client` separately.)
## Preventing client code regeneration	

To prevent regenerating the client, remove the `client:vuex` property from `config.yml`.	
