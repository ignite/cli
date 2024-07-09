---
description: Run Ignite CLI using a Docker container.
---

# Running inside a Docker container

You can run Ignite CLI inside a Docker container without installing the Ignite
CLI binary directly on your machine.

Running Ignite CLI in Docker can be useful for various reasons; isolating your
test environment, running Ignite CLI on an unsupported operating system, or
experimenting with a different version of Ignite CLI without installing it.

Docker containers are like virtual machines because they provide an isolated
environment to programs that runs inside them. In this case, you can run Ignite
CLI in an isolated environment.

Experimentation and file system impact is limited to the Docker instance. The
host machine is not impacted by changes to the container.

## Prerequisites

Docker must be installed. See [Get Started with
Docker](https://www.docker.com/get-started).

## Ignite CLI Commands in Docker

After you scaffold and start a chain in your Docker container, all Ignite CLI
commands are available. Just type the commands after `docker run -ti
ignite/cli`. For example:

```bash
docker run -ti ignitehq/cli -h
docker run -ti ignitehq/cli scaffold chain planet
docker run -ti ignitehq/cli chain serve
```

## Scaffolding a chain

When Docker is installed, you can build a blockchain with a single command.

Ignite CLI, and the chains you serve with Ignite CLI, persist some files. When
using the CLI binary directly, those files are located in `$HOME/.ignite` and
`$HOME/.cache`, but in the context of Docker it's better to use a directory
different from `$HOME`, so we use `$HOME/sdh`. This folder should be created
manually prior to the docker commands below, or else Docker creates it with the
root user.

```bash
mkdir $HOME/sdh
```

To scaffold a blockchain `planet` in the `/apps` directory in the container, run
this command in a terminal window:

```bash
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps ignitehq/cli:0.25.2 scaffold chain planet
```

Be patient, this command takes a minute or two to run because it does everything
for you:

- Creates a container that runs from the `ignitehq/cli:0.25.2` image.
- Executes the Ignite CLI binary inside the image.
- `-v $HOME/sdh:/home/tendermint` maps the `$HOME/sdh` directory in your local
  computer (the host machine) to the home directory `/home/tendermint` inside
  the container.
- `-v $PWD:/apps` maps the current directory in the terminal window on the host
  machine to the `/apps` directory in the container. You can optionally specify
  an absolute path instead of `$PWD`.

  Using `-w` and `-v` together provides file persistence on the host machine.
  The application source code on the Docker container is mirrored to the file
  system of the host machine.

  **Note:** The directory name for the `-w` and `-v` flags can be a name other
  than `/app`, but the same directory must be specified for both flags. If you
  omit `-w` and `-v`, the changes are made in the container only and are lost
  when that container is shut down.

## Starting a blockchain

To start the blockchain node in the Docker container you just created, run this
command:

```bash
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps -p 1317:1317 -p 26657:26657 ignitehq/cli:0.25.2 chain serve -p planet
```

This command does the following:

- `-v $HOME/sdh:/home/tendermint` maps the `$HOME/sdh` directory in your local
  computer (the host machine) to the home directory `/home/tendermint` inside
  the container.
- `-v $PWD:/apps` persists the scaffolded app in the container to the host
  machine at current working directory.
- `serve -p planet` specifies to use the `planet` directory that contains the
  source code of the blockchain.
- `-p 1317:1317` maps the API server port (cosmos-sdk) to the host machine to
  forward port 1317 listening inside the container to port 1317 on the host
  machine.
- `-p 26657:26657` maps RPC server port 26657 (tendermint) on the host machine
  to port 26657 in Docker.
- After the blockchain is started, open `http://localhost:26657` to see the
  Tendermint API.
- The `-v` flag specifies for the container to access the application's source
  code from the host machine, so it can build and run it.

## Versioning

You can specify which version of Ignite CLI to install and run in your Docker
container.

### Latest version

- By default, `ignite/cli` resolves to `ignite/cli:latest`.
- The `latest` image tag is always the latest stable [Ignite CLI
  release](https://github.com/ignite/cli/releases).

For example, if latest release is
[v0.25.2](https://github.com/ignite/cli/releases/tag/v0.25.2), the `latest` tag
points to the `0.25.2` tag.

### Specific version

You can specify to use a specific version of Ignite CLI. All available tags are
in the [ignite/cli
image](https://hub.docker.com/r/ignitehq/cli/tags?page=1&ordering=last_updated) on
Docker Hub.

For example:

- Use `ignitehq/cli:0.25.2` (without the `v` prefix) to use version `0.25.2`.
- Use `ignitehq/cli` to use the latest version.
- Use `ignitehq/cli:main` to use the `main` branch, so you can experiment with
  the upcoming version.

To get the latest image, run `docker pull`.

```bash
docker pull ignitehq/cli:main
```
