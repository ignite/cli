# Install Starport 

**Prerequisites:** If you want to install Starport locally, make sure to have [Golang >=1.14](https://golang.org/). The latest version of Starport also requires [Protocol Buffer compiler](https://grpc.io/docs/protoc-installation/) to be installed. [Node.js >=12.19.0](https://nodejs.org/) is used to build the welcome screen, block explorer and to run the web scaffold.

To install Starport:

```
curl https://get.starport.network/starport! | bash
```

This command will download the latest `starport` binary from Github and install it into `/usr/local/bin`. To learn more about how to install previous versions of the binary, refer to the [documentation](https://github.com/allinbits/starport-installer).

## macOS with Homebrew

```
brew install tendermint/tap/starport
```

## Build from source

```
git clone https://github.com/tendermint/starport && cd starport && make
```

This will build and install `starport` binary into `$GOBIN`.

Note: When building from source, it is important to have your `$GOPATH` set correctly.  When in doubt, the following should do:

```
mkdir ~/go
export GOPATH=~/go
```
