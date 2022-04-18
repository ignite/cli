---
description: Details on the Vue frontend app created by Ignite CLI.
order: 8
---

# Frontend overview

A Vue frontend app is created in the `vue` directory when a blockchain is scaffolded. To start the frontend app run `npm i && npm run serve` in the `vue` directory.

The frontend app is built using the `@ignt/vue`, `@ignt/client`, and `@ignt/plugins` packages. For details, see the [monorepo for Ignite CLI front-end development](https://github.com/ignite-hq/vue).

## Client code generation

A TypeScript (TS) client is automatically generated for your blockchain for custom and standard Cosmos SDK modules.

To enable client code generation, add the `client` entries to `config.yml`:

```yaml
client:
  typescript:
    path: "js"
```

A TS client is generated in the `vue/src/generated` directory.

## Client code regeneration

By default, the filesystem is watched and the clients are regenerated automatically. Clients for standard Cosmos SDK modules are generated after you scaffold a blockchain.

To regenerate all clients for custom and standard Cosmos SDK modules, run this command:

`ignite generate typescript`

## Preventing client code regeneration	

To prevent regenerating the client, remove the `client` property from `config.yml`.	
