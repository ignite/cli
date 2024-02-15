<div align="center">
  <h1> Ignite </h1>
</div>

<div align="center">
  <a href="https://github.com/ignite/cli/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/cosmos/cosmos-sdk.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/ignite/cli?tab=doc">
    <img alt="GoDoc" src="https://pkg.go.dev/badge/github.com/ignite/cli.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/ignite/cli">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/ignite/cli" />
  </a>
<!--
  <a href="https://codecov.io/gh/ignite/cli">
    <img alt="Code Coverage" src="https://codecov.io/gh/ignite/cli/branch/main/graph/badge.svg" />
  </a>
-->
</div>
<div align="center">
  <a href="https://github.com/ignite/cli">
    <img alt="Lines Of Code" src="https://tokei.rs/b1/github/ignite/cli" />
  </a>
    <img alt="Test Status" src="https://github.com/ignite/cli/workflows/Test/badge.svg" />
    <img alt="Lint Status" src="https://github.com/ignite/cli/workflows/Lint/badge.svg" />
</div>

![Ignite CLI](./assets/ignite-cli.png)

[Ignite CLI](https://ignite.com/cli) is the all-in-one platform to build,
launch, and maintain any crypto application on a sovereign and secured
blockchain. It is a developer-friendly interface to the [Cosmos
SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used
blockchain application framework. Ignite CLI generates boilerplate code for you,
so you can focus on writing business logic.

## Quick start

<<<<<<< HEAD
Open Ignite CLI [in your web
browser](https://gitpod.io/#https://github.com/ignite/cli/tree/v0.25.2) (or open
[nightly version](https://gitpod.io/#https://github.com/ignite/cli)), or
[install the latest release](https://docs.ignite.com/welcome/install).
=======
### Installation

You can install Ignite using [HomeBrew](https://formulae.brew.sh/formula/ignite) on macOS and GNU/Linux:

```sh
brew install ignite
```

Or using Snap on GNU/Linux:

```sh
snap install ignite --classic
```

Or manually using the following command:

```sh
curl https://get.ignite.com/cli! | bash
```

<details>
  <summary>Troubleshoot</summary>

If Ignite doesn't automatically move to your `/usr/local/bin` directory, use this command to do so:

```sh
sudo mv ignite /usr/local/bin
```

If you encounter an error, you might need to create the `/usr/local/bin` directory and set the necessary permissions as follows:

```sh
mkdir /usr/local/bin
sudo chown -R $(whoami) /usr/local/bin
```

</details>

For more options on installing and using Ignite, please see the following:

Open Ignite CLI [in your web browser](https://gitpod.io/#https://github.com/ignite/cli/tree/v28.1.1) (or open [nightly version](https://gitpod.io/#https://github.com/ignite/cli)), or [install the latest release](https://docs.ignite.com/welcome/install).
>>>>>>> 17f1a763 (chore: add homebrew formula (#3950))

To create and start a blockchain:

```bash
ignite scaffold chain mars

cd mars

ignite chain serve
```

## Documentation

To learn how to use Ignite CLI, check out the [Ignite CLI
docs](https://docs.ignite.com). To learn more about how to build blockchain apps
with Ignite CLI, see the [Ignite CLI Developer
Tutorials](https://docs.ignite.com/guide).

To install Ignite CLI locally on GNU, Linux, or macOS, see [Install Ignite
CLI](https://docs.ignite.com/welcome/install).

To learn more about building a JavaScript frontend for your Cosmos SDK
blockchain, see [ignite/web](https://github.com/ignite/web).

## Questions

For questions and support, join the official [Ignite
Discord](https://discord.gg/ignite) server. The issue list in this repo is
exclusively for bug reports and feature requests.

## Cosmos SDK compatibility

Blockchains created with Ignite CLI use the [Cosmos
SDK](https://github.com/cosmos/cosmos-sdk) framework. To ensure the best
possible experience, use the version of Ignite CLI that corresponds to the
version of Cosmos SDK that your blockchain is built with. Unless noted
otherwise, a row refers to a minor version and all associated patch versions.

| Ignite CLI  | Cosmos SDK  | IBC                  | Notes                                                         |
|-------------|-------------|----------------------|---------------------------------------------------------------|
| v28.x.x     | v0.50.x     | v8.0.0               | -                                                             |
| v0.27.1     | v0.47.3     | v7.1.0               | -                                                             |
| v0.26.0     | v0.46.7     | v6.1.0               | -                                                             |
| v0.25.2     | v0.46.6     | v5.1.0               | Bump Tendermint version to v0.34.24                           |
| v0.25.1     | v0.46.3     | v5.0.0               | Includes  Dragonberry security fix                            |
| ~~v0.24.0~~ | ~~v0.46.0~~ | ~~v5.0.0~~           | This version is deprecated due to a security fix in `v0.25.0` |
| v0.23.0     | v0.45.5     | v3.0.1               |                                                               |
| v0.21.1     | v0.45.4     | v2.0.3               | Supports Cosmos SDK v0.46.0-alpha1 and above                  |
| v0.21.0     | v0.45.4     | v2.0.3               |                                                               |
| v0.20.0     | v0.45.3     | v2.0.3               |                                                               |
| v0.19       | v0.44       | v1.2.2               |                                                               |
| v0.18       | v0.44       | v1.2.2               | `ignite chain serve` works with v0.44.x chains                |
| v0.17       | v0.42       | Same with Cosmos SDK |                                                               |

To upgrade your blockchain to the newer version of Cosmos SDK, see the
[Migration guide](https://docs.ignite.com/migration).

## Plugin system

Ignite CLI commands can be extended using plugins. A plugin is a program that
uses github.com/hashicorp/go-plugin to communicate with the ignite binary.

#### Use a plugin

Plugins must be declared in the `config.yml` file, using the following syntax:

```yaml
plugins:
  // path can be a repository or a local path
  // the directory must contain go code under a main package.
  // For repositories you can specify a suffix @branch or @tag to target a
  // specific git reference.
  - path: github.com/org/repo/my-plugin
    // Additional parameters can be passed to the plugin
    with:
      key: value
```

Once declared, the next time the ignite binary will be executed under this
configuration, it will fetch, build and run the plugin. As a result, more
commands should be available in the list of the ignite commands.

`ignite plugin` command allows to list the plugins and their status, and to
update a plugin if you need to get the latest version.

### Make a plugin

A plugin must implement `plugin.Interface`.

The easiest way to make a plugin is to use the `ignite plugin scaffold` command.
For example:

```
$ cd /home/user/src
$ ignite plugin scaffold github.com/foo/bar
```

It will create a folder `bar` under `/home/user/src` and generate predefined
`go.mod` and `main.go`. The code contains everything required to connect to the
ignite binary via `hashicorp/go-plugin`. What need to be adapted is the
implementation of the `plugin.Interface` (`Commands` and `Execute` methods).

To test your plugin, you only need to declare it under a chain config, for
instance:

```yaml
plugins:
  - path: /home/user/src/bar
```

Then run `ignite`, the plugin will compile and should be listed among the ignite
commands. Each time `ignite` is executed, the plugin is recompiled if the files
have changed since the last compilation. This allows fast and easy plugin
development, you only care about code and `ignite` handles the compilation.

## Contributing

We welcome contributions from everyone. The `main` branch contains the
development version of the code. You can create a branch from `main` and
create a pull request, or maintain your own fork and submit a cross-repository
pull request.

Our Ignite CLI bounty program provides incentives for your participation and
pays rewards. Track new, in-progress, and completed bounties on the [Bounty
board](https://github.com/ignite/cli/projects/5) in GitHub.

**Important** Before you start implementing a new Ignite CLI feature, the first
step is to create an issue on GitHub that describes the proposed changes.

If you're not sure where to start, check out [contributing.md](contributing.md)
for our guidelines and policies for how we develop Ignite CLI. Thank you to
everyone who has contributed to Ignite CLI!

## Community

Ignite CLI is a free and open source product maintained by
[Ignite](https://ignite.com). Here's where you can find us. Stay in touch.

* [ignite.com website](https://ignite.com)
* [@ignite\_dev on Twitter](https://twitter.com/ignite_dev)
* [ignite.com/blog](https://ignite.com/blog)
* [Ignite Discord](https://discord.com/invite/ignite)
* [Ignite YouTube](https://www.youtube.com/@ignitehq)
* [Ignite docs](https://docs.ignite.com)
* [Ignite jobs](https://ignite.com/careers)
