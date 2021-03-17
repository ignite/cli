# Install Starport on a local computer

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

As a developer, you will most likely want to install Starport to work in a local development environment.

## Prerequisites

Local Starport installation requires the follow software be installed and running:

- [Golang >=1.16](https://golang.org/)

  Make sure that your `$PATH` includes `$HOME/go/bin` with `export PATH=$PATH:$HOME/go/bin`.

- [Node.js >=12.19.0](https://nodejs.org/)

- [Protocol Buffer compiler](https://grpc.io/docs/protoc-installation/)


## Installing Starport with cURL

To install Starport with cURL:

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`.

To install previous versions of the precompiled `starport` binary, see [star-port installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### Write permission

Starport installation requires write permission to the `/usr/local/bin/` directory. If the installation fails because you do not have write permission to `/usr/local/bin/`, run the following command:

```
curl https://get.starport.network/starport | bash
```

Then run this command to move the `starport` executable to `/usr/local/bin/`:

```
sudo mv starport /usr/local/bin/
```

## Installing Starport on macOS with Homebrew

To install Starport on macOS with Homebrew:

```
brew install tendermint/tap/starport
```

## Build from source

You can build and install the precompiled `starport` binary into `$GOBIN`.

When building from source, you must set your `$GOPATH` environment variable.

To set `$GOPATH` environment variable:

```
mkdir ~/go
export GOPATH=~/go
```

To install the precompiled `starport` binary, run this command:

```
git clone https://github.com/tendermint/starport --depth=1
cd starport && make
```

## Summary

- To setup a local development environment, install Starport locally on your computer.
- Install Starport by fetching the binary using cURL, Homebrew, or by building from source.
- The latest version is installed by default. You can install previous versions of the precompiled `starport` binary.
