---
order: 3
description: Run Starport CLI using a Docker container.
---

# Run Starport in Docker

You can run Starport CLI inside a Docker container using a Docker CLI without installing the Starport binary directly on your machine.

Running Starport in Docker can be useful for various reasons; isolating your test environment, running Starport on an unsupported operating system, or experimenting with a different version of Starport CLI without installing it.

Docker containers are like virtual machines because they provide an isolated environment to programs that runs inside them. In this case, you can run Starport in an isolated environment.

Experimentation and file system impact is limited to the Docker instance. The host machine is not impacted by changes to the container.

## Prerequisites

Docker must be installed. See [Get Started with Docker](https://www.docker.com/get-started).

## Scaffolding a Chain

When Docker is installed, you can build a blockchain with a single command.

To scaffold a blockchain `planet` in a new directory, run this command in a terminal window:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps starport/cli:develop app github.com/hello/planet
```

Be patient, this command takes a minute or two to run because it does everything for you:

- Creates a container that runs from the `starport/cli` image.
- Executes the Starport binary inside the image. In this example, the `starport/cli:develop` image is on the `develop` branch so you can experiment with the next version.
- Uses `-v $HOME/sdh:/home/tendermint` to map the `$HOME/sdh` directory in your local computer (the host machine) to the `/home/tendermint` directory inside the container.

## Starting a Blockchain

To start the blockchain node in the Docker container you just created, run this command:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps -p 1317:1317 -p 26657:26657 starport/cli:develop serve -p planet
```

This command does the following:

- `serve -p planet` specifies to use the `planet` directory that contains the source code of the blockchain.
- `-p 26657:26657` maps port 26657 on the host machine to port 26657 in Docker. This mapping exposes ports from the container to the host machine.
- After the blockchain is started, open `http://localhost:26657` to see the Tendermint API.
- The `-v` flag specifies for the container to access the application's source code from the host machine so it can build and run it.

## Versioning

You can specify which version of Starport to install and run in your Docker container.

### Latest Version

- By default, `starport/cli` resolves to `starport/cli:latest`.
- The `latest` image tag is always the latest stable [Starport release](https://github.com/tendermint/starport/releases).

For example, if latest release is [v0.15.1](https://github.com/tendermint/starport/releases/tag/v0.15.1), the `latest` tag points to the `0.15.1` tag.

### Specific Version

You can set the version tag to use a specific version. All available tags are in the [starport/cli image](https://hub.docker.com/repository/docker/starport/cli/tags?page=1&ordering=last_updated) on Docker Hub.

For example:

- Use `starport/cli:0.15.1` (without the `v` prefix) to use version 0.15.1.
- Use `starport/cli:develop` points to the `develop` branch of Starport so you can experiment with the next version.

To get the latest image, run `docker pull`.

```
docker pull starport/cli:develop
```
