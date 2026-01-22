# Packaging and Distributing Ignite

Ignite CLI is distributed on package managers. This document describes how to package and distribute Ignite CLI.

## HomeBrew

Read the following resources to understand HomeBrew.

- <https://docs.brew.sh/>
- <https://docs.brew.sh/Formula-Cookbook>

```bash
HOMEBREW_NO_INSTALL_FROM_API=1 brew install --interactive ignite
brew audit --new-formula ignite
```

The formula is published in the [homebrew-core](https://github.com/homebrew/homebrew-core) repository: <https://github.com/Homebrew/homebrew-core/pull/161938>
