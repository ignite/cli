# Ignite CLI

![Ignite CLI](./assets/ignite-cli.png)

[Ignite CLI](https://ignite.com/cli) is the all-in-one platform to build, launch, and maintain any crypto application on a sovereign and secured blockchain. It is a developer-friendly interface to the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used blockchain application framework. Ignite CLI generates boilerplate code for you, so you can focus on writing business logic.

## Quick start

Open Ignite CLI [in your web browser](https://gitpod.io/#https://github.com/ignite/cli/tree/master) (or open [nightly version](https://gitpod.io/#https://github.com/ignite/cli/)), or [install the latest release](https://docs.ignite.com/guide/install.html). 

To create and start a blockchain:

```bash
ignite scaffold chain mars

cd mars

ignite chain serve
```

## Documentation

To learn how to use Ignite CLI, check out the [Ignite CLI docs](https://docs.ignite.com). To learn more about how to build blockchain apps with Ignite CLI, see the [Ignite CLI Developer Tutorials](https://docs.ignite.com/guide/). 

To install Ignite CLI locally on GNU, Linux, or macOS, see [Install Ignite CLI](https://docs.ignite.com/guide/install.html).

To learn more about building a JavaScript frontend for your Cosmos SDK blockchain, see [ignite/web](https://github.com/ignite/web).

## Questions

For questions and support, join the official [Ignite Discord](https://discord.gg/ignite) server. The issue list in this repo is exclusively for bug reports and feature requests.

## Cosmos SDK compatibility

Blockchains created with Ignite CLI use the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) framework. To ensure the best possible experience, use the version of Ignite CLI that corresponds to the version of Cosmos SDK that your blockchain is built with. Unless noted otherwise, a row refers to a minor version and all associated patch versions.

| Ignite CLI | Cosmos SDK | IBC                  | Notes |
| -------- | ---------- | -------------------- | ------------------------------------------------ |
| v0.21.1  | v0.45.4    | v2.0.3               | Supports Cosmos SDK v0.46.0-alpha1 and above |
| v0.21.0  | v0.45.4    | v2.0.3               | |
| v0.20.0  | v0.45.3    | v2.0.3              | |
| v0.19    | v0.44      | v1.2.2               | |
| v0.18    | v0.44      | v1.2.2               | `ignite chain serve` works with v0.44.x chains |
| v0.17    | v0.42      | Same with Cosmos SDK | |

To upgrade your blockchain to the newer version of Cosmos SDK, see the [Migration guide](https://docs.ignite.com/migration/).

## Contributing

We welcome contributions from everyone. The `develop` branch contains the development version of the code. You can create a branch from `develop` and create a pull request, or maintain your own fork and submit a cross-repository pull request.

Our [Ignite CLI bounty program](docs/bounty/index.md) provides incentives for your participation and pays rewards. Track new, in-progress, and completed bounties on the [Bounty board](https://github.com/ignite/cli/projects/5) in GitHub.

**Important** Before you start implementing a new Ignite CLI feature, the first step is to create an issue on GitHub that describes the proposed changes.

If you're not sure where to start, check out [contributing.md](contributing.md) for our guidelines and policies for how we develop Ignite CLI. Thank you to everyone who has contributed to Ignite CLI!

## Community

Ignite CLI is a free and open source product maintained by [Ignite](https://ignite.com). Here's where you can find us. Stay in touch.

- [ignite.com website](https://ignite.com)
- [@ignite_dev on Twitter](https://twitter.com/ignite_dev)
- [ignite.com/blog](https://ignite.com/blog/)
- [Ignite Discord](https://discord.com/invite/ignite)
- [Ignite YouTube](https://www.youtube.com/ignitehq)
- [Ignite docs](https://docs.ignite.com/)
- [Ignite jobs](https://ignite.com/careers)
