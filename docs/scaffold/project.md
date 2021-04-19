---
order: 2
description: Overview of a new Cosmos SDK blockchain project built with Starport.
---

# Project Scaffold Reference

The `starport app` command scaffolds a new Cosmos SDK blockchain project.

starport app github.com/hello/planet

This command will create a directory called `planet`, which contains all the files for your project. The `github.com` URL in the argument is a string that will be used for Go module's path. The repository name (`planet`, in this case) will be used as the project's name. A git repository will be initialized locally.

The project directory structure:

* `app`: files that wire the blockchain together
* `cmd`: blockchain node's binary
* `proto`: protocol buffer files for custom modules
* `x`: directory with custom modules
* `vue`: scaffolded web application (optional)
* `config.yml`: configuration file

Most of the logic of your application-specific blockchain is written in custom modules. Each module effectively encapsulates an independent piece of functionality. Custom modules are stored inside the `x` directory. By default, `starport app` scaffolds a module with a name that matches the name of the project. In our example, it will be `x/planet`.

Every Cosmos SDK module has protocol buffer files defining data structures, messages, queries, RPCs, etc. `proto` contains a directory with proto files per each custom module in `x`.

Global changes to your blockchain are defined in files inside the `app` directory. This includes importing third-party modules, defining relationships between modules, and configuring blockchain-wide settings.

`config.yml` is a file that contains configuration options that Starport uses to build, initialize and start your blockchain node in development.

## Address prefix

Account addresses on Cosmos SDK-based blockchains have string prefixes. For example, Cosmos Hub blockchain uses the default `cosmos` prefix, so that addresses look like this: `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`. 

When creating a new blockchain, pass a prefix as a value to the `--address-prefix` flag like so:

starport app github.com/hello/planet --address-prefix moonlight

Using this prefix, account addresses on your blockchain look like this: `moonlight12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

To change the prefix after the blockchain has been scaffolded, modify the `AccountAddressPrefix` in the `app/prefix.go` file.

1. Change the `AccountAddressPrefix` variable in the `/app/prefix.go` file. Be sure to preserve other variables in the file.
2. To recognize the new prefix, change the `VUE_APP_ADDRESS_PREFIX` variable in `/vue/.env`.

## Cosmos SDK version

By default, the `starport app` command creates a Cosmos SDK blockchain using the latest stable version of the SDK.

It is possible to use the legacy Cosmos SDK v0.39.2 (Launchpad). This legacy version has no active feature development and does not support IBC. You probably don't want to create a Launchpad blockchain, but if you do, use the `--sdk-version` flag with the value `launchpad`.

```
starport app github.com/hello/planet --sdk-version launchpad
```