---
order: 2
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

### Adding `~/go/bin` to `$PATH`

After Starport builds the chain, the resulting binary is installed in `~/go/bin`. For Starport to be able to access this binary the path to this location has to be set. To do so, add the following line to your shell config file (for example, `~/.bashrc`):

```
export $PATH=~/go/bin:$PATH
```

And apply the changes for the current terminal:

```
source ~/.bashrc
```

Please, follow the steps above if you get the following error: `: exec: “appd”: executable file not found in $PATH`.

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
