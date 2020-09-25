# Starport

Starport is the easiest way to build blockchains. It is a developer-friendly interface to the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used blockchain application framework. Starport generates boilerplate code for you, so you can focus on writing business logic.

![Banner](./assets/banner.jpeg)

Overview: https://www.youtube.com/watch?v=rmbPjCGDXek

## Install

### In browser

➡️ Check out <a href="https://gitpod.io/#https://github.com/tendermint/starport/" target="_blank">Starport in browser-based IDE</a>. Fastest way to get started! `starport` binary is already installed, just create an application and start hacking!

### NPM

```
npm i -g @tendermint/starport
```

### macOS with Homebrew

```
brew install tendermint/tap/starport
```

### Debian/Ubuntu with Snapcraft

```
snap install --classic node
```

Append your current working directory to the environment variable `PATH`:

```
export PATH=$PATH:$PWD/node_modules/.bin/
```

### Build from source

```
git clone https://github.com/tendermint/starport && cd starport && make
```

Requirements: Go 1.14 and Node.js (optional but highly recommended, used for the welcome screen and web app scaffolding).

## Documentation

The documentation can be found in the [`/docs`](/docs/README.md) directory [here](/docs/README.md).

0. [Quickstart](docs/README.md#quickstart-)
1. [Introduction](docs/01%20Introduction/README.md)     
2. [Using Starport](docs/02%20Using%20Starport/README.md)    
3. [Modules](docs/03%20Modules/README.md)  
4. [Use Cases](docs/04%20Use%20Cases/README.md)  
5. [Extras](docs/05%20Extras/README.md)

## More tutorials

- [Blog (video) tutorial](https://www.youtube.com/watch?v=rmbPjCGDXek): get started with your first blockchain
- [Poll tutorial](https://tutorials.cosmos.network/starport-polling-app/): build a voting application with a web-based UI
- [Smart contract tutorial](https://www.notion.so/Smart-contracts-with-CosmWasm-c6fbcd584b78437a843e738b922dc108): add smart contracts to your app with CosmWasm: build, upload, instantiate and run a smart contract
- [Blog (from scratch) tutorial](https://tutorials.cosmos.network/starport-blog/01-index.html): learn how Starport works by building a blog without scaffolding

## Questions & comments

For questions and support please join the #starport channel in the [Cosmos Community Discord](https://discord.com/invite/W8trcGV). The issue list of this repo is exclusively for bug reports and feature requests.

## Contributing

`develop` contains the development version. Find the last stable release under https://github.com/tendermint/starport/releases.

You can branch of from `develop` and create a Pull Request or maintain your own fork and submit a Pull Request from there.

## Stay in touch

Starport is a product built by [Tendermint](https://tendermint.com). Follow us to get the latest updates!

- [Twitter](https://twitter.com/tendermint_team)
- [Blog](https://medium.com/tendermint)
- [Jobs](https://tendermint.com/careers)
