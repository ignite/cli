# Frontend Overview

When scaffolding a blockchain, Starport also creates a Vue frontend application in `vue` directory. If this directory exists, `starport serve` automatically runs `npm run serve`.

The frontend app uses packages from [`tendermint/vue`](https://github.com/tendermint/vue).

## Client code generation

Starport automatically generates JavaScript, TypeScript and Vuex clients for your blockchain both for custom and standard Cosmos SDK modules. To enable client code generation add the following to `config.yml`:

```yaml
client:
  vuex:
    path: "js"
```

Starport will generate a Vuex client inside `js` directory. JS and TS clients are also generated, because they are dependencies of the Vuex client.

By default Starport watches the filesystem and regenerates clients for your custom modules automatically. Clients for standard Cosmos SDK modules are generated once, when you scaffold a blockchain. To regenerate all clients (for both custom and stadard Cosmos SDK modules), run `starport serve --reset-once --rebuild-proto-once`.