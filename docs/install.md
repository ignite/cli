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

