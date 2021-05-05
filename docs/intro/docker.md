---
order: 3
description: Run Starport CLI using a Docker container.
---

# Run Starport in Docker

You can run Starport CLI inside a Docker container using a Docker CLI without installing the Starport binary directly on your machine. This may be useful if you want to run Starport on an unsupported operating system or experiment with a different version of Starport CLI without installing it.

Docker containers are like virtual machines and they provide an isolated environment to programs that runs inside them (in our case it is Starport).

Changes to file system (adding, editing, removing files) inside a container does not have any effect on the host machine (user's machine). If starport creates a file inside a container, this file will not present on the host machine and it'll not be persisted in the container as well. All filesystem changes inside a container will be erased when container stops running. In our case, when `docker run` command finishes.

Isolation is not just for files, network is also isolated. Starport runs some servers that are accessible through http and tcp but when it's running inside the Docker container, host machine cannot reach to them. e.g. when you try to access to an http server from a browser on your host machine, you'll not be able to reach to the website.

Thus, we need to use different set of Docker flags when using different Starport commands so, their effects can be observed by the host machines as well.

Make sure you have [Docker installed](https://www.docker.com/get-started).

## Scaffolding a chain

The following command will scaffold a blockchain `planet` in a new directory:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps starport/cli:develop app github.com/hello/planet
```

This command creates a container that runs from `starport/cli` image and executes Starport binary inside the image. When the execution is done, container is stopped and deleted and `docker run` command is released.

`-v $HOME/sdh:/home/tendermint ` maps `$HOME/sdh` directory in the host machine to the `/home/tendermint` directory inside the container.

## Starting a blockchain

The following command will start the blockchain node in development:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps -p 1317:1317 -p 26657:26657 starport/cli:develop serve -p planet
```

`serve -p planet` tells Starport to use the `planet` directory, which contains the source code of the blockchain. `-p 26657:26657` exposes ports from the container to the host machine. Once the blockchain has started, open `http://localhost:26657` to see Tendermint's API. We need to use `-v` flag, because we want the container to access application's source code from the host machine so it can build and run it. 

### Versioning

`starport/cli` resolves to `starport/cli:latest` (`latest` is the image's tag name). And `latest` always points to the latest stable [release](https://github.com/tendermint/starport/releases) of Starport. E.g. if latest release is [v0.15.1](https://github.com/tendermint/starport/releases/tag/v0.15.1), `latest` tag will be pointing to the `0.15.1` tag.

Set the version tag to use a spesific version. e.g. `starport/cli:0.15.1` (without the `v` prefix). All available tags can be found here: https://hub.docker.com/repository/docker/starport/cli/tags?page=1&ordering=last_updated .

`starport/cli:develop` points to te `develop` branch of Starport. Make sure to run `docker pull` to get the latest image.

```
docker pull starport/cli:develop
```