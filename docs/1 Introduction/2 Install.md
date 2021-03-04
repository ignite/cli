# Install Starport on a local computer

You can run Starport in the GitPod IDE or you can install Starport on your local computer.

## Prerequisites

Local Starport installation requires the follow software be installed and running:

- [Golang >=1.14](https://golang.org/)

  Starport and Cosmos SDK modules are written in Go. Version 1.14 or later is required.

- [Protocol Buffer compiler](https://grpc.io/docs/protoc-installation/)

  <!-- purpose of this compiler? -->
- [Node.js >=12.19.0](https://nodejs.org/)

  Node.js is used to build the welcome screen, block explorer, and run the web scaffold. Version 12.1.9.0 or later is required.

## Installing Starport

To install Starport, run this command:

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`.

To install previous versions of the precompiled binary, see [star-port installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### /usr/local/bin/ requires user write permission

The installation requires user permission to write to `/usr/local/bin/`. If the installation fails because the user does not have write permission to `/usr/local/bin/`, run the following command:

```
curl https://get.starport.network/starport | bash
```

To move the `starport` executable to `/usr/local/bin/`, run this command:

```
sudo mv starport /usr/local/bin/
```

<!-- per <https://github.com/allinbits/starport-installer/blob/master/README.md> installing with Homebrew is not supported, so let's comment out from the doc ## macOS with Homebrew ``` brew install tendermint/tap/starport ``` -->

 ## Build from source

You can build and install the precompiled `starport` binary into `$GOBIN`.

When building from source, your `$GOPATH` environment variable must be set correctly.

To set `$GOPATH` environment variable, run this command:

```
mkdir ~/go
export GOPATH=~/go
```

To install the precompiled `starport` binary, run this command:

```
git clone https://github.com/tendermint/starport --depth=1
cd starport && make
```
