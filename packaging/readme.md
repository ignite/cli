# Packaging and Distributing Ignite

Ignite CLI is distributed on multiple platforms and package managers. This document describes how to package and distribute Ignite CLI.

## Snap

Read the following resources to understand Snap.

* <https://snapcraft.io/docs/snapcraft-overview>
* <https://snapcraft.io/docs/go-plugin#heading--core22>
* <https://snapcraft.io/docs/go-applications>

```bash
cd packaging/snap
sudo apt install snapd
sudo snap install multipass
multipass launch
SNAPCRAFT_BUILD_ENVIRONMENT=multipass snapcraft
snap install ignite_0.0.0_amd64.snap --dangerous --classic
```

For building, from snapcraft 8, use `SNAPCRAFT_REMOTE_BUILD_STRATEGY=force-fallback snapcraft remote-build`.

A [github action](../.github/workflows/release-binary.yml) is used to build and publish the Snap at each release.

## HomeBrew

Read the following resources to understand HomeBrew.

* <https://docs.brew.sh/>
* <https://docs.brew.sh/Formula-Cookbook>

```bash
HOMEBREW_NO_INSTALL_FROM_API=1 brew install --interactive ignite
brew audit --new-formula ignite
```

The formula is published in the [homebrew-core](https://github.com/homebrew/homebrew-core) repository: <https://github.com/Homebrew/homebrew-core/pull/161938>
