# Packaging and Distributing Ignite

Ignite CLI is distributed on multiple platforms and package managers. This document describes how to package and distribute Ignite CLI.

## Flatpak

Read the folowing resources to understand Flatpak.

* <https://docs.flathub.org/docs/category/for-app-authors>
* <https://docs.flatpak.org/en/latest/first-build.html>

```bash
cd packaging/flatpak
sudo apt install flatpak-builder
flatpak install org.freedesktop.Platform//23.08 org.freedesktop.Sdk//23.08 org.freedesktop.Sdk.Extension.golang//23.08
flatpak-builder build-dir com.ignite.Ignite.yml --force-clean
flatpak-builder --user --install--force-clean --repo=repo build-dir com.ignite.Ignite.yml
flatpak run com.ignite.Ignite
```

The Flatpak is published at <https://flathub.org/apps/com.ignite.Ignite>.
The update process is done manually at <https://github.com/flathub/com.ignite.Ignite> at the moment. At each release, edit `com.ignite.Ignite.yml` and the metainfo file to update the version.

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

Read the following resources to understand HomeBrew.

* <https://docs.brew.sh/>
* <https://docs.brew.sh/Formula-Cookbook>

```bash
HOMEBREW_NO_INSTALL_FROM_API=1 brew install --interactive ignite
brew audit --new-formula ignite
```

The formula is published in the [homebrew-core](https://github.com/homebrew/homebrew-core) repository: <https://github.com/Homebrew/homebrew-core/pull/161938>
