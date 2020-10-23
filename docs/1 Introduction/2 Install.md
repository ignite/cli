# Install Starport 

**Prerequisites:** If you want to install Starport locally, make sure to have [Golang >=1.14](https://golang.org/) and [Node.js >=12.19.0](https://nodejs.org/) installed on your system.

## Installation Options

### NPM

```
npm i -g @tendermint/starport
```

### macOS with Homebrew

```
brew install tendermint/tap/starport
```

<!-- ### Debian/Ubuntu with Snapcraft

```
snap install --classic node
```

Append your current working directory to the environment variable `PATH`:

```
export PATH=$PATH:$PWD/node_modules/.bin/
``` -->

### Build from source

```
git clone https://github.com/tendermint/starport && cd starport && make
```
