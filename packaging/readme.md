# Packaging and Distributing Ignite

Ignite CLI is distributed on multiple platforms and package managers. This document describes how to package and distribute Ignite CLI.

## Snap

Read the folowing resources to understand Snap.

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

A [github action](../.github/workflows/release-binary.yml) is used to build and publish the Snap at each release.

## HomeBrew

TBD.
