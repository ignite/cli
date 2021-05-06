---
order: 3
description: Run Starport CLI using a Docker container.
---

# Run Starport in Docker

You can run Starport CLI inside a Docker container using a Docker CLI without installing the Starport binary directly on your machine.

Running Starport in Docker can be useful if you want to run Starport on an unsupported operating system or experiment with a different version of Starport CLI without installing it.

Docker containers are like virtual machines because they provide an isolated environment to programs that runs inside them. In this case, you can run Starport in an isolated environment.

## Benefits of running Starport in a Docker container

Experimentation impact is limited to the Docker instance:

- Changes to the file system by adding, editing, and removing files do not have any effect on the host machine.

  - When a new file is created by Starport inside a Docker container, this file is not present on the host machine.

  - The new file is not be persisted in the Docker container.

  - All file system changes inside a container are erased when the container stops running. In this case, when the `docker run` command finishes.

## Network isolation

Isolation is not just for files, the network is also isolated. In a typical installation, Starport runs some servers that are accessible through HTTP and TCP. However, When Starport runs inside a Docker container, servers are not accessible using these methods.

Running Starport in Docker requires a different set of Docker flags that enable Starport to work with the host machines outside of the Docker container.

## Prerequisites

- Docker must be installed. See [Get Started with Docker](https://www.docker.com/get-started).
- A Docker container must be running.

## Scaffolding a chain

After Docker is confirmed to be installed and running, you can build a blockchain.

To scaffold a blockchain `planet` in a new directory, run this command in a terminal window:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps starport/cli:develop app github.com/hello/planet
```

This command creates a container that runs from the `starport/cli` image and executes the Starport binary inside the image. In this example, the `starport/cli:develop` image is on the `develop` branch so you can experiment with the next version.

After the execution is complete, the container is stopped and deleted and `docker run` command is released.

- `-v $HOME/sdh:/home/tendermint` maps the `$HOME/sdh` directory in the host machine to the `/home/tendermint` directory inside the container.

## Starting a blockchain

To start the blockchain node in development, run this command:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps -p 1317:1317 -p 26657:26657 starport/cli:develop serve -p planet
```

This command does the following:

- `serve -p planet` specifies to use the `planet` directory that contains the source code of the blockchain.
- `-p 26657:26657` exposes ports from the container to the host machine.
- After the blockchain is started, open `http://localhost:26657` to see the Tendermint API.
- The `-v` flag specifies for the container to access the application's source code from the host machine so it can build and run it.

## Versioning

You can specify which version of Starport to run in your Docker container.

### Latest

- By default, `starport/cli` resolves to `starport/cli:latest`. The `latest` image tag is the current release.
- `latest` always points to the latest stable [Starport release](https://github.com/tendermint/starport/releases).

For example, if latest release is [v0.15.1](https://github.com/tendermint/starport/releases/tag/v0.15.1), the `latest` tag points to the `0.15.1` tag.

### Specific version

You can set the version tag to use a specific version. All available tags are in the [starport/cli image](https://hub.docker.com/repository/docker/starport/cli/tags?page=1&ordering=last_updated) on Docker Hub.

For example:

- Use `starport/cli:0.15.1` (without the `v` prefix) to use version 0.15.1.
- Use `starport/cli:develop` points to the `develop` branch of Starport so you can experiment with the next version.

To get the latest image, run `docker pull`.

```
docker pull starport/cli:develop
```
