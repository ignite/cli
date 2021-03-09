# Quickstart for local Starport installation

Now that `starport` is installed on your machine, you can build your first blockchain! <!-- or do we want to write for gitpod too? -->

To build a blockchain:

```
starport app github.com/username/myapp && cd myapp
```

The `app` command creates the directory `myapp` and scaffolds a Cosmos SDK blockchain.

To run your blockchain:

```
starport serve
```

The `serve` command installs dependencies, builds, initializes, and starts your blockchain.

To update functionality for a custom type:

```
starport type post title body
```

The `type` command scaffolds functionality to create, read, update, and delete for a custom type.

<!-- the tutorial covers this quickstart command, here should we talk about a custom type and why we need it? -->
