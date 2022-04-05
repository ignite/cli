---
order: 1
description: Steps to install Starport on your local computer.
---

# Install Starport

You can run [Starport](https://github.com/ignite-hq/cli) in a web-based Gitpod IDE or you can install Starport on your local computer. 


## Prerequisites

Be sure you have met the prerequisites before you install and use Starport.

### Operating systems

Starport is supported for the following operating systems:

- GNU/Linux
- macOS
- Windows Subsystem for Linux (WSL)

### Go

Starport is written in the Go programming language. To use Starport on a local system:

- Install [Go](https://golang.org/doc/install) (**version 1.16** or higher)
- Ensure the Go environment variables are [set properly](https://golang.org/doc/gopath_code#GOPATH) on your system

## Verify your Starport version

To verify the version of Starport you have installed, run the following command:

```sh
starport version
```

## Installing Starport

To install the latest version of the `starport` binary use the following command.

```bash
curl https://get.starport.network/starport! | bash
```

This command invokes `curl` to download the install script and pipes the output to `bash` to perform the installation. The `starport` binary is installed in `/usr/local/bin`.

To learn more or customize the installation process, see the [installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### Write permission

Starport installation requires write permission to the `/usr/local/bin/` directory. If the installation fails because you do not have write permission to `/usr/local/bin/`, run the following command:

```bash
curl https://get.starport.network/starport | bash
```

Then run this command to move the `starport` executable to `/usr/local/bin/`:

```bash
sudo mv starport /usr/local/bin/
```

On some machines, a permissions error occurs:

```bash
mv: rename ./starport to /usr/local/bin/starport: Permission denied
============
Error: mv failed
```

In this case, use sudo before `curl` and before `bash`:

```bash
sudo curl https://get.starport.network/starport! | sudo bash
```

## Upgrading your Starport installation

Before you install a new version of Starport, remove all existing Starport installations.

To remove the current Starport installation:

1. On your terminal window, press `Ctrl+C` to stop the chain that you started with `starport chain serve`.
1. Remove the Starport binary with `rm $(which starport)`.
   Depending on your user permissions, run the command with or without `sudo`.
1. Repeat this step until all `starport` installations are removed from your system.

After all existing Starport installations are removed, follow the [Installing Starport with cURL](#installing-starport-with-curl) instructions. For details on version features and changes, see the [changelog.md](https://github.com/ignite-hq/cli/blob/develop/changelog.md) in the repo.

## Installing Starport on macOS with Homebrew

Using brew to install Starport is supported only for macOS machines without the M1 chip. 

For details on version features and changes, see the [changelog.md](https://github.com/tendermint/starport/blob/develop/changelog.md) in the repo.

## Build from source

To experiment with the source code, you can build from source:

```bash
git clone https://github.com/ignite-hq/cli --depth=1
cd starport && make install
```

## Summary

- Verify the prerequisites.
- To setup a local development environment, install Starport locally on your computer.
- Install Starport by fetching the binary using cURL or by building from source.
- The latest version is installed by default. You can install previous versions of the precompiled `starport` binary.
- Stop the chain and remove existing versions before installing a new version.
