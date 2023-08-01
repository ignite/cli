# Set up your Ignite CLI development environment

To ensure you have a successful experience developing with Ignite CLI, make sure that your local system meets these
technical requirements.

Ignite CLI is supported for the following operating systems:

* GNU/Linux
* macOS
* Windows Subsystem for Linux (WSL)

## Install Go

This installation method removes existing Go installations, installs Go, and sets the environment variables.

1. Install the latest version of [Go](https://golang.org/doc/install).

2. Download the release that is suitable for your system.

3. Follow the installation instructions.

**Note:** We recommend not using brew to install Go.

## Add the Go bin directory to your PATH

Ensure the Go environment variables are [set properly](https://golang.org/doc/gopath_code#GOPATH) on your system. Many
of the initial problems are related to incorrect environment variables.

1. Edit your `~/.bashrc` file and add `export PATH=$PATH:$(go env GOPATH)/bin`.
2. To apply the changes, run `source ~/.bashrc`.

## Remove existing Ignite CLI installations

Before you install a new version of Ignite CLI, remove all existing installations.

1. Remove the Ignite CLI binary with `rm $(which ignite)`

    Depending on your user permissions, run the command with or without `sudo`.

2. Repeat this step until all Ignite CLI installations are removed from your system.

```bash
curl https://get.ignite.com! | bash`
```

See [Install Ignite CLI](./docs/docs/01-welcome/02-install.md).

## Clone the Ignite CLI repo

`git clone --depth=1 git@github.com:ignite/cli.git`

## Run make install

1. After you clone the `cli` repo, change into the root directory `cd cli`.

2. Run `make install`.

## Verify your Ignite CLI version

To verify the version of Ignite CLI that is installed, run `ignite version`.
