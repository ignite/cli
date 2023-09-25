---
description: Using and Developing Ignite Apps
---

# Using Ignite Apps

Apps offer a way to extend the functionality of the Ignite CLI. There are two
core concepts within apps: `Commands` and `Hooks`. `Commands` extend the CLI's
functionality and `Hooks` extend existing CLI command functionality.

Apps are registered in an Ignite scaffolded blockchain project through the
`plugins.yml`, or globally through `$HOME/.ignite/plugins/plugins.yml`.

To use an app within your project execute the following command inside the
project directory:

```sh
ignite app install github.com/project/cli-app
```

The app will be available only when running `ignite` inside the project
directory.

To use an app globally on the other hand, execute the following command:

```sh
ignite app install -g github.com/project/cli-app
```

The command will compile the app and make it immediately available to the
`ignite` command lists.

## Listing installed apps

When in an ignite scaffolded blockchain you can use the command `ignite app
list` to list all Ignite Apps and there statuses.

## Updating apps

When an app in a remote repository releases updates, running `ignite app
update <path/to/app>` will update an specific app declared in your
project's `config.yml`.
