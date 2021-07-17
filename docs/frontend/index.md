---
order: 1
description: Details on the Vue frontend app created by Starport.
parent:
  order: 7
  title: Frontend
---

# Frontend Overview

A Vue frontend app is created in the `vue` directory when a blockchain is scaffolded. To start the frontend app run `npm i && npm run serve` in the `vue` directory.

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Starport front-end development](https://github.com/tendermint/vue).

## Client Code Generation

JavaScript (JS), TypeScript (TS), and Vuex clients are automatically generated for your blockchain for custom and standard Cosmos SDK modules. 

To enable client code generation, add the `client` entries to `config.yml`:

```yaml
client:
  vuex:
    path: "js"
```

A Vuex client is generated in the `js` directory. JS and TS clients are also generated because they are dependencies of the Vuex client.

## Client Code Regeneration

By default, the filesystem is watched and the clients are regenerated automatically. Clients for standard Cosmos SDK modules are generated after you scaffold a blockchain.

To regenerate all clients for custom and standard Cosmos SDK modules, run this command:

`starport chain serve --reset-once --rebuild-proto-once` 
