# Welcome to Starport ✨

Starport is an easy-to-use CLI tool for creating sovereign blockchains. Blockchains created with Starport use Cosmos SDK, the world's most widely used blockchain application framework.

In this browser-based development environment, the terminal window is in the lower part of the window. The `starport` binary is pre-installed and ready to use on the command line.

## Quick start

To create a blockchain and start a node in development:

```
starport scaffold chain github.com/user/hello

cd hello

starport chain serve
```

## Next steps

📺 **[Introduction to Starport](https://www.youtube.com/watch?v=5RqAIE0b8Kw)**: Watch an introductory video to learn about Starport.

🧑‍🏫 **[Tutorials](https://tutorials.cosmos.network)**: Learn by building a simple IBC-enabled module, nameservice, or a decentralized exchange (DEX).

📕 **[Starport Documentation](https://docs.starport.network)**: Explore the features of Starport.

📚 [Cosmos SDK Documentation](https://docs.cosmos.network): Learn about the framework for building application-specific blockchains.

⭐️ [Starport on Github](https://github.com/tendermint/starport): Submit an issue or contribute to the source code.

## Starport features

* Scaffold modules, messages, types with CRUD operations, IBC packets, and more
* Start a blockchain node in development with live reloading
* Connect to other blockchains with a built-in IBC relayer
* Use automatically generated TypeScript/Vuex clients to interact with your blockchain
* Use the Vue.js web app template with a set of components and Vuex modules

## Install Starport locally

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`. Learn more about [installing Starport](https://docs.starport.network/intro/install).

## Stay in touch

Starport is a free and open source product maintained by [Tendermint](https://tendermint.com). Follow us on [Twitter](https://twitter.com/tendermint_team) and [Medium](https://medium.com/tendermint) to get the latest updates!
