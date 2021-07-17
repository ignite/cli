---
order: 4
description: Generate code for new modules in your blockchain.
---

# Module Scaffold

Modules are the building blocks of Cosmos SDK blockchains. Modules encapsulate logic and allow sharing functionality between projects. To learn about building modules, see the [Introduction to SDK Modules](https://docs.cosmos.network/master/building-modules/index.html).

## Create Modules

Starport supports scaffolding Cosmos SDK modules.

```
starport scaffold module create [name] [flags]
```

The `name` parameter is the name of your new module. The module name must be unique within a project.

## Files and Directories

The following files and directories are created and modified by scaffolding a module:

- `proto/`: a directory that contains proto files for query and message services.
- `x`: common logic for a module.
- `app/app.go`: imports and initializes your module.

## Enable IBC Logic

The Inter-Blockchain Communication protocol (IBC) is an important part of the Cosmos SDK ecosystem.

To include all the logic for the scaffolded IBC modules, use `--ibc` flag.

## Create Module Example

This example command scaffolds an -IBC-enable module named blog.

```bash
starport scaffold module create blog --ibc
```

