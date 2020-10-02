# Install Starport 

If you want to install Starport locally, make sure to have Golang >=1.14 and a recent version of Node.js installed on your system.

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

Scaffold your application - [more](02%20Using%20Starport/01_using_starport/01_using_starport.md#your-blockchain-application)

## Start Hacking

```bash
starport app github.com/username/myapp
```

Serve the blockchain application - [more](02%20Using%20Starport/01_using_starport/01_using_starport.md#serve)

```bash
starport serve
```

Add a new transaction type to your application - [more](02%20Using%20Starport/01_using_starport/01_using_starport.md#how-to-use-types)

```bash
starport type post title body
```

## Learn More

To learn how to use Starport, continue to the [Starport Handbook](01%20Introduction/01_starport_introduction/introduction.md).
