---
order: 3
description: Steps to install Starport on your local computer.
---

# Install Starport

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

## Prerequisite

Starport is written in the Go programming language. To use Starport on a local installation, Go must be installed and running:

- [Golang >=1.16](https://golang.org/)

## Installing Starport with cURL

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`.

To install previous versions of the precompiled `starport` binary or customize the installation process, see [Starport installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### Write permission

Starport installation requires write permission to the `/usr/local/bin/` directory. If the installation fails because you do not have write permission to `/usr/local/bin/`, run the following command:

```
curl https://get.starport.network/starport | bash
```

Then run this command to move the `starport` executable to `/usr/local/bin/`:

```
sudo mv starport /usr/local/bin/
```

### Error while running command `protoc`

If you get errors related to `protoc`, make sure that you followed the [Install pre-compiled binaries (any OS)](https://grpc.io/docs/protoc-installation/#install-pre-compiled-binaries-any-os) instructions:

1. From [github.com/google/protobuf/releases](https://github.com/google/protobuf/releases), manually download the zip file that corresponds to your operating system and computer architecture (`protoc-<version>-<os><arch>.zip`).

2. Unzip the file under `$HOME/.local` or a directory of your choice. For example:

```
$ unzip protoc-3.15.5-linux-x86_64.zip -d $HOME/.local
```

1. Update your environment's `PATH` variable to include the path to the protoc executable. For example:

```
$ export PATH="$PATH:$HOME/.local/bin"
```

## Installing Starport on macOS with Homebrew

```
brew install tendermint/tap/starport
```

## Build from source

```
git clone https://github.com/tendermint/starport --depth=1
cd starport && make install
```

## Summary

- To setup a local development environment, install Starport locally on your computer.
- Install Starport by fetching the binary using cURL, Homebrew, or by building from source.
- The latest version is installed by default. You can install previous versions of the precompiled `starport` binary.
