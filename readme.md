# Starport

![Starport](./assets/starport.png)

[Starport](https://starport.com) is the all-in-one platform to build, launch, and maintain any crypto application on a sovereign and secured blockchain. It is a developer-friendly interface to the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used blockchain application framework. Starport generates boilerplate code for you, so you can focus on writing business logic.

* [**Build a blockchain with Starport in a web-based IDE** (stable)](https://gitpod.io/#https://github.com/tendermint/starport/tree/master) or use [nightly version](https://gitpod.io/#https://github.com/tendermint/starport/)
* [Check out the latest features in v0.19](https://starport.com/blog/starport-v0-19)

## Quick start

Open Starport [in your browser](https://gitpod.io/#https://github.com/tendermint/starport/tree/master), or [install it](https://docs.starport.network/guide/install.html). Create and start a blockchain:

```bash
starport scaffold chain github.com/cosmonaut/mars

cd mars

starport chain serve
```

## Documentation

To learn how to use Starport, check out the [Starport Documentation](https://docs.starport.com). To learn more about how to build blockchain apps with Starport, see the [Starport Developer Tutorials](https://docs.starport.com/guide/). 

To install Starport locally on GNU, Linux, or macOS, see [Install Starport](https://docs.starport.com/guide/install.html).

To learn more about building a JavaScript frontend for your Cosmos SDK blockchain, see [tendermint/vue](https://github.com/tendermint/vue).

## Questions

For questions and support, join the official [Starport Discord](https://discord.gg/starport) server. The issue list in this repo is exclusively for bug reports and feature requests.

## Cosmos SDK Compatibility

Blockchains created with Starport use the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) framework. To ensure the best possible experience, use the version of Starport that corresponds to the version of Cosmos SDK that your blockchain is built with. Unless noted otherwise, a row refers to a minor version and all associated patch versions.

| Starport | Cosmos SDK | Notes                                            |
| -------- | ---------- | ------------------------------------------------ |
| v0.19    | v0.44      | |
| v0.18    | v0.44      | `starport chain serve` works with v0.44.x chains |
| v0.17    | v0.42      | |

To upgrade your blockchain to the newer version of Cosmos SDK, see the [Migration guide](https://docs.starport.com/migration/).

## Plugin system

Starport commands can be extended using plugins. A plugin is a program that
uses github.com/hashicorp/go-plugin to communicate with the starport binary.

#### Use a plugin

Plugins must be declared in the `config.yml` file, using the following syntax:

```yaml
plugins:
  - name: my-plugin
    // path can be a repository or a local path
    // the directory must contain go code under a main package.
    // For repositories you can specify a suffix @branch or @tag to target a
    // specific git reference.
    path: github.com/org/repo/my-plugin
    // Additional parameters can be passed to the plugin
    with:
      key: value
```

Once declared, the next time the starport binary will be executed under this
configuration, it will fetch, build and run the plugin. As a result, more
commands should be available in the list of the starport commands.

`starport plugin` command allows to list the plugins and their status, and to
update a plugin if you need to get the latest version.

### Make a plugin

A plugin must implement `plugin.Interface`.

The easiest way to make a plugin is to use the `starport plugin scaffold`
command. For example:

```
$ cd /home/user/src
$ starport plugin scaffold github.com/foo/bar
```

It will create a folder `bar` under `/home/user/src` and generate predefined
`go.mod` and `main.go`. The code contains everything required to connect to the
starport binary via `hashicorp/go-plugin`. What need to be adapted is the
implementation of the `plugin.Interface` (`Commands` and `Execute` methods).

To test your plugin, you only need to declare it under a chain config, for
instance:

```yaml
plugins:
  - name: bar
    path: /home/user/src/bar
```

Then run `starport`, the plugin will compile and should be listed among the
starport commands. Each time `starport` is executed, the plugin is recompiled
if the files have changed since the last compilation. This allows fast and easy plugin
development, you only care about code and `starport` handles the compilation.

## Contributing

We welcome contributions from everyone. The `develop` branch contains the development version of the code. You can create a branch from `develop` and create a pull request, or maintain your own fork and submit a cross-repository pull request. 

**Important** Before you start implementing a new Starport feature, the first step is to create an issue on Github that describes the proposed changes.

If you're not sure where to start, check out [contributing.md](contributing.md) for our guidelines and policies for how we develop Starport. Thank you to everyone who has contributed to Starport!

## Community

Starport is a free and open-source product maintained by [Tendermint](https://tendermint.com). Here's where you can find us. Stay in touch.

- [Starport.com website](https://starport.com)
- [@StarportHQ on Twitter](https://twitter.com/StarportHQ)
- [Starport.com/blog](https://starport.com/blog/)
- [Starport Discord](https://discord.com/starport)
- [Starport YouTube](https://www.youtube.com/channel/UCXMndYLK7OuvjvElSeSWJ1Q)
- [Starport docs](https://docs.starport.com/)
- [Tendermint jobs](https://tendermint.com/careers)
