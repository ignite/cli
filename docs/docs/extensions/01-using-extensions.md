---
description: Using and Developing Ignite Extensions
---

# Using Ignite Extensions

Extensions offer a way to extend the functionality of the Ignite CLI. There are two
core concepts within extensions: `Commands` and `Hooks`. `Commands` extend the CLI's
functionality and `Hooks` extend existing CLI command functionality.

Extensions are registered in an Ignite scaffolded blockchain project through the
`extensions.yml`, or globally through `$HOME/.ignite/extensions/extensions.yml`.

To use an extension within your project execute the following command inside the
project directory:

```sh
ignite ext install github.com/project/cli-extension
```

The extension will be available only when running `ignite` inside the project
directory.

To use an extension globally on the other hand, execute the following command:

```sh
ignite ext install -g github.com/project/cli-extension
```

The command will compile the extension and make it immediately available to the
`ignite` command lists.

## Listing installed Extensions

When in an ignite scaffolded blockchain you can use the command `ignite extension
list` to list all Ignite Extensions and there statuses.

## Updating Extensions

When an extension in a remote repository releases updates, running `ignite extension
update <path/to/extension>` will update an specific extension declared in your
project's `config.yml`.
