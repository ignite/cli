---
order: 2
description: Steps to install Starport on your local computer.
---

# Install Starport

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

## Prerequisite

Starport is written in the Go programming language. To use Starport on a local installation, Go must be installed and running:

- [Golang >=1.16](https://golang.org/)

## Upgrading Your Starport Installation

Before you install a new version of Starport, remove all existing Starport installations. 

To remove the current Starport installation:

1. On your terminal window, press Ctrl C to stop the chain that you started with `starport serve`.
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

### Adding `~/go/bin` to `$PATH`

After Starport builds the chain, the resulting binary is installed in `~/go/bin`.  For Starport to access the `bin` subdirectory that contains the executable binary, you must set the path to this location. To set the PATH, add the following line to your shell config file (for example, `~/.bashrc`):

```
export $PATH=~/go/bin:$PATH
```

Apply the changes for the current terminal:

```
source ~/.bashrc
```

Your path is not set correctly if you get the following error: `: exec: “appd”: executable file not found in $PATH`.  Be sure to follow these steps. To use a different location and for more details, see [The GOPATH environment variable](https://golang.org/doc/gopath_code#GOPATH). 



### Write Permission

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
