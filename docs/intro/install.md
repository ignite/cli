---
order: 3
description: Steps to install Starport on your local computer.
---

# Install Starport

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

## Prerequisite

Starport is written in the Go programming language. To use Starport on a local system:

 - Install [Go](https://golang.org/doc/install) (**version 1.16** or higher)
 - Ensure the Go environment variables are [set properly](https://golang.org/doc/gopath_code#GOPATH) on your system
 
**Tip:** If you get the message `appd: command not found` then your Go environment variables might not be set correctly.

## Upgrading Your Starport Installation

Before you install a new version of Starport, remove all existing Starport installations. 

To remove the current Starport installation:

1. On your terminal window, press `Ctrl+C` to stop the chain that you started with `starport chain serve`.
1. Remove the Starport binary with `rm $(which starport)`.
   Depending on your user permissions, run the command with or without `sudo`.
1. Repeat this step until all `starport` installations are removed from your system.

After all existing Starport installations are removed, follow the [Installing Starport with cURL](#installing-starport-with-curl) instructions. For details on version features and changes, see the [changelog.md](https://github.com/tendermint/starport/blob/develop/changelog.md) in the repo. 

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
- Stop the chain and remove existing versions before installing a new version.
