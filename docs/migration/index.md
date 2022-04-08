---
order: 1
title: v0.19.2
parent:
  title: Migration
  order: 3
description: For chains that were scaffolded with Ignite CLI versions lower than v0.19.2, changes are required to use Ignite CLI v0.19.2. 
---

# Upgrading a Blockchain to use Ignite CLI v0.19.2

Ignite CLI _v0.19.2_ comes with IBC _v2.0.2_. To migrate your chain scaffold with _v0.19_ or _v0.18_ version of Ignite CLI, apply changes introduced in PR [#1975](https://github.com/ignite-hq/cli/pull/1975/files) to your chain.

With _v0.19.2_, contents of `tendermint/spm` moved to Starport's own repo. Upgrade your chain's `go.mod` by removing `tendermint/spm` and adding _v0.19.2_ version of `ignite-hq/cli`.
