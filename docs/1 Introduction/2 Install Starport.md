# Install Starport on a local computer

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

## Prerequisites

Local Starport installation requires the follow software be installed and running:

- [Golang >=1.16](https://golang.org/)

  <!-- is 1.16 correct? or revert back to 1.14? -->

   Starport and Cosmos SDK modules are written in Go. Version 1.16 or later is required.

- [Protocol Buffer compiler](https://grpc.io/docs/protoc-installation/)

  <!-- purpose of this compiler? we can be helpful and tell them why we are asking them to install protobuf -->
- [Node.js >=12.19.0](https://nodejs.org/)

  Node.js is used to build the welcome screen, block explorer, and run the web scaffold. Version 12.1.9.0 or later is required.

## Installing Starport with cURL

To install Starport with cURL:

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`.

To install previous versions of the precompiled binary, see [star-port installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### Write permission

The installation requires writer permission to the `/usr/local/bin/` directory. If the installation fails because the user does not have write permission to `/usr/local/bin/`, run the following command:

```
curl https://get.starport.network/starport | bash
```

Then run this command to move the `starport` executable to `/usr/local/bin/`:

```
sudo mv starport /usr/local/bin/
```

<!-- looks like we need to update <https://github.com/allinbits/starport-installer/blob/master/README.md> that says installing with Homebrew is not supported, I tried and it worked for me -->

 ## Installing Starport on macOS with Homebrew

To install Starport on macOS or Linux with Homebrew:

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

## Clone repo and install using make

To install the precompiled `starport` binary:

```
git clone https://github.com/tendermint/starport --depth=1
cd starport && make
```

<!-- what is different about local install? benefits? -->
