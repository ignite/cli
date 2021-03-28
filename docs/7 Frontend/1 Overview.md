# Frontend Overview

A Vue frontend app is created in `vue` directory when a blockchain is scaffolded. If the `vue` directory exists, the `starport serve` command automatically runs `npm run serve`.

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For  details, see the [monorepo for Starport front-end development](https://github.com/tendermint/vue).

## Client Code Generation

JavaScript (JS), TypeScript (TS), and Vuex clients are automatically generated for your blockchain for custom and standard Cosmos SDK modules. To enable client code generation, add the `client` entries to `config.yml`:

```yaml
client:
  vuex:
    path: "js"
```

Starport will generate a Vuex client inside `js` directory. JS and TS clients are also generated, because they are dependencies of the Vuex client.

By default Starport watches the filesystem and regenerates clients for your custom modules automatically. Clients for standard Cosmos SDK modules are generated once, when you scaffold a blockchain. To regenerate all clients (for both custom and stadard Cosmos SDK modules), run `starport serve --reset-once --rebuild-proto-once`.
